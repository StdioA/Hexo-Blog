title: 第一届 NUAACTF 非官方 Writeup
date: 2016-04-25 12:01:04
categories:
- 乱七八糟
tags:
- CTF
- 脑洞
- Python
toc: true

---

学校组织了第一届 NUAACTF，参加去玩玩，开开脑洞:joy:，嗯，还蛮好玩的~

<!-- more -->
# Web 0
**huan'ying'lai'dao**

签到题，F12 即可得到 flag.
Flag 藏在页面的某个 `js` 文件里，用 jsfuck 混淆了。

# Web 1
**百度一下，你就知道**

题目貌似改过一回。

改之前直接点击网页上的一个链接，跳转到某个网页，源代码里面有一段 php，具体内容不记得了，里面
有一串 MD5 `095a655fc809cbbdffa207717a5233f5`. Google 一下找到某白菜的博客，得到原文 `bnVhYWN0ZiU3Qi93ZWIyL2NlYmE2ZmJiZjBlZGU0MzI1MjY0MWNkMzM2ZTM2YTAzJTdE`.  
看起来像 `Base64`，于是 decode 一下，然后 `decodeURIComponents` 得到 flag.

后来题目改了，直接按照题目要求去百度就能直接找到某白菜的博客。

# Web 2
**不是 bug 是 feature**

Web 2 的入口点在 Web 1 的 flag 里。进去以后会跳转到某个 php, 然而某白菜又把 php 源码写进去了，思路跟上个题好像一样，也是去百度一下 MD5 得到密码，添加 `GET` 参数即可拿到 flag.

# Web 3
**笨笨的程序猿**

Web 3 的入口点在 Web 2 的 flag 里。进去以后跳转到某 php, 里面只有一个登录表单。看起来像是 SQL 注入，就用 `sqlmap` 扫了一下，把表 dump 出来，发现有两个账户 admin 和 user，密码都是一坨看起来什么都不像的东西，然后线索就断了。  
最后一个小时决定手动注入一下，然而并不怎么会 SQL 注入于是就找了几篇博客来看。注入了半天密码框均无果，最后试了一下注入用户名输入框 `admin' or '1'='1`，成功登录。  
登录以后 flag 一闪而过，本来想截屏然而懒得截了就抓了包，拿到 flag.

# Web 4 (未解出)
**你从哪里来，我的朋友。**

Web 3 成功登录后就跳转到 Web 4，网页写着"0nly welc0me pe0ple who c0me from http://cs.nuaa.edu.cn "。结合“你从哪里来”，想到了请求来源，但是我只想到了 `Referer`, 修改后无果。  
比赛结束后听说要改 `Origin` 和 `X-Forwarded-For`, 然而改了也没做出来，不知道哪里出了问题。

# Rev 1
**曾老师的 Android 课**

把 apk 下载下来，用 `jd-gui` 打开，一阵乱翻以后发现 flag 藏在 `MainActivity.class` 里。

# Rev 2 (未解出)
**奇妙的声音**

又看见安卓了，下下来以后一阵乱翻，没有头绪。

刚才看了别人的题解，发现 flag 在资源目录里面。  
拿出来 `res/raw/sound.wav`, 用随便什么鬼打开（来之前见识过这个脑洞，专门准备的 Sonic Visualiser, 然而卡住了没有用上），发现里面有四个声道，下面两个声道像是个方波。

![Rev 2](/pics/nuaactf/rev2.png)

然后一点一点数 01 数出来，得到：
> 01101110
> 01110101
> 01100001
> 01100001
> 01100011
> 01110100
> 01100110
> 01111011
> 01110011
> 01101000
> 00110000
> 01110010
> 01110100
> 01011111
> 01100110
> 00110001
> 01000001
> 01100111

数格子数的眼都要花了…  
每一行 8 位二进制转 ASCII 码，得到 flag `nuaactf{sh0rt_f1Ag`, 少了一个右括号:joy:  
真是幸亏 flag 短…

顺便说一句，那个音频真好听，听说是锤子手机的默认铃声:flushed:

# Rev 3
**不喜欢写界面的白菜哥**

`.NET` 逆向，直接拖到 IDA 里面，一阵乱翻，翻到了三段 base64 和一个摩尔斯电码表。把三段 base64 拼起来以后 decode, 得到一段摩尔斯电码 `-. ..- .- .- -.-. - ..-. 2D3f ... .- .--. .--. -.-- -.-. .-. .- -.-. -....- .toc: true

---- -. --. -.-. ... ... .- .-. .--. -...--.-` (2D3f 是什么鬼:joy:)，解码后得到 flag `NUAACTF{SAPPYCRACK1NGCSSARP}`，放到源程序里面，显示 flag 正确。

然后，提交以后说 flag 不对！  
为什么 flag 不是 Happy Cracking CSharp? 看了半天，又交了半天的 flag，没找出问题。后来发现程序编码表里面 H 和 S 的编码一样:joy: 所以里面所有的 S 换成 H，源程序都会提示 flag 正确…  
把 H 和 S 转换，一个一个试，最后试出来 flag 是 `NUAACTF{HAPPYCRACK1NGCHHARP}`，也是醉人。

# Pwn 1
**乱码！乱码！**

签到题。下载下来一个 txt，打开以后发现是 jsfuck，放到浏览器里运行得到 flag.

# Pwn 2 （未解出）
**回家的路**

CTF 竞赛考了算法…最短路…真是醉醉哒。  
一开始没有细想，瞎写了一个深度优先搜索，出来好多解，然而没有一个是对的，也是悲催。  
到最后也没有做。

# Misc 1
**奇怪的压缩包**

嘛，隐写的题都蛮好玩哒:stuck_out_tongue_winking_eye:

[misc1.rar](/pics/nuaactf/misc1.rar) 下载下来后无法打开，`file` 一下发现是个 PNG，改拓展名后打开，得到 flag。这个也算签到题吧。

# Misc 2
**奇怪的图片**

![misc2.png](/pics/nuaactf/misc2.png)

下载 `misc2.png`, `pngcheck` 提示 `additional data after IEND chunk`, 看来最后一块后面还有东西，用记事本打开，发现文件最后用明文写着 key.

刚刚在文件末尾发现了一个文件头 `PK\x03\x04`, 觉得应该是个 zip 文件，`binwalk` 一下然后用 `dd` 分开，得到一个 zip 文件，解压后得到 flag. 这应该是这个题的标准做法吧，只是为什么做 zip 文件的时候没有压缩:new_moon_with_face:

# Misc 3
**更奇怪的图片**（你们这起的都是什么名字）

![misc3.png](/pics/nuaactf/misc3.png)

下载下来（听说这是舰娘？），`pngcheck` 没有错误，卡了一会。后来用 PS 打开，把色阶拉低，发现图的左下角有个二维码，扫码得 `QlpoOTFBWSZTWXhAk1kAAAtfgAAQIABgAAgAAACvIbYKIAAigNHqNGmnqFMJpoDTEO0CXcIvl9SeOAB3axLQYn4u5IpwoSDwgSay`.  
解码后发现是乱码，不知道怎么做了，卡了很久:mask: 比赛的最后一个小时看见了解码后字符串开头的 BZ，突然意识到这可能是个文件头，于是将解码后的内容写入文件，`file` 一下发现果然是个 `bz` 压缩文件，解压后得到 flag.

# Misc 4 （未解出）
**讨厌的 APP**

觅动校园什么鬼？！


# 总结
比赛五个小时，解了 10 道题，前半小时解出来 5 道，后面各种卡，最后一个小时脑洞大开又解出来 3 道题:joy:
CTF 真好玩，考到了各种姿势各种脑洞，还是蛮有意思哒~  
然而五个小时的比赛真的太累了…血的教训，下次题目做不出来还是要跑出去歇一会再回来做\_(:зゝ∠)\_
