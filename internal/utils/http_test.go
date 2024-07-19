package utils

import (
	"testing"
	"time"
)

func TestCheckWebsocket(t *testing.T) {
	ws := "ws://127.0.0.1:5700/onebot/v11/ws"
	err := CheckWebsocket(ws, 1*time.Second)
	if err != nil {
		panic(err)
	}
}
