package model

type MetaData struct {
	MetaEventType string `json:"meta_event_type"`
	Time          int    `json:"time"`
	SelfID        int64  `json:"self_id"`
	PostType      string `json:"post_type"`
}

type HeartBeat struct {
	MetaData
	Interval        int              `json:"interval"`
	HeartBeatStatus *HeartBeatStatus `json:"status"`
}

type HeartBeatStatus struct {
	AppInitialized bool `json:"app_initialized"`
	AppEnabled     bool `json:"app_enabled"`
	AppGood        bool `json:"app_good"`
	Online         bool `json:"online"`
	Good           bool `json:"good"`
}

// PostType string `json:"post_type"`
