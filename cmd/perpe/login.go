package perp

import (
	"fmt"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/hook"
	"github.com/IUnlimit/perpetua/internal/logger"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

func Login() {
	config := conf.Config
	lgrFolder := config.ParentPath + "/" + conf.LgrFolder

	log.Info("Searching Lagrange.OneBot ...")
	err := initLagrange(lgrFolder, config)
	if err != nil {
		log.Fatalf("Lagrange.OneBot init error %v", err)
	}

	fileFolder := config.ParentPath + "/"
	exists := utils.FileExists(fileFolder + "appsettings.json")
	if !exists {
		log.Warn("Can't find `appsettings.json`, generating default configuration (See https://github.com/LagrangeDev/Lagrange.Core?tab=readme-ov-file#appsettingsjson-example)")
		err = conf.UpdateLgrConfig(fileFolder)
		if err != nil {
			log.Fatalf("Failed to update lgr config %v", err)
		}
		log.Info("Default configuration has been generated, please configure and restart perpetua")
		os.Exit(0)
	}
}

func Start() {
	log.Info("Lagrange.OneBot starting ...")
	err := runExec()
	if err != nil {
		log.Fatalf("file instance create failed %v", err)
	}
	log.Info("Lagrange.OneBot start success")
}

func runExec() error {
	execName := "Lagrange.OneBot"
	windows := utils.IsWinPlatform()
	if windows {
		execName += ".exe"
	}
	cmdDir := conf.Config.ParentPath
	execPath := conf.LgrFolder + execName

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
	go func() {
		_, err := io.Copy(in, os.Stdin)
		if err != nil {
			log.Fatalf("Failed to copy stdin %v", err)
		}
	}()

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

func initLagrange(lgrFolder string, config *model.Config) error {
	owner := "LagrangeDev"
	repo := "Lagrange.Core"
	exists := utils.FileExists(lgrFolder + "Lagrange.OneBot.pdb")

	if !exists || config.NTQQImpl.Update {
		zipPath := config.ParentPath + "/Lagrange.OneBot.zip"
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
	if exists && artifact.UpdatedAt.Before(conf.Config.NTQQImpl.UpdatedAt) {
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


