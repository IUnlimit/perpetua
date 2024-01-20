package perp

import (
	"github.com/IUnlimit/perpetua/internal/hook"
	"github.com/IUnlimit/perpetua/internal/logger"
)

func Init() {
	hook.Init()
	logger.Init()
}
