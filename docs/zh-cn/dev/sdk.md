# SDK 支持

Perpetua 提供 SDK 以简化客户端的接入开发。SDK 封装了 WebSocket 连接管理、API 调用、事件监听等功能。

## 可用 SDK

| 语言 / 平台 | 名称                                                                         | 简介     |
|---------|----------------------------------------------------------------------------|--------|
| Java    | [perpetua-sdk-for-java](https://github.com/IUnlimit/perpetua-sdk-for-java) | 官方社区实现 |

## 接入流程

使用 SDK 接入 Perpetua 的基本流程：

1. 通过 HTTP 请求 `GET /api/get_ws_port` 获取分配的 WebSocket 端口
2. 使用获取的端口建立 WebSocket 连接
3. 通过 `set_client_name` 设置客户端名称（推荐，便于识别和广播定向）
4. 监听事件推送，调用 OneBot API 或 Perpetua 拓展 API

## 开发自定义 SDK

如果您需要为其他语言开发 SDK，需要实现以下核心功能：

- **WebSocket 连接管理**：连接建立、心跳维持、断线重连
- **API 调用**：通过 WebSocket 发送 action 请求，处理 echo 回调
- **事件监听**：处理 OneBot 标准事件和 Perpetua 拓展事件
- **拓展 API 支持**：封装 `set_client_name`、`send_broadcast_data` 等接口

相关协议规范请参考：
- [OneBot 11 标准](https://github.com/botuniverse/onebot-11)
- [拓展 API](https://iunlimit.github.io/perpetua/#/zh-cn/user/enhance-api)
- [拓展 Event](https://iunlimit.github.io/perpetua/#/zh-cn/user/enhance-event)
- [拓展数据类型](https://iunlimit.github.io/perpetua/#/zh-cn/user/enhance-model)
