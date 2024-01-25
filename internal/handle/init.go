package handle

import (
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/gorilla/websocket"
	"net/http"
)

func Init() {
	echoMap = NewEchoMap()
	globalCache = NewCache(global.Config.MsgExpireTime)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// allow all sources conn
			return true
		},
	}
}
