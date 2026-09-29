package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestGetExternalReconnectInterval(t *testing.T) {
	var nilImpl *NTQQImpl
	assert.Equal(t, DefaultExternalReconnectInterval, nilImpl.GetExternalReconnectInterval())
	assert.Equal(t, DefaultExternalReconnectInterval, (&NTQQImpl{}).GetExternalReconnectInterval())
	assert.Equal(t, DefaultExternalReconnectInterval, (&NTQQImpl{ExternalReconnectInterval: -time.Second}).GetExternalReconnectInterval())

	var impl NTQQImpl
	err := yaml.Unmarshal([]byte("external-reconnect-interval: 3s"), &impl)
	assert.NoError(t, err)
	assert.Equal(t, 3*time.Second, impl.GetExternalReconnectInterval())
}
