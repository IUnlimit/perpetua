package configs

import "embed"

//go:embed config.yml
var Config embed.FS

var ConfigFileName = "config.yml"

//go:embed appsettings.json
var AppSettings embed.FS

var AppSettingsFileName = "appsettings.json"
