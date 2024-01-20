package main

import (
	perp "github.com/IUnlimit/perpetua/cmd/perpe"
	"github.com/IUnlimit/perpetua/internal/conf"
)

func main() {
	conf.Init()
	perp.Init()
	perp.Login()
	go perp.Start()

}
