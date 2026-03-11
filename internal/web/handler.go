package web

import (
	"net/http"
	"strconv"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/gin-gonic/gin"
)

// ConnectionsProvider is a callback to get connections without importing handle package
var ConnectionsProvider func() []map[string]interface{}

// GetConnections returns all active connections
func GetConnections(c *gin.Context) {
	var connections []map[string]interface{}
	if ConnectionsProvider != nil {
		connections = ConnectionsProvider()
	} else {
		connections = make([]map[string]interface{}, 0)
	}
	c.JSON(http.StatusOK, model.Response{
		Status:  "ok",
		RetCode: 0,
		Data:    connections,
	})
}

// GetPackets returns packets with optional handler filter and pagination
func GetPackets(c *gin.Context) {
	handlerID := c.Query("handler_id")
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "50")

	offset, _ := strconv.ParseInt(offsetStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var packets []*Packet
	var total int64
	var err error

	if handlerID != "" {
		packets, total, err = GetPacketsByHandler(handlerID, offset, limit)
	} else {
		packets, total, err = GetAllPackets(offset, limit)
	}

	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Status:  "failed",
			RetCode: -1,
			Msg:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  "ok",
		RetCode: 0,
		Data: gin.H{
			"packets": packets,
			"total":   total,
			"offset":  offset,
			"limit":   limit,
		},
	})
}

// GetSystemInfo returns system overview info
func GetSystemInfo(c *gin.Context) {
	var connections []map[string]interface{}
	if ConnectionsProvider != nil {
		connections = ConnectionsProvider()
	} else {
		connections = make([]map[string]interface{}, 0)
	}

	// Get total packet count from Redis
	var total int64
	if rdb != nil {
		total, _ = rdb.ZCard(ctx, "perpetua:all_packets").Result()
	}

	info := gin.H{
		"connections_count": len(connections),
		"packets_count":     total,
		"lifecycle":         global.Lifecycle,
		"heartbeat":         global.Heartbeat,
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  "ok",
		RetCode: 0,
		Data:    info,
	})
}

// DeletePackets deletes packets older than specified timestamp
func DeletePackets(c *gin.Context) {
	beforeStr := c.Query("before")
	if beforeStr == "" {
		c.JSON(http.StatusOK, model.Response{
			Status:  "failed",
			RetCode: -1,
			Msg:     "missing 'before' parameter (unix timestamp in ms)",
		})
		return
	}

	before, err := strconv.ParseFloat(beforeStr, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Status:  "failed",
			RetCode: -1,
			Msg:     "invalid 'before' parameter",
		})
		return
	}

	if rdb == nil {
		c.JSON(http.StatusOK, model.Response{
			Status:  "failed",
			RetCode: -1,
			Msg:     "Redis not initialized",
		})
		return
	}

	// Remove from global set
	removed, _ := rdb.ZRemRangeByScore(ctx, "perpetua:all_packets", "-inf", strconv.FormatFloat(before, 'f', -1, 64)).Result()

	// Remove from handler sets
	iter := rdb.Scan(ctx, 0, "perpetua:handler_packets:*", 100).Iterator()
	for iter.Next(ctx) {
		rdb.ZRemRangeByScore(ctx, iter.Val(), "-inf", strconv.FormatFloat(before, 'f', -1, 64))
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  "ok",
		RetCode: 0,
		Data: gin.H{
			"removed": removed,
		},
	})
}
