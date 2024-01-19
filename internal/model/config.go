package model

import "time"

type Config struct {
	NTQQImpl    *NTQQImpl
	ParentPath  string
	LogAging    time.Duration
	LogForceNew bool
	LogColorful bool
	LogLevel    string
}

type NTQQImpl struct {
	Update    bool
	ID        int64
	Platform  string
	UpdatedAt time.Time
}
