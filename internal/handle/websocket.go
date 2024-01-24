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
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// allow all sources conn
		return true
	},
}

// CreateWSInstance create new ws instance and wait client connection
func CreateWSInstance(port int) {
	var start bool
	ctx := context.Background()
	handleWebSocket := func(w http.ResponseWriter, r *http.Request) {
		// upgrade HTTP to WebSocket conn
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorf("Failed to upgrade connection to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		// check http header: AUTHORIZATION
		impl, _ := getForwardImpl()
		if len(impl.AccessToken) != 0 {
			auth := r.Header.Get("AUTHORIZATION")
			if auth != impl.AccessToken {
				log.Error("Connection verification failed with auth-field: ", auth)
				return
			}
		}

		start = true
		handler := NewHandler(ctx)
		handler.AddWait()
		handleList = append(handleList, handler)
		log.Infof("WebSocket connection established on port: %d with path: %s", port, r.URL.Path)

		// heartbeat
		gopool.Go(func() {
			handler.AddWait()
			err = doHeartbeat(conn, handler)
			if err != nil {
				log.Debugf("Error occurred when heartbeat on port-%d: %v", port, err)
			}
		})

		// write to client
		// perp -> client
		gopool.Go(func() {
			handler.AddWait()
			for {
				if handler.ShouldExit() {
					log.Debugf("Websocket write goroutine on port %d has exited", port)
					break
				}

				handler.GetMessage(func(data *map[string]interface{}) {
					gopool.Go(func() {
						if handler.ShouldExit() {
							return
						}
						err := conn.WriteJSON(data)
						if err != nil {
							log.Warnf("Failed to write message to WebSocket(port: %d): %v", port, err)
							handler.WaitExitAll()
						}
					})
				})

				<-handler.receive
			}
		})

		// read from client
		// perp <- client
		gopool.Go(func() {
			handler.AddWait()
			for {
				if handler.ShouldExit() {
					log.Debugf("Websocket read goroutine on port %d has exited", port)
					return
				}

				mType, message, err := conn.ReadMessage()
				if err != nil {
					log.Warnf("Failed to read message from WebSocket(port: %d): %v", port, err)
					handler.WaitExitAll()
					return
				}

				// TODO read and write to cache
				log.Debugf("Received message(type: %d) on port-%d: %s", mType, port, string(message))
			}
		})

		handler.WaitDone()
		log.Info("WebSocket connection closed on port: ", port)
	}

	server := &http.Server{
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

	gopool.Go(func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Warnf("WebSocket connection(port: %d) exit with error: %v", port, err)
		}
	})
}

func CreateNTQQWebSocket() error {
	impl, err := getForwardImpl()
	if err != nil {
		return err
	}

	request, _ := http.NewRequest("GET", "", nil)
	request.Header.Set("AccessToken", impl.AccessToken)
	wsUrl := fmt.Sprintf("ws://%s:%d/%s", impl.Host, impl.Port, impl.Suffix)

	log.Info("Start connecting to NTQQ websocket: ", wsUrl)
	<-waitNTQQStartup(impl.Host, impl.Port)
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, request.Header)
	if err != nil {
		return err
	}
	log.Info("NTQQ websocket connection successful")
	defer conn.Close()

	// write to NTQQ
	// NTQQ <- perp
	gopool.Go(func() {
		//for {
		//	// TODO
		//}
	})

	// read from NTQQ
	// NTQQ -> perp
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("Failed to read NTQQ message: %v", err)
			continue
		}

		log.Debug("Received NTQQ message: ", string(message))
		var event map[string]interface{}
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Errorf("Failed to unmarshal NTQQ message: %s", string(message))
			continue
		}

		// heartbeat
		if event["meta_event_type"] == "heartbeat" {
			status := event /*["status"].(map[string]interface{})*/
			global.Heartbeat = &status
			continue
		}

		// event
		uuid, err := globalCache.Append(event)
		if err != nil {
			log.Errorf("Failed to append global cache: %v", err)
			continue
		}
		// broadcast message
		for _, handler := range handleList {
			handler.AddMessage(uuid)
			handler.receive <- true
		}
	}
}

// wait NTQQ startup
func waitNTQQStartup(host string, port int) <-chan struct{} {
	done := make(chan struct{})

	gopool.Go(func() {
		for {
			err := utils.CheckPort(host, port, time.Second*1)
			if err != nil {
				log.Debugf("Wait NTQQ startup: %v", err)
				time.Sleep(time.Millisecond * 500)
				continue
			}
			break
		}
		close(done)
	})

	return done
}

// do heartbeat to client
func doHeartbeat(conn *websocket.Conn, handler *Handler) error {
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

		heartbeat := deepcopy.Copy(global.Heartbeat).(*map[string]interface{})
		(*heartbeat)["time"] = time.Now().UnixMilli()
		err := conn.WriteJSON(heartbeat)
		if err != nil {
			return err
		}

		select {
		case <-ticker.C:
			log.Debug("Heartbeat to client: ", conn.RemoteAddr())
		}
	}
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
