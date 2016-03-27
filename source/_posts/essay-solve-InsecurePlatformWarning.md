title: 随手记 - 解决InsecurePlatformWarning
date: 2015-12-7 23:01:00
categories:
- 随手记
tags:
- python
- openssl
- requests

---

第三次在VPS上面解决使用requests报InsecurePlatformWarning警告的问题。之前每次都要查资料折腾好久，这次决定把它记下来。

<!-- more -->

# 1. 干货
* Debian类系统
    1. `sudo apt-get install python-dev libssl-dev`
    2. `sudo pip install -U requests[security]`

* Redhat类系统
    1. `sudo yum install python-devel openssl-devel`
    2. `sudo pip install -U requests[security]`

# 2. 需求
自己的VPS系统有点老(Ubuntu 14.04 LTS), 所以python版本也比较落后(Python 2.7.3), 今天改代码需要用到requests新版本中提供的功能，但是requests升级后发送HTTPS请求时会报出InsecurePlatformWarning, 这是一个由openssl漏洞(Heartbleed)造成的警告，所以需要升级pyopenssl等模块。

# 3. 升级过程
pypi提供了一个升级包，叫做`requests[security]`, 用pip进行升级即可。输入`sudo pip install requests[security]`命令后，pip报错，才发现不能本地编译python包，遂安装`python-dev`. 然后再次安装时发现缺少`openssl/aes.h`头文件，又去安装openssl的开发包`libssl-dev`, 再次安装，安装成功。

# 4. 参考资料
1. [[原]pip安装模块警告InsecurePlatformWarning](http://blog.csdn.net/blog/henulwj/48131393)
2. [zsh - no matches found: requests[security]](http://stackoverflow.com/questions/30539798/zsh-no-matches-found-requestssecurity)
