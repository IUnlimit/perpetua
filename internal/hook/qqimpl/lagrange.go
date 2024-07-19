package qqimpl

import (
	"fmt"
	"os"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/hook"
	"github.com/IUnlimit/perpetua/internal/utils"
	log "github.com/sirupsen/logrus"
)

func InitLagrange(lgrFolder string, update bool) error {
	owner := "LagrangeDev"
	repo := "Lagrange.Core"
	exists := utils.FileExists(lgrFolder)

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
			"per_page": "32",
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
	if exists {
		if artifact.UpdatedAt.Before(global.Config.NTQQImpl.Update.UpdatedAt) {
			log.Info("Lagrange.OneBot is the latest version")
			return nil
		}
		log.Infof("Pulled the latest Lagrange.OneBot arch, time: %s", artifact.UpdatedAt.String())
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

	exists = utils.FileExists(lgrFolder)
	callback(exists)
	err = conf.UpdateConfig(artifact)
	if err != nil {
		return err
	}
	return nil
}
