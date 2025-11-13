package handle

import (
	"time"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/model"
	log "github.com/sirupsen/logrus"
)

// enhance event

func ClientOnlineStatusChangeEvent(trigger *Handler, online bool) {
	broadcast(trigger, handleSet.Iterator(), false, func(target *Handler) global.MsgData {
		return global.MsgData{
			"time":        time.Now().UnixMilli(),
			"self_id":     global.Heartbeat["self_id"],
			"post_type":   "notice",
			"notice_type": "client_status",
			"client": &model.Client{
				AppId:      trigger.GetId(),
				ClientName: trigger.GetName(),
			},
			"online":      online,
			"self_client": trigger == target,
		}
	})
}

func ClientBroadcastEvent(trigger *Handler, targets []interface{}, uuid string, data string) {
	broadcast(trigger, targets, true, func(target *Handler) global.MsgData {
		return global.MsgData{
			"time":             time.Now().UnixMilli(),
			"self_id":          global.Heartbeat["self_id"],
			"post_type":        "distributed",
			"distributed_type": "broadcast",
			"client": &model.Client{
				AppId:      trigger.GetId(),
				ClientName: trigger.GetName(),
			},
			"uuid": uuid,
			"data": data,
		}
	})
}

func ClientBroadcastEventCallback(trigger *Handler, target interface{}, uuid string, data string) {
	broadcast(trigger, []interface{}{target}, true, func(target *Handler) global.MsgData {
		return global.MsgData{
			"time":             time.Now().UnixMilli(),
			"self_id":          global.Heartbeat["self_id"],
			"post_type":        "distributed",
			"distributed_type": "broadcast_callback",
			"client": &model.Client{
				AppId:      trigger.GetId(),
				ClientName: trigger.GetName(),
			},
			"uuid": uuid,
			"data": data,
		}
	})
}

func broadcast(trigger *Handler, targets []interface{}, jumpTrigger bool, eventSuppiler func(target *Handler) global.MsgData) {
	for _, v := range targets {
		handler := v.(*Handler)
		msgData := eventSuppiler(handler)

		uuid, err := globalCache.Append(msgData)
		if err != nil {
			log.Errorf("[Enhance] Failed to append global cache: %v", err)
			trigger.WaitExitAll()
			return
		}

		if jumpTrigger && handler.GetId() == trigger.GetId() {
			continue
		}
		log.Debugf("[Enhance] Broadcast to channel-%s with event: %v", handler.GetId(), msgData)
		handler.AddMessage(uuid)
	}
}
