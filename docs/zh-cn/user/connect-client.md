# 连接客户端

Perpetua 与上游 NTQQ 实现建立连接后，下游用户客户端可通过以下三种方式接入 Perpetua：

```
NTQQ(Lagrange.OneBot) ←──上游──→ Perpetua ←──下游──→ 用户客户端
                                              ├── 正向 WebSocket（客户端主动连接）
                                              ├── 反向 WebSocket（Perpetua 主动连接）
                                              └── HTTP Post（Perpetua 主动推送）
```

## 正向 WebSocket

客户端主动连接到 Perpetua。这是最基本的接入方式。

### 配置

```yaml
# http 相关配置（提供 get_ws_port 等接口）
http:
  port: 8080

# websocket 相关配置
web-socket:
  # ws监听最长等待时间
  timeout: 15s
  # 指定范围 [start, end] 内随机监听端口
  range-port:
    # 是否开启功能
    enabled: false
    # 起始端口
    start: 8100
    # 终止端口
    end: 8110
```

### 连接流程

1. **获取 ws 端口**

    使用 `GET` 请求访问 `/api/get_ws_port`（端口为 `http.port`，默认 8080）。Perpetua 采取动态端口的监听方式来避免端口冲突和配置文件过度耦合等问题，以提高服务的可伸缩性

2. **建立 WebSocket 连接**

    根据获取的端口，与 Perpetua 建立 WebSocket 连接。目前仅支持通过 ['/' 接口](https://github.com/botuniverse/onebot-11/blob/master/communication/ws.md#-%E6%8E%A5%E5%8F%A3) 同时进行事件监听与 OneBot API 操作

> 若开启了 `range-port`，Perpetua 会在指定范围内随机选择可用端口；否则使用系统随机分配的端口。超过 `timeout` 时间无客户端连入，该端口自动释放。

## 反向 WebSocket

Perpetua 主动向指定地址发起 WebSocket 连接。适用于 nonebot 等需要被动接收连接的框架。

### 配置

```yaml
# 反向 websocket 相关配置
# 支持配置多个上报地址，Perpetua 会同时连接所有地址
reverse-web-socket:
  - url: 'ws://127.0.0.1:8800/onebot/v11/ws/'
    access-token: ''
```

| 配置项            | 说明                          |
|----------------|-----------------------------|
| `url`          | 上报地址，Perpetua 主动连接的目标 ws 地址 |
| `access-token` | 鉴权 token（可选）                |

### 连接行为

- Perpetua 启动后会自动连接到所有配置的反向 WebSocket 地址
- 连接时携带请求头：`X-Self-ID`（QQ号）、`X-Client-Role: Universal`、`Authorization: Bearer {token}`
- 断线后自动重连

### 具体接入示例

- [第三方支持 - nonebot](zh-cn/user/support-nonebot.md)

## HTTP Post

Perpetua 主动向指定地址推送事件。适用于仅需接收事件推送、不需要主动调用 API 的场景。

### 配置

```yaml
# http post 相关配置
# 支持配置多个上报地址，Perpetua 会同时向所有地址推送事件
http-post:
  - url: 'http://127.0.0.1:9000/callback'
    secret: ''
```

| 配置项      | 说明                                    |
|----------|---------------------------------------|
| `url`    | 事件上报地址                              |
| `secret` | HMAC 签名密匙，用于验证上报数据的完整性（可选） |

### 推送行为

- Perpetua 通过 HTTP POST 将事件数据以 JSON 格式推送到上报地址
- 请求头包含：`Content-Type: application/json`、`X-Self-ID`（QQ号）
- 上报地址的响应体可作为 API 调用回传

## 多端接入

三种连接方式可同时使用，且反向 WebSocket 和 HTTP Post 均支持配置多个地址。Perpetua 会将 NTQQ 的事件同时广播到所有已连接的客户端：

```yaml
# 多个正向 WS 客户端可各自调用 get_ws_port 获取独立端口
# 多个反向 WS 地址
reverse-web-socket:
  - url: 'ws://127.0.0.1:8800/onebot/v11/ws/'
    access-token: ''
  - url: 'ws://127.0.0.1:8801/onebot/v11/ws/'
    access-token: 'your-token'

# 多个 HTTP Post 地址
http-post:
  - url: 'http://127.0.0.1:9000/callback'
    secret: ''
  - url: 'http://127.0.0.1:9001/callback'
    secret: 'your-secret'
```
