package main

import (
	"fmt"
	perp "github.com/IUnlimit/perpetua/cmd/perpe"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/handle"
	"github.com/IUnlimit/perpetua/internal/hook"
	"github.com/IUnlimit/perpetua/internal/logger"
)

func main() {
	printBanner()
	conf.Init()
	hook.Init()
	logger.Init()
	handle.Init()
	perp.Bootstrap()
}

func printBanner() {
	fmt.Println(`	                                   __
	______   _________________   _____/  |_ __ _______
	\____ \_/ __ \_  __ \____ \_/ __ \   __\  |  \__  \
	|  |_> >  ___/|  | \/  |_> >  ___/|  | |  |  // __ \_
	|   __/ \___  >__|  |   __/ \___  >__| |____/(____  /
	|__|        \/      |__|        \/                \/ `)
}
