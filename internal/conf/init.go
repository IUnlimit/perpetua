package conf

import (
	"github.com/IUnlimit/perpetua/configs"
	global "github.com/IUnlimit/perpetua/internal"
	log "github.com/sirupsen/logrus"
)

func Init() {
	versionCheck()
	fileName := "config.yml"
	fileFolder := global.ParentPath + "/"
	_, err := LoadConfig(fileName, fileFolder, "yaml", configs.Config, &global.Config)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}
