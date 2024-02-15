package handle

import (
	"errors"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/utils"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// enhance onebot API impl

var mx sync.Mutex
var hookMap = map[string]func(global.MsgData, *Handler) (global.MsgData, error){
	"set_restart":                  setRestartHook,
	"set_client_name":              setClientName,
	"send_broadcast_data":          sendBroadcastData,
	"send_broadcast_data_callback": sendBroadcastDataCallback,
}
var emptyParams = make(map[string]interface{}, 0)

// TryTouchEnhanceHook check and handle api if exists
// @return data, exist, error
func TryTouchEnhanceHook(msgData global.MsgData, handler *Handler) (global.MsgData, bool, error) {
	params := msgData["params"]
	if _, ok := params.(map[string]interface{}); !ok {
		if params != nil {
			return nil, false, errors.New(fmt.Sprintf("unknown params field type: %s", params))
		}
		msgData["params"] = emptyParams
	}
	if hook := hookMap[msgData["action"].(string)]; hook != nil {
		data, err := hook(msgData, handler)
		if err != nil {
			return nil, true, err
		}
		return data, true, nil
	}
	return nil, false, nil
}

// reboot onebot instance
func setRestartHook(msgData global.MsgData, _ *Handler) (global.MsgData, error) {
	delay := msgData["params"].(map[string]interface{})["delay"]
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

// set the ws connect name
func setClientName(msgData global.MsgData, handler *Handler) (global.MsgData, error) {
	name := msgData["params"].(map[string]interface{})["name"]
	if _, ok := name.(string); !ok {
		return utils.BuildWSBadResponse(fmt.Sprintf("empty or unsupport name: %s", name), msgData["echo"].(string)), nil
	}

	handler.name = name.(string)
	return utils.BuildWSGoodResponse("ok", msgData["echo"].(string)), nil
}

// broadcast data
func sendBroadcastData(msgData global.MsgData, trigger *Handler) (global.MsgData, error) {
	params := msgData["params"].(map[string]interface{})
	clients := params["clients"]
	if clients == nil {
		clients = make([]interface{}, 0)
	} else if _, ok := clients.([]interface{}); !ok {
		return utils.BuildWSBadResponse(fmt.Sprintf("unknown clients list: %s", clients), msgData["echo"].(string)), nil
	} else if len(clients.([]interface{})) == 0 {
		return utils.BuildWSBadResponse(fmt.Sprintf("empty clients list"), msgData["echo"].(string)), nil
	}

	targets := make([]interface{}, 0)
	for _, v := range handleSet.Iterator() {
		h := v.(*Handler)
		for _, c := range clients.([]interface{}) {
			client := c.(map[string]interface{})
			appId := client["app_id"]
			clientName := client["client_name"]
			if appId != nil && appId == h.id {
				targets = append(targets, h)
			} else if clientName != nil && clientName == h.name {
				targets = append(targets, h)
			}
		}
	}

	id := uuid.NewString()
	data := utils.ReturnsNonNull(params["data"])
	ClientBroadcastEvent(trigger, targets, id, data)
	return utils.BuildWSGoodResponse("ok", msgData["echo"].(string),
		"uuid", id), nil
}

// broadcast data callback
func sendBroadcastDataCallback(msgData global.MsgData, trigger *Handler) (global.MsgData, error) {
	params := msgData["params"].(map[string]interface{})
	client := params["client"]
	if _, ok := client.(map[string]interface{}); !ok {
		return utils.BuildWSBadResponse(fmt.Sprintf("unknown client: %s", client), msgData["echo"].(string)), nil
	}
	id := params["uuid"]
	if _, ok := id.(string); !ok || len(id.(string)) == 0 {
		return utils.BuildWSBadResponse(fmt.Sprintf("broadcast callback api must specify uuid"), msgData["echo"].(string)), nil
	}

	var target interface{}
	appId := client.(map[string]interface{})["app_id"]
	clientName := client.(map[string]interface{})["client_name"]
	for _, v := range handleSet.Iterator() {
		h := v.(*Handler)
		if appId != nil && appId == h.id {
			target = h
		} else if clientName != nil && clientName == h.name {
			target = h
		}
	}
	if target == nil {
		return utils.BuildWSBadResponse(fmt.Sprintf("can't find target client: %s", client), msgData["echo"].(string)), nil
	}

	data := utils.ReturnsNonNull(params["data"])
	ClientBroadcastEventCallback(trigger, target, id.(string), data)
	return utils.BuildWSGoodResponse("ok", msgData["echo"].(string)), nil
}
