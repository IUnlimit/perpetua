package perp

import (
	"fmt"
	"github.com/IUnlimit/perpetua/configs"
	"github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/hook/qqimpl"
	"github.com/IUnlimit/perpetua/internal/logger"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/bytedance/gopkg/util/gopool"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

var isWindows bool

// Configure NTQQ settings using config.yml
func Configure() {
	isWindows = utils.IsWinPlatform()
	if isWindows { // TODO add env flag to jump confirm
		confirmShell()
	}

	config := global.Config
	lgrFolder := global.ParentPath + "/" + global.LgrFolder

	log.Info("Searching Lagrange.OneBot ...")
	err := qqimpl.InitLagrange(lgrFolder, config.NTQQImpl.Update)
	if err != nil {
		log.Fatalf("Lagrange.OneBot init error %v", err)
	}

	fileName := "appsettings.json"
	fileFolder := global.ParentPath + "/"
	exists, err := conf.LoadConfig(fileName, fileFolder, "json", configs.AppSettings, &global.AppSettings)
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
	err := runExec()
	if err != nil {
		log.Fatalf("file instance create failed: %v", err)
	}
	log.Info("Lagrange.OneBot start success")
}

func runExec() error {
	execName := "Lagrange.OneBot"
	if isWindows {
		execName += ".exe"
	}
	cmdDir := global.ParentPath
	execPath := global.LgrFolder + execName
	cmd := exec.Command(execPath)

	cmd.Dir = cmdDir
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	// 将错误输出与标准输出连接至同一管道
	cmd.Stderr = cmd.Stdout
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	// 将命令行输入复制到 stdin 管道中
	gopool.Go(func() {
		_, err := io.Copy(in, os.Stdin)
		if err != nil {
			log.Fatalf("Failed to copy stdin: %v", err)
		}
	})

	if err = cmd.Start(); err != nil {
		return err
	}

	var n int
	hook := logger.Hook
	bytes := make([]byte, 8*1024)
	for {
		n, err = out.Read(bytes)
		if err != nil {
			break
		}
		err = hook.ExecLogWrite(string(bytes[:n]))
		if err != nil {
			log.Warnf("Write exec log error: %v", err)
		}
	}

	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func confirmShell() {
	log.Warn("检测到 Windows 环境，请确保使用 powershell 或 bat 脚本运行程序。\n文档：https://iunlimit.github.io/perpetua\n是否继续？(y/n): ")
	var input string
	_, _ = fmt.Scanln(&input)
	if input == "n" || input == "N" {
		os.Exit(0)
	}
}
