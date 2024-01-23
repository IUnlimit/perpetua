package global

import (
	"github.com/IUnlimit/perpetua/internal/model"
)

const LgrFolder = "Lagrange.OneBot/"

// Config perpetua config.yml
var Config *model.Config

// AppSettings Lagrange.OneBot appsettings.json
var AppSettings *model.AppSettings

// Status NTQQ heartbeat status
var Status *model.HeartBeatStatus
