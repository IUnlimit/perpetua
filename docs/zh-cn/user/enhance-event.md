# 拓展 Event

`拓展 Event` 同样基于 `onebot` 协议实现。为了与原有事件类型作区分，拓展事件在部分归类与原有事件类型的基础上，额外分成以下几个大类：

- 通知事件，在原有 `onebot` 协议的基础上额外拓展的通知事件
- 分布式事件，包括各客户端实例间数据推送、回调等

## 通知事件

> 通知事件对应的 `post_type` 字段值为 `notice`

### 客户端在线状态变更

其他连接到 Perpetua 的客户端状态变更事件

**事件数据**

| 字段名           | 数据类型           | 可能的值            | 说明            |
|---------------|----------------|-----------------|---------------|
| `time`        | number (int64) | -               | 事件发生的时间戳      |
| `self_id`     | number (int64) | -               | 收到事件的机器人 QQ 号 |
| `post_type`   | string         | `notice`        | 上报类型          |
| `notice_type` | string         | `client_status` | 通知类型          |
| `client`      | Client*        | -               | 客户端信息         |
| `online`      | bool           | -               | 当前客户端是否在线        |
| `self_client` | bool           | -               | 当前客户端是否是自身        |

- Client 可在 [拓展数据类型](https://iunlimit.github.io/perpetua/#/zh-cn/user/enhance-model?id=client) 查看

## 分布式事件

> 通知事件对应的 `post_type` 字段值为 `distributed`

### 客户端广播

当一个客户端通过 API 向其他客户端发送广播传递数据时，接收消息的客户端会收到此事件

| 字段名                | 数据类型           | 可能的值          | 说明             |
|--------------------|----------------|---------------|----------------|
| `time`             | number (int64) | -             | 事件发生的时间戳       |
| `self_id`          | number (int64) | -             | 收到事件的机器人 QQ 号  |
| `post_type`        | string         | `distributed` | 上报类型           |
| `distributed_type` | string         | `broadcast`   | 分布式类型          |
| `client`           | Client*        | -             | 发送广播的客户端信息     |
| `uuid`             | string         | -             | 此次客户端广播事件的唯一id |
| `data`             | string         | -             | 广播的数据          |

### 客户端广播回调

当客户端接收到 `客户端广播` 事件时，可通过 API 向广播发起方回调事件，发起方接收到的事件定义如下

| 字段名                | 数据类型           | 可能的值                 | 说明             |
|--------------------|----------------|----------------------|----------------|
| `time`             | number (int64) | -                    | 事件发生的时间戳       |
| `self_id`          | number (int64) | -                    | 收到事件的机器人 QQ 号  |
| `post_type`        | string         | `distributed`        | 上报类型           |
| `distributed_type` | string         | `broadcast_callback` | 分布式类型          |
| `client`           | Client*        | -                    | 发起回调的客户端信息     |
| `uuid`             | string         | -                    | 此次客户端广播事件的唯一id |
| `data`             | string         | -                    | 回调的数据          |