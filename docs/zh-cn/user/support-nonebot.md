# 第三方支持 - nonebot

`perpetua` 与 `nonebot` 间采用反向 websocket 链接的方式进行通信，您可参考以下配置进行服务间的连接（nonebot 安装启动不再赘述）

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
# 反向 websocket 相关配置
reverse-web-socket:
  # 是否开启功能
  enabled: true
  # 上报地址 - 与 nonebot 配置对应
  url: 'ws://127.0.0.1:8800/onebot/v11/ws/'
  # AccessToken
  access-token: ''
```

## 连接标识

当连接成功后，您可以在 `perpetua` 控制台看到如下输出

> [Client] Start connecting to reverse-websocket: ws://127.0.0.1:8800/onebot/v11/ws/ with headers: map[X-Client-Role:[Universal] X-Self-Id:[3012218237]]

对应的，在 `nonebot` 控制台看到如下输出

> [INFO] nonebot | OneBot V11 | Bot 3012218237 connected