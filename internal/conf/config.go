package conf

import (
	"errors"
	"github.com/IUnlimit/perpetua/configs"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	"regexp"
	"time"
)

const LgrFolder = "Lagrange.OneBot/"

var Config *model.Config

func Init() {
	Config = &model.Config{
		ParentPath:  "perpetua/",
		LogAging:    time.Hour * 24,
		LogForceNew: false,
		LogColorful: true,
		LogLevel:    "info",
		NTQQImpl: &model.NTQQImpl{
			Update: false,
		},
	}
}

// UpdateLgrConfig update appsettings.json
func UpdateLgrConfig(fileFolder string) error {
	fileName := "appsettings.json"
	data, err := configs.AppSettings.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = utils.CreateFile(fileFolder+fileName, data)
	if err != nil {
		return err
	}
	return nil
}

// UpdateConfig update config.yml
func UpdateConfig(artifact *model.Artifact) error {
	platform := ""

	regx := regexp.MustCompile(`_(\w+)-`)
	match := regx.FindStringSubmatch(artifact.Name)
	if len(match) > 1 {
		platform = match[1]
	} else {
		return errors.New("can't match platform, artifact name: " + artifact.Name)
	}

	Config.NTQQImpl = &model.NTQQImpl{
		ID:        artifact.ID,
		Platform:  platform,
		UpdatedAt: artifact.UpdatedAt,
	}

	// TODO record
	return nil
}
