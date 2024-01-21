package model

import "time"

type Config struct {
	Log        *Log
	NTQQImpl   *NTQQImpl
	Http       *Http
	WebSocket  *WebSocket
	ParentPath string
}

type Log struct {
	LogForceNew bool
	LogLevel    string
	LogAging    time.Duration
	LogColorful bool
}

type NTQQImpl struct {
	Update    bool
	ID        int64
	Platform  string
	UpdatedAt time.Time
}

type Http struct {
	Port int
}

type WebSocket struct {
	Timeout time.Duration
}
