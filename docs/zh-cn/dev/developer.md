# 开发手册

## 项目结构

```
perpetua/
├── main.go                        # 程序入口
├── cmd/perpe/
│   ├── main.go                    # Bootstrap: Configure + Start + EnableAgent
│   ├── exec.go                    # 配置阶段：检测 EMBED/EXTERNAL 模式，初始化 Lagrange.OneBot
│   └── agent.go                   # 代理阶段：启动 Redis/Web面板/HTTP服务/反向WS/HttpPost/NTQQ WS
├── configs/
│   ├── embed.go                   # 嵌入默认配置文件
│   ├── config.yml                 # 默认配置模板
│   └── appsettings.json           # Lagrange.OneBot 默认配置模板
├── internal/
│   ├── global.go                  # 全局变量定义
│   ├── conf/                      # 配置加载、版本检查、配置更新
│   ├── handle/
│   │   ├── handler.go             # Handler 核心结构体与全局集合
│   │   ├── echo.go                # EchoMap：echo 字段路由分发
│   │   ├── cache.go               # 消息缓存（gcache，带过期时间）
│   │   ├── client.go              # Client 接口、读写循环、echo签名、拓展API拦截
│   │   ├── client_websocket.go    # WebSocket 客户端实现（正向/反向）
│   │   ├── client_http.go         # HTTP Post 客户端实现
│   │   ├── serve_websocket.go     # NTQQ WebSocket 连接与消息路由
│   │   ├── enhance_hook.go        # 拓展 API 实现（set_restart 等）
│   │   ├── enhance_event.go       # 拓展事件推送（客户端状态变更、广播等）
│   │   ├── enhance_http.go        # HTTP API 路由注册
│   │   └── init.go                # 初始化 echoMap、globalCache、upgrader
│   ├── hook/
│   │   ├── github.go              # GitHub API 交互
│   │   ├── init.go                # AES 解密 GitHub token
│   │   └── qqimpl/lagrange.go     # Lagrange.OneBot 下载、解压、版本检查
│   ├── model/
│   │   ├── config.go              # Config 配置模型
│   │   ├── api.go                 # ImplType 枚举、Client 结构体
│   │   ├── response.go            # 统一响应格式
│   │   ├── message.go             # 心跳/生命周期元数据
│   │   ├── appsettings.go         # Lagrange.OneBot 配置模型
│   │   └── workflow.go            # GitHub CI 模型
│   ├── logger/                    # 日志模块
│   └── web/
│       ├── server.go              # Web 管理面板服务
│       ├── handler.go             # 管理面板 API 处理器
│       ├── redis.go               # Redis 数据包存储与链路追踪
│       └── recorder.go            # 数据包记录器
├── web/                           # 嵌入式前端静态资源
│   ├── index.html
│   └── assets/
└── pkg/utils/                     # 工具函数
```

## 环境要求

- Go 1.21.4+
- Redis（Web 管理面板功能需要）

## 构建

```shell
go build -o perp .
```

## 核心机制

### Echo 路由

Perpetua 的核心在于利用 OneBot 协议中的 `echo` 字段实现多客户端消息路由。每个客户端发送的 API 请求会被添加唯一的 echo 前缀，当 NTQQ 返回响应时，通过 echo 前缀识别并路由到对应的客户端。

### 拓展 API 拦截

客户端发送的 API 请求在转发到 NTQQ 之前，会先经过 `hookMap` 拦截检查。如果 action 匹配已注册的拓展 API（如 `set_restart`、`set_client_name` 等），则由 Perpetua 直接处理而不转发。

### 事件广播

NTQQ 推送的事件会被广播到所有已连接的客户端（正向 WS、反向 WS、HTTP Post）。拓展事件（如客户端状态变更、分布式广播）同样遵循此机制。

## 贡献指南

1. Fork 本仓库
2. 创建特性分支：`git checkout -b feature/your-feature`
3. 提交更改：`git commit -m 'Add your feature'`
4. 推送分支：`git push origin feature/your-feature`
5. 提交 Pull Request

项目使用 AGPLv3 协议开源。
