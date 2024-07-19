package perp

import (
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/bytedance/gopkg/util/gopool"
)

func Bootstrap() {
	Configure()
	if global.ImplType == model.EMBED {
		gopool.Go(func() {
			Start()
		})
	}
	EnableAgent()
}
