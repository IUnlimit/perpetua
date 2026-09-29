package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestCreateWSInstance(t *testing.T) {
	conf.Init()
	//CreateWSInstance(25565)
}

func TestUnmarshall(t *testing.T) {
	msg := "{\"interval\":5000,\"status\":{\"app_initialized\":true,\"app_enabled\":true,\"app_good\":true,\"online\":true,\"good\":true},\"meta_event_type\":\"heartbeat\",\"time\":1705851315,\"self_id\":3012218237,\"post_type\":\"meta_event\"}"
	_ = json.Unmarshal([]byte(msg), &global.Heartbeat)
	fmt.Println(global.Heartbeat)

	selfID := strconv.Itoa(int(global.Heartbeat["self_id"].(float64)))
	fmt.Println(selfID)
}

// mockUpstream a forward websocket upstream that records received requests
type mockUpstream struct {
	server   *httptest.Server
	received chan global.MsgData
	conns    chan *websocket.Conn
}

func newMockUpstream(t *testing.T) *mockUpstream {
	m := &mockUpstream{
		received: make(chan global.MsgData, 16),
		conns:    make(chan *websocket.Conn, 1),
	}
	upgrader := websocket.Upgrader{}
	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		_ = conn.WriteJSON(global.MsgData{"post_type": "meta_event", "meta_event_type": "lifecycle", "self_id": 10001})
		m.conns <- conn
		for {
			var data global.MsgData
			if err := conn.ReadJSON(&data); err != nil {
				return
			}
			m.received <- data
		}
	}))
	return m
}

func (m *mockUpstream) url() string {
	return "ws" + strings.TrimPrefix(m.server.URL, "http")
}

func runNTQQWebSocket() <-chan error {
	result := make(chan error, 1)
	go func() {
		result <- CreateNTQQWebSocket()
	}()
	return result
}

// TestNTQQReconnect upstream disconnect must make CreateNTQQWebSocket return, and the next connection
// must receive the client requests (no stale writer of the previous connection steals them)
func TestNTQQReconnect(t *testing.T) {
	echoMap = NewEchoMap()
	globalCache = NewCache(time.Minute)
	global.ImplType = model.EXTERNAL
	global.Config = &model.Config{NTQQImpl: &model.NTQQImpl{}}

	handler := NewHandler(context.Background())
	handleSet.Add(handler)
	defer handleSet.Remove(handler)

	// first upstream: disconnect right after connected
	first := newMockUpstream(t)
	defer first.server.Close()
	global.Config.NTQQImpl.ExternalWebSocket = first.url()
	result := runNTQQWebSocket()
	conn := <-first.conns
	_ = conn.Close()
	select {
	case err := <-result:
		assert.Error(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("CreateNTQQWebSocket did not return after upstream disconnected")
	}

	// client request arrives while the upstream is disconnected
	sent := make(chan struct{})
	echoMap.JustPut(handler.GetId(), global.MsgData{"action": "send_private_msg", "echo": "queued"})
	go func() {
		echoMap.Receive <- true
		close(sent)
	}()

	// second upstream: queued request must be flushed to it
	second := newMockUpstream(t)
	defer second.server.Close()
	global.Config.NTQQImpl.ExternalWebSocket = second.url()
	result = runNTQQWebSocket()
	conn = <-second.conns
	select {
	case data := <-second.received:
		assert.Equal(t, "queued", data["echo"])
	case <-time.After(5 * time.Second):
		t.Fatal("queued request was not flushed to the reconnected upstream")
	}
	select {
	case <-sent:
	case <-time.After(5 * time.Second):
		t.Fatal("client notification was not consumed by the new writer")
	}

	// request after reconnected
	echoMap.JustPut(handler.GetId(), global.MsgData{"action": "get_login_info", "echo": "live"})
	echoMap.Receive <- true
	select {
	case data := <-second.received:
		assert.Equal(t, "live", data["echo"])
	case <-time.After(5 * time.Second):
		t.Fatal("request was not written to the reconnected upstream")
	}
	assert.Empty(t, first.received)

	_ = conn.Close()
	select {
	case <-result:
	case <-time.After(5 * time.Second):
		t.Fatal("CreateNTQQWebSocket did not return after upstream disconnected")
	}
}

func TestBuild(t *testing.T) {
	response := utils.BuildWSResponse("ok", 0, "echo",
		"uuid", "7657")
	fmt.Println(response)
}
