---
title: 由 TT-RSS 解析数据库地址失败引出的一个问题
author: David Dai
tags:
  - Docker
  - Alpine
  - 树莓派
categories:
  - 乱七八糟
date: 2021-02-23 20:51:14
toc: true
---

水一篇文章，主要用来告诫自己认真看文档。:new_moon_with_face:

<!--more-->

## 背景
下午随手在树莓派上升级了一下 [TT-RSS](https://tt-rss.org/) 的[镜像](https://hub.docker.com/r/wangqiru/ttrss)，然后它当场爆炸，看了容器日志告诉我 PHP 无法解析数据库的域名 `database.postgres`. 

## 尝试解决
进到容器里尝试手动解析一下，但是报错 `nslookup: clock_gettime(MONOTONIC) failed`.  
用自己的另一台运行 Debian testing 的 x86 机器运行了一下，无法复现这个问题。

Google 了一下找到 [Alpine 的一个 issue](https://gitlab.alpinelinux.org/alpine/aports/-/issues/12091)，简单看了一下发现是 Alpine 3.13 升级了 musl，使用了新的系统调用 `clock_gettime64`. 在容器里跑了下 date，结果如下：

```bash
$ docker run --rm -it alpine date
Sun Jan  0 00:100:4174038  1900
```

看起来，Alpine 3.13 需要 Docker 19.03.9. 然而我的 docker 版本已经是 20.10.3 了，但依然无法运行。  
Google 了一大圈发现树莓派上安装的 `libsecomp2` 太老（2.3.3-4），不支持 time64，尝试运行 `scmp_sys_resolver -a arm clock_gettime64` 返回 `-1` 也验证了这个观点。    
*PS：需要安装 `seccomp` 包。*

所以需要安装更新版本的 `libseccomp2`，但 Raspbian 不提供新版的包。所以，我从 [Debian 软件包目录](https://packages.debian.org/bullseye/libseccomp2) 找了新版（2.5.1-1）来安装，问题解决。

## 不仔细看英文文档的后果
Alpine 的 release notes 已经写明了问题：

> Therefore, Alpine Linux 3.13.0 requires the host Docker to be version 19.03.9 (which contains backported moby commit 89fabf0) or greater and the host libseccomp to be version 2.4.2 (which contains backported libseccomp commit bf747eb) or greater.

> Therefore, the following platforms are not suitable as Docker hosts for 32-bit Alpine Linux 3.13.0, due to containing out-of-date libseccomp: Amazon Linux 1 or 2, CentOS 7 or 8, Debian stable without debian-backports, Raspbian stable, Ubuntu 14.04 or earlier, and Windows. This applies regardless of whether the Linux distribution Docker packages or separate Docker package repositories are used.

树莓派的 Docker 能够满足条件，但因为运行的是 Raspbian stable，所以 `libseccomp` 的版本无法满足；Debian testing 的机器是 64 位的，所以可以正常运行。

## 结论
Alpine 升级了 musl → 使用了新的系统调用 → 如果系统是 32 位版本，且 Docker 和 `libsecccomp2` 版本较低，则无法在容器中正常获取时间，进而影响到各种功能。

1. 升级 Docker 到 19.03.9 及以上；
2. 升级 `libseccomp2` 到 2.4.2 及以上，如果官方软件源没有提供的话可以去 testing 库里找找；
3. 解决这个问题花了大半个小时，但如果仔细看文档的话估计 5 分钟就能搞定。**所以，要认真看英文文档。**

## 参考文档
* [Release Notes for Alpine 3.13.0](https://wiki.alpinelinux.org/wiki/Release_Notes_for_Alpine_3.13.0#time64_requirements)
* [alpine 3.13 在 Armhf docker 的网络问题](https://amefs.net/archives/2199.html)
* [Docker 使用 seccomp 无法获取系统时间的 bug 一则](https://harrychen.xyz/2020/08/07/docker-seccomp-time-bug/)
