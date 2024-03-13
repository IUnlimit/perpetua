package handle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func CreateNTQQWebSocket() error {
	var handle = NewHandler(context.Background())
	handle.AddWait()
	impl, err := utils.GetForwardImpl()
	if err != nil {
		return err
	}

	request, _ := http.NewRequest("GET", "", nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", impl.AccessToken))
	wsUrl := fmt.Sprintf("ws://%s:%d/%s", impl.Host, impl.Port, impl.Suffix)

	log.Info("[NTQQ] Start connecting to NTQQ websocket: ", wsUrl)
	<-utils.WaitNTQQStartup(impl.Host, impl.Port, func(err2 error) {
		log.Debugf("Wait NTQQ startup: %v", err2)
	})
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
			handler.AddMessage(uuid)
		}
	}
}
