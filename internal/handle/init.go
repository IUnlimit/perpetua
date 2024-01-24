package handle

import (
	global "github.com/IUnlimit/perpetua/internal"
)

func Init() {
	globalCache = NewCache(global.Config.MsgExpireTime)
}
