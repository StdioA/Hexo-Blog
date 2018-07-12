title: GitLab CI/CD 基础教程（三）
author: David Dai
tags:
  - DevOps
  - GitLab
  - CI/CD
categories:
  - DevOps
date: 2018-06-23 20:53:00
toc: true

---
前两篇我们讲了 GitLab CI/CD 的简单应用及部署方式，这一篇简单讲一下如何将 GitLab CI/CD 与日常开发部署流程结合。

<!-- more -->

# 0. TL;DR
看文档？其实就是简单应用吧。  
`.gitlab-ci.yml` 的完整配置定义可以见[第一篇博文](/2018/06/gitlab-cicd-fundmental/#2-cicd-流程配置)。

# 1. 测试阶段
测试阶段没什么好说的，只需要把 runner tag 打好（注册时使用 `--tag-list` 参数），基于 docker/k8s 把 Runner 搭起来，基本上就可以自动运行了。

`.gitlab-ci.yml` 配置如下：
```yaml
test_all:
  image: "pymicro"
  stage: test
  services:
    - name: mysql:5.6
      alias: mysql
      command: ["mysqld", "--character-set-server=utf8mb4", "--collation-server=utf8mb4_unicode_ci"]
  veriables:
    MYSQL_DATABASE: db
    MYSQL_ROOT_PASSWORD: password
  before_script:
    - pip install -U -r requirements.txt
  script:
    - flake8 app
    - pytest tests
```

这里定义的两个环境变量都是给 MySQL 服务用的，mysql 镜像会在容器启动时读取某些环境变量，来配置数据库。具体支持的环境变量可以参考 [MySQL 的 docker image 页面](https://hub.docker.com/_/mysql/)。  
我们可以在 service 中自定义启动命令，这里我将 MySQL 的默认字符集设置成了 `utf8mb4`，否则服务器中的数据库字符集会是 `latin1`.

需要注意的一点是，基于 Docker 部署的 Runner，可以使用服务别名, 也就是在跑测试的阶段中，可以通过 service alias 访问到对应的服务；而基于 Kubernetes 的 Runner 不支持，所以只能通过 `127.0.0.1` 访问。

# 2. 构建阶段
构建阶段中，我们会用 Docker 将工程打包成镜像，并推送到远端 registry.

## 2.1 基本配置

`.gitlab-ci.yml` 配置如下：
```yaml
build_image:
  image: "docker:17.11"
  stage: build
  services:
    - name: "docker:17.12.0-ce-dind"
      alias: dockerd
  variables:
    DOCKER_HOST: tcp://127.0.0.1:2375
    IMAGE: docker.registry/name/${CI_PROJECT_NAMESPACE}-${CI_PROJECT_NAME}
  before_script:
    - IMAGE_TAG=${IMAGE}:${CI_COMMIT_SHA:0:8}
  only:
    - master
  tags:
    - build
  script:
    - docker build -t ${IMAGE_TAG} -f Dockerfile .
    - docker push ${IMAGE_TAG}
```

在这个任务中，我们启用了一个 dind 作为 service，并使用 `DOCKER_HOST` 环境变量来让 docker 命令与我们的 dind 服务通信。  
任务执行时，会根据项目中的 `Dockerfile` 构建并推送镜像。

我们的镜像名称使用了项目组名 + 项目名的配置，tag 使用 commit SHA 前八位来构成。因为在 `variable` 字段中定义环境变量时，不能使用 `${CI_COMMIT_SHA:0:8}` 这种 shell 字符串操作，所以只好在 `before_script` 中来定义这个环境变量。

这里，我们需要在 docker 环境中启动一个 dind，来作为构建时所用的服务器。值得注意的一点是，如果你需要使用 dind，则 dind 所在的 container 应该具有[特权](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities)（[官方文档](https://docs.gitlab.com/runner/executors/kubernetes.html#using-docker-dind)也有讲到）。所以在 Runner 注册时，需要加上 `--docker-privileged` 或 `--kubernetes-privileged` 参数（具体视执行平台而定），来使 job 运行时所在的 container 拥有特权。不过，在部署 runner 时，Runner daemon 所需的 container 并不需要这个特权（其实可以机器上的 docker service 或者操控 pod 已经是一种特权了:joy:）。

## 2.2 dind 服务调（luan）优（gao）

### 2.2.1 dind 成为独立服务
上面定义的 job 中，dind 是一个 job service，也就是说，每次构建的时候都会从头开始构建。而 docker 构建提供了一套比较完善的缓存功能，如果 `Dockerfile` 某几层的构建命令完全一样（比如只是 `RUN apt-get install xxx`）的话，Docker 会在再次构建时自动使用之前已经构建好的层，这样可以减少构建时间。  
所以我在做完上面的那个流程之后立刻意识到了这一点，于是单独把 dind 从 CI 任务中抽离了出来，在 k8s namespace 中单独搭建了一个 dind 服务，并定义了 k8s service，而 job 中的 `DOCKER_HOST` 环境变量也改成了 `tcp://dockerd:2375`，因为 dind 并不在 job pod 里了，而是一个 k8s service，需要通过 DNS 来获取到具体的 IP.  
这样一来，在 job 结束后，dind 依然存在，并会保留前一次的构建层，这样下次构建的时候就可以跳过依赖安装步骤，大大缩短了构建所需的时间。

在 k8s 中搭建 dind 服务的内容不在本博文讲述范围内，想搭建的话，可以去 Google 一下。

### 2.2.2 尼玛… Node 存储空间满了？
在我们的 dind 服务运行起来一段时间后，就遇到了一个尴尬的问题：dind 服务占用了太多存储空间，导致 pod 的所在 node 出现了问题…  
这个问题有两种解决方案：一是单独做一个 node，并用 node selector 将 dind 单独放在那个 node 上，以避免影响其它服务；二是为 dind 单独挂一个 volume，使用 `PersistentVolume` 进行持久化存储。

而我司的解决方案简直是骚断腿：单独买一台 VPS，在上面搭一个 docker 服务器，然后把这个服务引入 k8s 集群中 :joy:  
 
所以，在搭建好服务器之后，修改 k8s 中的 `Service`，为 `Service` 添加 `Endpoint`，注意端点名称要与服务名称一致：

```yaml

kind: Service
metadata:
  name: dockerd
  namespace: cicd
  labels:
    app: dockerd
spec:
  ports:
  - protocol: TCP
    port: 2375
    targetPort: 2375
  clusterIP: Nonetoc: true

---
kind: Endpoints
apiVersion: v1
metadata:
  name: dockerd
  namespace: cicd
  labels:
    app: dockerd
subsets:
  - addresses:
      - ip: 192.168.8.45
    ports:
      - port: 2375
```

这样我们在集群中查看 `dockerd` 的 IP 地址时，就会得到那台 VPS 的 IP 了。  
然后…记得写个 cron job，定期清理 docker 服务器上的缓存，否则硬盘也是会满的\_(:з」∠)\_。

# 3. 部署阶段
部署阶段中，我们会使用 `kubectl set image` 命令，对特定 `Deployment` 触发一次滚动更新。

`.gitlab-ci.yml` 配置：
```yaml
deploy_production:
  image: "kubectl:1.8.1"
  stage: deploy
  variables:
    GIT_STRATEGY: none
  variables:
    IMAGE: docker.registry/name/${CI_PROJECT_NAMESPACE}-${CI_PROJECT_NAME}
  before_script:
    - IMAGE_TAG=${IMAGE}:${CI_COMMIT_SHA:0:8}
  only:
    - master
  when: manual
  tags:
    - deploy-production
  script:
    - kubectl set image deploy/myproject "app=${IMAGE_TAG}" --record
```

这里我们定义了一个 `GIT_STRATEGY` 环境变量，有了这个变量，在 CD 任务执行时，Runner 会跳过克隆代码的步骤，因为在这个阶段中我们并不需要项目代码。而 `when: manual` 属性表示这个任务需要手动触发。

这个阶段中，我们需要一个 `ServiceAccount` 来让 pod 使用 kubectl 与集群通信；同时为了保证 `set image` 命令的成功执行，我们还需要为这个账户赋予一些权限。

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: deployer
  namespace: cicd
imagePullSecrets:
- name: dockersecrettoc: true

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: deployer
rules:
  - apiGroups: ["extensions"]
    resources: ["deployments"]
    verbs: ["get", "patch"]toc: true

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cicd-deployer
subjects:
- kind: ServiceAccount
  name: deployer
  namespace: cicd
roleRef:
  kind: ClusterRole
  name: deployer
  apiGroup: rbac.authorization.k8s.io

```

同时，在我们注册这个 runner 时，需要加上 `--kubernetes-service-account deployer` 参数，这样在 job pod 启动时，集群将 deployer 账户的凭据注入进 pod，`kubectl` 命令才能正常使用。

至此，我们（应该）可以实现并部署一套完整的测试→构建→部署流程。  
这个系列到此也就结束了，后面还会有一篇，讲点周边的小工（玩）具。
