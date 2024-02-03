# perpetua

```
                                   __
______   _________________   _____/  |_ __ _______
\____ \_/ __ \_  __ \____ \_/ __ \   __\  |  \__  \
|  |_> >  ___/|  | \/  |_> >  ___/|  | |  |  // __ \_
|   __/ \___  >__|  |   __/ \___  >__| |____/(____  /
|__|        \/      |__|        \/                \/ 
```
<p align="center">
    本项目为 OneBot 协议实现下机器人实例与用户服务间第三方消息代理中间件。通过 WebSocket 代理与额外的 WebAPI 接口，提供单一QQ账号下多端接入，事件回调、消息互通等功能实现。并配备常规流量治理功能包括服务注册发现、熔断限流、负载均衡。
</p>

<p align="center">
    <a alt="Protocol" href="https://github.com/botuniverse/onebot-11"><image src="https://img.shields.io/badge/OneBot-v11-green"></image></a>
    <a alt="NTQQ" href="https://github.com/LagrangeDev/Lagrange.Core"><image src="https://img.shields.io/badge/Lagrange-OneBot-blue"></image></a>
</p>

<p align="center">
   <a alt="License" href="https://www.gnu.org/licenses/agpl-3.0.en.html"><image src="https://img.shields.io/badge/license-AGPLv3-4EB1BA.svg"></image></a>
   <a alt="Release" href="https://github.com/IUnlimit/lagrange-go-distributed/releases"><image src="https://img.shields.io/github/release/IUnlimit/lagrange-go-distributed.svg"></image></a>
</p>

<p align="center">
   <a href="https://iunlimit.github.io/perpetua/">[文档]</a>
</p>

### 注意事项

本项目内置了NTQQ实现（Lagrange.OneBot）的分发与运行，但[签名服务器配置](https://github.com/IUnlimit/perpetua/blob/main/configs/appsettings.json#L9)项部分需您自行寻找解决方案。详见[Lagrange.Core#known-problem](https://github.com/LagrangeDev/Lagrange.Core#known-problem)

### 适配平台

- [x] Linux
- [x] MacOS
- [x] Windows

### 支持协议

- [x] 正向 WebSocket

### 拓展支持

- [x] 服务注册与发现
- [ ] 自动重启与断点续传 (client -> NTQQ)
- [ ] 负载均衡
- [ ] 熔断限流

### 拓展实现

<details>
<summary>API</summary>

| API                | 功能          |
|--------------------|-------------|
| /get_ws_port       | [获取分配的ws端口] |

[获取分配的ws端口]: #

</details>

### 通信 SDK

| 语言 / 平台 | 名称                                                                         | 简介     | 通信协议支持   |
|---------|----------------------------------------------------------------------------|--------|------|
| Java    | [perpetua-sdk-for-java](https://github.com/IUnlimit/perpetua-sdk-for-java) | 官方社区实现 | Java |

