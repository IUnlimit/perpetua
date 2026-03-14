# Web 管理面板

Perpetua 内置了一个 Web 管理面板，提供连接状态监控、数据包记录查询和链路追踪功能。

## 前置依赖

Web 管理面板依赖 Redis 进行数据包记录的持久化存储，请确保 Redis 服务已启动并正确配置。

## 配置

在 `config.yml` 中配置 Redis 和 Web 面板：

```yaml
# Redis 配置（用于数据包记录持久化）
redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

# Web 管理面板配置
web:
  # 管理面板端口
  port: 9190
  # 数据包记录过期时间
  packet-expire: 24h
  # 过期数据清理间隔
  cleanup-interval: 1h
```

| 配置项                   | 说明                     | 默认值              |
|-----------------------|------------------------|------------------|
| `redis.addr`          | Redis 连接地址             | `127.0.0.1:6379` |
| `redis.password`      | Redis 密码               | 空                |
| `redis.db`            | Redis 数据库编号            | `0`              |
| `web.port`            | 管理面板监听端口               | `9190`           |
| `web.packet-expire`   | 数据包记录过期时间              | `24h`            |
| `web.cleanup-interval`| 过期数据自动清理间隔             | `1h`             |

## 访问

启动 Perpetua 后，在浏览器中访问 `http://127.0.0.1:9190` 即可打开管理面板。

## 功能

### 连接监控

查看当前所有活跃的客户端连接，包括客户端ID、名称、连接类型等信息。

### 数据包记录

浏览所有经过 Perpetua 的数据包记录，支持按处理器ID过滤和分页查询。每个数据包包含以下信息：

- 链路标识（`ntqq` 或 `client`）：标识数据包所在的通信链路
- 方向（`inbound` 或 `outbound`）：标识数据包相对于 Perpetua 的流向
- 关联的处理器和客户端信息
- 完整的数据包内容

### 链路追踪

通过 trace_id 关联同一请求链路下的所有数据包，实现 NTQQ ↔ Perpetua ↔ Client 的全链路追踪。

### 系统概览

查看系统运行状态，包括当前连接数、数据包总数、生命周期和心跳信息。

## API 接口

管理面板提供的 REST API 详见 [拓展 API - Web 管理面板 API](https://iunlimit.github.io/perpetua/#/zh-cn/user/enhance-api?id=web-管理面板-api)
