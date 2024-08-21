---
title: HomeLab 玩法简单分享
date: 2021-09-14 19:51:33
tags:
    - Homelab
    - NAS
    - Docker
categories:
    - 乱七八糟
toc: true
---

大学毕业之前一个冲动买了台式机，又一个冲动买了台 Linux 主机。到现在它已经运行了四年多了，简单分享下自己的玩法。
<!--more-->

> 注：本文于 2024 年 8 月进行过一次修订，根据自己的使用经验，对已有内容做了一些更新和补充。成段补充的部分会以引用的形式添加在文章中。

# 背景
毕业之前在公司附近租了房，再加上受到了网络的蛊惑，于是陷入了“买一台 NAS 来大幅提高生活质量”的念头之中。看了很多成熟的 NAS 方案（比如群辉或威联通），最后还是在高昂的价格面前望而却步。  
当时的我，傻乎乎地认为品牌 NAS 的平台只是“SMB + RAID + 媒体服务器”而已，那么既然成熟的方案那么贵，为什么不搞个 Linux 自己折腾呢？脱离了平台的束缚，反而可能有更多的可能性。  
最后，我决定自己买硬件搭一台 Linux 主机，做一个 HomeLab.

# 硬件选型
选择主机硬件的时候考虑了自己的需求，大致如下：
1. **价格便宜**：我只是个穷学生，看他们玩虚拟化的都上了 E3 E5，这么吃硬件的东西我还是不玩了吧 :new_moon_with_face:
2. 功耗低：24 小时开机，电费还得自己掏，所以买个 TDP 几十瓦的 i3-7100 或者 G4560 感觉好像有点烧钱。这条需求本质上可以划入上一条。
3. 有最基本的 IO 接口：对于我来说，有 3-4 个 SATA & 千兆网口就够了。
4. 不一定要做 RAID：硬盘容量是有限的，而备份重要的文件可以有很多方式。这条本质上还是第一条。
5. 可以运行 Linux：~~我承认这条是来凑数的~~

在[英特尔® 产品规范](https://ark.intel.com/content/www/cn/zh/ark.html)中翻看了两周后，我的选择范围从酷睿降到了奔腾，又降到了赛扬。最后，我选择了赛扬 J3455 来做这台主机的 CPU.  
配置清单如下，价格都是购买时的价格。

| 硬件 | 型号 | 价格（元） | 备注 |
| - | - | - | - |
| 主板 & CPU | 华擎 J3455-ITX | 494 | |
| 电源 | 金河田 480GT | 99 | |
| 内存 | 威刚 4G LPDDR3 | 189 |  |
| SSD | 22G  | 0 | 从上大学时买的超极本里拆下来的缓存盘 |
| 硬盘 | 2T 西数红盘 + 2T 西数蓝盘 | 599 + 405 | 贪便宜买了蓝盘 |
| 机箱 | 乔思伯 C2 | 139 | 2.5 寸盘位 * 1，3.5 寸盘位 * 2 |

后来又陆续升级了一些硬件：
| 硬件 | 型号 | 价格（元） | 备注 |
| - | - | - | - |
| 电源 | 航嘉 300W 电源 | 49.9 | 电源风扇出了点问题导致噪音很大，于是花了几十块钱从淘宝上买了块二手电源 |
| 内存 | 金士顿 8G LPDDR3 | 379 | 内存嘛，多多益善 |
| SSD | 256G 紫光 S100 | 179 | Docker 镜像太多，22G 硬盘快要装不下了 |
| SSD | 256G 闪迪 SSD Plus | 219 | 垃圾紫光买回来以后每个月都会出问题，只好把它换掉:unamused: |
| 硬盘 | 4T 日立企业盘 | 880 | 购买于 2018 年 10 月，替换了通电很久的蓝盘 |
| 机箱 | 乔思伯 N1 | 589 | 购买于 2021 年 11 月，改善外观并提升盘位数量 |
| SATA 拓展卡 | 杂牌 | 100 以内 | 主板的 SATA 接口不够用了，因此利用上了 PCIe x1 接口 |
| 硬盘 | 10T + 16T 西数企业盘 | 1730 元 | 10T 硬盘做主要存储，二手 16T 硬盘做定时热备 |

去除已淘汰的硬件，整机成本如下：
* 硬盘（2T + 4T + 10T + 16T）：3209 元
* 主机（除硬盘以外的部分，但包括 SSD 系统盘）：2138.9 元

写这篇文章的时候翻了下淘宝，发现这几年 Intel 推出了更多低功耗但性能更强的 CPU，然而这些 CPU 大多数都拿去做成了成套的 NAS、NUC 和工控机方案，貌似很难再买到 J3455-ITX 这种主板 + CPU 的组合了。

> 在 2024 年，由于市场对软路由的稳定需求，因此很多小厂商开发了搭载 N100 或 J4125 的 ITX 方案，目前选择比 2021 年要多了不少。

此外，由于主机的硬盘太过吵闹，有时为了享受更好的睡眠，我不得不关掉它的电源，所以我又买了一块可以安安静静挂机的树莓派 4（后来用上文淘汰下来的 22G SSD 重做了系统盘来代替性能捉急的 TF 卡），并将一部分我认为比较重要的软件挪到了它的上面。  
最近在机缘巧合之下又收了一块树莓派，不过目前处于在线但闲置的状态，跑了个 k3s 偶尔玩一玩。

附两张机器的配置图：
<img style="max-width: 45vw;" alt="NAS 配置" src="/pics/homelab/nas.png"></img>
<img style="max-width: 45vw;" alt="树莓派配置" src="/pics/homelab/raspberrypi.png"></img>

# 软件
> home lab有两种玩的方向。一种是 all in one，一种是 one by one。
>    —— Twitter [@riverscn](https://twitter.com/riverscn/status/1428712696615260163)

对于我来说，我更倾向于 one by one 的玩法。虽然我还没有那么多硬件，但一台 x86 主机，加上两块树莓派，已经足够避免主机层面的单点故障。

我的绝大部分服务都用 Docker 部署，启动脚本及配置文件全部放在一个 Git 仓库中，通过 Git 进行管理。这样以来，我就能够以极低的成本将某个服务在不同主机间迁移。  
不玩虚拟化，不玩软路由，Infrastructure as Code，多机多副本，虽然可用性没有那么强，但面对日常使用还是足够了。

## 基础平台
上大学时上过一门计算机体系结构的课，当时装过几台 Debian 的虚拟机，这使得我对这个发行版有了一些好感。于是我选择了 [Debian](https://debian.org/) 作为主机的操作系统。  
由于当时 SSD 空间有限，所以我选择了 Debian 的最小安装，并选用了轻量化的 LXDE 作为桌面系统（然而基本没有用过）。  
除此之外，为了使用更新的软件包，我将系统更新到了 Debian testing（最新版本代号是 `bookworm`）。不过使用新版系统就同样需要承担不稳定的风险：较新版本的内核会偶尔出现不稳定的情况（如网络接口自动断开），所以在遇到这样的问题时就需要手动回滚到旧版内核来使用。

主机买回来还是要当 NAS 用的，所以首先安装和配置的还是 [samba](https://www.samba.org/)。通过简单的配置我们就可以在 PC 中添加网络驱动器，在局域网内访问 NAS 中的文件，还可以通过[文件历史记录](https://support.microsoft.com/zh-cn/windows/windows-%E4%B8%AD%E7%9A%84%E6%96%87%E4%BB%B6%E5%8E%86%E5%8F%B2%E8%AE%B0%E5%BD%95-5de0e203-ebae-05ab-db85-d5aa0a199255)备份 PC 上的文件。

除了 NAS 之外，我们还需要搭建一个基础的平台用于托管各种应用，平台的组成大致如下：
1. Docker：用于容器托管。很多流行的应用都有官方或社区维护的 Docker 镜像，使用 Docker 部署应用能够大幅降低部署成本
2. [Gogs](https://gogs.io/)：用于管理 Git 仓库，后来换成了 [Gitea](https://gitea.io/zh-cn/)
3. [frps & frpc](https://github.com/fatedier/frp)：服务端搭建在腾讯云上，客户端通过 Docker 部署在主机中，实现 tcp 和 http 的内网穿透
4.  一个用来组织和管理服务脚本和配置的 Git 仓库，远端存储在 Gitea 中
5.  一些用于日常操作的基础工具，如 git, vim 和 tmux

有了这些组件，我就可以非常快捷地部署新的服务。

## 应用软件
在主机和树莓派上部署的应用可以称得上五花八门，绝大部分都是通过 Web UI 进行交互，使用 Docker 进行托管，通过 volume 完成目录共享。一部分不太适合使用 Docker 托管的服务（如需要使用 ssh 的 Gitea），早期我选择使用 [Supervisor](http://supervisord.org/) 托管，后来全部迁移到了 systemd 上。

部署方式大同小异，而且大部分部署在 Docker 上的应用都能够做到开箱即用，所以没有什么可以单独分享的。简单列举一下我在家部署的应用：

| 软件 | 简介 | 托管方式 |
| - | - | - |
| [Caddy](https://caddyserver.com/) | 配置简单的 HTTP 服务器，用做主机网关 | docker | |
| [acme.sh](https://github.com/acmesh-official/acme.sh) | 获取 Lets'Encrypt 的 HTTPS 证书，用于内网和外网的 HTTPS 访问 | cron |
| [Cockpit](https://cockpit-project.org/) | 主机监控面板 | systemd |
| [Portainer](https://www.portainer.io/) | 容器管理工具 | docker |
| [Prometheus](https://prometheus.io/) & AlertManager & [Grafana](https://grafana.com/) | 监控报警套件，偶尔用于主机问题排查 | docker |
| Node Exporter & cAdvisor | prometheus exporter | docker |
| [Loki](https://grafana.com/oss/loki/) & promtail | 日志收集平台及组件，负责聚合主机、容器和网关请求日志 | docker |
| [Drone](https://www.drone.io/) & Drone Runner | CI 平台，搭配 Gitea 使用 | docker |
| [Clash](https://github.com/Dreamacro/clash) | 代理服务器 | docker |
| [yacd](https://github.com/haishanh/yacd) | Clash 管理面板 | caddy 静态托管 |
| [aria2](https://aria2.github.io/) | 下载工具，支持 HTTP、BT 和磁力链接 | docker |
| [AriaNg](https://github.com/mayswind/AriaNg) | aria2 管理面板 | caddy 静态托管 |
| [Sync Home](https://www.resilio.com/individuals/) | 专有软件，P2P 文件分享，也可用于文件同步 | docker |
| [Kiwix](https://www.kiwix.org/en/) | 离线的维基百科服务，以备不时之需 | docker |
| [NextCloud](https://nextcloud.com/) | 网盘服务，从某种意义上可以代替百度网盘 | docker |
| [WebDAV](https://github.com/hacdias/webdav) | WebDAV 服务，托管 Joplin 中的笔记内容 | docker |
| [Tiny Tiny RSS](https://tt-rss.org/) | RSS 阅读器，用于替代 Feedly | docker，有人专门做了[容器化部署方案](https://ttrss.henry.wang/#deployment-via-docker) |
| [rss-proxy](https://git.stdioa.com/stdioa/rss-proxy) | 一个极为简陋的 HTTP 反向代理，方便 TT-RSS 抓取资源 | docker |
| [fava](https://github.com/beancount/fava) & [beancount-bot](https://github.com/StdioA/beancount-bot) | Beancount UI 和 Telegram bot，详见[《开始使用 Beancount》](/2020/09/using-beancount/) | systemd |
| [Snapdrop](https://snapdrop.net/) / [PairDrop](https://github.com/schlagmichdoch/pairdrop) | 基于 WebRTC 的本地网络文件传输工具 | docker |
| [Mattermost](https://https://www.navidrome.org/) | 音乐管理软件，私有云音乐 | docker | 
| [Navidrome](https://mattermost.com/) | 即时通讯软件，主要用做 Chatbot（Beancount 记账和 LLM 交互）平台 | docker | 
| [VaultWarden](https://github.com/dani-garcia/vaultwarden) | 与 BitWarden 兼容的密码管理器服务端 | docker | 
| [Stirling-PDF](https://github.com/Stirling-Tools/Stirling-PDF) | PDF 处理应用 | docker | 
| [Calibre](https://calibre-ebook.com/) & [Calibre Web](https://github.com/janeczku/calibre-web) | 图书管理 & 在线阅读（不过目前还是主要在用 iBooks）| docker | 

## 网络
服务部署可以通过 docker 轻松搞定，然而绝大部分服务都是通过 Web UI 进行交互的，所以我们还需要找到一个快捷的方案来从内网或外网访问这些服务。

在本文初次完成的 2021 年，我主要使用 frp 来完成内网穿透，但在 2024 年，能够选择的解决方案就非常多了，我目前使用过的方案有：
* Cloudflare IPv6 直连（可选 proxy）
* Cloudflare Tunnel
* frp
* Tailscale

### HTTP
处理 DNS 解析时，为了访问方便，我在内网中通过 `dnsmasq` 代理了顶级域名 `s.` 的解析请求：将 `*.s` 的域名全部解析到主机的固定 IP，针对某些希望暴露到公网的服务，我会将 `*.stdioa.com` 的域名解析到很久很久以前买的腾讯云学生机上。
针对 HTTP(S) 的路由，我在内网中使用 Caddy 作为网关；公网中我使用了 Nginx 作为网关，然后使用 frp 将公网的请求导入到内网中，再通过内网主机上的 Caddy 网关进行路由。

请求拓扑结构大致如下：
![请求拓扑结构](/pics/homelab/topo.png)

----
> 在 2024 年对本文进行修订时，笔者认为使用 CloudFlare 将内网服务暴露到公网也是一个不错的主意，这样可以省下一台 VPS 的钱，也无需再申请 IPv4 的公网 IP，只需 IPv6 的 IP 即可。不过某些小众宽带运营商和偏远地区的移动网络对 Cloudflare 的可访问性并不够高，因此需要结合使用地区的实际情况来选择合适的方案。
> 
> 使用 Cloudflare 主要有几种姿势：
> * DNS 直连（可以使用 Cloudflare 的 Proxy 来避免源站暴露，并为只有 V6 公网 IP 的家庭宽带添加双栈支持）
>   注：如果用 IPv6 直连，则大概率需要搭配 DDNS 使用，并确认路由器和光猫对 IPv6 防火墙的支持程度。可以参考[以前的文章](/2019/02/build-file-manager-on-ipv6/)。
> * Cloudflare Tunnel（无需暴露源站端口，但在国内的稳定性不佳）
> * Cloudflare Worker + 优选 IP（我尝试过但延迟过高，或许是使用姿势不对）
> 
> 关于这几种方案的具体使用方法，网上已有相当多的文章，本文不再赘述。目前我主要使用的是前两种方案，frp 虽然依然在运行，但几乎不再承接 HTTP 流量。
----

这样，Web 访问的问题就解决了，但偶尔还会有一些特殊的情况发生。

### SSH
在外偶尔会有一些针对 NAS 或树莓派的运维需求，或者可能只是连到树莓派上去编辑一下 beancount 的交易记录，此时就需要通过 SSH 来连接到主机，通过 shell 来进行操作。  
此时我同样用到了 frp：建立一个 TCP 代理，将 NAS 的 22 端口映射到腾讯云上的某端口，即可通过 SSH 来连接这个端口。

----
> 在运维方面，SSH 是非常重要的连接手段，为了保证链路的可用性，我构建了多条不同的链路：
> 1. 通过 Tailscale VPN 直接连接；
> 2. 通过 frp 建立一个 TCP 代理，将 NAS 的 22 端口映射到腾讯云上的某端口，通过腾讯云做跳板连接，通常用于使用别人的电脑（或工作机）时进行临时登录；
> 3. 通过 [Cloudflare Zero Trust](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/use-cases/ssh/#connect-to-ssh-server-with-cloudflared-access) 在浏览器中启动终端。
> 
> ~~虽然做了链路冗余，但停电和断网的风险并没有被考虑在内。~~
----

### 在外访问内网服务
出于安全考虑，我不会把涉及到隐私的服务（比如 Fava）暴露在公网上。但有时还是要访问一些不想暴露到公网的服务，或临时搭建的服务，所以需要想办法访问到内网主机的端口。

在本文初次撰写时，我想到了三种端口映射方案：
* 配置 frpc tcp 代理规则，临时将内网的端口暴露到公网
* 通过 SSH 隧道，将内网的端口映射到本地
* 如果你恰巧正在[使用 VSCode 通过 SSH 进行远程开发](https://code.visualstudio.com/docs/remote/ssh)，那也可以用 VSCode 来快速配置[端口映射](https://code.visualstudio.com/docs/remote/ssh#_forwarding-a-port-creating-ssh-tunnel)。

但其实我们还可以配置 VPN 作为终极解决方案。  
在写这篇文章时，我恰好接触到了 [Tailscale](https://tailscale.com/). 它是一个安全的虚拟组网产品，配置十分简单，且拥有相当高的 NAT 穿透成功率（此处推荐[一篇好文](https://tailscale.com/blog/how-nat-traversal-works)）。因此，如果需要高频访问内网服务的话，可以考虑直接使用 Tailscale 来完成。  
不过需要注意的是：tailscale 是一个商业产品。如果你对自己的数据安全十分担忧，可以考虑使用 [nebula](https://github.com/slackhq/nebula), [Zerotier](https://www.zerotier.com/) 或 [headscale](https://github.com/juanfont/headscale) 这样的自建组网方案作为替代。

----
> 在 2024 年重新回顾这篇文章时，由于 Tailscale 优秀的体验，我已经完全使用 Tailscale 来访问内网服务了。主要的方案如下：
> * 在 Tailscale 后台配置子网（subnets），将 NAS 设置家中内网网段的网络出口，这样家中设备（如路由器、NAS 和树莓派）均可以经由 NAS，通过内网 IP 直接访问；
> * 在电脑上安装 dnsmasq，在解析内网域名 `*.s` 时，使用家中路由器的地址作为上游 DNS 服务器。
>
> 使用这种方案，除了可以无痛访问内网的 Web 服务外，甚至还可以通过 NVIDIA GameStream 和 Moonlight 为游戏进行串流，实现在 MacBook Air 上打 PC 游戏的梦想。在 1080P 60 帧的分辨率下，游戏延迟只有 30ms，甚至连绝区零都可以勉强玩起来。虽然这套方案比不上米哈游自己的云游戏，但至少不用花钱。
----

## 数据备份
俗话说：“备份不做，十恶不赦”。但不得不承认我的经济实力和盘位都十分有限，没有办法做 RAID，所以只好选择性地备份某些比较重要的数据，如 Gitea Repo，WebDAV 中的笔记内容，以及 PC 内的文件等。

热备份的手段比较简陋：
1. PC 中的文件直接使用“文件历史记录”功能，通过 SMB 协议备份到 NAS 中；
2. 通过 systemd Timer 定时触发脚本，用 `rsync` 将本机或腾讯云上的数据增量同步到备份位置。

对于冷备份，我使用了 [Cobian Backup](https://www.cobiansoft.com/) 来定期将数据备份到一块离线硬盘。

# 如今
这台主机玩到现在，我觉得它已经不是一台 NAS 了，而是更像一个 HomeLab.  
对于我来说，HomeLab 是一个可以让人快速实现某个想法或需求的平台。有了 Linux 和 Docker，面对日常生活中某些天马行空的需求时，我可以快速实现、快速部署。在享受 HomeLab 给我带来便利的同时，我还可以享受折腾它带来的~~奇怪的~~成就感。  

举个例子，三月份的时候玩塞尔达发现了一个[旷野之息的地图网站](http://16p.top/)，然而由于网站的 CDN 带宽极为有限，导致这个服务的性能奇差，加载一张图片可能需要十几秒到一分钟的时间。于是我在某个晚上花了不到一个小时写个惰性的本地缓存服务，并将它部署在了内网中。此后，当我需要在海拉鲁大陆的角落寻找宝藏时，就可以享受到地图秒开的快感。

除此之外，HomeLab 的另一大价值在于它让我们拥有了数据和功能的**所有权**。  
现代互联网的商业魔爪已经伸向了每个用户，几乎所有的互联网服务都需要以金钱和/或隐私作为代价来换取使用权。如果美好的万维网早已不在，我们的力量也无法支撑我们颠覆商业规则，那我们是否有可能在互联网中搭建一个只属于自己或一小撮人的庇护所呢？  
对于我来说，通过 HomeLab 自建一些服务（如 Tiny Tiny RSS 或 BitWarden），将云服务私有化，让个人隐私留在自己的硬盘之中，可能是最好的解决方案了。

最后，我想为大家推荐 [awesome-selfhosted](https://github.com/awesome-selfhosted/awesome-selfhosted) ，大家可以在这里查找并选择合适的自建服务。

# 相关阅读
* [云服务器都99一年了，除了买来吃灰，你还能用来装这些免费云软件](https://wiki.ftqq.com/howto-use-your-99vps)
* [聊聊你的私有云 - Small talk](https://anchor.fm/tech-small-talk/episodes/ep-enlm2q)
