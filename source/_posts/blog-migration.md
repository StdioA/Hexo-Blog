title: 博客迁移记（一）
date: 2015-12-17 19:55:37
categories:
- 乱七八糟
tags:
- python
- Github
- nginx
toc: true

---

昨天买了一个VPS、一个域名，决定把博客迁移到VPS上。

<!-- more -->

# 1. 背景介绍
首先打个广告：借助“腾讯云+校园”计划，我成功地用上了1元/月的服务器和11元/月的.com域名。
有了国内访问速度极快的VPS，有了域名，就想自己鼓捣鼓捣，把自己的博客和其他站点托管到VPS上。

# 2. 博客托管
## 2.1 搭建静态文件服务器
因为自己的博客是静态博客，之前托管在[Github](https://github.com/StdioA/stdioa.github.io)上，依靠Github Pages实现部署，所以我先将文件clone了下来，存放在`~/blog`目录中。
下面要做的，是搭建一个静态服务器。

### 2.1.1 使用Python的SimpleHTTPServer
Python的SimpleHTTPServer是一个非常实用、方便的库，可以使用简单一条命令在当前目录创建一个HTTP文件服务器。所以输入`sudo python -m SimpleHTTPServer 80`命令，即可搭建一个静态文件服务器，实现从外网对静态博客的直接访问。
然而搭建好后，我发现了一个问题：因为我的blog文件夹本身是一个Git Repo, 所以我可以直接从外网访问.git文件夹，虽然我的.git目录没有保存任何设置及账户等，但这样会带来一定的安全隐患，所以要想办法禁止外部用户对.git文件夹的访问。

### 2.1.2 使用Nginx托管文件服务器
关于Nginx的介绍，请自行访问[官网](http://nginx.org/)与[维基百科](https://en.wikipedia.org/wiki/Nginx)。
之前用过Apache, 但是在接触Nginx后，我认为我对Nginx更有好感，所以采用了Nginx做为文件托管服务器。
首先安装Nginx，删掉`/etc/nginx/sites-enabled`目录（我的系统是Ubuntu Server 14.04 LTS），在`conf.d`目录中设置配置文件：
```
server {
    listen 80;                              # 监听端口
    server_name server;
    root /home/ubuntu/blog;                 # 托管目录
    index index.html;
    access_log /var/log/nginx/access_blog.log
    error_log /var/log/nginx/error_blog.log

    # 禁止对.git目录的访问
    location ^~ /.git {
        deny all;
    }
}
```
配置完成后，重启Nginx, 所有服务正常运行，访问.git目录时会返回403.

_SimpleHTTPServer在建立服务器的时候只能监听`0.0.0.0`而不能只监听`127.0.0.1`, 之前搭建服务器的时候用的SimpleHTTPServer+Nginx反向代理，现在看起来感觉我就是个傻逼…_

## 2.2 借助Webhook实现博客的自动部署
因为之前静态博客托管在Github Pages上，所以向Github Repo上面进行push操作之后会动态更新页面。但是如果将静态博客托管在VPS上，则需要每次执行`git pull`才能够将内容更新。所以我在VPS上写了一个脚本，能够在Push之后自动进行`git pull`来更新内容。

1. 在Github上的repo设置Webhook
在repo的设置页面可以设置Webhook, 可在该repo收到push之后，Github可以向一个特定的URL发送一个POST请求。所以我设置了一个Webhook, 在push之后可以向我的VPS的特定端口发送POST请求。

2. 设置服务器接收Webhook
在VPS中配置一个服务器，开放VPS的一个端口来接收来自Github的Webhook，这里使用bottle来搭建服务器框架。服务器代码:
```python
#!/usr/bin/env python
# coding: utf-8

import os
from bottle import *

@route("/push", method=["POST"])    # 监听Webhook
def pull():
    os.system("./auto_pull.sh")     # 执行git pull脚本
    return "OK"

run(host="0.0.0.0", port=23333)

```

3. 编写自动pull脚本
编写自动pull脚本`auto_pull.sh`:
```shell
#!/usr/bin/env sh

cd ~/blog
git pull
```
运行服务器，则可监听23333端口的POST请求，然后自动执行`git pull`更新博客内容。

# 3. 后记
至此，博客**部分**成功迁移到VPS.
为什么要说“部分”？因为我的域名在做备案啊QvQ备案好麻烦还要打印材料还要跑到市中心去照相还要发快递到北京QvQ都做好还得等待管局审核QvQ
吐槽时间结束。后面等域名备案好以后可能会鼓捣一阵子Nginx，为服务器启用HTTPS等…
**唉，还是先去复习吧\_(:зゝ∠)\_**

_在写博文时，去编译升级了一下Nginx, 鼓捣配置文件又弄了半天_╮(╯\_╰)╭
