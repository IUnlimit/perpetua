package handle

import (
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/bytedance/gopkg/util/gopool"
	log "github.com/sirupsen/logrus"
	"time"
)

// enhance event

func ClientOnlineStatusChangeEvent(trigger *Handler, online bool) {
	msgData := global.MsgData{
		"time":        time.Now().UnixMilli(),
		"self_id":     global.Heartbeat["self_id"],
		"post_type":   "notice",
		"notice_type": "client_status",
		"client": &model.Client{
			AppId:      trigger.GetId(),
			ClientName: trigger.GetName(),
		},
		"online": online,
	}
	broadcast(trigger, true, msgData)
}

func broadcast(trigger *Handler, jumpTrigger bool, msgData global.MsgData) {
	// msgData
	uuid, err := globalCache.Append(msgData)
	if err != nil {
		log.Errorf("[Enhance] Failed to append global cache: %v", err)
		trigger.WaitExitAll()
		return
	}

	log.Debugf("[Enhance] Broadcast event: %v", msgData)
	for _, v := range handleSet.Iterator() {
		handler := v.(*Handler)
		if jumpTrigger && handler.GetId() == trigger.GetId() {
			continue
		}
		gopool.Go(func() {
			handler.AddMessage(uuid)
			handler.Receive <- true
		})
	}
}
