package handle

import (
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/bluele/gcache"
	"github.com/google/uuid"
	"time"
)

var globalCache *Cache

type Cache struct {
	cache gcache.Cache
}

func NewCache(expireTime time.Duration) *Cache {
	cache := gcache.New(1024).Simple().Expiration(expireTime).Build()
	return &Cache{
		cache: cache,
	}
}

func (c *Cache) Append(data global.MsgData) (string, error) {
	id := uuid.NewString()
	err := c.cache.Set(id, data)
	if err != nil {
		return id, err
	}
	return id, nil
}
