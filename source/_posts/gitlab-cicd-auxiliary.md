title: 'GitLab CI/CD: 辅助工具'
author: David Dai
tags:
  - DevOps
  - GitLab
  - CI/CD
categories:
  - Devops
date: 2018-07-15 20:21:00
---
本文会讲一些在 GitLab CI/CD 中可能会用到的辅助工具，包括隐藏任务、依赖缓存、定时任务以及部署环境。

<!--more-->

# 0. TL;DR

* [Hidden keys (jobs)](https://docs.gitlab.com/ee/ci/yaml/#hidden-keys-jobs)
* [Cache dependencies in GitLab CI/CD](https://docs.gitlab.com/ee/ci/caching/index.html)
* [Pipeline Schedules](https://docs.gitlab.com/ee/user/project/pipelines/schedules.html)
* [Introduction to environments and deployments](https://docs.gitlab.com/ee/ci/environments.html)

# 1. 隐藏任务
先讲个简单的。

有的时候我们需要在 Pipeline 中跳过某些任务，通常情况下我们可以用任务定义中的 `when` 和 `except` 属性来控制任务是否显示。但是如果我们想暂时删掉这个任务怎么办？  
一种方法，是在 `.gitlab-ci.yml` 中删掉或注释掉这个任务；另一种做法是，直接在任务定义的 key 中加个点号(`.`)，就可以把这个任务隐藏起来。这种做法和 Linux 中隐藏文件的方法非常相似，也是 GitLab 官方推荐的做法。

# 2. 依赖缓存
之前我们的项目依赖是直接打在 Docker 镜像里的，但是后来技术更新后，单元测试使用的镜像变成了构建用的 `pymicro`，内部不包含任何依赖，需要在运行测试之前从头开始安装，为此会耗费大量时间。于是，我就想把这些依赖文件缓存起来。于是找到了[这个文档](https://docs.gitlab.com/ee/ci/caching/index.html#caching-python-dependencies)，在测试任务运行时，使用 `virtualenv` 将所有依赖放进项目目录下，并配置缓存，这样在任务运行**成功**后，CI 系统会将依赖缓存起来，保存在 Runner 中，下次运行时就不需要重新安装依赖了。  
经实测，加上依赖缓存可以使我们的测试任务运行时间从两分半缩短到一分半。于是很开心地省下了无数个 ~~1s~~ 一分钟。

配置文件如下：
```yaml
test_all:
  image: "/pymicro"
  stage: test_all
  variables:
    GIT_STRATEGY: fetch
    PIP_CACHE_DIR: "$CI_PROJECT_DIR/.cache"
  before_script:
    - "[ -e venv ] || ( pip install virtualenv -i https://mirrors.aliyun.com/pypi/simple && virtualenv venv )"
    - source venv/bin/activate
    - pip install -U -r test_requirements.txt -i https://mirrors.aliyun.com/pypi/simple
  script:
    - flake8 app
    - pytest tests
  cache:
    paths:
      - .cache/
      - venv/
    key: requirement-cache
```

需要注意的是，cache 只能缓存项目目录下的文件，不能缓存其它目录的文件，比如 `/opt/` 什么的。所以，我们必须用 `virtualenv` 或者 `pipenv` 将所有 `site-packages` 存在项目目录下。

再就 `before_script` 里面的第一句和第三句多说两句：
1. 一开始是照着文档来写的，但是那样的话每次都要重新安装 virtualenv，并且还要重复创建虚拟环境。这也是安装依赖啊…并且就算是用 `PIP_CACHE_DIR` 把依赖包缓存在本地，创建 `virtualenv` 是也要安装 `setuppools` 和 `pip` 等，依然很慢 :new_moon_with_face: 于是干脆加了个判断，如果有 `venv` 这个目录，就直接跳过创建虚拟环境的阶段。
2. 我们的 `requirements.txt` 是不锁依赖版本的，所以 `pip install -U -r` 可以在每次运行时对本地缓存的依赖进行更新，这样虽然缓存了依赖文件，但 `pip` 依然会和 registry 进行网络交互。去掉 `-U` 的话，依赖就不会被更新，但是 `pip install` 的执行时间会直接降低到一两秒钟。

# 3. 定时任务
更准确的叫法，应该叫**定时**流水线(Scheduling Pipelines)。

在项目的 CI/CD → Pipeline 菜单中，我们可以配置定时任务。定时的配置方式与 Crontab 的配置方式相同，还可以选择这个定时 pipeline 所使用的时区。

![定时任务](/pics/cicd/schedule.png)

配置好后就可以看到设置的定时任务，到时间就会在某个分支上自动触发。

![定时任务列表](/pics/cicd/schedule-pipeline.png)

到现在为止，我们平时的 Pipeline 和定时任务中执行的任务是一模一样的。但如果我们有一些特殊的任务需要只在定时任务中执行，可以在 job 的 `when` 属性中写入 `- schedules`；同样，如果某些任务不应该在定时任务中执行，配置一下 `except` 属性就可以了。

# 4. 部署环境
如果项目的 CD 流程在 GitLab 中进行的话，可以考虑在 `.gitlab-ci.yml` 中配置部署任务执行所在的环境：

```yaml
deploy_rpc:
  stage: deploy_production
  only:
    - master
  except:
    - schedules
  tags:
    - deploy-production
  when: manual
  environment:
    name: production/rpc
  script:
    - kubectl set image deploy/project-rpc "app=${IMAGE_TAG}"
```

配置完成后，在执行这个任务时，GitLab 就会从配置中读取环境配置，并记录当前环境部署时项目所在的 Git commit. 随后，我们就可以在 CI/CD → Environments 菜单中看到这个环境的部署情况。

![环境列表](/pics/cicd/environment.png)

点击环境右边的执行（:arrow_forward:）按钮，可以方便地将这个 commit 的代码部署到其它环境中；点击环境名，可以看到这个环境的部署历史，我们可以在这里方便地将环境中的代码回滚到之前的版本。

![部署历史](/pics/cicd/deploy-history.png)

GitLab 还提供了一个贴心的功能：我们可以直接在 Merge Request 的页面中看到当前 MR 的部署进度，以及该 MR 部署至每个环境的时间点。

![MR](/pics/cicd/environment-pr.png)
