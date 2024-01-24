package handle

import (
	"github.com/bluele/gcache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExpire(t *testing.T) {
	cache := gcache.New(1024).Simple().Expiration(time.Second * 1).Build()
	cache.Set("test", 666)
	time.Sleep(time.Second * 2)
	get, _ := cache.Get("test")
	assert.True(t, get == nil)
}
