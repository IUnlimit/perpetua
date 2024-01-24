package global

import "github.com/IUnlimit/perpetua/internal/model"

// MsgData websocket message data type
type MsgData map[string]interface{}

// ParentPath 配置文件目录路径
const ParentPath = "perpetua/"

// LgrFolder lgr文件存放路径
const LgrFolder = "Lagrange.OneBot/"

// Config perpetua config.yml
var Config *model.Config

// AppSettings Lagrange.OneBot appsettings.json
var AppSettings *model.AppSettings

// Heartbeat NTQQ heartbeat status
var Heartbeat *MsgData
