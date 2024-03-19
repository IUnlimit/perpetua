package handle

import (
	"context"
	"errors"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// ReadAndWriteClient read loop and write loop
type ReadAndWriteClient struct {
	addr    string
	conn    *websocket.Conn
	handler *Handler
}

func NewReadAndWriteClient(conn *websocket.Conn, handler *Handler) *ReadAndWriteClient {
	addr := conn.LocalAddr().String()
	return &ReadAndWriteClient{
		addr:    addr,
		conn:    conn,
		handler: handler,
	}
}

func (raw *ReadAndWriteClient) writeFunc(data global.MsgData) error {
	return raw.conn.WriteJSON(data)
}

func (raw *ReadAndWriteClient) readFunc() ([]byte, error) {
	_, bytes, err := raw.conn.ReadMessage()
	return bytes, err
}

func (raw *ReadAndWriteClient) getHandler() *Handler {
	return raw.handler
}

func (raw *ReadAndWriteClient) getUrl() string {
	return raw.addr
}

var upgrader websocket.Upgrader

// TryReverseWebsocket try to establish reverse ws client connection
func TryReverseWebsocket(wsUrl string, accessToken string) error {
	impl, err := utils.GetForwardImpl()
	if err != nil {
		return err
	}
	<-utils.WaitNTQQStartup(impl.Host, impl.Port, nil)
	<-utils.WaitCondition(time.Duration(2000), func() error {
		if global.Lifecycle == nil {
			return errors.New("not init yet")
		}
		return nil
	}, nil)

	request, _ := http.NewRequest("GET", wsUrl, nil)
	if len(accessToken) != 0 {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}
	request.Header.Set("X-Self-ID", strconv.Itoa(int(global.Lifecycle["self_id"].(float64))))
	request.Header.Set("X-Client-Role", "Universal")
	log.Infof("[Client] Start try to report events to wsUrl: %s with headers: %s", wsUrl, request.Header)

	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, request.Header)
	if err != nil {
		return err
	}
	log.Info("[Client] Reverse-websocket connection successful")
	defer conn.Close()

	handler := NewHandler(context.Background())
	client := NewReadAndWriteClient(conn, handler)
	ConfigureRAWClientHandler(client)
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
		client := NewReadAndWriteClient(conn, handler)
		ConfigureRAWClientHandler(client)
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
