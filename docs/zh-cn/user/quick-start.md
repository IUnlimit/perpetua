# 快速开始

> 鉴于不同平台用户群体分布不同，本文以 Windows 用户视角给出软件的使用流程示例

## 运行

Windows 用户请使用 `powershell` 运行本项目，以避免直接运行 `exe` 文件或使用 `cmd` 运行所导致的无法查看程序退出日志、控制台字符打印异常的部分问题。

推荐的运行脚本内容如下：

```bat
# start.bat
./perp.exe
```

> Windows 下可使用 `faststart` 参数跳过 shell 确认提示：`./perp.exe faststart`

## 基础配置

Perpetua 的配置文件路径为 `./config/config.yml`，首次运行会自动生成。以下为基础配置项：

```yaml
# 日志项配置
log:
  # 是否每次启动新建log文件
  force-new: false
  # 日志等级
  #   trace debug info warn error
  level: "debug"
  # 日志存活时间，到期新建log文件
  aging: 24h
  # 是否开启控制台颜色
  colorful: true

# 接收消息的最大缓存时间
msg-expire-time: 30m
```

如果希望在发生 bug 时进行更详细的溯源，可以将日志等级调整为 `debug`。

## 下一步

Perpetua 作为中间件，需要分别配置上游（NTQQ 实现）和下游（用户客户端）的连接：

1. **[连接服务端](zh-cn/user/connect-server.md)** — 配置 Perpetua 与 NTQQ 实现（Lagrange.OneBot）的连接
2. **[连接客户端](zh-cn/user/connect-client.md)** — 配置用户客户端如何接入 Perpetua

```
NTQQ(Lagrange.OneBot) ←──上游──→ Perpetua ←──下游──→ 用户客户端
```
