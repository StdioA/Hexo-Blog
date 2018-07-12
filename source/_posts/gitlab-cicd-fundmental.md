title: GitLab CI/CD 基础教程（一）
author: David Dai
categories:
  - DevOps
tags:
  - DevOps
  - GitLab
  - CI/CD
date: 2018-06-06 15:46:00
toc: true

---
最近玩了 GitLab CI/CD 平台，通过搭建这个平台也收获了一些关于 DevOps 的基本技能，打算通过几篇文章来讲述一下 GitLab CI/CD 平台的构建及应用。本文对 GitLab CI/CD 以及 CI/CD 流程定义文件的写法做了简要介绍。

<!-- more -->
前几个月公司技术改进，某些业务部署在 k8s 集群中，于是我们开始通过 GitLab 自带的 CI/CD 功能来实现服务的测试、构建及部署，所以才有了这篇文章。

# 1. 基本概念

## 1.1 CI/CD
CI，Continuous Integration，为持续集成。即在代码构建过程中持续地进行代码的集成、构建、以及自动化测试等；有了 CI 工具，我们可以在代码提交的过程中通过单元测试等尽早地发现引入的错误；  
CD，Continuous Deployment，为持续交付。在代码构建完毕后，可以方便地将新版本部署上线，这样有利于快速迭代并交付产品。

## 1.2 GitLab CI/CD
GitLab CI/CD（后简称 GitLab CI）是一套基于 GitLab 的 CI/CD 系统，可以让开发人员通过 `.gitlab-ci.yml` 在项目中配置 CI/CD 流程，在提交后，系统可以自动/手动地执行任务，完成 CI/CD 操作。而且，它的配置非常简单，CI Runner 由 Go 语言编写，最终打包成单文件，所以只需要一个 Runner 程序、以及一个用于运行 jobs 的[执行平台](https://docs.gitlab.com/runner/executors/README.html)（如裸机+SSH，Docker 或 Kubernetes 等，我推荐用 Docker，因为搭建相当容易）即可运行一套完整的 CI/CD 系统。

下面针对 Gitlab CI 平台的一些基本概念做一个简单介绍：

1. Job  
Job 为任务，是 GitLab CI 系统中可以独立控制并运行的最小单位。  
在提交代码后，开发者可以针对特定的 commit 完成一个或多个 job，从而进行 CI/CD 操作。

2. Pipeline  
Pipeline 即流水线，可以像流水线一样执行多个 Job.  
在代码提交或 MR 被合并时，GitLab 可以在最新生成的 commit 上建立一个 pipeline，在同一个 pipeline 上产生的多个任务中，所用到的代码版本是一致的。

3. Stage  
一般的流水线通常会分为几段；在 pipeline 中，可以将多个任务划分在多个阶段中，只有当前一阶段的所有任务都执行成功后，下一阶段的任务才可被执行。  
注：如果某一阶段的任务均被设定为“允许失败”，那这个阶段的任务执行情况，不会影响到下一阶段的执行。

![CI Pipeline](/pics/cicd/pipeline-demonstration.png)
上图中，整条流水线从左向右依次执行，每一列均为一个阶段，而列中的每个可操控元素均为任务。  
左边两个阶段的任务是自动执行的任务，在 commit 提交后即可自动开始运行，执行成功或失败后，可以点击任务右边的按钮重试；而右边两个是手动触发任务，需要人工点击右边的“播放”按钮来手动运行。

# 2. CI/CD 流程配置
## 2.1 完整定义
GitLab 允许在项目中编写 `.gitlab-ci.yml` 文件，来配置 CI/CD 流程。

下面，我们来编写一个简单的测试→构建→部署的 CI/CD 流程。

首先，可以定义流程所包含的阶段。我们的流程包含三个阶段：测试、构建和部署。  
在 `.gitlab-ci.yml` 的开头，定义好所有阶段、以及执行每个任务之前所需要的环境变量以及准备工作，然后定义整个流程中包含的所有任务：

```yaml
stages:
  - test
  - build
  - deploy

variables:
  IMAGE: docker.registry/name/${CI_PROJECT_NAMESPACE}-${CI_PROJECT_NAME}

before_script:
  - IMAGE_TAG=${IMAGE}:${CI_COMMIT_SHA:0:8}

test_all:
  image: "pymicro"
  stage: test
  services:
    - name: mysql:5.6
      alias: mysql
  veriables:
    MYSQL_DATABASE: db
    MYSQL_ROOT_PASSWORD: password
  before_script:
    - pip install -U -r requirements.txt
  script:
    - flake8 app
    - pytest tests

build_image:
  image: "docker:17.11"
  stage: build
  services:
    - name: "docker:17.12.0-ce-dind"
      alias: dockerd
  variables:
    DOCKER_HOST: tcp://dockerd:2375
  only:
    - master
  tags:
    - build
  script:
    - docker build -t ${IMAGE_TAG} -f Dockerfile .
    - docker push ${IMAGE_TAG}

deploy_production:
  stage: deploy
  variables:
    GIT_STRATEGY: none
  only:
    - master
  when: manual
  tags:
    - deploy-production
  script:
    - kubectl set image deploy/myproject "app=${IMAGE_TAG}" --record


```

在每个任务中，通常会包含 `image`, `stage`,`services`, `script` 等字段。  
其中，`stage` 定义了任务所属的阶段；`image` 字段指定了执行任务时所需要的 docker 镜像；`services` 指定了执行任务时所需的依赖服务（如数据库、Docker 服务器等）；而 `script` 直接定义了任务所需执行的命令。

下面简单介绍一下每个阶段中的任务。
## 2.2 测试
在测试任务中，我们启动了 `MySQL` 服务，并通过环境变量注入了 MySQL 的初始数据库以及 Root 密码，在服务启动后，Runner 会运行 `before_script` 中的命令来安装所需依赖；安装成功后就会运行 `script` 属性中的命令来进行代码风格检查以及单元测试；  
可以注意到，我们的 `MySQL` 服务下有一个 `alias` 属性标识服务别名。如果你的 Runner 运行在 Docker 平台下，你可以直接通过服务别名访问到该测试环境中对应的服务。比如在这个任务中，我们就可以用 `mysql://root:password@mysql/db` 来访问测试数据库。

## 2.3 构建
在构建任务中，我们会用 `Dockerfile` 注入依赖，将工程打包成 Docker 镜像并上传；  
我们为这个任务定义了一些额外的属性：`tag` 属性可以标记这个任务将在含有特定 tag 的 CI Runner 上运行；而 `only` 属性表示只有这个 commit 在特定的分支下（如 master）时，才可以在此 commit 上运行这个任务。  
`only` 和 `except` 支持很多种环境条件判断，详细的用法可以参考[官方文档](https://docs.gitlab.com/ee/ci/yaml/README.html#only-and-except-simplified)。  
另外，我们在 `before_scripts` 中，通过环境变量拿到了项目所属的组，以及项目名称。GitLab 会在运行任务前，向环境中注入很多环境变量，来表明运行环境以及上下文。所有的环境变量列表可以看[文档](https://docs.gitlab.com/ee/ci/variables/)。

## 2.4 部署
在部署任务中，我们会用 `kubectl set image` 命令将我们刚刚构建的镜像发布到生产环境。  
这个任务中的 `when` 表示运行该任务所需要的必要条件，如前一阶段任务全部成功。`when: manual` 表示该操作只允许手动触发。该属性具有四个选项，具体请见[文档](https://docs.gitlab.com/ee/ci/yaml/README.html#when)。

至此，我们在 `.gitlab-ci.yml` 中定义了一套完整的测试→构建→部署流程。
