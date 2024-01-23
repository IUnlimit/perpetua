package model

import "time"

type Config struct {
	Log           *Log          `json:"log"`
	NTQQImpl      *NTQQImpl     `json:"ntqq-impl"`
	Http          *Http         `json:"http"`
	WebSocket     *WebSocket    `json:"web-socket"`
	ParentPath    string        `json:"parent-path"`
	MsgExpireTime time.Duration `json:"msg-expire-time"`
}

type Log struct {
	ForceNew bool          `json:"force-new,omitempty"`
	Level    string        `json:"level,omitempty"`
	Aging    time.Duration `json:"aging,omitempty"`
	Colorful bool          `json:"colorful,omitempty"`
}

type NTQQImpl struct {
	Update    bool      `json:"update,omitempty"`
	ID        int64     `json:"id,omitempty"`
	Platform  string    `json:"platform,omitempty"`
	UpdatedAt time.Time `json:"updated-at"`
}

type Http struct {
	Port int `json:"port,omitempty"`
}

type WebSocket struct {
	Timeout time.Duration `json:"timeout,omitempty"`
}
