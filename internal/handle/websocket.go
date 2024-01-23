package handle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
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

// CreateWSInstance create new ws instance
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

		start = true
		handler := NewHandler(ctx)
		handleList = append(handleList, handler)
		log.Info("WebSocket connection established on port: ", port)

		// write goroutine
		// go pool: /api read and wait to respond
		// go pool: /event push events
		gopool.Go(func() {
			handler.AddWait()
			for {
				if handler.ShouldExit() {
					log.Debugf("Websocket weite goroutine on port %d has exited", port)
					break
				}

				// TODO consume and write to conn
				//var resp []byte
				//err = conn.WriteMessage(messageType, resp)

				if err != nil {
					log.Warnf("Failed to write message to WebSocket(port: %d): %v", port, err)
					break
				}
			}
		})

		// read goroutine
		gopool.Go(func() {
			handler.AddWait()
			for {
				if handler.ShouldExit() {
					log.Debugf("Websocket read goroutine on port %d has exited", port)
					break
				}

				mType, message, err := conn.ReadMessage()
				if err != nil {
					log.Warnf("Failed to read message from WebSocket(port: %d): %v", port, err)
					break
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
	defer conn.Close()

	log.Info("NTQQ websocket connection successful")
	for {
		// NTQQ读取消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("Failed to read NTQQ message: %v", err)
			continue
		}

		log.Debug("Received NTQQ message: ", string(message))
		var event *model.HeartBeat
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Errorf("Failed to unmarshal NTQQ message: %s", string(message))
			continue
		}

		// heartbeat
		if event.MetaEventType == "heartbeat" {
			global.Status = event.HeartBeatStatus
			continue
		}

		// event
		err = globalCache.Append(&event.MetaData, message)
		if err != nil {
			log.Errorf("Failed to append global cache: %v", err)
			continue
		}
	}
}

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

func getForwardImpl() (*model.Implementation, error) {
	for _, impl := range global.AppSettings.Implementations {
		if impl.Type == "ForwardWebSocket" {
			return impl, nil
		}
	}
	return nil, errors.New("can't find forward websocket impl")
}
