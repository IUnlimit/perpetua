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
    本项目为OneBot协议下机器人实例与用户服务间第三方中间件。基于NTQQ实现，提供了额外的分布式支持，包括WebSocket代理与请求分配，并提供额外的API接口
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
   <a href="https://iunlimit.github.io/perpetua/">Document</a>
</p>

### 注意事项

本项目内置了NTQQ实现（Lagrange.OneBot）的分发与运行，但[签名服务器配置](https://github.com/IUnlimit/perpetua/blob/main/configs/appsettings.json#L9)项部分需您自行寻找解决方案。详见[Lagrange.Core#known-problem](https://github.com/LagrangeDev/Lagrange.Core#known-problem)