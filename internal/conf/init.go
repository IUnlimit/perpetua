package conf

import (
	"github.com/IUnlimit/perpetua/configs"
	global "github.com/IUnlimit/perpetua/internal"
	log "github.com/sirupsen/logrus"
)

func Init() {
	versionCheck()
	fileFolder := global.ParentPath + "/"
	_, err := LoadConfig(configs.ConfigFileName, fileFolder, "yaml", configs.Config, &global.Config)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Info("Current perpetua instance version: ", Version)
}
