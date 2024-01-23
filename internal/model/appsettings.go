package model

type AppSettings struct {
	Logging         *Logging          `json:"Logging"`
	SignServerURL   string            `json:"SignServerUrl"`
	Account         *Account          `json:"Account"`
	Message         *Message          `json:"Message"`
	Implementations []*Implementation `json:"Implementations"`
}

type LogLevel struct {
	Default                  string `json:"Default"`
	Microsoft                string `json:"Microsoft"`
	MicrosoftHostingLifetime string `json:"Microsoft.Hosting.Lifetime"`
}

type Logging struct {
	LogLevel *LogLevel `json:"LogLevel"`
}

type Account struct {
	Uin              int    `json:"Uin"`
	Password         string `json:"Password"`
	Protocol         string `json:"Protocol"`
	AutoReconnect    bool   `json:"AutoReconnect"`
	GetOptimumServer bool   `json:"GetOptimumServer"`
}

type Message struct {
	IgnoreSelf bool `json:"IgnoreSelf"`
}

type Implementation struct {
	Type              string `json:"Type"`
	Host              string `json:"Host"`
	Port              int    `json:"Port"`
	Suffix            string `json:"Suffix"`
	ReconnectInterval int    `json:"ReconnectInterval"`
	HeartBeatInterval int    `json:"HeartBeatInterval"`
	AccessToken       string `json:"AccessToken"`
}
