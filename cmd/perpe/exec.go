package perp

import (
	"fmt"
	"github.com/IUnlimit/perpetua/configs"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/hook/qqimpl"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	log "github.com/sirupsen/logrus"
	"os"
)

// Configure NTQQ settings using config.yml
func Configure() {
	isWindows := utils.IsWinPlatform()
	if isWindows && !utils.ContainsArgs("faststart") {
		confirmShell()
	}

	config := global.Config
	lgrFolder := global.ParentPath + "/" + global.LgrFolder
	global.ImplType = model.EMBED

	// check impl type
	lgrWS := config.NTQQImpl.ExternalWebSocket
	if lgrWS != "" {
		<-utils.WaitExternalNTQQStartup(lgrWS, 5, func(alive bool) {
			if alive {
				log.Info("External NTQQ connection successful: ", lgrWS)
				global.ImplType = model.EXTERNAL
			}
		}, func(err2 error) {
			log.Debugf("Wait External-NTQQ startup: %v", err2)
		})

		if global.ImplType == model.EXTERNAL {
			config := global.Config.NTQQImpl
			host, port, suffix, err := utils.ParseWebSocketURL(config.ExternalWebSocket)
			if err != nil {
				log.Fatalf("Parse external websocket url: %v", err)
				return
			}
			global.AppSettings = &model.AppSettings{
				Implementations: []*model.Implementation{
					{
						Type:              "ForwardWebSocket",
						Host:              host,
						Port:              port,
						Suffix:            suffix,
						ReconnectInterval: 5000,
						HeartBeatInterval: 5000,
						AccessToken:       config.ExternalAccessToken,
					},
				},
			}
			return
		}
		log.Warn("External NTQQ connect failed, try to start EMBED")
	}

	log.Info("Searching Lagrange.OneBot ...")
	err := qqimpl.InitLagrange(lgrFolder, config.NTQQImpl.Update.Enable)
	if err != nil {
		log.Fatalf("Lagrange.OneBot init error %v", err)
	}

	fileFolder := global.ParentPath + "/"
	exists, err := conf.LoadConfig(configs.AppSettingsFileName, fileFolder, "json", configs.AppSettings, &global.AppSettings)
	if err != nil {
		log.Fatalf("Failed to load lgr config: %v", err)
	}
	if !exists {
		log.Info("Default `appsettings.json` has been generated, please configure and restart perpetua (See https://github.com/LagrangeDev/Lagrange.Core?tab=readme-ov-file#signserver)")
		os.Exit(0)
	}
}

// Start the exec file
func Start() {
	log.Info("Lagrange.OneBot starting ...")
	err := utils.RunExec(nil)
	if !global.Restart && err != nil {
		log.Errorf("File instance create failed: %v", err)
	}
}

func confirmShell() {
	log.Warn("检测到 Windows 环境，请确保使用 powershell 或 bat 脚本运行程序（使用命令行参数 faststart 即可跳过本提示）。\n文档：https://iunlimit.github.io/perpetua\n是否继续？(y/n): ")
	var input string
	_, _ = fmt.Scanln(&input)
	if input == "n" || input == "N" {
		os.Exit(0)
	}
}
