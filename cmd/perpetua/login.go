package perp

import (
	"errors"
	"fmt"
	"github.com/IUnlimit/perpetua/internal/conf"
	"github.com/IUnlimit/perpetua/internal/hook"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
)

func Login() {
	config := conf.Config
	lgrFolder := config.ParentPath + "/" + conf.LgrFolder

	log.Info("Searching Lagrange.OneBot ...")
	err := initLagrange(lgrFolder, config)
	if err != nil {
		log.Fatalf("Lagrange.OneBot init error %v", err)
	}

	exists := utils.FileExists("appsettings.json")
	if !exists {
		log.Warn("Can't find `appsettings.json`, generating default configuration (See https://github.com/LagrangeDev/Lagrange.Core?tab=readme-ov-file#appsettingsjson-example)")
		err = conf.UpdateLgrConfig()
		if err != nil {
			log.Fatalf("Failed to update lgr config %v", err)
		}
		log.Info("Default configuration has been generated, please configure and restart perpetua")
		os.Exit(0)
	}

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
	updateConfig(artifact)
	return nil
}

func updateConfig(artifact *model.Artifact) error {
	platform := ""

	regx := regexp.MustCompile(`_(\w+)-`)
	match := regx.FindStringSubmatch(artifact.Name)
	if len(match) > 1 {
		platform = match[1]
	} else {
		return errors.New("can't match platform, artifact name: " + artifact.Name)
	}

	conf.Config.NTQQImpl = &model.NTQQImpl{
		ID:        artifact.ID,
		Platform:  platform,
		UpdatedAt: artifact.UpdatedAt,
	}

	// TODO record
	return nil
}
