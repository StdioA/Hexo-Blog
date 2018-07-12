title: 搭建私有KMS服务器
date: 2015-09-02 22:26:00
categories:
- 乱七八糟
tags:
- KMS
toc: true

---

最近win10又频繁提示“许可即将过期”，之前找到的kms服务器又挂掉了，于是就决定自己搭一个自己用。

<!-- more -->

# 1. 引子
好久没更新日志了，暑假在家每天过的乱七八糟，就没怎么研究技术，于是博客也没怎么更新…今天随便写点什么凑个数\_(:зゝ∠)\_

# 2. 科普
KMS (Key Management Service), 密钥管理服务，是一种对Windows及Office产品进行批量授权的服务，通常被部署在大型企业局域网中，用于对批量授权版（即VOL版）Windows系统进行大批量的激活。KMS服务器的作用是给局域网中的所有计算机的操作系统提供一个有效的产品序列号，然后计算机里面的KMS服务就会自动将系统激活。每一个由KMS Server提供的序列号的有效期只有180天，而不是其他版本的永久使用一个序列号。所以操作者必须在快到期的时候在此手动连接KMS服务器让它提供一个新的序列号，否则180天以后就会回到试用版本状态。由于KMS系统部署较为容易，所以在国内很多人通过MSDN等渠道下载VOL版本的软件，然后通过KMS服务进行激活，已达到盗版的目的。

# 3. 搭(luan)建(gao)过程
昨天在满大街乱找野生KMS服务器的时候发现了一个帖子：[使用KMS激活windows系统及VL-office系列](http://wrlog.com/activate-kms-vlmcsd.html), 里面提供了一个[链接](http://forums.mydigitallife.info/threads/50234-Emulated-KMS-Servers-on-non-Windows-platforms)指向一个帖子，帖子中提供了Python和C版本的KMS服务器模拟器，可以在自己的服务器中进行KMS服务器部署。于是把代码搞了下来（为此还注册了个账号, python版的代码在[这里](https://mega.co.nz/#F!6pIGEbhQ!DE2twA7dVG5C4knjAq56zQ))，用VSFTP将代码传到了自己的树莓派上，然后运行`python server.py`进行部署。
    
    python server.py
    TCP server listening at 0.0.0.0 on port 1688.

然后在Windows系统中打开具有管理员权限的命令提示符，输入`slmgr -skms 192.168.155.2:1688`设置KMS服务器地址（地址可以更换），然后输入`slmgr -ato'进行系统激活，此时服务器端显示：

    Connection accepted: 192.168.56.1:13023
    Received V6 request on Wed Sep  2 22:59:55 2015.
    Connection closed: 192.168.56.1:13023

Windows系统提示“成功地激活了产品”，激活成功。

在树莓派上部署成功以后，随手在VPS上部署了一份以备用。帖子中提了一种设置Linux系统启动项来使KMS服务器开机自动部署的方法，不过自己没有这个需求就没搞。

# 4. 总结
这篇文章好像有点水，主要是因为自己实在没什么东西写了…自己的暑假过的乱七八糟，浪费了很多时间在游戏上，没什么心情研究技术。现在开学了，有更多时间来钻研技术了，收收心找找状态，以后博客会定期更新的，我对树莓派电源灯发誓→_→

# 5. 各种Key

最后附上Office 2016和Windows 10的VOL版激活码，其它版本软件激活码请自行百度。

> Office Professional Plus 2016 - XQNVK-8JYDB-WJ9W3-YJ8YR-WFG99
> Office Standard 2016 - JNRGM-WHDWX-FJJG3-K47QV-DRTFM
> Project Professional 2016 - YG9NW-3K39V-2T3HJ-93F3Q-G83KT
> Project Standard 2016 - GNFHQ-F6YQM-KQDGJ-327XX-KQBVC
> Visio Professional 2016 - PD3PC-RHNGV-FXJ29-8JK7D-RJRJK
> Visio Standard 2016 - 7WHWN-4T7MP-G96JF-G33KR-W8GF4
> Access 2016 - GNH9Y-D2J4T-FJHGG-QRVH7-QPFDW
> Excel 2016 - 9C2PK-NWTVB-JMPW8-BFT28-7FTBF
> OneNote 2016 - DR92N-9HTF2-97XKM-XW2WJ-XW3J6
> Outlook 2016 - R69KK-NTPKF-7M3Q4-QYBHW-6MT9B
> PowerPoint 2016 - J7MQP-HNJ4Y-WJ7YM-PFYGF-BY6C6
> Publisher 2016 - F47MM-N3XJP-TQXJ9-BP99D-8K837
> Skype for Business 2016 - 869NQ-FJ69K-466HW-QYCP2-DDBV6
> Word 2016 - WXY84-JN2Q9-RBCCQ-3Q3J3-3PFJ6
> 
> Windows 10 Home - TX9XD-98N7V-6WMQ6-BX7FG-H8Q99
> Windows 10 Home N - 3KHY7-WNT83-DGQKR-F7HPR-844BM
> Windows 10 Home Single Language - 7HNRX-D7KGG-3K4RQ-4WPJ4-YTDFH
> Windows 10 Home Country Specific - PVMJN-6DFY6-9CCP6-7BKTT-D3WVR
> Windows 10 Professional - W269N-WFGWX-YVC9B-4J6C9-T83GX
> Windows 10 Professional N - MH37W-N47XK-V7XM9-C7227-GCQG9
> Windows 10 Education - NW6C2-QMPVW-D7KKK-3GKT6-VCFB2
> Windows 10 Education N - 2WH4N-8QGBV-H22JP-CT43Q-MDWWJ
> Windows 10 Enterprise - NPPR9-FWDCX-D2C8J-H872K-2YT43
> Windows 10 Enterprise N - DPH2V-TTNVB-4X9Q3-TJR4H-KHJW4
> Windows 10 Enterprise 2015 LTSB - WNMTR-4C88C-JK8YV-HQ7T2-76DF9
> Windows 10 Enterprise 2015 LTSB N - 2F77B-TNFGY-69QQF-B8YKP-D69TJ

_可以考虑有空搞一个Office 2016来玩玩。_

# 6. 参考文档
1. [KMS - 互动百科](http://www.baike.com/wiki/KMS)
2. [使用KMS激活windows系统及VL-office系列](http://wrlog.com/activate-kms-vlmcsd.html)
3. [Emulated KMS Servers on non-Windows platforms](http://forums.mydigitallife.info/threads/50234-Emulated-KMS-Servers-on-non-Windows-platforms)
