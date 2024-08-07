package global

import (
	"os"
	"regexp"

	"github.com/IUnlimit/perpetua/internal/model"
)

// MsgData websocket message data type
type MsgData map[string]interface{}

// ParentPath perp files path
const ParentPath = "config/"

// LgrFolder lgr bin directory
const LgrFolder = "Lagrange.OneBot/"

// EchoPrefix is prefix for generating echos
const EchoPrefix = "perp"

// The NTQQ impl connection type
var ImplType model.ImplType

// Restart marks whether the end status of the process is restarted
var Restart bool

// OneBotProcess the currently running onebot process
var OneBotProcess *os.Process

// EchoRegx to match ${EchoPrefix}#${uuid}#client-echo
var EchoRegx = regexp.MustCompile(`([^#]+)#([^#]+)#(.+)`)

// Config perpetua config.yml
var Config *model.Config

// AppSettings Lagrange.OneBot appsettings.json
var AppSettings *model.AppSettings

// Lifecycle NTQQ lifecycle metadata
var Lifecycle MsgData

// Heartbeat NTQQ heartbeat status
var Heartbeat MsgData
