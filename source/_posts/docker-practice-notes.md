title: 《Docker 实践》阅读笔记
author: David Dai
tags:
  - Docker
  - 读书
categories:
  - DevOps
date: 2018-10-30 21:14:00
toc: true

---
这几天看了《Docker 实践》，写了一点自己不知道或者想记录下来的内容，所以是一份笔记，但不是一份基础教程。

<!--more-->

# 1. 第一部分：Docker 基础
## Docker 的优势
* 通过将环境打包成镜像的方式来标准化系统环境，需要使用这个环境的人可以直接使用镜像，无须重头配置环境。所以，Docker 在很多情况下可以作为虚拟机的替代使用。
* 对 Linux 用户而言，Docker 镜像没有依赖，所以非常适合用于打包软件。

## 关键概念：镜像和容器
简而言之，容器运行着由镜像定义的系统，而镜像本质上是一个文件系统，由一个或多个层加上一些 Docker 的元数据组成。  
我们可以从一个镜像中生成多个容器，这些容器完全隔离，其行为不会相互影响。

一个巧妙的类比：镜像和容器的关系，就相当于类和对象的关系。

创建 Docker 镜像有四种标准的方式：
* Docker run & docker commit: 手工创建镜像
* Dockerfile
* Dockerfile 及配置管理（configuration management）工具
* 从头创建镜像并导入一组文件（FROM scratch & ADD sth）

Docker 容器修改文件时会使用**写时复制**（copy-on-write）的方式：  
容器的最顶层是一个可写层，当容器需要修改文件时，docker 会将该文件从下面的只读层复制到可写层，再在可写层对文件进行修改。  
在 `docker commit` 时，这个可写层将会冻结，变为一个具有自身标识符的只读层。

## 使用技巧
1. 以守护进程方式运行容器：
    `-d` 参数会让镜像在后台运行；`--restart` 参数指定了容器重启的条件：
    | 策略 | 描述 |
    | - | - |
    | no | 容器退出时不重启 | 
    | always | 容器退出时每次都会重启 | 
    | on-failure[:max-retry] | 只在失败时（返回非 0 状态码）时重启 |

2. 如果想要移动 Docker 存储数据的位置，则在启动 `docker daemon` 时，使用 `-g` 参数并指定新位置；
3. 实现容器间通信：在 `docker run` 时使用 `--link <hostport>:<container>:<containerport>` 参数可以将另外一个容器的某个端口映射到当前容器的端口中以；  
    实现原理是更改当前容器的 `hosts` 文件；  
    前提条件：构建镜像时必须用 `EXPOSE` 命令暴露容器的端口。
4. 在线查找镜像：使用 `docker search` 功能。

# 2. 第二部分：Docker 与开发

## 用 Docker 代替虚拟机

可以考虑使用 Docker 来代替虚拟机，但由于缺少 systemd 等工具，所以可以考虑用 `supervisord` 托管服务。

Docker 和虚拟机的差异：
* Docker 面向应用，而虚拟机面向操作系统
* Docker 容器和其它容器共享操作系统，而每个虚拟机独享一个操作系统
* Docker 被设计成只运行一个主要进程，而不是管理多组进程

## 构建镜像
### Dockerfile 的使用
* Dockerfile 的用途：从给定镜像开始，为 Docker 指定一系列的 shell 命令和元指令，从而产出最终所需的镜像
* `ADD` 和 `COPY` 命令的区别：`ADD` 会自动在镜像内解压归档文件（如 `.tar` 或 `.tar.gz`），但 `COPY` 只会单纯复制文件。按需使用。
* `ADD` 命令可以将一个 URL 对应的文件添加到容器，但通过 URL 下载的文件不会自动解压。
* 在 `RUN` 命令中使用命令链，有助于减小镜像层数，缩小容器体积。而且将 `apt-get update` 和 `apt-get install` 命令连起来，可以保证每次构建时所装的软件都是最新的，而不会从之前缓存的索引中安装一个旧版本软件。
* 如果希望手动清除某一层的缓存，可以在命令后面加一条注释，如 `ADD a /a # bust the cache`
* `ENTRYPOINT` 指定了镜像的入口点，用户在 `docker run` 时所写的命令都是入口点执行文件的参数。  
    如果不想用镜像的入口点，则需要在 `docker run` 的时候添加 `--entrypoint=xxxx` 选项以重载入口点。

* `ENTRYPOINT` 和 `CMD` 的区别：
    `ENTRYPOINT` 指定了容器入口点，而 `CMD` 指定了入口点程序的默认参数。  
    假设 Dockerfile 为：
    ```dockerfile
    FROM ...
    ENTRYPOINT ['/entrypoint.sh']
    CMD ['xxx', 'yyy']
    ```
    `docker run <image>` 时，会执行 `/entrypoint.sh xxx yyy`；
    `docker run <image> a b` 时，`a b` 会覆盖掉 `CMD` 的值，而不会覆盖入口点，所以会执行 `/entrypoint.sh a b`

* `ENTRYPOINT` 和 `CMD` 命令的参数形式：
    这两个命令的参数有两种形式，一种为字符串类型 `CMD /entrypoint.sh a b`，一种为数组类型 `CMD ['/entrypoint.sh', 'a', 'b']`，其中字符串类型的参数在实际执行前会在前面加上 `bash -c` 命令变成 `bash -c '/entrypoint.sh a b'`，但数组类型的参数则不会改变，直接运行 `/entrypoint.sh`.  
    两种方法有利有弊，按需使用。

### 对镜像的操作
* 扁平化镜像
    如果想将镜像中的多层合为一层（如在某层中添加了密钥又在后面删除），则可以在运行容器之后，使用 `docker export <container> | docker import some-image` 来讲容器的**目录结构**导出为 tar 文件，然后再以此重新制作镜像。这样的镜像只会有一层。

* 对容器进行逆向工程
    书里有个脚本，但是不能用；从 [StackOverFlow](https://stackoverflow.com/questions/19104847/how-to-generate-a-dockerfile-from-an-image) 上找了一个可以用，但是都不如我在 [Portainer](https://portainer.io/) 里看的全:joy:  
    用这些方法可以逆向出一部分命令，比如 `MAINTAINER` `EXPOSE` `RUN` 等，但由于**构建上下文**的缺失，`ADD` 命令只能显示出添加文件的哈希值和容器内路径，并不能知道具体添加的文件是什么样子的。

### 减小镜像体积的方法
上文提到的“扁平化镜像”方法可以有效地减少构建时镜像分层所带来的开销；除此之外，还有一些方法可以减小容器的体积：
1. 使用一个更小的基础镜像：ubuntu 有数十 MB，而 alpine 只有几 MB
2. 自己事后清理：可以在装完软件包以后用 `apt clean` 等命令删除缓存和软件包索引
3. 将一系列命令设置为一行，这样可以减少层数
4. 编写一个脚本来完成安装：原理同 3，只不过不需要在 Dockerfile 中写太多代码
5. 删除不必要的软件包和文档文件：进入容器中，删除所有用不到的文件（甚至基础的可执行文件），并将容器导出（至于这么拼嘛 :new_moon_with_face:）
6. 特殊情况——系统只需要一个带静态链接的二进制文件（如 go 编译后的文件）：用 `scratch` 就可以了，绝对小

进行静态编译，并将可执行文件放入另一个容器中：
书中做出了 `CMD ["cat", "/go/bin/go-web-server"]` `docker run go-web-server > go-web-server` 这样的的操作用来跨镜像复制文件。  
但自从 17.05 版本引入多阶段构建（multi-stage build）后，这个繁琐的过程已经不需要了，构建程序和添加程序的操作可以在一个 Dockerfile 中完成，具体可以参见 [Docker 文档](https://docs.docker.com/develop/develop-images/multistage-build/)。

## 运行容器
### 容器中的服务
> 在 Docker 的世界里，公认的最佳实践是尽可能多地把系统拆分开，直到在每个容器上都只运行一个“服务”，并且所有容器都通过链接相互连通。

如果想在容器中管理多个进程，可以考虑用 `supervisord`，或者使用 `phusion/baseimage`。参见[这篇文章](https://www.summershrimp.com/2018/08/run-multi-service-in-one-container/)。

### 在 Docker 中使用外部数据卷  
除了在 `docker run` 时使用 `-v` 参数以外，我们还可以定义数据容器，然后在运行其它容器时使用 `--volumes-from` 标志。  
使用数据容器可以在多个容器共享数据卷时更方便管理数据卷。  
例：需要改变其中一个容器的挂在路径时，如果不使用数据容器，则需要在多个容器的启动脚本中修改 `-v` 参数的值，而使用数据容器后，只需要更改数据容器就可以了。  
使用数据容器中的卷并不需要让容器处在运行状态，所以可以在运行时使用 `/bin/true` 等命令，让数据容器创建后立即退出。  
注意：多个容器共享数据容器时，同时写入同一文件可能会导致数据卷中的数据**被覆盖或截断**。

PS: 刚刚遇到了一个宿主机文件更改但未同步至容器的问题，可以参考[这个帖子](https://forums.docker.com/t/modify-a-file-which-mount-as-a-data-volume-but-it-didnt-change-in-container/2813/14)最后面的解释。

### 删除数据卷  
为了保证数据安全，Docker 在删除容器时不会自动删除容器锁关联的数据卷，用户可以选择手动将这些数据卷清除  
如果希望删除容器时自动删除数据卷，可以在 `docker rm` 中加入 `-v` 标志。

### 解绑（detach）容器
如果想要从一个容器的交互会话中退出，可以按 `Ctrl+P Ctrl+Q`，Docker 检测到这个按键序列后，就会自动解绑容器，但同时容器依旧会在后台运行。  
如果想重新回到容器中，可以用 `docker attach` 命令。  
这个操作和 `docker run -d` 然后 `docker exec` 有点相似，但上面的方法操纵的是镜像内 PID 为 1 的“主进程”，而 `exec` 命令会新启动一个新的进程给当前 tty 使用。

### 在运行的容器里执行一些命令
如果容器主进程不是 shell 程序而是一些别的，可以用 `docker exec` 命令进入容器，这样 Docker 会在容器中新开一个进程给用户来使用。  
`docker exec` 有三种“模式”：
* 基本的运行模式，同步运行命令，成功后退出，如 `docker exec ps`；
* 守护进程模式，立即退出，命令在后台执行，如 `docker exec -d nginx -g daemon off`；
* 交互模式，就是 `-it` 的样子啦，允许与进程进行交互，如 `docker exec -it bash`

## 使用技巧
* 如果想让镜像立刻完成任务退出，可以使用 `/bin/true` 作为镜像启动命令，也可以用 `touch /somefile`，我更喜欢用第一个；
* 如果想让镜像启动后立即挂起，可以使用 `sleep infinity`，或 `tail -f /etc/hosts` 等作为启动命令；

# 3. 第三部分：Docker 与 DevOps

第三部分主要涉及到将 Docker 应用至 DevOps 流水线中，并在本地利用 Docker 模拟一些生产环境的网络条件（如高延迟、丢包等）来对服务的健壮性进行测试。  
由于还没有对这一部分进行实践，所以这部分的内容只会进行一些摘抄和总结。

又是基本概念：
> 持续集成：持续集成是指用于加快流水线的一个软件生命周期策略。在每次代码库发生重大修改时，通过自动重新运行测试，可以获得更快且稳定的交付，因为被交付的软件具有一个基础层次的稳定性。
> 
> Docker 的可移植性和轻量性，使其成为 CI 从节点（一台供 CI 主服务器连接以便执行构建的机器）的理想选择。与虚拟机从节点相比，Docker CI 从节点向前迈了一大步（相对构建裸机更是一个飞跃）。它可以使用一台宿主机在多种环境上进行构建、快速销毁并创建整洁的环境来确保不受污染的构建，来使用所有熟悉的 Docker 工具来管理构建环境。

## CI 技巧
如果是开源项目，可以考虑用 Docker Hub 工作流完成自动构建；  
如果是本地构建，可以为包管理器安装一个 Squid 代理，通过缓存软件包来加快软件下载速度，同时节省流量。

## CI/CD 流水线
> CD 背后的关键思想之一是构建提升。构建提升是指流水线的每个场景（用户验收测试、集成测试以及性能测试）只有在前一个场景成功时才能触发下一个场景。

## “Docker 契约”
在 CD 全过程中，从 CI 产出的镜像必须是最终的、不可修改的。  
这样，在不同团队、不同环境中运行的代码和依赖才可以被彻底固化，有利于问题的复现及排查。

## 微服务架构
etcd 可作为环境的中央配置存储，服务发现可以用 etcd、confd 及 nginx 的组合来实现。

## 网络模拟：无痛的现实环境测试
### Docker Compose: 管理容器间链接
之前写到可以用链接的方式链接容器从而实现容器间通信，但链接的配置比较繁琐，而且出现问题难以恢复（需要依序重启所有容器才能重置所有链接）。  
所以，如果需要启动一组互相连接的容器，可以使用 Docker Compose.  
Docker Compose 的 YAML 配置可以使容器的管理变得十分简单，它把编排容器的复杂事务从手工且易出错的过程变成了可通过源代码控制的更安全和自动化的过程。

```yml
echo-server:
    image: server
    expose:
    - "2000"
client:
    build: .                # 使用 ./Dockerfile 来构建镜像
    links:
    - echo-server:talkto    # 这里的参数与 --link 的参数一致
```

`docker-compose up` 后，client 容器就可以通过 `talkto` 的 host 与 echo-server 通信。

但是，这里的 host 解析是静态的，如果希望在容器内使用可动态配置的 DNS，可以引入 `resolvable`.

### 网络测试
* 想要为单个容器应用不同的网络状况，可以用 Comcast
* 想要对大量容器进行网络状况编排设置，可以用 Blockade
* 想要跨宿主机进行容器间无缝通信，可以使用 Weave 构建基底网络
* `docker network` 提供试验性的网络构建功能

啊…随着容器编排和 Service Mesh 框架的出现，貌似这些问题都可以更轻松地解决了 \_(:з」∠)\_

