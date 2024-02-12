# 拓展 API

!> 下文中的 API 调用分为 `http` 调用与 `websocket` 调用两种，两种方式均需要以 onebot 协议规定的格式发起请求，此处不再追述。详见 [onebot-11/communication](https://github.com/botuniverse/onebot-11/tree/master/communication)

## 新增支持 - http 调用

> `新增支持 API` 是 `perpetua` 基于自身业务场景额外提供的 API 接口，大部分是为提供分布式服务治理支持所服务，故如您仅有单机需求时，只需关注少数关键 `新增支持 API` 即可。下文中 `响应数据` 字段为 `基础响应数据` 中 `data` 片段数据，其 `基础响应数据` 格式如下

```json
{
    "status": "状态, 表示 API 是否调用成功, 如果成功, 则是 OK, 其他的在下面会说明",
    "retcode": 0,
    "msg": "错误消息, 仅在 API 调用失败时有该字段",
    "data": {
        "响应数据名": "数据值",
        "响应数据名2": "数据值"
    }
}
```

### `get_ws_port` 获取分配的ws端口

- uri: `/api/get_ws_port`
- method: `GET`

**参数**

无

**响应数据**

| 字段名    | 数据类型         | 说明           |
|--------|--------------|--------------|
| `port` | number (int) | 开放监听的 ws 端口号 |

### `get_online_clients` 获取当前在线客户端列表

- uri: `/api/get_online_clients`
- method: `GET`

**参数**

无

**响应数据**

| 字段名       | 数据类型     | 说明      |
|-----------|----------|---------|
| `clients` | Client[] | 在线客户端列表 |

### `set_client_name` 设置当前客户端名称

- uri: `/api/set_client_name`
- method: `POST`

**参数**

| 字段名    | 数据类型   | 默认值 | 说明            |
|--------|--------|-----|---------------|
| `name` | string | -   | 当前客户端名称，需全局唯一 |

**响应数据**

无

### `send_broadcast_data` 发送客户端广播数据

从当前客户端向其他客户端广播数据，可指定目标客户端

- uri: `/api/send_broadcast_data`
- method: `POST`

**参数**

| 字段名       | 数据类型     | 默认值 | 说明       |
|-----------|----------|-----|----------|
| `clients` | Client[] | -   | 指定的客户端列表 |
| `data`    | string   | -   | 需要广播的数据  |

**响应数据**

无

### `send_broadcast_data_callback` 发送客户端广播数据回调

接收到其他客户端广播数据的客户端，可调用此 API 回调发送方客户端，传递数据

- uri: `/api/send_broadcast_data_callback`
- method: `POST`

**参数**

| 字段名      | 数据类型   | 默认值 | 说明              |
|----------|--------|-----|-----------------|
| `client` | Client | -   | 接收回调的客户端，从事件中获取 |
| `data`   | string | -   | 需要回调的数据         |

**响应数据**

无

## 功能增强 - ws 调用

> `功能增强 API` 是 `perpetua` 在符合原 onebot 协议的基础上，在 NTQQ 实现无法满足用户需求的情景下，额外对协议中规定的部分 API 进行的实现、优化、拓展。即基于原有 `onebot` 协议规范，在 `websocket` 连接中进行调用。下文中 `响应数据` 字段为 `基础响应数据` 中 `data` 片段数据，其 `基础响应数据` 格式如下

```json
{
  "status": "状态, 表示 API 是否调用成功, 仅在成功时返回 OK",
  "retcode": 0,
  "echo": "'回声', 如果指定了 echo 字段, 那么响应包也会同时包含一个 echo 字段, 它们会有相同的值",
  "data": {
    "响应数据名": "数据值",
    "响应数据名2": "数据值"
  }
}
```

### `set_restart` 重启 OneBot 实现

由于重启 OneBot 实现同时需要重启 API 服务，这意味着当前的 API 请求会被中断，因此需要异步地重启，接口返回的 `status` 是 `async`。

**参数**

| 字段名     | 数据类型   | 默认值 | 说明                                    |
|---------|--------|-----|---------------------------------------|
| `delay` | number | `0` | 要延迟的毫秒数，如果默认情况下无法重启，可以尝试设置延迟为 2000 左右 |

**响应数据**

无

<hr>