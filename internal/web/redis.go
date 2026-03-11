package web

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	global "github.com/IUnlimit/perpetua/internal"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var rdb *redis.Client
var ctx = context.Background()

// Packet represents a captured data packet
type Packet struct {
	ID         string                 `json:"id"`
	Timestamp  int64                  `json:"timestamp"`
	Direction  string                 `json:"direction"` // "ntqq->client" | "client->ntqq"
	HandlerID  string                 `json:"handler_id"`
	ClientName string                 `json:"client_name"`
	Data       map[string]interface{} `json:"data"`
}

// InitRedis initializes the Redis client
func InitRedis() error {
	cfg := global.Config.Redis
	if cfg == nil {
		log.Warn("[Web] Redis config not found, using defaults")
		rdb = redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   0,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		})
	}
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}
	log.Info("[Web] Redis connection established")
	return nil
}

// SavePacket saves a packet record to Redis
func SavePacket(p *Packet) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	pipe := rdb.Pipeline()

	// Store packet data with expiration
	expire := 24 * time.Hour
	if global.Config.Web != nil {
		expire = global.Config.Web.PacketExpire
	}
	key := fmt.Sprintf("perpetua:packet:%s", p.ID)
	pipe.Set(ctx, key, data, expire)

	// Add to handler's packet sorted set (score = timestamp)
	handlerKey := fmt.Sprintf("perpetua:handler_packets:%s", p.HandlerID)
	pipe.ZAdd(ctx, handlerKey, redis.Z{
		Score:  float64(p.Timestamp),
		Member: p.ID,
	})
	pipe.Expire(ctx, handlerKey, expire)

	// Add to global packet sorted set
	pipe.ZAdd(ctx, "perpetua:all_packets", redis.Z{
		Score:  float64(p.Timestamp),
		Member: p.ID,
	})
	pipe.Expire(ctx, "perpetua:all_packets", expire)

	_, err = pipe.Exec(ctx)
	return err
}

// GetPacketsByHandler returns packets for a specific handler with pagination
func GetPacketsByHandler(handlerID string, offset, limit int64) ([]*Packet, int64, error) {
	handlerKey := fmt.Sprintf("perpetua:handler_packets:%s", handlerID)

	total, err := rdb.ZCard(ctx, handlerKey).Result()
	if err != nil {
		return nil, 0, err
	}

	// Get packet IDs in reverse chronological order
	ids, err := rdb.ZRevRange(ctx, handlerKey, offset, offset+limit-1).Result()
	if err != nil {
		return nil, 0, err
	}

	return getPacketsByIDs(ids, total)
}

// GetAllPackets returns all packets with pagination
func GetAllPackets(offset, limit int64) ([]*Packet, int64, error) {
	total, err := rdb.ZCard(ctx, "perpetua:all_packets").Result()
	if err != nil {
		return nil, 0, err
	}

	ids, err := rdb.ZRevRange(ctx, "perpetua:all_packets", offset, offset+limit-1).Result()
	if err != nil {
		return nil, 0, err
	}

	return getPacketsByIDs(ids, total)
}

func getPacketsByIDs(ids []string, total int64) ([]*Packet, int64, error) {
	if len(ids) == 0 {
		return []*Packet{}, total, nil
	}

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = fmt.Sprintf("perpetua:packet:%s", id)
	}

	results, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, 0, err
	}

	packets := make([]*Packet, 0, len(results))
	for _, r := range results {
		if r == nil {
			continue
		}
		var p Packet
		if err := json.Unmarshal([]byte(r.(string)), &p); err != nil {
			continue
		}
		packets = append(packets, &p)
	}

	return packets, total, nil
}

// CleanupExpiredPackets removes expired packet references from sorted sets
func CleanupExpiredPackets() {
	expire := 24 * time.Hour
	if global.Config.Web != nil {
		expire = global.Config.Web.PacketExpire
	}
	cutoff := float64(time.Now().Add(-expire).UnixMilli())

	// Clean global set
	removed, _ := rdb.ZRemRangeByScore(ctx, "perpetua:all_packets", "-inf", fmt.Sprintf("%f", cutoff)).Result()
	if removed > 0 {
		log.Infof("[Web] Cleaned up %d expired packet references", removed)
	}

	// Clean handler-specific sets
	iter := rdb.Scan(ctx, 0, "perpetua:handler_packets:*", 100).Iterator()
	for iter.Next(ctx) {
		rdb.ZRemRangeByScore(ctx, iter.Val(), "-inf", fmt.Sprintf("%f", cutoff))
	}
}

// StartCleanupTask starts a periodic cleanup goroutine
func StartCleanupTask() {
	interval := time.Hour
	if global.Config.Web != nil && global.Config.Web.CleanupInterval > 0 {
		interval = global.Config.Web.CleanupInterval
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			CleanupExpiredPackets()
		}
	}()
	log.Infof("[Web] Packet cleanup task started (interval: %v)", interval)
}
