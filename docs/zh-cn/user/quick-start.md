# 快速开始

> 鉴于不同平台用户群体分布不同，本文以 Windows 用户视角给出软件的使用流程示例

## 运行

Windows 用户请使用 `powershell` 运行本项目，以避免直接运行 `exe` 文件或使用 `cmd` 运行所导致的无法查看程序退出日志、控制台字符打印异常的部分问题。

推荐的运行脚本内容如下：

```bat
# start.bat
./perp.exe
```

## 初始化

第一次运行程序，会自动获取最新 NTQQ 实现列表供您选择下载（目前仅为 `Lagrange.OneBot`），您需要选择符合您机器架构的构建版本（输入不同版本前`[]`内的数字）进行下载。

> 注意：下载内容为 [Github Action](https://github.com/LagrangeDev/Lagrange.Core/actions/workflows/Lagrange.OneBot-build.yml) 构建的最新版本，可能存在需要代理才能访问的情况。您也可以手动下载并解压到 `./config/Lagrange.OneBot/` 路径下

此阶段日志输出内容如下：

<details>
<summary><b>点击展开</b></summary>

```text
[PERP] [INFO] [2024-02-01 21:20:47]: Searching Lagrange.OneBot ...
[PERP] [INFO] [2024-02-01 21:20:48]: Please choose the Lagrange.OneBot software suitable for your platform (send the number before option)
[0] Lagrange.OneBot_win-x86
[1] Lagrange.OneBot_win-x64
[2] Lagrange.OneBot_osx-x64
[3] Lagrange.OneBot_osx-arm64
[4] Lagrange.OneBot_linux-x64
[5] Lagrange.OneBot_linux-arm64
[6] Lagrange.OneBot_linux-arm
```

</details>

下载完成后，因为缺少 `appsettings.json`（Lagrange.OneBot 配置）文件，程序会自动生成初始配置文件并退出

## Lagrange.OneBot 配置

下为 `appsettings.json` 的默认配置，本文就关键内容给出注释说明

<details>
<summary><b>点击展开</b></summary>

```json
{
  "Logging": {
    "LogLevel": {
      "Default": "Information",
      "Microsoft": "Warning",
      "Microsoft.Hosting.Lifetime": "Information"
    }
  },
  // 验证服务器，必填项（v0.0.4 后自动配置
  "SignServerUrl": "",
  "Account": {
    // qq 账户，若不填则使用扫码登陆
    "Uin": 0,
    // qq 密码，若不填则使用扫码登陆
    "Password": "",
    // 协议类型，目前仅支持 Linux
    "Protocol": "Linux",
    "AutoReconnect": true,
    "GetOptimumServer": true
  },
  "Message": {
    "IgnoreSelf": true
  },
  // 连接配置
  "Implementations": [
    {
      // 连接类型：正向 WebSocket 连接
      // Perpetua 将自动读取使用第一个 ForwardWebSocket 连接配置
      "Type": "ForwardWebSocket",
      "Host": "127.0.0.1",
      "Port": 5700,
      "Suffix": "/onebot/v11/ws",
      "ReconnectInterval": 5000,
      "HeartBeatInterval": 5000,
      "AccessToken": ""
    }
  ]
}
```

</details>

## 权限配置

### Windows

Windows 用户只需在程序启动时批准相关权限弹窗即可

### Linux

Linux 用户需要在运行前对 `./config/Lagrange.OneBot/` 路径下的 `Lagrange.OneBot` 可执行程序赋予运行权限

```shell
chmod +x ./config/Lagrange.OneBot/Lagrange.OneBot
```

## Perpetua 配置

Perpetua 默认生产的配置即可满足绝大部分需求，其路径为 `./config/config.yml`。如果希望在发生bug时进行更详细的溯源，可以将日志等级调整为 `debug`。

<details>
<summary><b>点击展开</b></summary>

```yaml
#	                                    __
#	______   _________________   _____/  |_ __ _______
#	\____ \_/ __ \_  __ \____ \_/ __ \   __\  |  \__  \
#	|  |_> >  ___/|  | \/  |_> >  ___/|  | |  |  // __ \_
#	|   __/ \___  >__|  |   __/ \___  >__| |____/(____  /
#	|__|        \/      |__|        \/                \/
#
# Notice
#   perpetua 固定连接第一个 ForwardWebSocket 配置项

# 日志项配置
log:
  # 是否每次启动新建log文件
  force-new: false
  # 日志等级
  #   trace debug info warn error
  level: "info"
  # 日志存活时间，到期新建log文件
  aging: 24h
  # 是否开启控制台颜色
  colorful: true

# 本配置项自动更新，无需手动
ntqq-impl:
  update: false
  id: 0
  platform: ""
  updated-at: "0001-01-01T00:00:00Z"

# http 相关配置
http:
  # 监听端口
  port: 8080

# websocket 相关配置
web-socket:
  # ws监听最长等待时间
  timeout: 15s

# 接收消息的最大缓存时间
msg-expire-time: 30m
```

</details>




