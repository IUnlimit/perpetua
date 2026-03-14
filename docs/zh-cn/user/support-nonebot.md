# 第三方支持 - nonebot

`perpetua` 与 `nonebot` 间采用反向 websocket 链接的方式进行通信，您可参考以下配置进行服务间的连接（nonebot 安装启动不再赘述）

> 关于反向 WebSocket 的详细说明请参考 [连接客户端 - 反向 WebSocket](zh-cn/user/connect-client.md?id=反向-websocket)

## nonebot 配置

```
# 选择 fastapi 监听端口
DRIVER=~fastapi

# 监听链接配置
# ws://127.0.0.1:8800/onebot/v11/ws/
HOST=127.0.0.1
PORT=8800
# 超级管理员配置
SUPERUSERS=["765743073"]
```

## perpetua 配置

```yaml
reverse-web-socket:
  - url: 'ws://127.0.0.1:8800/onebot/v11/ws/'
    access-token: ''
```

## 连接标识

当连接成功后，您可以在 `perpetua` 控制台看到如下输出

> [Client] Start connecting to reverse-websocket: ws://127.0.0.1:8800/onebot/v11/ws/ with headers: map[X-Client-Role:[Universal] X-Self-Id:[3012218237]]

对应的，在 `nonebot` 控制台看到如下输出

> [INFO] nonebot | OneBot V11 | Bot 3012218237 connected
