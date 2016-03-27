title: 博客迁移记（三）
date: 2016-1-20 11:40:00
categories:
- 乱七八糟
tags:
- nginx
- 七牛
- git

---

一个月没更新，期末复习的时候光鼓捣网站却懒得写东西，期末考完了，把这一个月鼓捣的东西记录一下，比如利用七牛进行静态文件托管。

<!-- more -->

# 1. 使用七牛托管静态文件
## 1.1 背景
博客的事情处理完成后，我又做了一个网站主页，其中包括一个100KB+的背景图片。不过因为我的VPS出站带宽只有1Mb/s, 所以背景图片加载时间过长，导致网站访问速度较慢。所以我将绝大部分的资源全部挂到了七牛上，在服务器端对静态资源请求进行302跳转，将流量转移到七牛的节点，提高访问速度。

## 1.2 操作
### 1.2.1 同步文件
我需要托管的文件包括博客和个人主页的所有图片文件。首先，我在七牛上建立了空间`cdn-stdioa`并绑定了个人域名、申请了HTTPS域名。为方便将本地文件与七牛空间同步，七牛提供了命令行同步工具`qrsync`. 查看[文档](http://developer.qiniu.com/docs/v6/tools/qrsync.html)后，我新建并修改了配置文件`sync_conf.json`，内容如下：
```json
{
    "src":          "blog/source/pics",
    "dest":         "qiniu:access_key=<your_access_key>&secret_key=<your_secret>key>&bucket=cdn-stdioa&key_prefix=blog/pics/",
    "debug_level":  1
}
```
若添加资源，则执行`qrsync sync_conf.json`, 即可完成静态资源与七牛空间的自动同步。

### 1.2.2 在服务器端设置跳转
为了提高访问速度，需要在nginx端将所有指向图片的请求全部跳转到七牛的链接上。学习了一下location重定向规则，直接上配置文件吧。
```nginx
location ^~ /pics {
    return https://xxx.qnssl.com/blog$request_uri;
}
```
这样，就可以将所有指向`/pics`的请求全部重定向到七牛的链接。顺便，在七牛设置了一下防盗链。

## 1.3 效率提升——Git hook
因为我的博客和个人主页都是使用Coding+webhook部署的，所以每次更改页面后要推代码，还要同步静态资源。有没有方法可以把这两条操作简化一下呢？一开始写了个批处理脚本，后来觉得一定有更好的办法，于是翻了翻，遇到了`Git hook`这种神器。
`Git hook`跟webhook类似，都是在某个操作上挂一个“钩子”，使得在进行某操作发生时自动触发自定义脚本来达到某些目的，实现快捷操作。所有的钩子均在`.git/hooks`目录下，在该目录下设置特殊文件名的脚本文件来设置钩子。因为我需要在推送代码之前进行资源同步，所以我需要设置`pre-push`脚本。脚本内容如下：
```bash
#!/bin/sh

echo -e "Sync up the static resources with Qiniu..."
cd /f/websites/blog/
qrsync sync_conf.json 2>>sync.log
echo -e "Done!"
```
这样，我就可以在每次执行`git push`之前自动同步静态资源，少敲一条命令，提升工作效率→_→

# 2. 谱站搭建
内容不是太多，就写在这里啦。
我有一个乐谱分享站，里面存放一些自己收集的钢琴谱。一年以前写了一个程序，用来为每个文件夹生成index页面。前几天翻开那个程序，一下子被自己写的连环replace吓到了（那个时候连正则还都不会，不过转念一想要是会了正则，写出来的东西会多可怕😂），于是怒用`Jinja`模板渲染引擎重写了一个，把原来最核心部分的十几行代码变为了一行渲染语句，代码看起来清爽多了。 
嗯，具体链接可以在左侧或顶栏的“友情链接”中找。

# 3. 后记
历经一个月，博客的迁移工作基本完成（刚刚写博客的工夫还给它添加了`twemoji`支持😶），以后还要做个人页面，不过可能不会有太多可写的地方了。就酱，后面看一阵子go以后可以来写写golang的东西。

# 4. 参考资料
[Nginx重定向规则详细介绍](http://www.seo-guwen.com/post-17.html)
[自定义 Git - Git 钩子](http://git-scm.com/book/zh/v2/%E8%87%AA%E5%AE%9A%E4%B9%89-Git-Git-%E9%92%A9%E5%AD%90)
