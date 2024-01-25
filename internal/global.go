package global

import (
	"github.com/IUnlimit/perpetua/internal/model"
	"regexp"
)

// MsgData websocket message data type
type MsgData map[string]interface{}

// ParentPath perp files path
const ParentPath = "config/"

// LgrFolder lgr bin directory
const LgrFolder = "Lagrange.OneBot/"

// EchoPrefix is prefix for generating echos
const EchoPrefix = "perp"

// EchoRegx to match ${EchoPrefix}#${uuid}#client-echo
var EchoRegx = regexp.MustCompile(`([^#]+)#([^#]+)#(.+)`)

// Config perpetua config.yml
var Config *model.Config

// AppSettings Lagrange.OneBot appsettings.json
var AppSettings *model.AppSettings

// Heartbeat NTQQ heartbeat status
var Heartbeat *MsgData
