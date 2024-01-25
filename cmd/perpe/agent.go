package perp

import (
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/handle"
	"github.com/IUnlimit/perpetua/internal/handle/api"
	"github.com/bytedance/gopkg/util/gopool"
	log "github.com/sirupsen/logrus"
)

func EnableAgent() {
	var config = global.Config.Http
	gopool.Go(func() {
		api.EnableHttpService(config.Port)
	})

	err := handle.CreateNTQQWebSocket()
	if err != nil {
		log.Fatalf("Failed to connect to NTQQ websocket: %v", err)
	}
}
