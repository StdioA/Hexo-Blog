title: 使用 supervisor 及 gunicorn 部署 Web 应用
date: 2017-03-19 17:00
categories:
- Web
tags:
- Python
- supervisor
- gunicorn
- Django
toc: true

---

很久之前就想尝试一下用 supervisor 部署 Web 应用，几个月前把 Python 应用的服务器都换成了 gunicorn，今天终于把进程管理服务换成了 supervisor. 看我的拖延症。

<!-- more -->

# 1. 前言
之前一直在用 tmux 来托管各种 Web 应用进程，感觉这种想法真的很蠢，于是今天把托管方式换成了更专业的 supervisor，并用它托管了三个 Django APP，一个 flask APP，还有一个 node APP.

简单看看今天要用的东西：
* [gunicorn](http://gunicorn.org)，一个 Python 实现的 WSGI 服务器；
* [supervisor](http://supervisord.org/)，一个进程管理工具。
* Django, Flask, Node.js，不多说。

# 2. supervisor 安装及基础配置
`pip2 install supervisor` 即可。注意，supervisor 不支持 Python 3.

supervisor 提供了一个配置生成程序 `echo_supervisord_conf`，可以用它来直接生成一个实例配置文件。直接输入 `echo_supervisord_conf > /etc/supervisor/supervisord.conf`，将配置写入文件。

随后，修改配置文件，开启 http 管理服务：
```ini
[inet_http_server]            ; inet (TCP) server disabled by default
port = 127.0.0.1:9001         ; (ip_address:port specifier, *:port for all iface)
username = stdio              ; (default is no username (open server))
password = password           ; (default is no password (open server))
```

添加选项，包含子目录 `conf.d` 下的所有配置文件：
```ini
[include]
files = conf.d/*.conf
```

随后，使用 `sudo supervisord -c /etc/supervisor/supervisord.conf`，启动 supervisor 守护进程。

配置一下 Nginx，登录管理页面，可以看到 supervisor 正在运行，不过现在还没有配置任何服务。  
同理，可以使用 `supervisorctl` 程序，来查看服务运行状态。

# 3. 服务托管
## 3.1 托管一个 Django 应用
进入 `conf.d` 文件夹，创建配置文件（如 `baybook.conf`）。

```ini
[program:baybook]
command = gunicorn baybook.wsgi -b 127.0.0.1:8002 -n baybook      ; 运行命令
directory = /home/stdio/websites/baybook                          ; 运行路径
user = stdio
autostart = true
autorestart = true
startsecs = 5
startretries = 3
stdout_logfile = /var/log/supervisor/baybook_stdout.log
stderr_logfile = /var/log/supervisor/baybook_stderr.log
environment=DJANGO_SETTINGS_MODULE="baybook.settings.production"
```

其中，`environment` 选项可以配置运行环境的环境变量，比如在此处更改了 Django 的配置文件选项；`startsecs` 选项表示正常启动所需的时间，比如，当程序已持续运行超过 5 秒时，则视为程序启动成功。

具体的配置选项可以查看[文档](http://supervisord.org/configuration.html)。

保存文件后，运行 `supervisorctl update` 使配置生效。  
随后可以输入 `supervisorctl status` 来查看配置状态。

```shell
$ sudo supervisorctl status
baybook                          RUNNING   pid 7658, uptime 1:19:55
```

supervisor 提供了一个命令行界面，直接输入 `supervisorctl` 即可进入，随后可以输入一系列命令，如 `start`, `stop`, `status`, `restart` 来查看和控制服务运行状态。

当然，也可以在 Web 管理页面中，查看及控制服务的运行状态。

![supervisor 管理页面](/pics/supervisor-status.png)

## 3.2 托管 Node 应用
Node 和 Django 的配置文件基本相同，不过，因为我的 node 应用 的启动速度比较慢，所以我把 `startsecs` 调高到了 30 秒。

## 3.3 托管 Flask 应用
因为我的 Flask 应用运行在 `virtualenv` 创建的虚拟环境中，所以 `command` 命令要稍微改一下，将 `gunicorn` 可执行文件的路径改为虚拟环境中的 `gunicorn` 绝对路径，如 `/home/stdio/websites/crypt/venv/bin/gunicorn`.

# 4. 后记
好像可写的就这么多…  
以后研究一下 gunicorn，再补充 gunicorn 的相关内容吧。
