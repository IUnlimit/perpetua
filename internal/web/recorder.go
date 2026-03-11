package web

import (
	"time"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// RecordPacket records a packet to Redis
// direction: "ntqq->client" or "client->ntqq"
func RecordPacket(direction string, handlerID string, clientName string, data global.MsgData) {
	if rdb == nil {
		return
	}

	p := &Packet{
		ID:         uuid.NewString(),
		Timestamp:  time.Now().UnixMilli(),
		Direction:  direction,
		HandlerID:  handlerID,
		ClientName: clientName,
		Data:       data,
	}

	if err := SavePacket(p); err != nil {
		log.Warnf("[Web] Failed to record packet: %v", err)
	}
}
