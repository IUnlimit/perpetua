package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

const mockToken = "snowluma-token"

// newMockSnowLumaServer mimics the upgrade verification of SnowLuma forward websocket (ws-default):
// path must be "/" (400 otherwise) and Bearer / access_token query must match (401 otherwise)
func newMockSnowLumaServer(t *testing.T) *httptest.Server {
	upgrader := websocket.Upgrader{}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/api" && r.URL.Path != "/event" {
			http.Error(w, "Bad path for WebSocket", http.StatusBadRequest)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+mockToken && r.URL.Query().Get("access_token") != mockToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade failed: %v", err)
			return
		}
		_ = conn.Close()
	}))
}

func wsURL(server *httptest.Server, path string) string {
	return "ws" + strings.TrimPrefix(server.URL, "http") + path
}

func TestCheckWebsocket(t *testing.T) {
	server := newMockSnowLumaServer(t)
	defer server.Close()

	err := CheckWebsocket(wsURL(server, "/"), mockToken, time.Second)
	assert.NoError(t, err)
}

func TestCheckWebsocketWithoutToken(t *testing.T) {
	server := newMockSnowLumaServer(t)
	defer server.Close()

	err := CheckWebsocket(wsURL(server, "/"), "", time.Second)
	assert.ErrorIs(t, err, websocket.ErrBadHandshake)
	assert.Contains(t, err.Error(), "401")
	assert.Contains(t, err.Error(), "external-access-token")
}

func TestCheckWebsocketWrongPath(t *testing.T) {
	server := newMockSnowLumaServer(t)
	defer server.Close()

	err := CheckWebsocket(wsURL(server, "/onebot/v11/ws"), mockToken, time.Second)
	assert.ErrorIs(t, err, websocket.ErrBadHandshake)
	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "Bad path for WebSocket")
}

func TestDialWebsocketUnreachable(t *testing.T) {
	server := newMockSnowLumaServer(t)
	url := wsURL(server, "/")
	server.Close()

	_, err := DialWebsocket(url, mockToken, time.Second)
	assert.Error(t, err)
	assert.NotErrorIs(t, err, websocket.ErrBadHandshake)
}
