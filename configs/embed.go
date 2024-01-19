package configs

import "embed"

//go:embed config.yml
var Config embed.FS

//go:embed appsettings.json
var AppSettings embed.FS
