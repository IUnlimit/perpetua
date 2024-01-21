package perp

import (
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/handle"
)

func EnableAgent() {
	var config = conf.Config.Http
	go handle.EnableHttpService(config.Port)

	select {}
}
