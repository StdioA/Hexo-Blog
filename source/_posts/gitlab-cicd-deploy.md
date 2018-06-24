title: GitLab CI/CD 基础教程（二）
author: David Dai
categories:
  - DevOps
tags:
  - DevOps
  - GitLab
  - CI/CD
date: 2018-06-18 15:46:00
---

本文是 GitLab CI/CD 系列的第二篇，主要介绍 GitLab CI Runner 在 Docker 和 Kubernetes 环境下的部署方式。

<!-- more -->

# 0. TL;DR
文档在[这儿](https://docs.gitlab.com/runner/)。

# 1. GitLab Runner 的运行环境及执行环境选择
GitLab Runner 用 Go 语言写成，最后打包成单文件进行分发，所以可以在很多平台下快速运行，包括 Windows / GNU Linux / MacOS 等，同时也提供 Docker 镜像，方便在 Docker / Kubernetes 环境中部署。

但除了 Runner 运行外，Runner 还需要一个环境来运行 jobs. 这个环境称之为执行环境（executor）。GitLab Runner 支持多种执行环境，包括 SSH，Docker，VirtualBox 等。  
不同执行环境对 GitLab CI/CD 不同功能的支持情况，可以看官方文档中的[兼容性表格](https://docs.gitlab.com/runner/executors/README.html#compatibility-chart)。

由于我司主要用 Docker 或 Kubernetes 来托管服务，并且在进行测试时需要 service 的支持，所以自然只剩下了两种选择，Docker 和 Kubernetes. 在下文中会讲解 GitLab 在 Docker 和 Kubernetes 中的部署方式。  
尽管运行环境和执行环境可以相互独立，但为了方便起见，我更推荐在同一个环境中运行 runner daemon 和 jobs.

# 2. GitLab Runner 部署

简单来讲，GitLab Runner 的部署方式分为两步：运行、注册。

## 2.1 注册 runner
在一个 runner daemon 进程中，我们可以同时注册并运行多个 runner，来并行完成多个场景下的不同任务。

注册的步骤在不同平台下大同小异：运行 `gitlab-runner register` 命令，会输出一个交互界面，在里面依次输入 GitLab 实例地址、CI Token、Runner 描述、标签列表（用于区分不同类型的 Runner，使不同阶段的 job 在不同的 Runner 中运行）、执行环境类型，如果选择基于 Docker 的执行环境，则需要再输入一个缺省的 job image.  
注册之后，Runner 会将配置写入 `/etc/gitlab-runner/config.toml` 文件中，如果文件内容不丢失，`gitlab-runner` 程序会自动读取配置内容并运行，无需重复注册。  
当然，如果你需要一些更高级的配置，则可以直接修改 `config.toml`. 具体的配置可以见[官方文档](https://docs.gitlab.com/runner/configuration/advanced-configuration.html)。在配置文件更改后，runner daemon 会自动重新加载配置，无需重启。  
此外，`gitlab-runner register` 非常鬼畜的一点在于，`config.toml` 中的部分配置可以通过环境变量注入，且**所有**配置都可以通过注册命令的参数传入。也就是说，如果你希望不接触 `config.toml`，只使用一条命令来配置并注册单个 runner，`gitlab-runner register` 命令完全可以满足你的要求。:new_moon_with_face:

## 2.2 在 Docker 环境中部署
一句话：

```bash
docker run -d --name gitlab-runner --restart always \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /data/gitlab-runner/conf:/etc/gitlab-runner \
    --net=host \
    gitlab/gitlab-runner:latest
```

注意：我们在这里将 `/var/run/docker.sock` 挂载进了 gitlab-runner 容器，这也就意味着我们将 Docker 环境的控制权交给了 runner daemon，这样 runner daemon 可以在收到任务指令时，使用当前 Docker 环境作为执行环境，在里面运行容器以执行任务。  
当然，这样挂载是一种相当危险的做法。如果你比较担心安全问题，可以考虑做一个 docker compose，并在里面运行一个 dind （Docker in Docker）来作为 Runner 的执行环境。

## 2.3 在 Kubernetes 环境中部署
在 k8s 环境中的部署要稍微复杂一些（因为配置文件太长:joy:），大体需要配置以下五部分：
* 一个或者两个 `Namespace`（取决于你是否要把 job pod 和 daemon 放在同一个 namespace 里）
* 一个 `ConfigMap`，用于注入 runner 配置
* 一个 `Deployment`，用于运行 runner daemon
* 一个 `ServiceAccount`，用于给 daemon 使用，来启动 pod，为运行任务提供环境（当然，用 `default`）也不是不可以
* 一个 `Role` + `RoleBinding`，为上面的 `ServiceAccount` 赋予 `pods` 和 `pods/exec` 权限

runner 的配置注入有两种方式：
1. 先在 GitLab 中注册好 runner，然后在 `ConfigMap` 中写好配置文件，并在 `Deployment` 中作为一个 Volume 挂载到配置目录下。这是官方文档中推荐的做法；
2. 在 `ConfigMap` 中定义环境变量，并在 pod template 的 `envFrom` 属性中定义一个 `configMapRef` 来注入环境变量，并在 pod 启动时当场注册一个 runner；在 pod 被停止前再调用命令将这个 runner 注销掉（推荐，否则 GitLab 后台会看到很多离线的 runner）。

由于编辑器的高亮规则通常不会处理一个文件里出现两种不同语言的情况，在 `ConfigMap` 里写 toml 会很挣扎，所以我使用的是第二种注入方式。

完整的配置文件：
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cicd
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: executor
  namespace: cicd
imagePullSecrets:
- name: dockersecret
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: cicd
  name: executor-role
rules:
  # runner 要新建 pod，所以为它赋予 pod 相关的权限
  - apiGroups: [""]
    resources: ["pods", "pods/exec"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: cicd
  name: executor-rolebinding
subjects:
- kind: ServiceAccount
  name: executor
  namespace: cicd
roleRef:
  kind: Role
  name: executor
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: ConfigMap
metadata:
  namespace: cicd
  labels:
    app: gitlab-deployer
  name: gitlab-runner-cm
data:
  # 具体可用的参数配置以及环境变量配置可以运行 gitlab-runner register --help 查看
  REGISTER_NON_INTERACTIVE: "true"
  REGISTER_LOCKED: "false"
  CI_SERVER_URL: "https://gitlab.com/ci"
  METRICS_SERVER: "0.0.0.0:9100"
  RUNNER_CONCURRENT_BUILDS: "4"
  RUNNER_REQUEST_CONCURRENCY: "4"
  RUNNER_TAG_LIST: "tag1,tag2"
  RUNNER_EXECUTOR: "kubernetes"
  KUBERNETES_NAMESPACE: "cicd"
  KUBERNETES_SERVICE_ACCOUNT: "executor"
  KUBERNETES_CPU_LIMIT: "100m"
  KUBERNETES_MEMORY_LIMIT: "100Mi"
  KUBERNETES_SERVICE_CPU_LIMIT: "100m"
  KUBERNETES_SERVICE_MEMORY_LIMIT: "100Mi"
  KUBERNETES_HELPER_CPU_LIMIT: "100m"
  KUBERNETES_HELPER_MEMORY_LIMIT: "100Mi"
  KUBERNETES_PULL_POLICY: "if-not-present"
  KUBERNETES_TERMINATIONGRACEPERIODSECONDS: "10"
  KUBERNETES_POLL_INTERVAL: "5"
  KUBERNETES_POLL_TIMEOUT: "360"
  KUBERNETES_IMAGE: "kubectl:1.8.1"
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: runner
  namespace: cicd
  labels:
    app: runner
spec:
  replicas: 1
  selector:
    matchLabels:
      app: runner
  template:
    metadata:
      labels:
        app: runner
    spec:
      containers:
      - name: ci-builder
        image: gitlab/gitlab-runner:v10.6.0
        command:
        # 命令有点长，做了以下几步：注销当前的 runner name 以防止 runner 冲突；注册新的 runner；启动 runner daemon
        - /bin/bash
        - -c
        - "/usr/bin/gitlab-runner unregister -n $RUNNER_NAME || true; /usr/bin/gitlab-runner register; exec /usr/bin/gitlab-runner run"
        imagePullPolicy: IfNotPresent
        envFrom:
        # 通过 ConfigMap 注入 runner 配置
        - configMapRef:
            name: gitlab-runner-cm
        env:
        # 通过 Secret 注入与 GitLab 实例进行交互所用的 CI Token
        # runner 命令会自动从环境变量中读取这个 token，用于注册 runner
        - name: CI_SERVER_TOKEN
          valueFrom:
            secretKeyRef:
              name: gitlab-ci-token
              key: token
        # 动态注入环境变量，使用 pod name 作为 runner name
        # 刚查了一下文档，如果不通过环境变量指定 runner name 的话，会用当前环境的 hostname，也就是 pod name 来做 runner name
        # 那完全没必要把这个 pod name 注册进去嘛…
        - name: RUNNER_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        # gitlab-runner 自带 Prometheus metrics server，通过上面的 METRICS_SERVER 环境变量配置
        # 强的一比！
        ports:
        - containerPort: 9100
          name: http-metrics
          protocol: TCP
        resources:
          limits:
            cpu: "100m"
            memory: "100Mi"
          requests:
            cpu: "100m"
            memory: "100Mi"
        lifecycle:
          # 在 pod 停止前，注销这个 runner
          preStop:
            exec:
              command:
                - /bin/bash
                - -c
                - "/usr/bin/gitlab-runner unregister -n $RUNNER_NAME"
      restartPolicy: Always
```

这里注意两点：
1. Kubernetes executor 不支持在运行 job 是使用 service alias，所以访问服务时都要通过 `127.0.0.1` 来访问；
2. 你可能注意到，配置里出现了三种 resource limit 配置：`KUBERNETES_CPU_LIMIT`, `KUBERNETES_SERVICE_CPU_LIMIT`, `KUBERNETES_HELPER_CPU_LIMIT`，这里三个配置将应用于 kubernetes job pod 中的三类容器，第一类用于实际执行命令，第二类（service 容器）用于启动 job 所需的 service，第三类（helper 容器）用于任务执行之前的代码拉取，以及任务执行之后构建产物（artifact）的上传。我们可以通过配置来为三种容器赋予不同的资源限制。

至此，我们可以在 Docker 和 Kubernetes 环境下部署 GitLab Runner.
