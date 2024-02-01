package qqimpl

import (
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/hook"
	"github.com/IUnlimit/perpetua/internal/utils"
	log "github.com/sirupsen/logrus"
	"os"
)

func InitLagrange(lgrFolder string, update bool) error {
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
