# 拓展 Event

拓展 Event 同样基于 onebot 协议实现。为了与原有事件类型作区分，拓展事件在部分归类与原有事件类型的基础上，额外分成以下几个大类：

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
| `online`      | bool           | -               | 当前是否在线        |

- Client 可在 [API - 获取当前在线客户端列表](https://iunlimit.github.io/perpetua/#/zh-cn/user/enhance-api?id=get_online_clients-%e8%8e%b7%e5%8f%96%e5%bd%93%e5%89%8d%e5%9c%a8%e7%ba%bf%e5%ae%a2%e6%88%b7%e7%ab%af%e5%88%97%e8%a1%a8) 查看

## 分布式事件

> 通知事件对应的 `post_type` 字段值为 `distributed`

### 客户端广播

### 客户端广播回调