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
   <a alt="Actions" href="https://github.com/IUnlimit/perpetua/actions"><image src="https://github.com/IUnlimit/perpetua/workflows/CI/badge.svg"></image></a>
</p>

<p align="center">
   <a href="https://iunlimit.github.io/perpetua/">[文档]</a>
</p>

### 注意事项

本项目内置了NTQQ实现（Lagrange.OneBot）的分发与运行，并临时提供了签名服务器配置项部分。若您在使用时有任何疑问，随时欢迎进群咨询。Group: [863522624](https://qm.qq.com/cgi-bin/qm/qr?k=Xby1-vbC43Hgv4TXd8LcI889zEhwkq_a&jump_from=webapi&authKey=SmcLCk3eBSQyC0ylq9CiwTafuDk7ls+5QrNDB2//hjTZY6sCTdCz/RKzRwVRrN4J)

### 适配平台

- [x] Linux
- [x] MacOS
- [x] Windows

### OneBot 协议支持

> 为了将 perpetua API 与 OneBot 实现作区分，perpetua 仅提供正向 WebSocket 端口作 OneBot 协议接口，拓展 API 实现使用 WebAPI 的形式供给调用。详见 [拓展实现-API](https://github.com/IUnlimit/perpetua?tab=readme-ov-file#%E6%8B%93%E5%B1%95%E5%AE%9E%E7%8E%B0)

- [x] 正向 WebSocket

### 拓展支持

- [x] 服务注册与发现
- [ ] 自动重启与断点续传 (client -> NTQQ)
- [ ] 负载均衡
- [ ] 熔断限流

### 拓展实现

<details>
<summary>新增支持 API</summary>

| API          | 功能                                                                                                 |
|--------------|----------------------------------------------------------------------------------------------------|
| /get_ws_port | [获取分配的ws端口](https://iunlimit.github.io/perpetua/#/zh-cn/user/enhance-api?id=get_ws_port-获取分配的ws端口) |

[获取分配的ws端口]: #

</details>

<details>
<summary>功能增强 API</summary>

| API          | 功能                                                                                                       |
|--------------|----------------------------------------------------------------------------------------------------------|
| /set_restart | [重启 OneBot 实现](https://iunlimit.github.io/perpetua/#/zh-cn/user/enhance-api?id=set_restart-重启-onebot-实现) |

[重启 OneBot 实现]: #

</details>

### 通信 SDK

| 语言 / 平台 | 名称                                                                         | 简介     | 通信协议支持 |
|---------|----------------------------------------------------------------------------|--------|--------|
| Java    | [perpetua-sdk-for-java](https://github.com/IUnlimit/perpetua-sdk-for-java) | 官方社区实现 | Java   |


### 致谢
- 感谢原机器人社区的贡献者：†白可乐、Alan Zhao、[@fred913](https://github.com/fred913)
- 感谢 [@Thexiaoyuqaq](https://github.com/Thexiaoyuqaq)、小豆子、阿丽塔、polar、一口小雨、黑土、仔仔 等用户在测试、策划方面提供的帮助与支持
