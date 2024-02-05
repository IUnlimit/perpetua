package model

type Response struct {
	/**
	 * 状态, 表示 API 是否调用成功
	 * - ok	    api 调用成功
	 * - async  api 调用已经提交异步处理, 具体 api 调用是否成功无法得知
	 * - failed api 调用失败
	 * */
	Status string `json:"status"`

	/**
	 * 响应码
	 * - 0   调用成功
	 * - 1   已提交 async 处理
	 * - 其他 操作失败, 具体原因可以看响应的 msg 字段
	 * */
	RetCode int32 `json:"retcode"`

	/**
	 * 错误消息, 仅在 API 调用失败时存在该字段
	 * */
	Msg string `json:"msg,omitempty"`

	/**
	 * 响应数据
	 * - key: 响应数据名
	 * - value: 数据值
	 * */
	Data map[string]any `json:"data,omitempty"`
}
