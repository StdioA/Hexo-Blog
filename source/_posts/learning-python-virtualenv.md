title: Python学习之virtualenv
date: 2015-10-22 10:15:00
categories:
- Python
tags:
- python
- virtualenv
toc: true

---

用过virtualenv的人都说好，可是我没有具体使用过，所以尝试了一下，用完我也说好~233333

<!-- more -->

# 1. 简介
virtualenv是一个python库，用于创建独立python开发及运行环境。一般linux环境下如果在全局用pip安装模块时需要使用sudo命令，可是在共享主机上将root权限交给一般用户是不显示而且不安全的。可是有了virtualenv, 普通用户就可以创建一个虚拟环境，然后在虚拟环境中以普通用户权限安装模块，更改环境变量，进行开发和运行python程序而不会影响系统环境的环境变量和Python模块。

# 2. 创建virtualenv环境
输入`virtualenv venv`创建名为venv的虚拟环境。
创建虚拟环境的常用选项：

> --no-site-packages 不使用系统中的site packages
> --system-site-package 使用系统中的site packages *(据说是默认，但在我这默认是不使用的)*
> -p PYTHON_EXE, --python=PYTHON_EXE 使用指定的python解释器，这里的PYTHON_EXE在Windows下需要用绝对路径，比如C:\Python27\python.exe\

当然，也可以使用虚拟环境和系统配置文件来设置virtualenv默认创建选项，详情见[这里](http://virtualenv-chinese-docs.readthedocs.org/en/latest/#id7)。

输出：
```shell
G:\>virtualenv venv
New python executable in venv\Scripts\python.exe
Installing setuptools, pip, wheel...done.
```

# 3. 进入及退出虚拟环境
进入venv目录，Linux下输入`bin/activate`, Windows下输入`Scripts\activate`进入虚拟环境。
```shell
G:\venv>Scripts\activate
(venv) G:\venv>
```
进入环境后，输入`deactivate`退出。
```shell
(venv) G:\venv>deactivate
G:\venv>
```

# 4. 使用虚拟环境
## 4.1 安装第三方模块
安装过程与平时相符（比如使用`pip install`), 只不过安装后的包会存储在虚拟环境中。

## 4.2 设置环境变量
Linux下输入`export VAR1="value1"`, Windows下输入`set VAR1=value1`来设置虚拟环境的环境变量
```shell
(venv) C:\Users\Stdio\Desktop\temp\venv>set VAR1=value1

(venv) C:\Users\Stdio\Desktop\temp\venv>echo %VAR1%
value1
```
环境变量设置成功后，即可在Python程序中利用虚拟环境的环境变量进行程序配置（比如flask中的`app.config.from_envvar("FLASK_SETTINGS")`）。该环境变量只在虚拟环境中有效，退出虚拟环境后环境变量即消失。

# 5. 参考文档
1. [Virtualenv - virtualenv 13.1.2 documentation](http://virtualenv.readthedocs.org/en/latest/index.html)
2. [Virtualenv - virtualenv 1.7.1.2.post1 documentation 中文版](http://virtualenv.readthedocs.org/en/latest/index.html)
3. [virtualenv入门教程](http://flask123.sinaapp.com/article/39/)
4. [virtualenv -- python虚拟沙盒](http://www.cnblogs.com/tk091/p/3700013.html)
5. [用virtualenv建立多个Python独立开发环境](http://www.nowamagic.net/academy/detail/1330228)

# 6. 后记
又填完一个坑，本来打算昨天就写好的，结果昨天没写完…
