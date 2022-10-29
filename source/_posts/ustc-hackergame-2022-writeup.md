title: USTC Hackergame 2022 玩耍记录
date: 2022-10-29 22:50:00
categories:
- 乱七八糟
tags:
- CTF
- 脑洞
toc: true

---

上周日晚上，偶然看到一个原神玩家群（？）里面有人发了一张图，说是 USTC 的 CTF. 上一次玩 CTF 还是六年前，不过这次一时兴起打算玩玩看。  
由于上班比较忙，所以只玩了一天多一点，做了一些简单题。


<!-- more -->

# 注册
周日的半夜，在群里看到了一张 CTF 的图，于是兴冲冲地跑去注册，没想到直接吃了闭门羹:new_moon_with_face:。

![关门了](/pics/hackergame2022/door_closed.png)

周一早上八点，准时开干。

# Binary
## Flag 自动机
题目在[这里](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/Flag%20%E8%87%AA%E5%8A%A8%E6%9C%BA/src/flag_machine.exe)。  
打开之后可以看到一个对话框，鼠标一碰到“狠心夺取”的按钮，它就会跑掉。  
拖到 IDA 里反编译，可以看到 `sub_401510` 函数中存在输出 flag 的代码，但生成 flag 的代码被混淆过，很难看懂。  
在这个过程中，看到了出题人给的注释：“Hint: You don't need to reverst the decryption logis itself.”，那看来是需要看看别的方法。

仔细看下 `sub_401510`:
![](/pics/hackergame2022/flag_machine.png)

其中 `2u` 的 case 是“放手离开”的触发逻辑，而 `1u` 是“狠心夺取”。中间的 `0x111u` 分支需要想办法触发才能拿到 flag. 这个程序加了反动态调试，因此也无法通过打断点的方式来进入这个分支。（当然官方题解也讲了绕过动态调试的方法）。  
一番搜索（感谢浩子哥哥），找到了 Windows 的 [PostMessageA](https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-postmessagea) 函数，通过它可以向指定窗口发送事件消息。  
那么创建一个 Windows C++ 项目，运行 `PostMessageA(HWND_BROADCAST, 0x111, 3, 114514)`，即可让 flag 自动机弹窗并输出 flag.

不过，老司机告诉我可以直接用 Python:
```python
import win32gui
import win32con

hwnd = win32gui.FindWindow("flag 自动机", "flag 自动机")
win32gui.SetForegroundWindow(hwnd)
win32con.SendMessage(hwnd, 0x111, 3, 114514)
```

OK，这就是我唯一能做的逆向题了。

# General
General 就是 misc，是想当年最喜欢做的题目种类。

## 猫咪问答喵
> 参加猫咪问答喵，参加喵咪问答谢谢喵。

点进去以后发现是 6 道问答题，答对一半给一个 flag，答对全部给两个。

1. > 中国科学技术大学 NEBULA 战队（USTC NEBULA）是于何时成立的喵？

    Google “USTC NEBULA 成立时间”搜到一条[新闻](https://cybersec.ustc.edu.cn/2022/0826/c23847a565848/page.htm)，可以得知成立于 2017 年 3 月。
2. > 2022 年 9 月，中国科学技术大学学生 Linux 用户协会（LUG @ USTC）在科大校内承办了软件自由日活动。除了专注于自由撸猫的主会场之外，还有一些和技术相关的分会场（如闪电演讲 Lightning Talk）。其中在第一个闪电演讲主题里，主讲人于 slides 中展示了一张在 GNOME Wayland 下使用 Wayland 后端会出现显示问题的 KDE 程序截图，请问这个 KDE 程序的名字是什么？

    Google 一番找到了当时的[活动记录](https://lug.ustc.edu.cn/wiki/lug/events/sfd/)，看了一眼 slides 感觉是 `Dolphin`，然鹅并不是它。

    于是去看了 B 站上的[视频回放](https://www.bilibili.com/video/BV11e411M7t9/)，把这段演讲听了两三遍 :new_moon_with_face:。在视频的 2:42:06 处，主讲人念了一个单词，最后 Google 搜索联想告诉我它是 Kdenlive.
3. > 22 年坚持，小 C 仍然使用着一台他从小用到大的 Windows 2000 计算机。那么，在不变更系统配置和程序代码的前提下，Firefox 浏览器能在 Windows 2000 下运行的最后一个大版本号是多少？

    放飞自我，随便 Google，最后发现 [12 能用](https://support.mozilla.org/gl/questions/937250)，但是 [13 就不能了](
    https://groups.google.com/g/mozilla.support.firefox/c/PaBc5BH7ZB4?pli=1)。所以答案是 12.
4. > 你知道 PwnKit（CVE-2021-4034）喵？据可靠谣传，出题组的某位同学本来想出这样一道类似的题，但是发现 Linux 内核更新之后居然不再允许 argc 为 0 了喵！那么，请找出在 Linux 内核 master 分支（torvalds/linux.git）下，首个变动此行为的 commit 的 hash 吧喵！

    在 `torvalds/linux.git` 下搜索 CVE-2021-4034，得到 [Commit dcd46d897adb70d63e025f175a00a89797d31a43](https://github.com/torvalds/linux/commit/dcd46d897adb70d63e025f175a00a89797d31a43).

5. > 通过监视猫咪在键盘上看似乱踩的故意行为，不出所料发现其秘密连上了一个 ssh 服务器，终端显示 ED25519 key fingerprint is MD5:e4:ff:65:d7:be:5d:c8:44:1d:89:6b:50:f5:50:a0:ce.，你知道猫咪在连接什么域名吗？

    一开始看到这个题我有点蒙。全球那么多域名那么多服务器，怎么知道服务器的 key 是什么？做完上面的题回来看，发现可以直接 Google，看到 [Zeek 的文档](https://docs.zeek.org/en/master/logs/ssh.html) 用这个 key 做了例子，而目的 IP 是可以直接访问的。
    
    通过网页标题重新 Google，发现答案是 `sdf.org`.
6. > 中国科学技术大学可以出校访问国内国际网络从而允许云撸猫的“网络通”定价为 20 元一个月是从哪一天正式实行的？

    一开始查到了一篇[介绍收费标准的文章](https://netfee.ustc.edu.cn/faq/index.html#fee)，然而 `2011-11-01` 不是正确答案。

    看了看以前的题解，这个题是可以爆破的，于是按天对题目进行了爆破，最后得到答案 `2003-03-01`.

## 家目录里的秘密
### VS Code 里的 flag
把 home 目录下下来，在 `.config/Code` 里全局搜索 "flag"，就能拿到第一个 flag.

### Rclone 里的 flag
rclone 没有用过，不是很熟。不过看了一眼 `.bash_history`，最后有一行 `cat .config/rclone/rclone.conf`.  
hmm… 提示可以说是非常明显了。

查看 `rclone.conf`，看到密码是 `tqqTq4tmQRDZ0sT_leJr7-WtCiHVXSMrVN49dWELPH1uce-5DPiuDtjBUN3EI38zvewgN5JaZqAirNnLlsQ`. 搜索“rclone decrypt password”，找到了一个[解密脚本](https://forum.rclone.org/t/how-to-retrieve-a-crypt-password-from-a-config-file/20051)，跑一下拿到 flag.

## HeiLang
>  来自 Heicore 社区的新一代编程语言 HeiLang，基于第三代大蟒蛇语言，但是抛弃了原有的难以理解的 `|` 运算，升级为了更加先进的语法，用 `A[x | y | z] = t` 来表示之前复杂的 `A[x] = t; A[y] = t; A[z] = t`。

没啥好说的，把代码下下来做个字符串替换，然后塞回源代码里，跑一下就能拿到 flag.

## 旅行照片 2.0
社工题，有点意思，不过做完这道题的第二天，我就把我的个人签名改成了“记得取消 FR24 订阅”。

题目给了这样一张[图片](https://raw.githubusercontent.com/USTC-Hackergame/hackergame2022-writeups/master/official/%E6%97%85%E8%A1%8C%E7%85%A7%E7%89%87%202.0/src/travel-photo-2.jpg)，需要从中来挖掘有关信息。

### 第一题：照片分析
1. 图片所包含的 EXIF 信息版本是多少？
2. 拍照使用手机的品牌是什么？
3. 该图片被拍摄时相机的感光度（ISO）是多少？
4. 照片拍摄日期是哪一天？
5. 照片拍摄时是否使用了闪光灯？

题目很明显都与照片的 EXIF 信息有关。而 Mac 自带的查看功能并不能看到全部信息（比如 ISO），所以找了一个 [Python EXIF](https://pypi.org/project/exif/) 库来协助完成。

### 第二题：社工实践
#### 酒店
1. 请写出拍照人所在地点的邮政编码，格式为 3 至 10 位数字，不含空格或下划线等特殊符号（如 230026、94720）。
2. 照片窗户上反射出了拍照人的手机。那么这部手机的屏幕分辨率是多少呢？（格式为长 + 字母 x + 宽，如 1920x1080）


仔细看图，体育馆的一楼有一行字，但看不太清，像是 “zozomandie stadium”。于是万能的 Google 搜索纠正告诉我它是 “zozo marine stadium”，也就是千叶海洋球场。那么，根据拍照人的位置，可以推测出拍照人住在附近的酒店，而且与球场只隔一条街。  
根据 Google 地图提供的酒店地址（〒261-0021 千葉県千葉市美浜区ひび野２丁目３），我们就能够知道邮编就是“2610021”。

而拍照人的手机，同样是通过 EXIF 信息获得。搜索“xiaomi sm6115”，前面几个搜索结果似乎没有提供很多有价值的信息，但后面的结果告诉我这个手机应该是「Redmi 9T」，或者「Redmi Note 9」. 搜索“红米9T”，可以直接查到它的屏幕分辨率是 `2340x1080`.

#### 航班
仔细观察，可以发现照片空中（白色云上方中间位置）有一架飞机。你能调查出这架飞机的信息吗？
1. 起飞机场（IATA 机场编号，如 PEK）
2. 降落机场（IATA 机场编号，如 HFE）
3. 航班号（两个大写字母和若干个数字，如 CA1813）

EXIF 信息告诉我这张照片拍摄于 2022 年 5 月 14 日 18:23:35，一开始不太确定手机的时区是否是日本时区（UTC+9），所以通过 2022-05-14 千叶县的日落时间比对了一下，确定拍摄时间所在的时区就是 UTC+9. 那么我们要查的就是那个时刻经过的航班了。  
查这种东西我只能想到 FlightRadar24，但免费用户只能查询 7 天内的航班。于是被迫开了 7 天试用，查到了当时那家飞机的航班号（NH683），以及起降机场（HND → HIJ）。

## 猜数字
看看代码：
```java
var guess = Double.parseDouble(event.asCharacters().getData());

var isLess = guess < this.number - 1e-6 / 2;
var isMore = guess > this.number + 1e-6 / 2;

var isPassed = !isLess && !isMore;
```

浮点数，**不大于**也**不小于**，为什么不来个 `NaN` 呢？  
输入 `NaN`，拿到 flag. 因为前端输入框限制了只能输入数字和小数点，所以还要构造下请求才能顺利提交。

## 线路板
![](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/%E7%BA%BF%E8%B7%AF%E6%9D%BF/src/circuit_boards.png?raw=true)
题目给了一套[PCB 的生产文件](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/%E7%BA%BF%E8%B7%AF%E6%9D%BF/src/fabrication_files.zip)，需要找到 PCB 版中的 flag.

下载下来之后，看到里面有 `.gbr` 文件，虽然是纯文本格式，但它似乎并不能直接读出来。  
于是下载了一套 KiCad 工具，用 [PCBNew](https://docs.kicad.org/master/zh/pcbnew/pcbnew.html) 打开电路板，在啥都不懂的情况下一番鼓捣，最后看到了 Flag 的轮廓，是 `flag{8_1ayER_rogeRS_81ind_V1a}`.

![藏在 PCB 中的 flag](/pics/hackergame2022/PCB.png)

为了它我下了一套 15GB 的 App，用了 10 分钟就删掉了，略微有点离谱 :joy:

## 光与影
题目文件在[这里](https://github.com/USTC-Hackergame/hackergame2022-writeups/tree/master/official/%E5%85%89%E4%B8%8E%E5%BD%B1/src)，是一个使用 OpenGL 编写的动画。而我们就是需要找到被挡住的 flag.

![光与影](/pics/hackergame2022/light_and_shadow.png)

下下来以后看到了一堆完全不认识的代码，而经验告诉我核心的渲染逻辑都在 `fragment-shader.js` 中，而这个 js 里面使用了一种没见过但很像 C 的语言。 
翻到文件的最后，在主函数看到了一个 `isTerrain` 变量，顺着 `isTerrain` 为 `false` 的逻辑一路往上看，遇到不懂的地方就改改代码再运行一下看看效果，最后找到了关键代码。

```c
float sceneSDF(vec3 p, out vec3 pColor) {
    pColor = vec3(1.0, 1.0, 1.0);
    
    vec4 pH = mk_homo(p);
    vec4 pTO = mk_trans(35.0, -5.0, -20.0) * mk_scale(1.5, 1.5, 1.0) * pH;
    
    float t1 = t1SDF(pTO.xyz);
    float t2 = t2SDF((mk_trans(-45.0, 0.0, 0.0) * pTO).xyz);
    float t3 = t3SDF((mk_trans(-80.0, 0.0, 0.0) * pTO).xyz);
    float t4 = t4SDF((mk_trans(-106.0, 0.0, 0.0) * pTO).xyz);
    float t5 = t5SDF(p - vec3(36.0, 10.0, 15.0), vec3(30.0, 5.0, 5.0), 2.0);
    
    float tmin = min(min(min(min(t1, t2), t3), t4), t5);
    return tmin;
}
```

只要把 `t5` 从 `tmin` 的判断中去掉，就能看到 flag: `flag{SDF-i3-FuN!}`.

虽然改出来了，然而我依然没有看懂这堆代码。[官方题解](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/%E5%85%89%E4%B8%8E%E5%BD%B1/README.md)中详细地讲了一下渲染方法的实现，有兴趣的朋友可以去看看。

# Web
## 签到
需要手写识别 2022，但它把识别结果写在 URL 里了。  
把 URL `http://202.38.93.111:12022/?result=????` 里的 result 改成 2022，拿到 flag.

## XCaptcha
题目源代码在[这里](https://github.com/USTC-Hackergame/hackergame2022-writeups/tree/master/official/Xcaptcha/src)。  
出题人给了一个网页，需要你通过在一秒内算出三个大数加法并提交，来证明你是机器人。

写个自动提交的爬虫吧，网页结构不是很复杂，所以直接用正则即可，不需要判断 DOM.

```python
sess = requests.Session()
sess.cookies.set("session", "<some_jwt>")
response = sess.request("GET", url)

res = re.findall(r"(\d+)\+(\d+)", response.text)
result = []
for x, y in res:
    result.append(int(x)+int(y))
for k, v in response.cookies.items():
    sess.cookies.set(k, v)

response = sess.request("POST", url, data={
    "captcha1": result[0],
    "captcha2": result[1],
    "captcha3": result[2],
})
print(response.text)
```

## LaTeX 机器人
渲染 LaTeX 图片的脚本在[这里](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/LaTeX%20%E6%9C%BA%E5%99%A8%E4%BA%BA/src/latexbot-backend.zip)。

靠 Google 只会做一半：`\input{/flag1}` 拿到第一个 flag.  
由于构建脚本里关掉了 `shell escape`，所以其它基于命令执行的方案都会失败，后面卡了很久。

某天晚上，浩哥哥发来了[一个页面](https://tex.stackexchange.com/questions/630044/how-can-i-input-multiple-txt-files-with-special-characters-into-my-main-tex-file)，讲了如何去在文章里引用包含特殊符号的内容，而里面提到了 `\catcode` 命令。  
`\catcode` 的常见用法是将某个字符定义成某种宏。而我们通过它就可以将特殊字符视为一个 [Catcode 12](https://www-sop.inria.fr/marelle/tralics/doc-symbols.html#catcode12) 的字符，而这类字符是无法参与指令控制的。

最后输入 `{ \catcode``#=12  \catcode\``_=12  \input{/flag2} }` 得到 flag：`flag{latex_bec_0_m##es_co__#ol_2a1fd66cfe}`.

## Flag 的痕迹
题目写明了服务的版本：
> （题目 Dokuwiki 版本基于 2022-07-31a "Igor"）

这大概已经算是明示了？去搜一发 CVE，看到了 [CVE-2022-3123](https://nvd.nist.gov/vuln/detail/CVE-2022-3123#vulnCurrentDescriptionTitle)，通过构造 XSS 即可打开历史记录页面。

构造请求：
```
POST /doku.php?id=start HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Referer: http://202.38.93.111:15004/doku.php?id=start
Cookie: DokuWiki=57vk0n23v486p8vdjqc15oigpu; DOKU_PREFS=show_changes%23both%23difftype%23sidebyside
Content-Length: 139
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Encoding: gzip,deflate,br
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36
Host: 202.38.93.111:15004
Connection: Keep-alive

difftype=sidebyside'"()%26%25<zzz><ScRiPt%20>alert("1")</ScRiPt>&do=diff&do[diff]=1&id=start&rev2[0]=1&rev2[1]=0&sectok=1
```
看到历史记录，就可以拿到 flag.

---

当然以上都是误打误撞。

比赛结束后看到了[官方题解](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/Flag%20%E7%9A%84%E7%97%95%E8%BF%B9/README.md)，才发现真正的关键是请求参数中的 `do=diff`，其它的包括 CVE, XSS payload 等等都不重要。

因此只要访问 `http://202.38.93.111:15004/doku.php?id=start&do=diff`，然后在页面中查看版本变更即可。

## 微积分计算小练习
XSS 基础练习。

题目给了一个[练习页面](https://github.com/USTC-Hackergame/hackergame2022-writeups/tree/master/official/%E5%BE%AE%E7%A7%AF%E5%88%86%E8%AE%A1%E7%AE%97%E5%B0%8F%E7%BB%83%E4%B9%A0/src/web)，以及一个提交练习成绩的[后端代码](https://github.com/USTC-Hackergame/hackergame2022-writeups/tree/master/official/%E5%BE%AE%E7%A7%AF%E5%88%86%E8%AE%A1%E7%AE%97%E5%B0%8F%E7%BB%83%E4%B9%A0/src/bot)。

<img src="/pics/hackergame2022/calculus.png" style="height:30em">

练习页面长这样，先 X 一发，发现姓名字段是可以注入的，没有任何限制。  
然后再去看[后端代码](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/%E5%BE%AE%E7%A7%AF%E5%88%86%E8%AE%A1%E7%AE%97%E5%B0%8F%E7%BB%83%E4%B9%A0/src/bot/bot.py)，可以看到它使用 selenium 打开了练习结果页面，但打开结果之前，它把 flag 注入到了当前浏览器的 cookie 中。  
那么，直接在练习页面来一发 `<img src=# onerror="alert(document.cookie)">` 并将结果提交，就可以看到后端抛出了异常：

```
- Logining...
 Putting secret flag...
- Now browsing your quiz result...
ERROR <class 'selenium.common.exceptions.UnexpectedAlertPresentException'>
selenium.common.exceptions.UnexpectedAlertPresentException: Alert Text: flag=flag{xS5_1OI_is_N0t_SOHARD_abb4def144}
Message: unexpected alert open: {Alert text : flag=flag{xS5_1OI_is_N0t_SOHARD_abb4def144}}
  (Session info: headless chrome=106.0.5249.119)
Stacktrace:
...
```

后来看了[官方题解](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/%E5%BE%AE%E7%A7%AF%E5%88%86%E8%AE%A1%E7%AE%97%E5%B0%8F%E7%BB%83%E4%B9%A0/README.md)，他们给的 payload 是把 cookie 放到了名字的 DOM 中，这样看起来似乎更优雅一些。


# Math
我的数学是真的很垃圾，所以只能做出一道半…

## 蒙特卡罗轮盘赌
好家伙，头一次见到随机数种子碰撞的题。

题目描述见[这里](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/%E8%92%99%E7%89%B9%E5%8D%A1%E7%BD%97%E8%BD%AE%E7%9B%98%E8%B5%8C/README.md)。

看了下代码，主要逻辑是使用 `(unsigned)time(0) + clock()` 初始化随机数种子，并根据实验计算出 π 的值。而我们如果要答对题，就需要找到或猜到服务运行时所用的种子。  
那么我的思路主要靠猜：通过大量收集服务端给出的结果，并对其进行爆破算出每组数据实际使用的 `clock()`，统计出 `clock()` 的大致范围，然后用一个固定 clock 去随机碰撞答案。

题目的服务限制了连接频率，每个用户 10 秒钟只能建立一个连接。收集程序跑了一中午拿了两百多组数据，而我的 NAS 平均一分半才能算碰撞出一组结果。算了大概半个小时，发现 `clock()` 的值大多落在 700-900 之间。  
于是简单写了个脚本，用 800 作为 clock 的值开始蒙，没想到只用了一分半就蒙出来了！然而更让我没想到的，是我在撞完之后没有把服务端的所有输出都打出来就退出了，这导致我只拿到了结果，但没拿到 flag，直接哭死。

![](/pics/hackergame2022/lost_flag.jpg)

改了一下代码，使用了 `interactive()` 保证所有输出都被打印，重新开蒙。这次蒙了将近 70 分钟才蒙到。

---
最后看了[题解](https://github.com/USTC-Hackergame/hackergame2022-writeups/tree/master/official/%E8%92%99%E7%89%B9%E5%8D%A1%E7%BD%97%E8%BD%AE%E7%9B%98%E8%B5%8C)，发现完全不需要碰撞，只需要先拿两组算出种子，然后提交三个正确答案即可。  
果然还是绕了远路 :new_moon_with_face:

## 企鹅拼盘
嗯，这个题真的很可爱。

我的数学实在是太差，所以解法正如题目描述那样“大力出奇迹”。  
为了优化遍历速度，还用 Go 重写了代码，并对模拟算法做了常数级别的优化，就差加个 channel 搞并行了。想看爆破代码的可以戳[这里](/pics/hackergame2022/traverse.go)。

各位还是移步[官方题解](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/%E4%BC%81%E9%B9%85%E6%8B%BC%E7%9B%98/README.md)好啦。

# 后续
周六比赛结束后，官方发了题解，于是去看了下想做但没思路的[《安全的在线测评》](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/%E5%AE%89%E5%85%A8%E7%9A%84%E5%9C%A8%E7%BA%BF%E6%B5%8B%E8%AF%84/README.md)》和[《杯窗鹅影》](https://github.com/USTC-Hackergame/hackergame2022-writeups/blob/master/official/%E6%9D%AF%E7%AA%97%E9%B9%85%E5%BD%B1/README.md)，又涨了不少姿势。

太久没玩 CTF 了，周一玩了大半天，晚上靠着手速冲到了 58 名，虽然最后掉到了 134 名，不过还是蛮有成就感的。  
听浩子哥哥说，这几年的 CTF 的题基本都是赛棍特供，没想到 USTC 的大佬们能够设计出这么精彩的题目，在此感谢各位出题的同学。

最后放两张图留个纪念吧，希望明年这个时候我还能记得 hackergame 2023。

![](/pics/hackergame2022/rank58.png)
![](/pics/hackergame2022/final.png)
