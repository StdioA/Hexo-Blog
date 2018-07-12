title: 博客迁移记（二）
date: 2015-12-18 15:55:37
categories:
- 乱七八糟
tags:
- nginx
- Github
- Let's Encrypt
toc: true

---

域名依然在备案中，依然想将博客部署到新买的VPS的我和腾讯斗智斗勇（雾），为博客添加了HTTPS>_< 
同时，把自己的密码学课设也转移到VPS上来了。

<!-- more -->

# 1. 使用Nginx针对多个host部署服务器
## 1.1 概述
如果一个VPS只能搭一个网站，那未免太浪费了，所以我们可以通过配置Nginx的方式来将针对多个域名的访问请求分开，从而进行不同的处理。例如，我在DNS配置时将`blog.stdioa.com`与`crypt.stdioa.com`同时指向到我的腾讯云的IP地址，用户访问这两个域名时，都会向我的VPS发送请求，我要做的是将这两种针对不同域名的访问请求分开。而这两种请求的区别在HTTP请求头的Host字段，所以我只需要针对不同的Host使用不同的处理方式即可。

## 1.2 操作
我的VPS上撘有两个网站，其中`blog.stdioa.com`域名指向的是我的博客——一个静态服务器，而`crypt`域名指向的是一个使用flask搭建的网站，所以要在Nginx端进行反向代理，将请求转发到本地的5000端口。
具体配置文件如下：
```nginx
server {
     listen       80;
     server_name  stdioa.com blog.stdioa.com;       # 通往博客的请求直接通过文件服务器返回

     access_log  /var/log/nginx/access_blog.log;
     error_log /var/log/nginx/error_blog.log;

     root   /home/stdio/blog;
     index  index index.html;

     location / {
     }
     location ^~ /.git {                            # 禁止访问.git文件夹
         deny all;
     }
     error_page  404     /404.html;
}
server {
     listen 80;
     server_name crypt.stdioa.com;                  # host为crypt的请求转发本地端口

     access_log /var/log/nginx/access_crypt.log;
     error_log /var/log/nginx/error_crypt.log;

     location / {
         proxy_pass http://127.0.0.1:5000/;
     }
}
```

配置完成后，重启Nginx, 可以看到访问不同域名时，请求会交给不同的程序处理。

# 2. 屏蔽来自特定域名的请求
## 2.1 背景
网站部署好后，又发现了一个来自奇怪域名的请求；更坑爹的是，这个域名指向自己的VPS；更更坑爹的是，来自这个域名的请求有好多，直接把我的日志刷爆了…所以我需要将来自这个域名的所有请求拒绝掉。

## 2.2 操作
新建一个Server就好啦。具体配置文件如下：
```nginx
server {
    listen 80;
    server_name bailigo.com *.bailigo.com;
    location / {
        return 410;         # 410 Gone, 使用了这个状态码，不知道能不能不再让爬虫爬这个网页
    }
}
```

配置完成，重新载入配置文件，成功。
_本来想用418, 但是Nginx没有418的返回页面，想了想还是算了吧_

# 3. 为博客部署HTTPS服务器
## 3.1 背景
`blog.stdioa.com`博客上午还可以访问，中午吃顿饭就发现访问被截断了，原因与以前一样——域名未完成备案。经查看，腾讯对访问博客的请求进行了301跳转，于是想了想，给博客配置了HTTPS, 让你们再阻断→_→（好吧再阻断的话我真的不知道该怎么弄了

## 3.2 使用Let's Encrypt生成网站证书
之前就看中了[Let's Encrypt](https://letsencrypt.org/)，它提供主流浏览器认证的免费证书，只可惜当时没有域名无法体验。现在有了域名，加上腾讯对未备案的域名查的很紧，所有未备案域名下的网站搭起来半天就被封掉，想了想，还是生成一个证书，配个HTTPS撑一阵子吧╮(╯\_╰)╭
按照官方的[指南](https://community.letsencrypt.org/t/quick-start-guide/1631)以及[此指南](http://www.vpser.net/build/letsencrypt-free-ssl.html)在将[Let's Encrypt的Repo](https://github.com/letsencrypt/letsencrypt) clone到本地之后，输入`./letsencrypt-auto certonly --email 邮箱 -d 域名 --agree-tos`来生成证书。由于腾讯云到Let's Encrypt的服务器的链接极其不稳定，通常需要重试很多次才能正常跟Let's Encrypt的服务器通信。成功后会显示：
> IMPORTANT NOTES:
> \- Congratulations! Your certificate and chain have been saved at
>   /etc/letsencrypt/live/blog.stdioa.com/fullchain.pem. Your cert
>   will expire on 2016-03-17. To obtain a new version of the
>   certificate in the future, simply run Let's Encrypt again.
> \- If you like Let's Encrypt, please consider supporting our work by:
> 
>   Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
>   Donating to EFF:                    https://eff.org/donate-le

证书已存储在上面信息指示的目录。

## 3.3 配置Nginx, 搭建HTTPS服务器
证书已生成，下面该配置Nginx了，添加下列配置：
```nginx
server {
    listen  443 ssl;
    server_name  blog.stdioa.com;

    ssl_certificate /etc/letsencrypt/live/blog.stdioa.com/fullchain.pem;        # 添加证书
    ssl_certificate_key /etc/letsencrypt/live/blog.stdioa.com/privkey.pem;      # 添加密钥
    ssl_ciphers "EECDH+AESGCM:EDH+AESGCM:AES256+EECDH:AES256+EDH";
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;

    # 后面和普通服务器一样
    access_log  /var/log/nginx/access_blog.log;
    error_log /var/log/nginx/error_blog.log;

    root   /home/stdio/blog;
    index  index index.html;

    location ^~ /.git {
        deny all;
    }

    error_page  404 /404.html;
}
```

重新加载配置，访问<https://blog.stdioa.com> , 成功~

小插曲：
昨天晚上熄灯之前发现Ubuntu上使用`apt-get`安装的Nginx版本太老了，于是在VPS上重新编译升级了Nginx，结果配置好HTTPS以后重启Nginx时发现启动失败，原因为缺少`ngx`组件。经Google后发现没有编译该组件，所以重新编译安装Nginx, 在配置时加入`--with-http_ssl_module`选项。

小插曲2：
HTTPS配置好后访问博客，Chrome提示“存在不安全的内容”，选择加载后，地址栏左边的HTTPS会变成红色，极其不好看（雾），于是看了博客的模板，发现有两个js在加载时选择使用HTTP方式加载，于是将`http://.../*.js`改为`//.../*.js`, 使浏览器可以根据当前协议自动选择JS文件的加载协议。更新博客模板，重新访问，所有js均使用HTTPS方式加载，问题解决。附图:
<center>![哈哈哈HTTPS](/pics/hhh_HTTPS.jpg)</center>

小插曲3：
上面那张图片的链接之前来自七牛，采用HTTP协议而不是HTTPS加载。发布这篇博文后，我发现该文章页面中地址栏左侧的HTTPS标志变为了白色…还是不好看！打开Chrome的控制台，在控制台中看到，如果图片链接协议为HTTPS，则Chrome依然会提示“不安全”，但是个人感觉这只是一张图片而已啊，又不是JS\_(:зゝ∠)\_
解决方案：将图片链接改回本地，待域名备案后再想办法在七牛那边解决域名绑定及HTTPS的问题。

# 4. 参考资料
[免费SSL安全证书Let's Encrypt安装使用教程(附Nginx/Apache配置)](http://www.vpser.net/build/letsencrypt-free-ssl.html)
http://blog.csdn.net/donghustone/article/details/25797727
[nginx SSL error - Server Fault](http://serverfault.com/questions/438889/nginx-ssl-error)
[网站存在不安全因素的解决办法](http://www.zzidc.com/main/help/showHelpContent/id_473.html)

# 5. 后记
快递已发出，相片已审核成功，腾讯收到资料之后会尽快报管局审理…希望2016年之前能够备案成功QvQ
嗯，Google真好用
**说好的准备考试呢！**

---

**Update @ 16:38**
上午才把资料用快递发走，下午腾讯云就说我的纸质资料已收到…打电话问了一下，客服说腾讯为了加快审核速度，帮我准备了一份材料上交管局，估计他们是用我的扫描件打印了一份交上去了吧，也是不错2333 最后希望通信交通管理局快一点\_(:зゝ∠)\_

---

**Update @ 19:59, Dec. 24th, 2015**
腾讯云访问Github的速度太慢了，于是将博客的Repo迁移到了Coding上。

1. 建立私有项目
2. 设置部署公钥
    因为是私有项目所以无法通过一个__不含用户名密码__的链接访问Repo来进行Pull，所以还是通过SSH访问比较舒服一些。
3. 改掉VPS上的origin链接，完成

曾经听到过一种说法：用HTTPS访问Git Repo要比SSH更好，然而一直不知道好在哪…记得用HTTPS访问Repo的时候账户的用户名和密码是要附在链接里面的，感觉好不安全
域名备案看来真的要奔着20天的样子去了…

**Update @ 11:07, Jan. 5th, 2015**
域名备案在2016年的第一个工作日通过啦！庆祝一下，期末之后开始准备制作个人网站~
