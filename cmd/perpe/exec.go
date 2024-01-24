package perp

import (
	"fmt"
	"github.com/IUnlimit/perpetua/configs"
	"github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/hook"
	"github.com/IUnlimit/perpetua/internal/logger"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/bytedance/gopkg/util/gopool"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

// Configure NTQQ settings using config.yml
func Configure() {
	config := global.Config
	lgrFolder := global.ParentPath + "/" + global.LgrFolder

	log.Info("Searching Lagrange.OneBot ...")
	err := initLagrange(lgrFolder, config.NTQQImpl.Update)
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
	windows := utils.IsWinPlatform()
	if windows {
		execName += ".exe"
	}
	cmdDir := global.ParentPath
	execPath := global.LgrFolder + execName

	var cmd *exec.Cmd
	if windows {
		// TODO
		cmd = nil
	} else { // unix
		cmd = exec.Command(execPath)
	}

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

func initLagrange(lgrFolder string, update bool) error {
	owner := "LagrangeDev"
	repo := "Lagrange.Core"
	exists := utils.FileExists(lgrFolder + "Lagrange.OneBot.pdb")

	if !exists || update {
		zipPath := global.ParentPath + "/Lagrange.OneBot.zip"
		err := updateNTQQImpl(owner, repo, zipPath, lgrFolder, exists, func(existNow bool) {
			if !existNow {
				log.Fatal("Lagrange.OneBot file can't be found, initialization failed")
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func updateNTQQImpl(owner string, repo string, zipPath string, lgrFolder string, exists bool, callback func(exists bool)) error {
	artifacts, err := hook.GetArtifactsUrls(
		owner, repo,
		map[string]string{
			"per_page": "7",
		})
	if err != nil {
		err = fmt.Errorf("artifact get excption: %v", err)
		return err
	}

	log.Info("Please choose the Lagrange.OneBot software suitable for your platform (send the number before option)")
	for i, artifact := range artifacts {
		fmt.Printf("[%d] %s\n", i, artifact.Name)
	}

	var selectIndex int8
	_, err = fmt.Scanf("%d", &selectIndex)
	if err != nil {
		err = fmt.Errorf("failed to capture input: %v", err)
		return err
	}
	log.Info("Start downloading ...")

	artifact := artifacts[selectIndex]
	// check lgr version
	if exists && artifact.UpdatedAt.Before(global.Config.NTQQImpl.UpdatedAt) {
		log.Info("Lagrange.OneBot is the latest version")
		return nil
	}

	err = hook.GetAuthorizedFile(artifact.ArchiveDownloadURL, zipPath, -1)
	defer os.Remove(zipPath)
	if err != nil {
		err = fmt.Errorf("artifact download excption: %v", err)
		return err
	}

	log.Info("Success download, unzipping ...")
	err = utils.Unzip(zipPath, lgrFolder)
	if err != nil {
		err = fmt.Errorf("artifact unzip excption: %v", err)
		return err
	}

	exists = utils.FileExists(lgrFolder + "Lagrange.OneBot.pdb")
	callback(exists)
	conf.UpdateConfig(artifact)
	return nil
}
