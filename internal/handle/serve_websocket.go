package handle

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/IUnlimit/perpetua/internal/web"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// ntqqHandshakeTimeout websocket handshake timeout for the NTQQ upstream connection
const ntqqHandshakeTimeout = 10 * time.Second

// ntqqLink the lifecycle of one NTQQ upstream websocket connection
type ntqqLink struct {
	conn *websocket.Conn
	done chan struct{}
	once sync.Once
}

func newNTQQLink(conn *websocket.Conn) *ntqqLink {
	return &ntqqLink{
		conn: conn,
		done: make(chan struct{}),
	}
}

// stop closes the connection and notifies all goroutines bound to it, safe to call multiple times
func (l *ntqqLink) stop() {
	l.once.Do(func() {
		close(l.done)
		_ = l.conn.Close()
	})
}

func (l *ntqqLink) closed() bool {
	select {
	case <-l.done:
		return true
	default:
		return false
	}
}

// CreateNTQQWebSocket connects to the NTQQ upstream and blocks until the connection is lost.
// All goroutines bound to the connection have exited when it returns, so it's safe to call again to reconnect.
func CreateNTQQWebSocket() error {
	var wsUrl string
	var accessToken string
	if global.ImplType == model.EXTERNAL {
		config := global.Config.NTQQImpl
		wsUrl = config.ExternalWebSocket
		accessToken = config.ExternalAccessToken
		log.Info("[NTQQ] Start connecting to external NTQQ websocket: ", wsUrl)
	} else { // EMBED
		impl, err := utils.GetForwardImpl()
		if err != nil {
			return err
		}
		wsUrl = fmt.Sprintf("ws://%s:%d/%s", impl.Host, impl.Port, impl.Suffix)
		accessToken = impl.AccessToken

		log.Info("[NTQQ] Start connecting to NTQQ websocket: ", wsUrl)
		<-utils.WaitNTQQStartup(impl.Host, impl.Port, func(err2 error) {
			log.Debugf("Wait NTQQ startup: %v", err2)
		})
	}

	conn, err := utils.DialWebsocket(wsUrl, accessToken, ntqqHandshakeTimeout)
	if err != nil {
		return err
	}
	log.Info("[NTQQ] Websocket connection successful")

	link := newNTQQLink(conn)
	writerExited := make(chan struct{})

	// write to NTQQ
	// NTQQ <- perp
	gopool.Go(func() {
		defer close(writerExited)
		write2NTQQLoop(link)
	})

	// read from NTQQ
	// NTQQ -> perp
	err = readFromNTQQLoop(conn)
	link.stop()
	// make sure the writer bound to this connection has exited before reconnecting,
	// otherwise it would steal requests and write them to the closed connection
	<-writerExited
	return err
}

func write2NTQQLoop(link *ntqqLink) {
	// flush requests queued while the upstream was disconnected
	flush2NTQQ(link)
	for {
		select {
		case <-link.done:
			return
		case <-echoMap.Receive:
		}
		if link.closed() {
			// keep the data queued, the writer of next connection will flush it
			return
		}
		flush2NTQQ(link)
	}
}

// flush2NTQQ writes all queued client requests to NTQQ
func flush2NTQQ(link *ntqqLink) {
	for _, v := range handleSet.Iterator() {
		id := v.(*Handler).GetId()
		echoMap.JustGet(id, func() bool {
			return !link.closed()
		}, func(data global.MsgData) {
			// TODO 断点续传 echo赋值错误
			log.Debugf("[NTQQ<-] Write to channel(id: %s) with message: %v", id, data)
			// Extract and remove trace_id before sending to NTQQ
			traceID, _ := data["_trace_id"].(string)
			delete(data, "_trace_id")
			// Record packet: perpetua -> NTQQ (outbound on ntqq link, same trace)
			if traceID != "" {
				web.RecordNTQQPacketWithTrace(traceID, "outbound", data)
			} else {
				web.RecordNTQQPacket("outbound", data)
			}
			err := link.conn.WriteJSON(data)
			if err != nil {
				log.Errorf("[NTQQ<-] Channel(id: %s) write to NTQQ failed: %v", id, err)
				link.stop()
			}
		})
	}
}

// readFromNTQQLoop dispatches NTQQ messages and only returns when the connection is broken.
// Errors of a single message are logged and skipped without breaking the connection.
func readFromNTQQLoop(conn *websocket.Conn) error {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("[NTQQ->] Failed to read NTQQ message: %v", err)
			return fmt.Errorf("read from NTQQ: %w", err)
		}

		var msgData global.MsgData
		err = json.Unmarshal(message, &msgData)
		if err != nil {
			log.Errorf("[NTQQ->] Failed to unmarshal NTQQ message: %s", string(message))
			continue
		}

		if msgData["meta_event_type"] == "heartbeat" {
			global.Heartbeat = msgData
			continue
		} else if msgData["meta_event_type"] == "lifecycle" {
			global.Lifecycle = msgData
			continue
		}

		// Record packet: NTQQ -> perpetua (inbound on ntqq link)
		traceID := web.RecordNTQQPacket("inbound", msgData)

		// msgData
		uuid, err := globalCache.Append(msgData)
		if err != nil {
			log.Errorf("[NTQQ->] Failed to append global cache: %v", err)
			continue
		}

		// broadcast message
		receivers := make([]interface{}, 0)
		echo, _ := msgData["echo"].(string)
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
					log.Errorf("[NTQQ->] Unable to match handler's echo value: %s", echo)
					continue
				}
				id = matches[2]
				msgData["echo"] = matches[3]
			}
			handler := FindHandler(id)
			if handler == nil {
				// the client may have disconnected before the response arrived
				log.Warn("[NTQQ->] Unknown handler id, response dropped: ", id)
				continue
			}
			log.Debugf("[NTQQ->] Received NTQQ message: %s", msgData)
			receivers = append(receivers, handler)
		}
		// when closed, staying dispatch
		for _, v := range receivers {
			handler := v.(*Handler)
			// Record packet: perpetua -> client (outbound on client link, same trace)
			web.RecordClientPacket(traceID, "outbound", handler.GetId(), handler.GetName(), msgData)
			handler.AddMessage(uuid)
		}
	}
}
