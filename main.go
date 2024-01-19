package main

import (
	perp "github.com/IUnlimit/perpetua/cmd/perpetua"
	"github.com/IUnlimit/perpetua/internal/conf"
)

func main() {
	conf.Init()
	perp.Init()
	perp.Login()
	perp.Start()

}
