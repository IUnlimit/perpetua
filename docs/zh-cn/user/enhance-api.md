# 拓展 API

!> 下文中的 API 有 `http` 调用与 `websocket` 两种调用方式，两种方式均需要以 onebot 协议规定的格式发起请求，此处不再追述。详见 [onebot-11/communication](https://github.com/botuniverse/onebot-11/tree/master/communication)

Perpetua 提供的 API 一共分为两种类型

1. **新增支持**  
  是 `perpetua` 基于自身业务场景额外提供的 API 接口，大部分是为提供分布式服务治理支持所服务，故如您仅有单机需求时，只需关注少数关键 `新增支持 API` 即可。
2. **功能增强**  
  是 `perpetua` 在符合原 onebot 协议的基础上，在 NTQQ 实现无法满足用户需求的情景下，额外对协议中规定的部分 API 进行的优化、拓展。基于原有 `onebot` 协议规范实现。

下文中 `响应数据` 字段为 `基础响应数据` 中 `data` 片段数据，其 `基础响应数据` 格式将分别在片段开头给出

## http 调用

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
- type: `新增支持`

**参数**

无

**响应数据**

| 字段名    | 数据类型         | 说明           |
|--------|--------------|--------------|
| `port` | number (int) | 开放监听的 ws 端口号 |

### `get_online_clients` 获取当前在线客户端列表

- uri: `/api/get_online_clients`
- method: `GET`
- type: `新增支持`

**参数**

无

**响应数据**

响应内容为 JSON 数组，每个元素定义见 [Client](https://iunlimit.github.io/perpetua/#/zh-cn/user/enhance-model?id=client)

## ws 调用

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

- type: `功能增强`

**参数**

| 字段名     | 数据类型   | 默认值 | 说明                                    |
|---------|--------|-----|---------------------------------------|
| `delay` | number | `0` | 要延迟的毫秒数，如果默认情况下无法重启，可以尝试设置延迟为 2000 左右 |

**响应数据**

无

### `set_client_name` 设置当前客户端名称

- type: `新增支持`

**参数**

| 字段名    | 数据类型   | 默认值 | 说明            |
|--------|--------|-----|---------------|
| `name` | string | -   | 当前客户端名称，需全局唯一 |

**响应数据**

无

### `send_broadcast_data` 发送客户端广播数据

从当前客户端向其他客户端广播数据，需指定目标客户端

- type: `新增支持`

**参数**

| 字段名       | 数据类型     | 默认值 | 说明       |
|-----------|----------|-----|----------|
| `clients` | Client[] | -   | 指定的客户端列表 |
| `data`    | string   | -   | 需要广播的数据  |

**响应数据**

| 字段名    | 数据类型   | 说明             |
|--------|--------|----------------|
| `uuid` | string | 此次客户端广播事件的唯一id |

### `send_broadcast_data_callback` 发送客户端广播数据回调

接收到其他客户端广播数据的客户端，可调用此 API 回调发送方客户端，传递数据

- type: `新增支持`

**参数**

| 字段名      | 数据类型   | 默认值 | 说明                    |
|----------|--------|-----|-----------------------|
| `client` | Client | -   | 接收回调的客户端，从事件中获取       |
| `uuid`   | string | -   | 此次客户端广播事件的唯一id，从事件中获取 |
| `data`   | string | -   | 需要回调的数据               |

**响应数据**

无

## Web 管理面板 API

Web 管理面板提供了一组 REST API 用于查看连接状态、数据包记录和链路追踪。这些 API 运行在独立端口上（默认 `9190`，可通过 `web.port` 配置）。

> 管理面板 API 依赖 Redis 进行数据包持久化，请确保 Redis 配置正确

### `get_connections` 获取所有活跃连接

- uri: `/api/web/connections`
- method: `GET`

**参数**

无

**响应数据**

响应内容为 JSON 数组，每个元素包含连接的详细信息（客户端ID、名称、连接类型等）

### `get_packets` 查询数据包记录

- uri: `/api/web/packets`
- method: `GET`

**参数**

| 字段名          | 数据类型   | 默认值  | 说明                          |
|--------------|--------|------|-----------------------------|
| `handler_id` | string | -    | 按处理器ID过滤（可选）                |
| `offset`     | number | `0`  | 分页偏移量                       |
| `limit`      | number | `50` | 每页数量，最大 200                 |

**响应数据**

| 字段名       | 数据类型     | 说明       |
|-----------|----------|----------|
| `packets` | Packet[] | 数据包列表    |
| `total`   | number   | 数据包总数    |
| `offset`  | number   | 当前偏移量    |
| `limit`   | number   | 当前每页数量   |

### `get_packet_trace` 查询数据包链路追踪

根据数据包ID查询关联的完整链路追踪，返回同一 trace_id 下的所有数据包

- uri: `/api/web/packets/trace`
- method: `GET`

**参数**

| 字段名  | 数据类型   | 默认值 | 说明    |
|------|--------|-----|-------|
| `id` | string | -   | 数据包ID |

**响应数据**

| 字段名       | 数据类型     | 说明              |
|-----------|----------|-----------------|
| `source`  | Packet   | 源数据包            |
| `related` | Packet[] | 同一链路下的关联数据包列表   |

### `get_system_info` 获取系统概览

- uri: `/api/web/system`
- method: `GET`

**参数**

无

**响应数据**

| 字段名                 | 数据类型   | 说明        |
|---------------------|--------|-----------|
| `connections_count` | number | 当前活跃连接数   |
| `packets_count`     | number | 数据包记录总数   |
| `lifecycle`         | object | 生命周期元数据   |
| `heartbeat`         | object | 最近一次心跳数据  |

### `delete_packets` 删除数据包记录

删除指定时间戳之前的所有数据包记录

- uri: `/api/web/packets`
- method: `DELETE`

**参数**

| 字段名      | 数据类型   | 默认值 | 说明                |
|----------|--------|-----|-------------------|
| `before` | number | -   | Unix 时间戳（毫秒），删除此时间之前的数据包 |

**响应数据**

| 字段名       | 数据类型   | 说明       |
|-----------|--------|----------|
| `removed` | number | 删除的数据包数量 |

<hr>
