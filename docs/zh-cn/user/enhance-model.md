# 拓展数据类型

## Client

| 字段名           | 数据类型   | 说明    |
|---------------|--------|-------|
| `app_id`      | string | 客户端ID |
| `client_name` | string | 客户端名称 |

## Response

所有 API 调用的统一响应格式

| 字段名       | 数据类型   | 说明                                                                 |
|-----------|--------|--------------------------------------------------------------------|
| `status`  | string | 状态：`ok` 调用成功 / `async` 已提交异步处理 / `failed` 调用失败 |
| `retcode` | number | 响应码：`0` 成功 / `1` 已提交异步处理 / 其他值表示失败                  |
| `msg`     | string | 错误消息，仅在调用失败时存在                                           |
| `data`    | object | 响应数据                                                             |

## Packet

Web 管理面板中的数据包记录类型，用于链路追踪

| 字段名           | 数据类型   | 说明                                                        |
|---------------|--------|-----------------------------------------------------------|
| `id`          | string | 数据包唯一ID                                                  |
| `trace_id`    | string | 链路追踪ID，同一链路下的数据包共享相同的 trace_id                    |
| `timestamp`   | number | 时间戳（毫秒）                                                 |
| `link`        | string | 链路标识：`ntqq`（NTQQ ↔ Perpetua）或 `client`（Perpetua ↔ 客户端） |
| `direction`   | string | 方向：`inbound`（流入 Perpetua）或 `outbound`（流出 Perpetua）      |
| `handler_id`  | string | 处理器ID（可选）                                                |
| `client_name` | string | 客户端名称（可选）                                                |
| `data`        | object | 数据包内容                                                     |
