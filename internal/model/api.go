package model

type Client struct {
	// 客户端ID
	AppId string `json:"app_id"`
	// 客户端名称
	ClientName string `json:"client_name,omitempty"`
}
