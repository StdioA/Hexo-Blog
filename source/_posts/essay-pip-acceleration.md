title: 随手记 - 用国内镜像加速pip
date: 2015-09-29 23:01:00
categories:
- 随手记
tags:
- python
- pip
toc: true

---

原来Windows和Linux更改镜像源的方式是不一样的啊。

<!-- more -->

# 1. 引子
偶然发现USTC有一个pypi的源([在这里](http://mirrors.ustc.edu.cn/pypi/))，照着USTC给的[镜像使用帮助](https://lug.ustc.edu.cn/wiki/mirrors/help/pypi)更改镜像源无果。今天闲来无事就多搜了一下。

# 2. 过程
USTC的镜像使用帮助里说，将`index-url = https://pypi.mirrors.ustc.edu.cn/simple`添加到`~/.pip/pip.conf`文件中，按照此思路看，如果我在Windows下使用的话，应该将配置信息添加到`C:\Users\Stdio\.pip\pip.conf`文件中。然而添加后并没有什么卵用。今天看见一个[博文](http://blog.csdn.net/sasoritattoo/article/details/10020547)说文件设置路径应为`%HOME%\pip\pip.ini`，于是试了一下（绝对路径为`C:\Users\Stdio\pip\pip.ini`），成功。用了USTC的源，装个软件速度简直飞起︿(￣︶￣)︿  
弄好以后多了个心眼，去看看官方怎么说，于是找到了pip的Documentation（[在这里](https://pip.pypa.io/en/stable/user_guide/#configuration)）里面讲述了配置文件所在位置。

```markdown
Config file
pip allows you to set all command line option defaults in a standard ini style config file.

The names and locations of the configuration files vary slightly across platforms. You may have per-user, per-virtualenv or site-wide (shared amongst all users) configuration:

Per-user:

On Unix the default configuration file is: $HOME/.config/pip/pip.conf which respects the XDG_CONFIG_HOME environment variable.
On Mac OS X the configuration file is $HOME/Library/Application Support/pip/pip.conf.
On Windows the configuration file is %APPDATA%\pip\pip.ini.
There are also a legacy per-user configuration file which is also respected, these are located at:

On Unix and Mac OS X the configuration file is: $HOME/.pip/pip.conf
On Windows the configuration file is: %HOME%\pip\pip.ini
You can set a custom path location for this config file using the environment variable PIP_CONFIG_FILE.

Inside a virtualenv:

On Unix and Mac OS X the file is $VIRTUAL_ENV/pip.conf
On Windows the file is: %VIRTUAL_ENV%\pip.ini
```
于是自己将配置文件放在了`%APPDATA%\pip\pip.ini`，即`C:\Users\Stdio\AppData\Roaming\pip\pip.ini`下，经实验，加速成功。
经测试，系统层面的全局设置会影响到virtualenv建立的虚拟环境设置，所以可以通过设置虚拟环境的配置文件来更改虚拟环境设置，设置文件就放在虚拟环境根目录下就好了。

# 3. 瞎写
改了软件源以后为了做测试升级了numpy，然后编译用了半天中间电脑卡到死，现在它还在我的腿上发烫真是伤不起- -  
写的东西越来越没有营养了，就当是写着玩顺带积累一下知识吧。如果不出意外的话，下一篇应该是关于Python Qt Creater和Qt GUI设计的小文章0.0也许是virtualenv的随手记吧。
