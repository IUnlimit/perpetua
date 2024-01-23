package handle

import (
	"github.com/IUnlimit/perpetua/internal/model"
	"github.com/allegro/bigcache/v3"
)

// handleList stores handle.Handler for client websocket
var handleList []*Handler

var globalCache *Cache

type Cache struct {
	cache *bigcache.BigCache
}

func NewCache(cache *bigcache.BigCache) *Cache {
	return &Cache{
		cache: cache,
	}
}

func (c *Cache) Append(meta *model.MetaData, entry []byte) error {
	var key string // TODO
	err := c.cache.Append(key, entry)
	if err != nil {
		return err
	}
	return nil
}
