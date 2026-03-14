# 连接服务端

Perpetua 作为中间件，需要先与上游 NTQQ 实现（Lagrange.OneBot）建立连接。支持两种模式：

## 内置模式（EMBED）

默认模式。Perpetua 自动下载并管理 `Lagrange.OneBot` 实例。

### 初始化

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

### Lagrange.OneBot 配置

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
  // 验证服务器，必填项
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

### 权限配置

#### Windows

Windows 用户只需在程序启动时批准相关权限弹窗即可

#### Linux

Linux 用户需要在运行前对 `./config/Lagrange.OneBot/` 路径下的 `Lagrange.OneBot` 可执行程序赋予运行权限

```shell
chmod +x ./config/Lagrange.OneBot/Lagrange.OneBot
```

### Perpetua 配置

内置模式下，`ntqq-impl` 配置段无需手动修改，程序会自动维护：

```yaml
ntqq-impl:
  # 以下配置项除 enable 外自动更新，无需手动变动
  update:
    enable: false
    id: 0
    platform: ""
    updated-at: "0001-01-01T00:00:00Z"
```

设置 `update.enable` 为 `true` 可在启动时自动检查 Lagrange.OneBot 更新。

## 外置模式（EXTERNAL）

当您已有独立运行的 OneBot 实现时（如 Docker 部署、自行管理的 Lagrange 实例），可使用外置模式。Perpetua 将直接连接到指定的正向 WebSocket 地址，不再管理 NTQQ 实现的生命周期。

### 适用场景

- Docker 容器化部署，NTQQ 实现运行在独立容器中
- 已有现成的 OneBot 实现实例，无需 Perpetua 管理
- 需要自定义 NTQQ 实现的启动参数或版本

### Perpetua 配置

在 `config.yml` 中指定 `external-web-socket` 即可启用外置模式：

```yaml
ntqq-impl:
  # 指定外置 OneBot 实现的正向 ws 地址
  external-web-socket: "ws://127.0.0.1:5700/onebot/v11/ws"
  # 外置 OneBot 实现的 AccessToken（可选）
  external-access-token: ""
```

> 注意：`external-web-socket` 地址需与 OneBot 实现的 `ForwardWebSocket` 配置一致
