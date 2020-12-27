---
title: protobuf 升级后带来的一些坑
author: David Dai
tags:
  - Golang
  - Protobuf
categories:
  - Golang
date: 2020-12-27 14:41:53
toc: true
---
前段时间把公司某项目依赖的 `github.com/golang/protobuf` 的版本从 v1.3.3 升级到了 v1.4.2，本文记录了升级过程中遇到的一些问题。

<!--more-->

Google 对 Go 的 protobuf 库的底层进行了大的改进，新版本的包路径转移到了 `google.golang.org/protobuf`.  
同时，这些改进也被带进了 `github.com/golang/protobuf`：从 `v1.4` 版本起，`github.com/golang/protobuf` 会在 `google.golang.org/protobuf` 的基础上实现，但会保证接口兼容，这也表明当前依赖 `github.com/golang/protobuf` 的项目可以直接升级版本，而无需对上层代码进行改动。

然而，新版的 `protobuf-gen-go` 使用了 `google.golang.org/protobuf/protoreflect`，编译出的 message 结构体与之前完全不同，这给我们的升级工作带来了一些麻烦。

## 1. 代码中对 `XXX_Unmarshal` 的直接调用
老版的 `protoc-gen-go` 会暴露一个 `XXX_Unmarshal` 接口，用于在 `proto.Unmarshal` 时进行调用，所以有一些同事选择会直接调用 `message.XXX_Unmarshal` 方法。新版的 proto 通过 `ProtoReflect` 接口暴露 message 内部信息，编译 `pb.go` 时也没有了 `XXX_Unmarshal` 方法，所以会导致编译时报错 `message.XXX_Unmarshal undefined`.

解决方案很简单：改用 `proto.Unmarshal` 即可。

## 2. 结构体内部结构变化导致测试出错
针对同一个 message，老版本编译后的结构体结构如下:

```golang
type HealthCheckResponse struct {
    Status               HealthCheckResponse_ServingStatus `protobuf:"varint,1,opt,name=status,proto3,enum=liulishuo.common.health.v1.HealthCheckResponse_ServingStatus" json:"status,omitempty"`
    XXX_NoUnkeyedLiteral struct{}                          `json:"-"`
    XXX_unrecognized     []byte                            `json:"-"`
    XXX_sizecache        int32                             `json:"-"`
}
```

而新版本编译后的结构如下：
```golang
type HealthCheckResponse struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Status HealthCheckResponse_ServingStatus `protobuf:"varint,1,opt,name=status,proto3,enum=liulishuo.common.health.v1.HealthCheckResponse_ServingStatus" json:"status,omitempty"`
}
```
可以看到，新版本中添加了三个未导出字段，而这三个字段为我们的测试代码带来了一些麻烦。

1. cmp.Equal 时 panic
我们的测试中使用了 `github.com/google/go-cmp/cmp.Equal` 来对 proto 结构体进行比较，而结构体中的未导出字段会让 `cmp.Equal` 和 `cmp.Diff` panic: 

```
panic: cannot handle unexported field at {*pkg.SomeRequest}.state:
    ".../services_go_proto".SomeRequest
consider using a custom Comparer; if you control the implementation of type, you can also consider using an Exporter, AllowUnexported, or cmpopts.IgnoreUnexported [recovered]
```

`go-cmp` 推荐的方式是使用 `IgnoreUnexported`，但这种方式需要传递所有需要忽略的类型，对含有多层嵌套的 message 非常不友好。  
经过一番搜索，发现 [`protocmp.Transform`](https://pkg.go.dev/google.golang.org/protobuf/testing/protocmp#Transform) 可以将所有的 protobuf message 转换成自定义的 `map[string]interface{}` 进行比较，所以我们可以用 `Transform()` 来解决问题：

```golang
import "google.golang.org/protobuf/testing/protocmp"

// ...
opt := protocomp.Trnasform()
if !cmp.Equal(exp, got, opt) {
    t.Error(exp, got, opt)
}
```

2. assert 卡死并占满内存
相比上面的问题，下面的问题更加奇怪：使用 `github.com/stretchr/testify/assert.Equal` 比较某些特殊的 proto message 时会卡死，同时内存占用会暴涨。

尝试用 pprof 取样，取出来的 CPU 和堆采样图长这样：

![cpu pprof](/pics/protobuf/cpu-pprof.png)

![memory pprof](/pics/protobuf/memory-pprof.png)

可以看到 `spew.Dump` 中存在无限递归，这导致了程序卡死以及持续的内存分配。

随后搜到了[testify 的 issue](https://github.com/stretchr/testify/issues/930)，相关评论中提出了几种绕过的方案，然而这个问题至今没有解决。

个人推荐的解决方式有两种：
1. 使用 `marshalTextString()` 将 message 转换成 proto text，然后再进行比较；
2. 使用 `cmp.Equal`，结合 `protocmp.Transform`.

## 3. lint 报错 copylocks
处理完业务代码处理测试，处理完测试代码还有 lint 要处理。

我们的项目在升级完后，go vet 会报 copylocks 错误：`assignment copies lock value to got: .../message_go_protos.Message contains google.golang.org/protobuf/internal/impl.MessageState contains sync.Mutex`

解决方式也比较简单：将所有 proto message 改为指针传递即可。
