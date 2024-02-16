package handle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/IUnlimit/perpetua/pkg/deepcopy"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

var upgrader websocket.Upgrader

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
		impl, _ := getForwardImpl()
		if len(impl.AccessToken) != 0 {
			auth := r.Header.Get("AUTHORIZATION")
			if auth != impl.AccessToken {
				log.Error("[Client] Connection verification failed with auth-field: ", auth)
				return
			}
		}

		start = true
		handler = NewHandler(ctx)
		handleSet.Add(handler)
		handler.AddWait()
		log.Infof("[Client] WebSocket connection established on port: %d with path: %s", port, r.URL.Path)
		ClientOnlineStatusChangeEvent(handler, true)

		// heartbeat
		gopool.Go(func() {
			handler.AddWait()
			err := doHeartbeat(handler, conn)
			if err != nil {
				log.Debugf("[Client] Error occurred when heartbeat on port-%d: %v", port, err)
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
		log.Info("[Client] WebSocket connection closed on port: ", port)
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

func CreateNTQQWebSocket() error {
	var handle = NewHandler(context.Background())
	handle.AddWait()
	impl, err := getForwardImpl()
	if err != nil {
		return err
	}

	request, _ := http.NewRequest("GET", "", nil)
	request.Header.Set("AccessToken", impl.AccessToken)
	wsUrl := fmt.Sprintf("ws://%s:%d/%s", impl.Host, impl.Port, impl.Suffix)

	log.Info("[NTQQ] Start connecting to NTQQ websocket: ", wsUrl)
	<-waitNTQQStartup(impl.Host, impl.Port)
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, request.Header)
	if err != nil {
		return err
	}
	log.Info("[NTQQ] Websocket connection successful")
	defer conn.Close()

	// write to NTQQ
	// NTQQ <- perp
	gopool.Go(func() {
		handle.AddWait()
		write2NTQQLoop(handle, conn)
	})

	// read from NTQQ
	// NTQQ -> perp
	err = readFromNTQQLoop(handle, conn)
	if err != nil {
		return err
	}
	return nil
}

// wait NTQQ startup
func waitNTQQStartup(host string, port int) <-chan struct{} {
	done := make(chan struct{})

	gopool.Go(func() {
		for {
			err := utils.CheckPort(host, port, time.Second*1)
			if err != nil {
				log.Debugf("Wait NTQQ startup: %v", err)
				time.Sleep(time.Millisecond * 1000)
				continue
			}
			break
		}
		close(done)
	})

	return done
}

// do heartbeat to client
func doHeartbeat(handler *Handler, conn *websocket.Conn) error {
	impl, err := getForwardImpl()
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
			log.Debug("Heartbeat to client: ", conn.RemoteAddr())
		}
	}
}

func write2NTQQLoop(handle *Handler, conn *websocket.Conn) {
	for {
		if handle.ShouldExit() {
			return
		}
		<-echoMap.Receive
		for _, v := range handleSet.Iterator() {
			handler := v.(*Handler)
			id := handler.GetId()
			echoMap.JustGet(id, func(data global.MsgData) {
				// TODO 断点续传,NTQQ重连尝试 echo赋值错误
				log.Debugf("[NTQQ<-] Write to channel(id: %s) with message: %v", handler.GetId(), data)
				err := conn.WriteJSON(data)
				if err != nil {
					log.Errorf("[NTQQ<-] Channel(id: %s) write to NTQQ failed: %v", handler.GetId(), err)
					handle.WaitExitAll()
				}
			})
		}
	}
}

func readFromNTQQLoop(handle *Handler, conn *websocket.Conn) error {
	for {
		if handle.ShouldExit() {
			return errors.New("exception interrupt")
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("[NTQQ->] Failed to read NTQQ message: %v", err)
			handle.WaitExitAll()
			continue
		}

		var msgData global.MsgData
		err = json.Unmarshal(message, &msgData)
		if err != nil {
			log.Errorf("[NTQQ->] Failed to unmarshal NTQQ message: %s", string(message))
			handle.WaitExitAll()
			continue
		}

		// heartbeat
		if msgData["meta_event_type"] == "heartbeat" {
			global.Heartbeat = msgData
			continue
		}

		// msgData
		uuid, err := globalCache.Append(msgData)
		if err != nil {
			log.Errorf("[NTQQ->] Failed to append global cache: %v", err)
			handle.WaitExitAll()
			continue
		}

		// broadcast message
		receivers := make([]interface{}, 0)
		echo := utils.GetDefault(msgData["echo"], "")
		if len(echo) == 0 { // global
			log.Debug("[NTQQ->] Received global NTQQ message: ", string(message))
			receivers = append(receivers, handleSet.Iterator()...)
		} else { // response
			var id string
			if len(echo) == len(global.EchoPrefix)+2+36 {
				// ${EchoPrefix}#${uuid}#
				log.Debugf("[NTQQ->] Received NTQQ message to specified handler(id: %s)", string(message))
				id = strings.Split(echo, "#")[1]
				msgData["echo"] = ""
			} else {
				// ${EchoPrefix}#${uuid}#client-echo
				matches := global.EchoRegx.FindStringSubmatch(echo)
				if len(matches) != 4 {
					log.Errorf("[NTQQ->] Unable to match handler's(id: %s) echo value: %s", id, echo)
					handle.WaitExitAll()
					continue
				}
				id = matches[2]
				msgData["echo"] = matches[3]
			}
			handler := FindHandler(id)
			if handler == nil {
				log.Error("[NTQQ->] Unknown handler id: ", id)
				handle.WaitExitAll()
				continue
			}
			log.Debugf("[NTQQ->] Received NTQQ message: %s", msgData)
			receivers = append(receivers, handler)
		}
		// when closed, staying dispatch
		for _, v := range receivers {
			handler := v.(*Handler)
			gopool.Go(func() {
				handler.AddMessage(uuid)
				handler.Receive <- true
			})
		}
	}
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

// get first forwardImpl
func getForwardImpl() (*model.Implementation, error) {
	for _, impl := range global.AppSettings.Implementations {
		if impl.Type == "ForwardWebSocket" {
			return impl, nil
		}
	}
	return nil, errors.New("can't find forward websocket impl")
}
