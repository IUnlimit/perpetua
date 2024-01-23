package conf

import (
	"embed"
	"encoding/json"
	"errors"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
)

// LoadConfig creat and load config, return exists(file)
func LoadConfig(fileName string, fileFolder string, fs embed.FS, config any) (bool, error) {
	filePath := fileFolder + fileName
	exists := utils.FileExists(filePath)
	if !exists {
		log.Warnf("Can't find `%s`, generating default configuration", fileName)
		data, err := fs.ReadFile(fileName)
		if err != nil {
			return false, err
		}
		err = utils.CreateFile(filePath, data)
		if err != nil {
			return false, err
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return exists, err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return exists, err
	}
	return exists, nil
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

	global.Config.NTQQImpl = &model.NTQQImpl{
		ID:        artifact.ID,
		Platform:  platform,
		UpdatedAt: artifact.UpdatedAt,
	}

	// TODO record
	return nil
}
