package perp

import (
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/handle"
	"github.com/bytedance/gopkg/util/gopool"
	log "github.com/sirupsen/logrus"
	"time"
)

func EnableAgent() {
	gopool.Go(func() {
		config := global.Config.Http
		handle.EnableHttpService(config.Port)
	})
	gopool.Go(func() {
		config := global.Config.ReverseWebSocket
		if !config.Enabled {
			return
		}
		for {
			err := handle.TryReverseWSInstance(config.Url, config.AccessToken)
			if err != nil {
				log.Infof("Failed to establish reverse websocket connection: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	})

	for {
		err := handle.CreateNTQQWebSocket()
		if err != nil {
			log.Errorf("Failed to connect to NTQQ websocket(: %v), will try again later", err)
		}
		if !global.Restart {
			break
		}
		global.Restart = false
	}

	if global.OneBotProcess != nil {
		_ = global.OneBotProcess.Kill()
	}
}
