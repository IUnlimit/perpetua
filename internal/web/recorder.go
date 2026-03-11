package web

import (
	"time"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// RecordNTQQPacket records a packet on the NTQQ <-> perpetua link
// direction: "inbound" (ntqq -> perpetua) or "outbound" (perpetua -> ntqq)
func RecordNTQQPacket(direction string, data global.MsgData) {
	recordPacket("ntqq", direction, "", "", data)
}

// RecordClientPacket records a packet on the perpetua <-> client link
// direction: "inbound" (client -> perpetua) or "outbound" (perpetua -> client)
func RecordClientPacket(direction string, handlerID string, clientName string, data global.MsgData) {
	recordPacket("client", direction, handlerID, clientName, data)
}

func recordPacket(link, direction, handlerID, clientName string, data global.MsgData) {
	if rdb == nil {
		return
	}

	p := &Packet{
		ID:         uuid.NewString(),
		Timestamp:  time.Now().UnixMilli(),
		Link:       link,
		Direction:  direction,
		HandlerID:  handlerID,
		ClientName: clientName,
		Data:       data,
	}

	if err := SavePacket(p); err != nil {
		log.Warnf("[Web] Failed to record packet: %v", err)
	}
}
