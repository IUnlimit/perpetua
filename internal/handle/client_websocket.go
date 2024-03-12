package handle

import (
	"context"
	"encoding/json"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/IUnlimit/perpetua/pkg/deepcopy"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var upgrader websocket.Upgrader

// TryReverseWSInstance try to establish reverse ws client connection
func TryReverseWSInstance(wsUrl string, accessToken string) error {
	request, _ := http.NewRequest("GET", "", nil)
	if len(accessToken) != 0 {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}
	log.Infof("[Client] Start connecting to reverse-websocket: %s with headers: %s", wsUrl, request.Header)

	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, request.Header)
	if err != nil {
		return err
	}
	log.Info("[Client] Reverse-websocket connection successful")
	defer conn.Close()

	handler := NewHandler(context.Background())
	configureClientHandler(handler, conn)
	return nil
}

// CreateWSInstance create new ws instance and wait client connection
func CreateWSInstance(port int) {
	var start bool
	var handler *Handler
	var server *http.Server
	ctx := context.Background()
	handleWebSocket := func(w http.ResponseWriter, r *http.Request) {
		// upgrade HTTP to WebSocket conn
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorf("[Client] Failed to upgrade connection to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		// check http header: AUTHORIZATION
		impl, _ := utils.GetForwardImpl()
		if len(impl.AccessToken) != 0 {
			auth := r.Header.Get("AUTHORIZATION")
			if auth != impl.AccessToken {
				log.Errorf("[Client] Connection verification failed with auth-field: %s", auth)
				return
			}
		}

		start = true
		handler = NewHandler(ctx)
		configureClientHandler(handler, conn)
		_ = server.Shutdown(ctx)
	}

	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handleWebSocket),
	}

	// wait to connect
	config := global.Config.WebSocket
	gopool.Go(func() {
		timer := time.After(config.Timeout)
		<-timer
		if !start {
			_ = server.Shutdown(ctx)
		}
	})

	// exit logic
	gopool.Go(func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Warnf("[Client] WebSocket connection(port: %d) exit with error: %v", port, err)
		}
		// remove handler
		handleSet.Remove(handler)
		if start {
			ClientOnlineStatusChangeEvent(handler, false)
		}
	})
}

// configure and link handler to connection
func configureClientHandler(handler *Handler, conn *websocket.Conn) {
	addr := conn.LocalAddr().String()
	addressParts := strings.Split(addr, ":")
	portStr := addressParts[len(addressParts)-1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Errorf("[Client] Failed to obtain the link port, address: %s", addr)
		return
	}

	handleSet.Add(handler)
	handler.AddWait()
	log.Infof("[Client] WebSocket connection established on %s", addr)
	ClientOnlineStatusChangeEvent(handler, true)

	// heartbeat
	gopool.Go(func() {
		handler.AddWait()
		err := doHeartbeat(handler, conn)
		if err != nil {
			log.Debugf("[Client] Error occurred when heartbeat on %s: %v", addr, err)
		}
	})

	// write to client
	// perp -> client
	gopool.Go(func() {
		handler.AddWait()
		write2ClientLoop(handler, conn, port)
	})

	// read from client
	// perp <- client
	gopool.Go(func() {
		handler.AddWait()
		readFromClientLoop(handler, conn, port)
	})

	handler.WaitDone()
	log.Infof("[Client] WebSocket connection closed on port: %d", port)
}

func write2ClientLoop(handler *Handler, conn *websocket.Conn, port int) {
	for {
		if handler.ShouldExit() {
			log.Debugf("[->Client] Websocket write goroutine on port %d has exited", port)
			return
		}

		<-handler.Receive
		handler.GetMessage(func(data global.MsgData) {
			if handler.ShouldExit() {
				return
			}
			log.Debugf("[->Client] Try to send message to client(id-%s, name-%s)", handler.id, handler.name)
			handler.Lock.Lock()
			err := conn.WriteJSON(data)
			handler.Lock.Unlock()
			if err != nil {
				log.Warnf("[->Client] Failed to write message to WebSocket(port: %d): %v", port, err)
				handler.WaitExitAll()
			}
		})
	}
}

func readFromClientLoop(handler *Handler, conn *websocket.Conn, port int) {
	for {
		if handler.ShouldExit() {
			log.Debugf("[<-Client] Websocket read goroutine on port %d has exited", port)
			return
		}

		mType, message, err := conn.ReadMessage()
		if err != nil {
			log.Warnf("[<-Client] Failed to read message from WebSocket(port: %d): %v", port, err)
			handler.WaitExitAll()
			return
		}
		log.Debugf("[<-Client] Received message(type: %d) on port-%d: %s", mType, port, string(message))

		var msgData global.MsgData
		err = json.Unmarshal(message, &msgData)
		if err != nil {
			log.Errorf("[<-Client] Failed to unmarshal client message: %s", string(message))
			continue
		}

		if msgData["echo"] == nil {
			msgData["echo"] = ""
		}

		// intercept hooked requests
		exist, err := interceptHookedRequests(msgData, handler)
		if exist {
			// touch hook method error, no need to break
			if err != nil {
				log.Errorf("[<-Client] Failed to intercept hooked requests: %v", msgData)
			}
			continue
		}

		// sign with echo field
		id := handler.GetId()
		echo := msgData["echo"].(string)
		echo = fmt.Sprintf("%s#%s#%s", global.EchoPrefix, id, echo)
		msgData["echo"] = echo
		log.Debugf("[<-Client] Update client(port-%d) message echo: %s", port, echo)
		echoMap.JustPut(id, msgData)
		echoMap.Receive <- true
	}
}

// @return continue loop
func interceptHookedRequests(msgData global.MsgData, handler *Handler) (bool, error) {
	resp, exist, err := TryTouchEnhanceHook(msgData, handler)
	if err != nil {
		return true, err
	}
	// todo jump able
	if !exist {
		return false, nil
	}

	uuid, err := globalCache.Append(resp)
	if err != nil {
		handler.WaitExitAll()
		return true, err
	}

	handler.AddMessage(uuid)
	handler.Receive <- true
	return true, nil
}

// do heartbeat to client
func doHeartbeat(handler *Handler, conn *websocket.Conn) error {
	impl, err := utils.GetForwardImpl()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Millisecond * time.Duration(impl.HeartBeatInterval))
	defer ticker.Stop()

	for {
		if handler.ShouldExit() {
			return nil
		}

		// update heartbeat time
		heartbeat := deepcopy.Copy(global.Heartbeat).(global.MsgData)
		heartbeat["time"] = time.Now().UnixMilli()
		handler.Lock.Lock()
		err := conn.WriteJSON(heartbeat)
		handler.Lock.Unlock()
		if err != nil {
			return err
		}

		select {
		case <-ticker.C:
			log.Debugf("Heartbeat to client: %s", conn.RemoteAddr().String())
		}
	}
}
