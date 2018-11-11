title: 《Prometheus Book》阅读笔记
author: David Dai
tags:
  - Prometheus
  - 读书
categories:
  - DevOps
date: 2018-11-4 16:38:00
toc: true

---
看了一本在线的小书，叫[《Prometheus Book》](https://yunlzheng.gitbook.io/prometheus-book/)，做了一点摘抄和笔记。

<!--more-->


# 第 1 章：天降奇兵

第一章对 Prometheus 的架构和用法做了简单的介绍。

## 基本介绍
通过建立完善的监控体系，我们可以达到以下目的：长期趋势分析、对照分析、告警、故障分析与定位、数据可视化。  
Prometheus是一个开源的完整监控解决方案，其对传统监控系统的测试和告警模型进行了彻底的颠覆，形成了基于**中央化**的规则计算、统一分析和告警的新模型。  
个人理解，这里的中央化主要是指数据存储中央化：各个数据源将指标暴露出来，由 Prometheus 服务器采集后**统一存储、统一分析**，方便聚合查询。

Prometheus 的基本结构如图所示：

![Prometheus 架构](https://winderresearch.com/img/blog/2017/prometheus/prometheus-architecture.svg)

其中 Prometheus Server 是整个组件中的核心部分，负责实现对监控数据的获取、存储及查询；PushGateway 用于在某些环境中被动收集用户指标，并提供给 Server；在 Server 中定义的报警规则如果被触发，则 AlertManager 会负责该报警的后续处理流程（如通知用户）。

Prometheus Server 通过两种方式来获取数据：
1. 通过 HTTP API 直接访问实例以获取数据
2. 当 Prometheus Server 无法直接访问到实例时，实例可以主动推送数据到 PushGateway，Server 访问 PushGateway 来获取数据。

Prometheus 的基本数据模型：  
一个指标，由它的名字和值组成。  
指标的名字由指标名称以及多对描述样本特征的 KV 标签构成；而指标的值是一个序列，记录了每个时间点的指标值。  
Prometheus 将这个基本指标进行了组合和拓展，构成了四种指标数据结构：Counter、Gauge、Summary 和 Histogram，使用者可以使用这四种数据结构来实现多样化的指标定义和分析模式。

## 安装运行
Prometheus 采用 Golang 编译，不存在第三方依赖，可以直接运行；或者 Prometheus 官方也提供了 Docker 镜像。  
官方提供了 node_exporter 用于采集主机的运行指标。运行后，可以通过配置 `prometheus.yml` 的 `scrape_configs` 来添加任务：

```yml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
```

重启后，查询 `up` 指标，可以看到 `up{job="node"}` 的值为 1，说明来自改端点的指标可以正常获取。


## 杂项
* Prometheus 自带的 UI 界面功能比较简单，此时可以考虑使用第三方的可视化的工具如 Grafana.
* 在 Prometheus 中，每个暴露监控样本数据的 HTTP 服务称为一个实例，一组出于相同采集目的的实例，可以通过一个任务（Job）进行分组管理。

# 第 2 章：探索 PromQL

第二章主要讲解了查询中所用到的数据类型；监控指标类型；以及基本的 PromQL 使用。

## 基本概念
* 时间序列：一个序列，其中的每个值由一个时间戳和一个浮点数构成，表示某一时刻监控样本的具体值
* 样本：时间序列的一个点，由三部分组成：
    * 指标（metric）：指标名，标签组；
    * 时间戳（timestamp）
    * 样本值（value）
* 指标名反映了监控样本的含义；标签组反映了样本的特征维度，可用于过滤或聚合。*指标名其实也是标签，key 为 `__name__`*

## 指标类型
Prometheus 为实例应用定义了四种不同的指标类型：Counter、Gauge、Summary 和 Histogram.

Counter 只增不减，Gauge 可加可减；Summary 和 Histogram 是组合指标，用于统计和分析样本的分布情况。

> 与Summary类型的指标相似之处在于Histogram类型的样本同样会反应当前指标的记录的总数（以`_count`作为后缀）以及其值的总量（以`_sum`作为后缀）。不同在于Histogram指标直接反应了在不同区间内样本的个数，区间通过标签len进行定义。  
> 同时对于Histogram的指标，我们还可以通过histogram_quantile()函数计算出其值的分位数。不同在于Histogram通过histogram_quantile函数是在服务器端计算的分位数。 而Sumamry的分位数则是直接在客户端计算完成。因此对于分位数的计算而言，Summary在通过PromQL进行查询时有更好的性能表现，而Histogram则会消耗更多的资源。反之对于客户端而言Histogram消耗的资源更少。在选择这两种方式时用户应该按照自己的实际场景进行选择。

## 查询
### 查询方法
查询时，首先要提供指标名，然后根据 label 来筛选指标。PromQL 提供四种 label 筛选的方式：`k=v`, `k!=v`, `k=~regex`, `k!~regex`

查询时间范围：使用方括号将时间段括住，可以查询最近一段时间内的数据，如 `http_requrest_total{}[5m]`  
时间偏移：使用 `offset` 将查询时间段向前推移如 `http_requrest_total{}[5m] offset 1h`  

聚合筛选操作：用 `by` 可以对某些标签进行聚合，用 `without` 可以将某些标签排除，对剩下的标签进行聚合。如 `sum(up) by (instance)`

### 数据结构表示
每个实例的 `metrics` API 产出的数据结构只有一个由多个数据指标构成的瞬时向量；这里的数据结构为输出数据结构，用于查询时进行运算及输出。

* 标量（Scalar）：一个浮点型数据值
* 字符串（String）
* 时间序列（Time Series）：一个序列，其中的每个点都由一个时间戳和一个数据值构成
* 向量：由多个指标和其对应的时间序列构成
* 瞬时向量（Instant Vector）：一个向量，其中每个指标只包含一个点
* 区间向量：使用区间向量表达式 `[5m]` 得到的结果，每个指标包含多个点；用于表示一段时间内的变化值

栗子：

```
# 瞬时向量
Query: checkin_rpc_request_time_count

checkin_rpc_request_time_count{endpoint="CreateFCL"}	1
checkin_rpc_request_time_count{endpoint="CreateUCL"}	3

# 区间向量
Query: checkin_rpc_request_time_count[2m]

checkin_rpc_request_time_count{endpoint="CreateFCL"}	1 @1541421409.712
                                                        1 @1541421469.712
checkin_rpc_request_time_count{endpoint="CreateUCL"}	3 @1541421437.904
                                                        3 @1541421497.904
```


### 运算
#### 数学运算
PromQL 支持的数学运算符有 `+`, `-`, `*`, `/`, `%`, `^`（幂运算）

当两个向量相加时，对应时间点的值将会相加并返回；若某点在另一向量中不存在，则该点会被丢弃；  
当一个向量加一个标量时，会将向量中每个点的值加上这个标量并返回。

#### 条件运算
PromQL 支持的条件运算符有 `==`, `!=`, `>`, `<`, `>=`, `<=`

`a > 2`，会筛选出所有值大于 2 的点，不满足条件的点将会被**丢弃**；  
`a > bool 2`，会对所有点进行逻辑运算，得出结果 **0 或 1** 并返回，不会将点丢弃。

#### 集合运算
PromQL 支持瞬时向量间的集合运算， `and`, `or`, `unless` 分别对应交集、并集和差集。

#### 向量匹配模式
匹配模式用于在向量运算时，对左边的指标名进行标签匹配。  
PromQL 有两种典型的匹配模式：一对一和一对多  
* 一对一用于两边表达式的标签一致的情况；  
    若不一致，可以使用 `on` 或 `ignoreing` 筛选出一些标签再进行匹配
* 一对多 / 多对一用于两边向量长度不一致的情况：  
    修正标签后，可能一边的向量中只有一个值，而另一边有多个；  
    此时，可以使用 `group_left` 或 `group_right` 来进行匹配，其中 `group_left` 表示左边的向量基数更大，也就是多对一；反之亦然。

#### 聚合操作
PromQL 可以对瞬时向量中的多个指标进行聚合，生成另外一个瞬时向量。  
支持的聚合函数太多，见[文档](https://prometheus.io/docs/prometheus/latest/querying/operators/#aggregation-operators)。  
在聚合时，可以通过 `by`/`without` 语句选择根据/不根据那些标签进行聚合。有点像 SQL 里的 `GROUP BY`.  

栗子：

```
Query: request_time_count
request_time_count{endpoint='/', instance='1'} 1
request_time_count{endpoint='/', instance='2'} 1
request_time_count{endpoint='/', instance='3'} 1
request_time_count{endpoint='/index.html', instance='1'} 2
request_time_count{endpoint='/index.html', instance='2'} 3
request_time_count{endpoint='/index.html', instance='3'} 4

Query: sum(request_time_count)
sum(request_time_count) 12

Query: sum(request_time_count) by (instance)
request_time_count{instance='1'} 3
request_time_count{instance='2'} 4
request_time_count{instance='3'} 5
```

#### 内置函数
更多了，见[文档](https://prometheus.io/docs/prometheus/latest/querying/functions/)。

`rate` 和 `irate` 的区别：两个函数都能够表示区间向量中各个时间的变化率，不过 `irate` 比 `rate` 具有更高的灵敏度。  
所以，观察指标时可以用 `irate`，但设置报警规则时应该用 `rate`，以免一些瞬时变化产生误报。

## 监控指标的设计实践
Prometheus 鼓励大家设计多层指标，从多个维度监控到所有东西。  
书中列举的一些常用的监控维度。

|   级别               | 监控什么                                              |    Exporter                     | 
|--------             |---------                                             |                       ----------|
|   网络               | 网络协议：http、dns、tcp、icmp；网络硬件：路由器，交换机等  | BlockBox Exporter;SNMP Exporter |
|   主机               | 资源用量                                              |     node exporter               |
|   容器               | 资源用量                                              |     cAdvisor                    |
|   应用(包括Library)   |  延迟，错误，QPS，内部状态等                             |     代码中集成Prmometheus Client  |
|   中间件状态          |  资源用量，以及服务状态                                 |     代码中集成Prmometheus Client  |
|   编排工具           |  集群资源用量，调度等                                    |     Kubernetes Components       |

除了监控的不同级别以外，还有一些常见的监控模式，用于从不同角度检测应用状态：
* 四个黄金指标：延迟、通讯量、错误发生速率、饱和度
* RED 方法：请求速率、请求错误（每秒失败的请求数）、请求耗时
* USE 方法（用于甄别系统性能问题）：使用率、饱和度、错误计数


# 第 3 章：Prometheus 告警处理

## 基本概念
Prometheus Server 负责存储告警规则并产生告警，而告警的具体处理方式（如发邮件通知用户）则交由 Alertmanager 来做。

> Alertmanage 作为一个独立的组件，负责**接收并处理**来自 Prometheus Server(也可以是其它的客户端程序)的告警信息。Alertmanager 可以对这些告警信息进行进一步的处理，比如**消除重复的告警信息**，**对告警信息进行分组并且路由到正确的接受方**。Prometheus 内置了对邮件，Slack 等通知方式的支持，同时还支持与 Webhook 的通知集成，以支持更多的可能性，例如可以通过 Webhook 与钉钉或者企业微信进行集成。同时 AlertManager 还提供了**静默和告警抑制机制**来对告警通知行为进行优化。

## 告警规则定义
在 Prometheus 配置文件中通过 `rule_files` 指定一组告警规则文件的访问路径。  
设置规则后，Prometheus 会根据 `global.evaluation_interval` 定义的时间周期计算报警规则中定义的 PromQL 表达式。如果表达式能够找到匹配的时间序列，则会对每条序列产生一个告警实例。  
Prometheus 告警配置规则如下：

```yaml
groups:
- name: example
  rules:
  - alert: HighErrorRate
    expr: job:request_latency_seconds:mean5m{job="myjob"} > 0.5
    for: 1m
    labels:
      severity: page
    annotations:
      summary: High request latency for {{ $labels.job }}：{{ $value }}
      description: description info
```

其中，`annotations` 定义的标注部分内容可以进行模板化：`{{ $labels.labelname }}` 可以访问告警实例中标签的值，而 `{{ $value }}` 可以访问表达式算出的样本值。

假设 Prom Server 每 10s 计算一次，则当第一个满足条件的报警出现时，Prom Server 会产生一条报警，但会处于 `PENDING` 状态，而不会触发；在经过由 `for` 指定的一段时间后，如果报警条件依然满足，则报警状态变为 `FIRING`，触发报警。

## Alertmanager 使用方法

### 部署
和 Prometheus Server 一样，Alertmanager 也是单文件可执行程序，也有官方提供的 Docker 镜像。  
Alertmanager 的配置主要包含两部分：路由（route）和接收器（Receiver）。所有的告警信息都会从配置中的顶级路由进入路由树，根据路由规则将告警信息发送给相应的接收器。  
配置成功后，在 Prom Server 的配置中添加内容，将 Prom Server 和 Alertmanager 关联起来。

```yaml
alerting:
  alertmanagers:
    - static_configs:
        targets: ['localhost:9093']
```

### 配置路由
Alertmanager 的路由为树状结构。  
在匹配规则时，可以通警告名称（`alertname`）和报警标签（`labelname`）来进行匹配，匹配方式支持完全匹配和正则匹配。  
如果 `continue` 的值为 `true`，则会将报警交由下一层节点，否则将会直接在当前节点进行处理。  
路由配置格式如下：

```yaml
[ receiver: <string> ]
[ continue: <boolean> | default = false ]

match:
  [ <labelname>: <labelvalue>, ... ]

match_re:
  [ <labelname>: <regex>, ... ]

# 分组规则，如果多条报警的的某些标签值相等，则这些报警将会被归为一组
[ group_by: '[' <labelname>, ... ']' ]
[ group_wait: <duration> | default = 30s ]
[ group_interval: <duration> | default = 5m ]
[ repeat_interval: <duration> | default = 4h ]

routes:       # 路由树的子节点
  [ - <route> ... ]
```

### 配置接收器
Alertmanager 默认支持多种接收器，如邮件，Slack 等，具体看[文档](https://prometheus.io/docs/alerting/configuration/)，同时也有多种基于 Webhook 的继承方式（如 Telegram Bot 和钉钉机器人），具体看[文档](https://prometheus.io/docs/operating/integrations/#alertmanager-webhook-receiver)，如果想自己制作 Webhook 接收器，只需要了解 [Webhook 格式](https://prometheus.io/docs/alerting/configuration/#webhook_config)即可。。

配置报警具体消息时，可以使用 Go 模板来进行自定义。

### 配置抑制 / 静默机制
抑制机制的作用：用户在收到一条报警通知后，可以通过规则来屏蔽掉后续的其它报警，以免收到过多垃圾信息，影响问题分析。  
栗子：如果集群炸了，那我们只需要收到“集群爆炸了”的报警消息就足够了，后面的一万个“A 服务不可用”“B 服务响应延迟 xxx 秒”的报警都不需要再报出来了。  
抑制规则是长期规则，需要在 Alertmanager 配置文件中进行配置。

如果只想临时关掉某些报警，管理员可以通过 Alertmanager 的 UI 来临时屏蔽满足规则的报警通知。  
Alertmanager 的 UI 中可以创建静默规则，管理员可以配置报警匹配规则以及静默时间。

## 杂项
如果某些语句 / 报警规则的计算成本很高，那现场计算可能会导致 Prometheus 响应超时。此时可以通过配置 [Recording Rules](https://prometheus.io/docs/practices/rules/) 来在后台提前对结果进行运算。  
这个功能有点像数据库索引… 

# 第 4 章：使用 Exporter
## 基本概念
广义上讲，所有可以向 Prometheus 提供监控样本数据的程序都可以称作做 Exporter. 一个运行 Exporter 的实例称为一个 Target.  
Prometheus 社区提供了非常丰富的 Exporter 实现，涵盖了从基础设施，中间件以及网络等各个方面的监控功能。

Exporter 应通过 HTTP API 为 Prometheus Server 提供符合格式规范的内容。如：

```
# HELP HTTP Response amount
# TYPE http_res_amount counter
http_res_amount{code="200"} 123
http_res_amount{code="400"} 12

# HELP Example Histogram
# TYPE example_histogram histogram
example_histogram_bucket{le="0.1"}  1
example_histogram_bucket{le="+inf"}  2
example_histogram_sum  5
example_histogram_count  3

# HELP Example Summary
# TYPE example_sumary sumary
example_sumary{quantile="0.1"}  1
example_sumary{quantile="0.5"}  2
example_sumary_sum  5
example_sumary_count  3
```

## Exporter 举例
书中详细提到了三个 Exporter，分别是监控容器状态的 `cAdvidor`、监控 MySQL 运行状态的 `MySQLD Exporter` 和监控网络状态的 `Blackbox Exporter`.

这里不详细展开讲，可以直接看[文档](https://prometheus.io/docs/instrumenting/exporters/)~~看花眼~~。

## 集成 Exporter
除了单独运行的 Exporter 程序外，我们还可以在自己的应用中集成一个 Exporter.  
Prometheus 官方和社区为多种语言提供了[集成支持](https://prometheus.io/docs/instrumenting/clientlibs/)，使用这些库，则可以方便地在应用中提供 HTTP API，并在程序运行过程中进行指标收集。

# 第 5 章：可视化一切
Prometheus 提供了一个 Console Template，可以通过 Go 模板来配置任意控制台界面。  
不过，这种配置方法非常麻烦，所以我们更推荐使用 Grafana 来创建美观的 Dashboard.

## Grafana 设置
基本概念：
* 数据源（Data Source）
* 仪表盘（Dashboard）
* 行（Row）
* 面板（Panel）

### 创建 Dashboard
Grafana 支持多种 Panel，我常用的有 Singlestat 和 Graph.

Singlestat 用于展示当前状态，如 CPU 使用率等，而 Graph 用于展示指标随时间的变化。  
Graph 创建时，如果设置的查询条件返回了多个指标，则可以画出多个图表，这一点非常棒。

# 第 6 章：集群与高可用

# 第 7 章：Prometheus 服务发现

# 第 8 章：监控 Kubernetes
