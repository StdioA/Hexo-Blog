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
