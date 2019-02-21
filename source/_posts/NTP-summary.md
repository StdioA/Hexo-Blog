title: NTP 简介
author: David Dai
tags:
  - Linux
  - NTP
categories:
  - 乱七八糟
date: 2019-02-12 19:50:00
toc: true
---
昨天遇到了一个神奇的问题，最后发现是服务器的 ntpd 没开导致本地时间没有同步:joy: 正好了解一下 NTP.

<!--more-->

## NTP 协议
NTP 协议用于在网络之中通过分组交换进行时钟同步。基于 UDP，使用 123 端口。  

### 协议实现
客户端和服务器间会通过修改版的 [Marzullo 算法](https://en.wikipedia.org/wiki/Marzullo%27s_algorithm) 完成时间同步。

在传递时间时，服务器会给出 64 位的时间戳，浮点精度为 32 位。这个时间戳每 2^32 秒会翻转一次，理论分辨率为 2^-32 秒。时间戳以 1900 年 1 月 1 日作为开始时间。

NTP 时间源会进行分层，通过阶层 n 同步的服务器将运行在阶层 n+1. 分层机制用来防止循环请求。阶层 0 的服务器与高精度计时设备（如原子钟）相连，也成为基准时钟。

## 使用 NTP 同步 Linux 系统时间
### ntpd
ntpd 是某些 Linux 发行版自带的 NTP 同步工具，它通常在后台运行，与授时服务器进行时间同步。  
ntpd 默认使用 `/etc/ntp.conf` 作为配置文件。配置方法及范例可以参考 [Linux System Administrators Guide](https://www.tldp.org/LDP/sag/html/basic-ntp-config.html)

不过需要注意的是，启动 ntpd 并不会立即纠正本地时间，而是会缓慢地进行时间同步。

> Be patient! A simple test is to change your system clock by 10 minutes before you go to bed and then check it when you get up. The time should be correct.

### ntpdate
ntpdate 是一款使用 NTP 协议同步本地时间的工具。如果需要快速纠正时间，可以使用 ntpdate 进行手动同步。

安装：
```bash
apt install ntpdate
```

使用：
```bash
$ sudo ntpdate cn.ntp.org.cn
12 Feb 13:19:05 ntpdate[27659]: adjust time server 119.28.183.184 offset -0.007086 sec
```

注意：由于 NTP 协议使用固定端口，在使用 ntpdate 时，需要关闭 ntpd 服务。

### systemd-timesyncd
[systemd-timesyncd](https://www.freedesktop.org/software/systemd/man/systemd-timesyncd.service.html) 是 [timedated](https://www.freedesktop.org/wiki/Software/systemd/timedated/) 提供的时钟同步守护软件。它可以通过 `systemd-timesyncd.service` 服务启动。

`timesyncd` 的配置文件位于 `/etc/systemd/timesyncd.conf`，格式如下：
```
[Time]
NTP=cn.ntp.org.cn 0.cn.pool.ntp.org # 主要 NTP 服务器
FallbackNTP=0.arch.pool.ntp.org 1.arch.pool.ntp.org 2.arch.pool.ntp.org 3.arch.pool.ntp.org # 备用 NTP 服务器
RootDistanceMaxSec=5
PollIntervalMinSec=32
PollIntervalMaxSec=2048
```

配置完成后，运行 `sudo timedatectl set-ntp true` 启用时间同步服务 `timedatectl status` 可以查看当前同步设置，`timedatectl timesync-status` 可以查看当前时间同步服务的运行状态，包括时间延迟、误差、至今同步次数等。

## 国内常用的 NTP 地址（IP 池）
* [ntp.org.cn](http://ntp.org.cn/index.php)：cn.ntp.org.cn
* [NTP Pool Project](https://www.ntppool.org/zone/cn)：0.cn.pool.ntp.org
* 阿里云 NTP：time.pool.aliyun.com

更多地址可参考[国内常用NTP服务器地址](https://www.jianshu.com/p/28864ab7fdd9)

## 其它参考文档
https://wiki.archlinux.org/index.php/Systemd-timesyncd_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)
https://zh.wikipedia.org/wiki/%E7%B6%B2%E8%B7%AF%E6%99%82%E9%96%93%E5%8D%94%E5%AE%9A
https://wiki.archlinux.org/index.php/Systemd-timesyncd_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)
https://www.freedesktop.org/software/systemd/man/systemd-timesyncd.service.html
