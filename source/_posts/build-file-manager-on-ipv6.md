title: 瞎玩IPv6——在公网搭建文件管理器
author: David Dai
tags:
  - IPv6
  - NAS
  - Linux
categories:
  - 网络
date: 2019-02-21 19:28:00
toc: true
---
IPv6 是个好东西，希望人人都有一个。

<!--more-->

# IPV6

## IPv6 是啥？
新一代的 IP 协议，解决了 IPv4 地址枯竭的问题。具体可以见 [Wikipedia](https://zh.wikipedia.org/wiki/IPv6).  
IPv6 的 IP 长度为 128 位，总量非常非常多，不用担心用不完，所以接入 IPv6 的客户端都会分到一**段** IP，比如 `240e:1c:ce8:fd00::/64`，然后客户端又可以把这段 IP 继续分段，下发到下面的所有子网中。  
不过需要注意的是，虽然客户端会分到一段 IP 的所有权，不过客户端本身还是会有至少一个确定的 IP，以确定自己的位置。  
_其实发现 ISP 分给自己 2^64 个 IP 地址的时候还是感觉很奢侈…_:joy:

## IPv6 跟我有什么关系？
感谢去年工信部发布了[《工业和信息化部关于贯彻落实〈推进互联网协议第六版（IPv6）规模部署行动计划〉的通知》](http://www.miit.gov.cn/n1146295/n1652858/n1652930/n3757020/c6154756/content.html)，IPv6 现在也飞入寻常百姓家了~

ISP 下发的 IP，都是公网 IP ，我们再也不用躲在层层 NAT 后面，想要一个公网 IP 都要去打电话跟运营商扯皮了。  
年前意外看到自己的路由器界面上有了公网的 IPv6 IP，电信 4G 网络也支持 IPv6 了，于是打算折腾一下，把这个公网 IP 利用起来，比如给家里的 NAS 搭个公网可以访问的云盘什么的:flushed:

# 配置
## DHCP
如果要用公网 IP，首先你要有一个公网 IP. 如果 `ifconfig` 看到的 IP 只有 `fe80` 开头的地址，那是用于链路本地通信的保留地址，是不能在公网使用的。

打开路由器的后台，看一下 WAN6 的接口，如果能看到公网 IP，那就可以配置 `dhcp` 下发了。

编辑 `/etc/config/dhcp` 打开 DHCPv6 的中继模式：
```
config dhcp 'lan'
        option interface 'lan'
        option start '100'
        option limit '150'
        option leasetime '12h'
        option ndp 'relay'
        option dhcpv6 'relay'
        option ra 'relay'

config dhcp 'wan'
        option interface 'wan'
        option start '100'
        option limit '150'
        option ignore '1'
        option leasetime '12h'
        option ra 'relay'
        option dhcpv6 'relay'
        option ndp 'relay'

config dhcp 'wan6'
        option interface 'wan'
        option ra 'relay'
        option dhcpv6 'relay'
        option ndp 'relay'
        option master '1'
```

配置好后，用 `/etc/init.d/odhcpd restart` 重启 DHCP 服务器；等待一会，应该就可以在客户端网口上看到一个公网的 ipv6 地址了。

```bash
$ ifconfig enp1s0
enp1s0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.168.10.209  netmask 255.255.255.0  broadcast 192.168.10.255
        inet6 240e:17:ce8:fd00:7285:c2ff:fe38:2a56  prefixlen 64  scopeid 0x0<global>
        inet6 fe80::7285:c2ff:fe38:2a56  prefixlen 64  scopeid 0x20<link>
```

参考资料：[OpenWRT IPv6 三种配置方式](http://blog.kompaz.win/2017/02/22/OpenWRT%20IPv6%20%E9%85%8D%E7%BD%AE/)

## 防火墙
打开以后，可以随便找一个外网的机器来 ping 一下自己的机器地址，如果顺利的话，应该可以 ping 通了，但是 SSH/HTTP 访问应该还是不可以的，因为路由器的防火墙拒绝了来自外网的 TCP 请求。

在路由器的防火墙界面配置流量规则，允许目标地址为服务器的 IP 地址，端口为 22/80/443 端口的请求：

![防火墙规则配置](/pics/firewall-config.png)

或者直接使用 `ip6tables` 命令配置路由器的 iptables：
```bash
ip6tables -A zone_wan_forward -d 240e:17:ce8:fd00:7285:c2ff:fe38:2a56/128 -p tcp -m tcp --dport 22 -m comment --comment "!fw3: Allow-v6-forward" -j zone_lan_dest_ACCEPT
ip6tables -A zone_wan_forward -d 240e:17:ce8:fd00:7285:c2ff:fe38:2a56/128 -p tcp -m tcp --dport 80 -m comment --comment "!fw3: Allow-v6-forward" -j zone_lan_dest_ACCEPT
ip6tables -A zone_wan_forward -d 240e:17:ce8:fd00:7285:c2ff:fe38:2a56/128 -p tcp -m tcp --dport 443 -m comment --comment "!fw3: Allow-v6-forward" -j zone_lan_dest_ACCEPT
```

然后本地开个 HTTP 服务器，再尝试访问一下 `http://[240e:17:ce8:fd00:7285:c2ff:fe38:2a56]/`，如果顺利的话，就可以自己做网站了~

测试的时候遇到了一个小插曲：我在 NAS 的 6000 端口上随手部署了一个 HTTP 服务器，但用手机的 Chrome 访问的时候报出了 `ERR_UNSAFE_PORT` 的错误，才发现 Chrome 还有这种奇怪的端口访问限制…具体可见[这篇文章](https://blog.csdn.net/testcs_dn/article/details/39186225)。

## 文件服务器
文件服务器的选型，我使用了 Nginx + FileRun，用 Docker 进行托管。  
Nginx 是之前在内网搭服务的时候就搭好的，[FileRun](https://www.filerun.com/) 是一款长得和 Google Drive 有点像~~（哪里像了）~~的文件服务器，支持文件分享、图片/视频在线预览，还可以使用 Office Web/Google Docs 在线编辑文件。  
它是收费软件，不过免费功能也够用。功能强大，长得又很漂亮，就用它了~

FileRun 支持使用 Docker 部署，方法可见[官方文档](https://docs.filerun.com/docker)。  
其中，容器中的 `/var/www/html` 目录会存放 FileRun 首次启动时生成的 PHP 文件，而 `/user-files` 目录是 FileRun 程序进行文件操作的根目录。

在部署时，我的 NAS 上已经有一个 MySQL 服务了，不想再另起一个，所以通过 `external_link` 设置让 FileRun 容器直接访问现有的 MySQL 容器：
```yaml
version: '2'
services:
  web:
    image: afian/filerun
    environment:
      FR_DB_HOST: mysql
      FR_DB_PORT: 3306
      FR_DB_NAME: filerun
      FR_DB_USER: filerun
      FR_DB_PASS: password
      APACHE_RUN_USER: www-data
      APACHE_RUN_USER_ID: 33
      APACHE_RUN_GROUP: www-data
      APACHE_RUN_GROUP_ID: 33
    network_mode: bridge
    external_links:
        - mysql
    ports:
      - "8030:80"
    volumes:
      - /data/filerun/html:/var/www/html
      - /mnt:/user-files
    restart: always
```

Nginx 配个反代，DNS 配一条 AAAA 记录，[acme.sh](https://github.com/Neilpang/acme.sh) 签个 HTTPS 证书，这些都不细讲了。

登录时，发现 JS 里有个奇怪的 `base_url` 指向了 `localhost:8030`，导致页面脚本不能正常工作。  
查到[这个帖子](https://www.reddit.com/r/FileRun/comments/68vucn/running_filerun_behind_nginx_reverse_proxy_xpost/)，新建文件 `/user-files/customizables/config.php`，写入如下内容：
```php
<?php $config['url']['root'] = 'https://filerun.somedomain.com'; $_SERVER['HTTPS'] = 'on';?>
```

重启之后就可以正常使用了。

来张图感受一下:smile:

![filerun](/pics/filerun.png)

## 还剩下一个小问题：DDNS
路由器重启或重新拨号后，NAS 会拿到新的 IP，网站就无法正常访问了。  
于是写了个脚本用 CloudFlare API 更改 DNS 记录，可以定时运行脚本，来更新 IP 地址。

```python
#!/usr/bin/env python3
import socket
import requests

KEY = "<Cloudflare API Key>"
EMAIL = "<Cloudflare Email>"
ZONE_ID = "<Zone ID>"

records = [["<DNS Record ID>", "<Domain>"]]

def get_ip_6(host, port=0):
     sock = socket.socket(socket.AF_INET6, socket.SOCK_DGRAM)
     sock.connect((host, port))
     return sock.getsockname()[0]

def main():
    ip = get_ip_6('ipv6.google.com')

    api_url = "https://api.cloudflare.com/client/v4/zones/{}/dns_records/{}"
    headers = {
        "X-Auth-Email": EMAIL,
        "X-Auth-Key": KEY
    }
    for record_id, domain in records:
        url = api_url.format(ZONE_ID, record_id)
        res = requests.put(url, json={
            "type": "AAAA",
            "name": domain,
            "content": ip,
            "ttl": 120,
            "proxied": False
        }, headers=headers)
        print(res.status_code, res.json())


if __name__ == '__main__':
    main()
```

Ref: [Finding local IP addresses using Python's stdlib](https://stackoverflow.com/questions/166506/finding-local-ip-addresses-using-pythons-stdlib)

差不多就写这些了。
