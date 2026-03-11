package web

import (
	"time"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// RecordNTQQPacket records a packet on the NTQQ <-> perpetua link.
// Returns the trace_id for correlation with downstream packets.
func RecordNTQQPacket(direction string, data global.MsgData) string {
	traceID := uuid.NewString()
	recordPacket(traceID, "ntqq", direction, "", "", data)
	return traceID
}

// RecordNTQQPacketWithTrace records a packet on the NTQQ link with an existing trace_id.
func RecordNTQQPacketWithTrace(traceID string, direction string, data global.MsgData) {
	recordPacket(traceID, "ntqq", direction, "", "", data)
}

// RecordClientPacket records a packet on the perpetua <-> client link.
// traceID ties it to the originating NTQQ-side packet.
func RecordClientPacket(traceID string, direction string, handlerID string, clientName string, data global.MsgData) {
	if traceID == "" {
		traceID = uuid.NewString()
	}
	recordPacket(traceID, "client", direction, handlerID, clientName, data)
}

func recordPacket(traceID, link, direction, handlerID, clientName string, data global.MsgData) {
	if rdb == nil {
		return
	}

	p := &Packet{
		ID:         uuid.NewString(),
		TraceID:    traceID,
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
