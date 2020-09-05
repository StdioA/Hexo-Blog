title: 开始使用 Beancount
author: David Dai
tags:
  - 记账
  - Beancount
  - fava
  - Python
categories:
  - 乱七八糟
date: 2020-09-05 15:15:00
toc: true
---
使用 Beancount 记账已经有将近两个月了，简单写一写我都做了什么。

<!--more-->
注：本文只是一个流水账，并不是一个 Beancount 使用教程，如果想详细了解 Beancount 的话，可以参考下面提到的那些文章。

## 一些背景
### Beancount 是什么
如上文所说，[Beancount](https://beancount.github.io/) 是一个记账工具，更准确些来讲，是一个[**复式记账**](https://zh.wikipedia.org/zh-cn/%E5%A4%8D%E5%BC%8F%E7%B0%BF%E8%AE%B0)工具。但直到我写这篇文章的时候才发现官方将其定义为“一种复式记账计算机语言”。

简单来讲，它可以让你以纯文本方式记账，并通过一种[类 SQL 的语言](https://beancount.github.io/docs/beancount_query_language.html)来对交易进行查询。记账文件还可以配合 Git 进行版本控制。  
此外，Beancount 官方提供了一个名叫 [fava](https://beancount.github.io/fava/) 的图形化管理工具，它基于 Web，能够提供比原生页面更加丰富的内容，一般记账所需要的信息一目了然。想体验的同学可以在[官方提供的 Demo](https://fava.pythonanywhere.com/example-beancount-file/income_statement/) 中简单感受下。

### 为什么要用它
我从 17 年起就在用[口袋记账](https://www.qeeniao.com/) App 来随手记账。但毕竟人家也是商业公司，所以 App 里逐渐出现了不少的理财产品 :new_moon_with_face:，而且 App 打开越来越卡了。于是就开始考虑其它的记账方案。  
想了想，我的记账需求有：
1. **快速**对刚刚产生的交易进行记录，一定要比口袋记账快！
2. 查看各个粒度的收支统计
3. 查账方便
4. 数据安全（要能够备份）

之前有了解过 GNUCash，但是不太感冒。在这期间也多次想过自己去写一个记账工具，但一直懒得动手造轮子。之前在网上看到一个玩笑说“每一个程序员都想过去写一个记账工具”，我也不例外。  
直到七月份，看到赵神在他的 Channel 里[提到了 Beancount](https://t.me/TripleZsChannel/1370)，在他的[安利](https://pastebin.pl/view/raw/3d2d4756)下，我发现它非常轻便，并且能够满足我的需求。于是在某个周日的傍晚的冲动之下，决定把口袋记账的数据转移到 Beancount 上，从此用它来记账。

## 看了些东西
在动手之前也顺着赵神的安利去看了一下他人写的博客，并了解了一下会计恒等式，以及复式记账的基本知识。以下是我看过的博文。

* BYVoid 写的 Beancount 系列文章：
  * [Beancount复式记账（一）：为什么](https://byvoid.com/zhs/blog/beancount-bookkeeping-1/)
  * [Beancount复式记账（二）：借贷记账法](https://byvoid.com/zhs/blog/beancount-bookkeeping-2/)
  * [Beancount复式记账（三）：结余与资产](https://byvoid.com/zhs/blog/beancount-bookkeeping-3/)
  * [Beancount复式记账（四）：项目管理](https://byvoid.com/zhs/blog/beancount-bookkeeping-4/)
* [Beancount —— 命令行复式簿记](https://wzyboy.im/post/1063.html)
* [Beancount 最佳实践](https://xwartz.xyz/blog/beancount-best-practice/)
* [Beancount使用经验 —— 通过Beancount导入支付宝&微信csv账单](http://lidongchao.com/2018/07/20/has_header_in_csv_Sniffer/)

## 做了点工作
看了些东西以后，对 Beancount 的使用，以及记账项目管理有了基本的概念，就开始动手了。

我决定把 Beancount 安装在我的树莓派上，用 fava 做统计分析，用 Git 做版本控制及数据备份。  
Git 远端同样在树莓派上，是一个 [Gitea](https://gitea.io/zh-cn/) 实例。  
*写到这还想起来去年从 Gogs 迁移到了 Gitea，当时记了下过程但是内容太少了，就没有发出来。*:joy:

### 安装及部署
`python3 -m pip install beancount fava`，没什么特别的。  
之前很多服务都是用 Docker 部署的，但是我的 Python 服务一直没有用 Docker 部署，而是直接装在了系统里，用 systemd 托管，可能是我的某种执念吧。

Beancount 本身不需要守护进程，因为记账文件是直接用文本存储在系统中的，beancount 只是用来做查询。不过 fava 倒是可以以守护进程的方式部署起来。

```systemd
[Unit]
Description=Fava
Documentation=https://beancount.github.io/fava/
After=network-online.target
Wants=network-online.target systemd-networkd-wait-online.service

[Service]
Restart=on-abnormal
User=pi
Group=pi
Environment=HOME=/home/pi
ExecStart=/home/pi/.local/bin/fava -p 6666 main.bean
WorkingDirectory=/home/pi/projects/accounting

[Install]
WantedBy=multi-user.target
```

### 数据导入
万幸口袋记账提供了导出功能，可以将用户的全部交易记录以 xls 的格式导出。  
导出之后，我写了一个极丑无比的 Python 脚本，将每条记录渲染成 Beancount 的交易格式，然后写入到 `main.bean` 文件中。

*最后平了下账，平出一比巨款来，可见之前记的账是多么多么不靠谱*:new_moon_with_face:

在导入并与口袋记账上的统计数据作比对，确认数据无误后，我对项目结构做了一些调整：
* 将账户声明语句按照账户类型组织，单独存放在 `accounts/{类型}.bean` 中；
* 将前几年的交易进行归档，存放在 `txs/{年份}.bean` 中；
* 在 `main.bean` 中用 `include` 语句导入 `accounts/*.bean` 和 `txs/*.bean`。

这样，就有了以下的项目结构：

```
.
├── accounts             // 账目定义文件
│   ├── assets.bean
│   ├── equity.bean
│   ├── expences.bean
│   └── income.bean
├── document
│   └── initialize       // 初始化导入的脚本及口袋记账数据
│       ├── initial.py
│       └── ...
├── txs                  // 历史交易数据
│   ├── 2017.bean
│   ├── 2018.bean
│   └── 2019.bean
├── main.bean            // 主文件，包含当年（2020）的交易数据
└── Makefile
```

### Telegram bot
用 vim 记了几笔账之后感觉还是不怎么方便，毕竟在外面玩的时候，需要频繁用手机 ssh 回家记账也不是那么回事。  
于是，我就写了一个简单的 Telegram bot，通过与 bot 交互来编辑树莓派上的 `main.bean` 文件。平时消费时，只需要按照之前在口袋记账里记账的思路，将几个关键词填好，bot 就可以渲染出这笔账对应的交易文本，并将其追加到 `main.bean` 文件中。这样就能涵盖绝大多数日常在外的使用场景了。不过，去超市或 AA 这种复杂一些的账，还是没办法用 bot 来记，不过也无所谓，我会简单记下来，然后回家用 vim 或 fava 把账补好。  
此外，还实现了两个 bot 命令，方便在外查看自己花了多少钱（虽然平时基本用不上 :joy:）。

最后的实现效果：

<img style="max-height: 60vh;" src="/pics/beancount-bot.jpg" />

~~顺便安利一下，FAT PHO 的河粉相当好吃~~

### 除此之外
每天睡觉前会将当日的更改提交，并做好备份；顺便做了 `pre-commit` 的 Git 钩子，用于规范文件格式；  
每个月结束的时候都会为所有资产账户添加断言，如果断言的数值与 beancount 里计算的数值不一致的话，beancount 就会抛出错误。不过报错也是在所难免，因为之前记账记得够细，所以对账也相当方便，偶尔还可以看到一些因为 bot 太笨而闹的笑话：

```diff
2020-08-30 * "自行车" ""
  Assets:Digital:支付宝:Deposit                               -3.50 CNY
-  Assets:Bank:交通银行:Card:Deposit
+  Expenses:Daily:交通出行
```

## 总结
这套基于 Beancount 的记账系统已经运行了快两个月了。从口袋记账迁移到 Telegram 记账后，整个记账流程都顺滑和快捷了很多。之前等待口袋记账 App 打开、卡死、看广告所浪费的时间，放到现在可以让我打开 Telegram 记完一笔账。  
所以，还是要吹一波开源社区才对。:white_check_mark:

如果有什么问题，欢迎与我交流。
