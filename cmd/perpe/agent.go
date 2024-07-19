package perp

import (
	"time"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/handle"
	"github.com/bytedance/gopkg/util/gopool"
	log "github.com/sirupsen/logrus"
)

func EnableAgent() {
	gopool.Go(func() {
		config := global.Config.Http
		handle.EnableHttpService(config.Port)
	})

	gopool.Go(func() {
		reverseWSList := global.Config.ReverseWebSocket
		if len(reverseWSList) == 0 {
			return
		}
		for i, reverseWS := range reverseWSList {
			if len(reverseWS.Url) == 0 {
				log.Infof("Empty reverse-websocket config ignored with index: %d", i)
				continue
			}
			gopool.Go(func() {
				for {
					err := handle.TryReverseWebsocket(reverseWS.Url, reverseWS.AccessToken)
					if err != nil {
						log.Infof("Attempt to establish http post connection failed: %v", err)
					}
					time.Sleep(1 * time.Second)
				}
			})
		}
	})

	gopool.Go(func() {
		httpPostList := global.Config.HttpPost
		if len(httpPostList) == 0 {
			return
		}
		for i, httpPost := range httpPostList {
			if len(httpPost.Url) == 0 {
				log.Infof("Empty http-post config ignored with index: %d", i)
				continue
			}
			gopool.Go(func() {
				for {
					err := handle.TryPostHttp(httpPost.Url, httpPost.Secret)
					if err != nil {
						log.Infof("Attempt to establish http post connection failed: %v", err)
					}
					time.Sleep(1 * time.Second)
				}
			})
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
