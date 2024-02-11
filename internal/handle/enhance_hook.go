package handle

import (
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/bytedance/gopkg/util/gopool"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// enhance onebot API impl

var mx sync.Mutex
var hookMap = map[string]func(global.MsgData) (global.MsgData, error){
	"set_restart": setRestartHook,
}

// TryTouchEnhanceHook check and handle api if exists
// @return data, exist, error
func TryTouchEnhanceHook(msgData global.MsgData) (global.MsgData, bool, error) {
	if hook := hookMap[msgData["action"].(string)]; hookMap != nil {
		data, err := hook(msgData)
		if err != nil {
			return nil, true, err
		}
		return data, true, nil
	}
	return nil, false, nil
}

// reboot onebot instance
func setRestartHook(msgData global.MsgData) (global.MsgData, error) {
	delay := msgData["delay"]
	if _, ok := delay.(int); ok {
		// sleep before locking so that other threads that end delay earlier can seize the lock
		time.Sleep(time.Duration(delay.(int)) * time.Millisecond)
	} // else delay = 0, run directly

	if !mx.TryLock() {
		return utils.BuildWSBadResponse("bot instance is restarting", msgData["echo"].(string)), nil
	}
	global.Restart = true
	process := global.OneBotProcess
	_ = process.Kill()
	global.OneBotProcess = nil
	gopool.Go(func() {
		err := utils.RunExec(&mx)
		if err != nil {
			log.Fatalf("[Enhance] File instance recreate failed: %v", err)
		}
		log.Info("[Enhance] Lagrange.OneBot restart success")
	})

	return utils.BuildWSGoodResponse("async", msgData["echo"].(string)), nil
}
