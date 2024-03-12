package perp

import (
	"fmt"
	"github.com/IUnlimit/perpetua/configs"
	"github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/hook/qqimpl"
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

	log.Info("Searching Lagrange.OneBot ...")
	err := qqimpl.InitLagrange(lgrFolder, config.NTQQImpl.Update)
	if err != nil {
		log.Fatalf("Lagrange.OneBot init error %v", err)
	}

	fileFolder := global.ParentPath + "/"
	exists, err := conf.LoadConfig(configs.AppSettingsFileName, fileFolder, "json", configs.AppSettings, &global.AppSettings)
	if err != nil {
		log.Fatalf("Failed to load lgr config: %v", err)
	}
	if !exists {
		log.Info("Default `appsettings.json` has been generated, please configure and restart perpetua (See https://github.com/LagrangeDev/Lagrange.Core?tab=readme-ov-file#appsettingsjson-example)")
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
