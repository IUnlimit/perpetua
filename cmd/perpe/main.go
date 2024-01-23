package perp

import "github.com/bytedance/gopkg/util/gopool"

func Bootstrap() {
	Configure()
	gopool.Go(func() {
		Start()
	})
	EnableAgent()
}
