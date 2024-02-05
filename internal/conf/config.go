package conf

import (
	"embed"
	"encoding/json"
	"errors"
	"github.com/IUnlimit/perpetua/configs"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/IUnlimit/perpetua/internal/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

// LoadConfig creat and load config, return exists(file)
// kind: json / yaml
func LoadConfig(fileName string, fileFolder string, kind string, fs embed.FS, config any) (bool, error) {
	filePath := fileFolder + fileName
	exists := utils.FileExists(filePath)
	if !exists {
		log.Warnf("Can't find `%s`, generating default configuration", fileName)
		data, err := fs.ReadFile(fileName)
		if err != nil {
			return false, err
		}
		err = os.MkdirAll(fileFolder, os.ModePerm)
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

	if kind == "json" {
		err = json.Unmarshal(data, config)
	} else if kind == "yaml" {
		err = yaml.Unmarshal(data, config)
	} else {
		err = errors.New("unknown file type: " + kind)
	}
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

	config := global.Config
	config.NTQQImpl = &model.NTQQImpl{
		ID:        artifact.ID,
		Platform:  platform,
		UpdatedAt: artifact.UpdatedAt,
	}

	bytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	filePath := global.ParentPath + "/" + configs.ConfigFileName
	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
