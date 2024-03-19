package model

import "time"

type Config struct {
	Log              *Log                `yaml:"log"`
	NTQQImpl         *NTQQImpl           `yaml:"ntqq-impl"`
	Http             *Http               `yaml:"http"`
	WebSocket        *WebSocket          `yaml:"web-socket"`
	ReverseWebSocket []*ReverseWebSocket `yaml:"reverse-web-socket"`
	HttpPost         []*HttpPost         `yaml:"http-post"`
	MsgExpireTime    time.Duration       `yaml:"msg-expire-time,omitempty"`
}

type Log struct {
	ForceNew bool          `yaml:"force-new,omitempty"`
	Level    string        `yaml:"level,omitempty"`
	Aging    time.Duration `yaml:"aging,omitempty"`
	Colorful bool          `yaml:"colorful,omitempty"`
}

type NTQQImpl struct {
	Update    bool      `yaml:"update,omitempty"`
	ID        int64     `yaml:"id,omitempty"`
	Platform  string    `yaml:"platform,omitempty"`
	UpdatedAt time.Time `yaml:"updated-at"`
}

type Http struct {
	Port int `yaml:"port,omitempty"`
}

type WebSocket struct {
	Timeout   time.Duration `yaml:"timeout,omitempty"`
	RangePort *RangePort    `yaml:"range-port"`
}

type ReverseWebSocket struct {
	Url         string `yaml:"url,omitempty"`
	AccessToken string `yaml:"access-token,omitempty"`
}

type HttpPost struct {
	Url    string `yaml:"url,omitempty"`
	Secret string `yaml:"secret,omitempty"`
}

type RangePort struct {
	Enabled bool `yaml:"enabled,omitempty"`
	Start   int  `yaml:"start,omitempty"`
	End     int  `yaml:"end,omitempty"`
}
