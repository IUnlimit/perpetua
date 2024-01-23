package handle

import (
	"context"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/allegro/bigcache/v3"
)

func Init() {
	config := bigcache.DefaultConfig(global.Config.MsgExpireTime)
	cache, _ := bigcache.New(context.Background(), config)
	globalCache = NewCache(cache)
}
