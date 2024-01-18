package main

import (
	"github.com/IUnlimit/perpetua/cmd/perpetua"
	"github.com/IUnlimit/perpetua/internal/logger"
)

func main() {
	logger.Init()
	perpetua.Login()
}
