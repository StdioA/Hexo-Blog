title: 在Windows中使用Wget
date: 2015-09-18 20:07:00
categories:
- 乱七八糟
tags:
- Windows
- wget

---

一直心心念念想在Windows里用wget，今天随便搞了搞，在Windows里面用上了wget，终于可以下载一些乱七八糟的小文件辣。

<!-- more -->

# 1. 引子
一直觉得Windows里面应该有一个类似wget的工具。今天想从github上下一个文件，结果文件直接以文本方式返回了(响应头的Content-Type为text/plain而不是application神马的)，自己很怨念，觉得要是Windows能像linux一样有个用来wget神马的下载文件的命令就好了…所以自己上网搜了一下，让自己用上了wget.

# 2. 乱搞
随便百度，找到一篇博文，又找到了一个叫做GnuWin(<http://sourceforge.net/projects/gnuwin32/>)的项目，该项目主要提供Windows下的GNU工具，比如sed, grep, wget等等。于是下载安装，安装程序默认将程序安装在了C:\Program Files (x86)\GnuWin32\bin目录下。然而即使这样，我还是无法直接在命令行里使用。

Windows中有一个“PATH环境变量”，将某目录（比如C:\Python27）添加进PATH环境变量中，就可以直接在命令行中输入命令打开该目录下的文件。如果想在命令行中使用wget的话，把C:\Program Files (x86)\GnuWin32\bin目录添加进PATH变量中当然可行。但是这样的话，随着工具越来越多，PATH变量里面的路径也就越来越多，维护起来也会更加困难，所以自己使用了一种特殊的方法：在C盘的目录下建了一个Command\_line\_programs文件夹用来放一些自己常用的命令行工具（比如sqlite），然后将C:\Command\_line\_programs（以下简称C.L.P.）添加到PATH变量中，这样自己就能使用C.L.P.文件夹中的工具。可是，自己已经将程序安装到Program Files里了，所以需要想个办法在C.L.P.文件中新建一个文件来连接到wget程序。_经过尝试，使用快捷方式的方法失败了。_所以自己想到了**新建批处理文件**。于是在C.L.P.文件夹中新建wget.bat, 内容为`C:\"Program Files (x86)"\GnuWin32\bin\wget.exe %*` （`%*`表示在命令中嵌入程序所有参数，详情请百度或Google Windows批处理编程），保存，在命令行中尝试使用wget，成功。

# 3. 然而这并没有什么卵用啊！
自己搞完以后，盯着GnuWin这几个字，突然想到了git里面自带了很多GNU工具可以直接使用，于是就去git/bin的目录翻了翻，果然，里面没有wget，然而我看到了curl。然后…就没有然后了。![](/pics/hands_up_crying.jpg)

# 4. 后记
这篇文章是胡乱凑数的，最近学习兴趣不高，又有一堆乱七八糟的事情要忙，所以托更了好久T^T
好了好了要好好学习，说好了多去体验框架呢QAQ
