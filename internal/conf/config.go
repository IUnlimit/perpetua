package conf

import (
	"github.com/IUnlimit/perpetua/configs"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	"time"
)

const LgrFolder = "Lagrange.OneBot/"

var Config *model.Config

func Init() {
	Config = &model.Config{
		ParentPath:  "perpetua/",
		LogAging:    time.Hour * 24,
		LogForceNew: true,
		LogColorful: true,
		LogLevel:    "info",
		NTQQImpl: &model.NTQQImpl{
			Update: false,
		},
	}
}

func UpdateLgrConfig() error {
	fileName := "appsettings.json"
	data, err := configs.AppSettings.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = utils.CreateFile(fileName, data)
	if err != nil {
		return err
	}
	return nil
}
