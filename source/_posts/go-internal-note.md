---
title: 《深入解析Go》笔记
author: David Dai
tags:
  - Golang
categories:
  - Golang
date: 2019-06-24 16:22
toc: true
---

在 GitHub 上找到一本解读 Go 实现细节的好书，名叫[《深入解析 Go》](https://github.com/tiancaiamao/go-internals/)。  
大致看了一遍，简单做了些笔记。

<!--more-->

这本书的代码来自 Go 1.3，所以还有一部分由 C 语言写成。  
这份笔记里的代码来自 Go 1.12.5，数据结构全部由 Go 语言实现。

## 数据结构
string 和 slice 都是引用类型，可能开在栈上，也可能开在堆上；  
channel 和 map 是引用类型，但一定开在堆上，栈中只有指针。

### string 底层结构
`src/reflect/value.go`
```go
type StringHeader struct {
    Data uintptr
    Len  int
}
```
string 是不可变数据结构，任何对 string 的操作都会产生一个新的 string. 因此，需要拼接的时候尽量使用 `bytes.Buffer`，或 `strings.Join` 这种用了 Buffer 的函数。

### slice 底层结构
`src/reflect/value.go`
```go
type SliceHeader struct {
    Data uintptr
    Len  int
    Cap  int
}
```

注意，两个数据结构的底层数据全部共享。  
关于 slice 的扩容，可以参见[《深入解析 Go 中 Slice 底层实现》](https://halfrost.com/go_slice/)。  
有一个小用法：对切片做 slice 操作时，末位的下标可以超过原 slice 的 len，但不能超过 cap.

```go
a := make([]int, 2, 10)
b := a[:5]
```

### map 底层结构
`src/runtime/map.go`
```go
type hmap struct {
    // Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
    // Make sure this stays in sync with the compiler's definition.
    count     int // # live cells == size of map.  Must be first (used by len() builtin)
    flags     uint8
    B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
    noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
    hash0     uint32 // hash seed

    buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
    oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
    nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

    extra *mapextra // optional fields
}

type bmap struct {
    // tophash generally contains the top byte of the hash value
    // for each key in this bucket. If tophash[0] < minTopHash,
    // tophash[0] is a bucket evacuation state instead.
    // 用于 hash 的快速比较，当 hash 相等的时候还会跟原 key 进行一次匹配。
    tophash [bucketCnt]uint8
    // Followed by bucketCnt keys and then bucketCnt values.
    // NOTE: packing all the keys together and then all the values together makes the
    // code a bit more complicated than alternating key/value/key/value/... but it allows
    // us to eliminate padding which would be needed for, e.g., map[int64]int8.
    // Followed by an overflow pointer.
}
```
注意看上面的 NOTE 里的 KV 排列结构，这样做有利于内存对齐。  
`bmap` 后面的内存分配未在结构体中定义，需要拿 KV 的时候要通过 `unsafe.Pointer` 根据 offset 去拿。  
见 `mapaccess1` 或 `mapaccess2`.

哈希表在每次扩容时，容量会增大到原来的两倍，也就是从 `2^B`  到 `2^(B+1)`。为了保证运行效率，会将 key 的搬迁操作平摊到每一次写操作（insert & remove）上，每次操作时迁移 1-2 个键值对。查找时会先在 `old—buckets` 中找，找不到再去新的 `buckets` 中找。

map 使用快速的 murmurhash 作为哈希算法。

`map` 解决冲突的方法是链地址法的改进形式。  
在创建 `bmap` 时，会分配一个数组，都可以容纳 8 个键值对。发生哈希冲突时，都会将新的键值对放入添加到数组里。如果数组放满了，则会新建一个 `bmap`，通过 `overflow` 指针来链接到当前 bucket 节点后面。

关于查找、插入和删除的细节，请见[该书相关章节](https://tiancaiamao.gitbooks.io/go-internals/content/zh/02.3.html)。  
一个需要注意的点：在 bmap 链中，相同的 Key 可能会存在于两个 bucket 里，而前面 bucket 的值会直接覆盖后面 bucket 的值。在进行更新操作时，如果前面 bucket 不存在该键，但是数组包含空位，则直接在该 bucket 中插入，而不会再去链表后面的 bucket 中查找并更新。查找操作亦然。*怎么说，这个机制有点像 docker 镜像中的不同层的文件覆盖机制。*

### channel
`src/runtime/chan.go`
```go
type hchan struct {
    qcount   uint           // 环形队列的总数量
    dataqsiz uint           // 环形队列大小
    buf      unsafe.Pointer // 环形队列
    elemsize uint16
    closed   uint32
    elemtype *_type // 元素类型
    sendx    uint   // 缓存区的发送指针（也就是缓冲区尾）
    recvx    uint   // 缓冲区的接收索引（也就是缓冲区头）
    recvq    waitq  // 因接收而阻塞的等待队列
    sendq    waitq  // 因发送而阻塞的等待队列

    // lock protects all fields in hchan, as well as several
    // fields in sudogs blocked on this channel.
    //
    // Do not change another G's status while holding this lock
    // (in particular, do not ready a G), as this can deadlock
    // with stack shrinking.
    lock mutex
}
type waitq struct {
    first *sudog
    last  *sudog
}
type sudog struct {
    // The following fields are protected by the hchan.lock of the
    // channel this sudog is blocking on. shrinkstack depends on
    // this for sudogs involved in channel ops.

    g *g

    // isSelect indicates g is participating in a select, so
    // g.selectDone must be CAS'd to win the wake-up race.
    isSelect bool
    next     *sudog
    prev     *sudog
    elem     unsafe.Pointer // data element (may point to stack)
                            // 存储 goroutine 的

    // The following fields are never accessed concurrently.
    // For channels, waitlink is only accessed by g.
    // For semaphores, all fields (including the ones above)
    // are only accessed when holding a semaRoot lock.

    acquiretime int64
    releasetime int64
    ticket      uint32
    parent      *sudog // semaRoot binary tree
    waitlink    *sudog // g.waiting list or semaRoot
    waittail    *sudog // semaRoot
    c           *hchan // channel
}
```

以前 channel 的缓冲区会在 hchan 之后的内存中创建，现在将它们分开了。  
对 channel 进行读写操作的时候，会根据缓冲区的状态、以及读 / 写链表来决定是否要阻塞当前 goroutine，阻塞的话会将当前的 G 挂载链表中。

### interface
```go
// 带方法的接口
type iface struct {
    tab  *itab
    data unsafe.Pointer
}

type itab struct {
    inter *interfacetype
    _type *_type
    hash  uint32 // copy of _type.hash. Used for type switches.
    _     [4]byte
    fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}
// fun 中会保存接口值背后具体类型的方法

// 空接口
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
```

`Type` 的 `UncommonType` 会记录某个**具体类型**实现的方法；  
`interface` 的 `Itab` 的 `InterfaceType` 中的方法表，记录了接口所声明的方法；    
`itab` 的 `fun` 数组，记录了具体的函数指针，用作接口值的方法缓存。

在接口值赋值时，会将 `UncommonType` 和 `InterfaceType` 中的方法表进行比对；如果比对成功的话，会将 `UncommonType` 中的方法指针拷贝到 `itab` 的 `fun` 中，方便方法调用时对方法进行查找。

### 方法调用
`a.F(b)` 会在编译时直接转换为 `A.F(a, b)`；  
当一个类型被匿名嵌入结构体时，它的方法会被拷贝到嵌入结构体的 Type 的方法表中。这个过程在编译时就可以完成。

接口的方法调用要通过 `itab.fun` 中的函数指针来确定具体调用的方法。方法拷贝的过程是在运行时完成的，所以接口的方法调用的成本要略高于一般方法调用。


## 函数调用协议
Go 把返回值放在上一个栈帧最后的内存中，这样调用链前后的两个函数都可以触及到这段内存，以此来实现多值返回。  
```
返回值2
返回值1
参数3
参数2
参数1 <- SP
```

go 关键字是个语法糖，`go f(args)` 可以看做 `runtine.newproc(size, f, args)`；  
同样 defer 也是语法糖，通过 `runtine.deferproc` 和 `runtine.deferreturn` 来实现。

### 连续栈
一个程序中可能会有非常多的 goroutine，为了节省内存，每个 goroutine 一开始只会得到非常小的一块栈。

使用可变栈时，每次函数调用时，都会通过 SP 和 stackguard 检查栈的使用情况。  
当栈不够大时，会进行栈扩张，开一块新栈，并把旧栈的内容复制过去。如果栈里有指针，而指针指向的是栈中的变量，那么在复制时会对指针的值加一个偏移，来保证指针指向的对象是被迁移过后的对象。  
gc 时，如果检测到栈只用了不到 1/4 时，会将栈缩小为原来的 1/2.

### 闭包
Go 通过 escape analyze 来检查逃逸的值，从而确定该值是否应该在堆上创建，而不是在栈中。

```go
func New() *T {
    var t T
    return &t
}
```

返回闭包时并不是单纯返回一个函数，而是返回了一个结构体，记录下函数返回地址和引用的环境中的变量地址。书中是有这样一个结构体的，但在我的源码里，这个数据结构应该是用汇编实现了。

## Go 程序初始化
* 系统初始化：初始化栈、设置本地线程存储（g）
* 调度器初始化：`runtime.schedinit` 
    根据 `GOMAXPROCS` 决定可用线程数；把 `runtime.main` 放入就绪线程队列；  
    调用 `runtime.mstart`，`mstart` 调用 `schedule`，也就是一直运行的调度器。
* `schedule` 选中 `runtime.main`
    * 通过 `newm(sysmon, nil)` 启动一个线程运行 `sysmon`，用于处理网络 epoll 
    * 通过 `runtime.newproc` 启动一个 G 执行 scavenger，用于垃圾回收
    * 调用 `main.main` 进入用户代码。

会单独启动一个线程（M）用于 `poll()`。GC 都是 goroutine，这个任务的地位是高于 GC 的。

## 调度模型
### Go 如何实现并发？
Go 通过 goroutine 和 channel 来实现 CSP 并发模型，从而实现并发。  
* goroutine 是 Go 语言中并发的执行单位。可以理解为轻量级的“线程”。
* channel 是 Go 语言中各个并发结构体 (goroutine) 之前的通信机制。 通俗的讲，就是各个 goroutine 之间通信的“管道”，有点类似于 Linux 中的管道。

再深入一点，Go 线程模型的实现依靠 MPG 以及 Sched 结构体。  
M 是 Machine，是对机器的抽象，一个 M 会关联一个物理线程；  
P 是 Processor，代表 Go 代码执行是所需要的资源；  
G 是 Goroutine，代表 Goroutine 的控制结构；
Sched 是调度实现中使用的数据结构。

多个 G 会以队列的形式挂靠在 P 上；当 G 与 M 绑定时，才能够执行 Go 代码。  
调度时会采用抢占式调度模型以避免一个 goroutine 运行太长时间；某个 P 的队列变空时会从其它的 P 队列上偷 G，然后继续运行。

http://morsmachine.dk/go-scheduler

![](https://i6448038.github.io/img/csp/GMPrelation.png)

### 系统调用细节
当某个 goroutine 发起一次系统调用时，会调用 `runtime.entersyscall`。  
调度器会将 G 的状态设置为 `Gsyscall` 后放入就绪队列；  
此时，p 会和 m 进行剥离，p 的状态被设为 `Psyscall`，而 m 会去执行系统调用。  
如果系统调用时非阻塞的，那么 m 会很快返回。返回时会调用 `runtime.exitsyscall`，这个时候会去检查当前 m 的 P，如果 P 处于 `Psyscall` 且队列非空，则重新将 p 和 m 绑定，恢复 g 的状态为 `Grunning`，继续运行。

如果 goroutine 发起的是阻塞的系统调用，则会调用 `runtime.entersyscallblock`。  
与 `entersyscall` 不同的是，`entersyscallblock` 会调用 `releasep` 和 `handoffp`。  
`releasep` 将当前的 M 与 M 关联的 P 剥离，M 会负责去执行系统调用；  
执行 `handoffp` 会让 P 尝试挂靠到其它空闲的 M 上继续执行。  
如果 P 上没有 G 了，P 会被设置为 `Pidle`；如果没有空闲的 M 了，则会调用 `startm` 来让 P 与新的 M 绑定后继续执行。  
当系统调用完成时，要让发起调用的 G 来继续执行。这时 G 会去找可用的 P。如果当前不存在 `Pidle` 的 P，调度器将会把 G 变成 `Grunnable`，将它挂到全局的就绪 G 队列中，然后停止当前 m 并调用 `schedule` 函数。  
*换句话说，block 的时候，可能 m 上的 P 会空闲，等 m 返回后还可以继续挂载执行，此时 M-P-G 的绑定原样恢复。整个逻辑就有点像非阻塞调用。如果 m 上的 P 去干别的了（比如又找了个新的 M 继续执行），那么当前 m 会将信息传递给 G，改变 G 的状态，然后 m 自己退出（因为所有的 P 都有 M 了）*

## 内存管理
### 内存池
每个线程都会有自己的本地内存，当线程内存不够时会向全局分配链中申请内存。

Go 会为每个 M 在 `MCache` 中存储一些空闲的小内存块；作为备用分配存储。当 `MCache` 用完后，会从 `MCentral` 自由链拿**一些**对象进行补充；`MCentral` 为空时，又会从 `Mheap` 中拿一些对象进行补充。这样的多级批量补充机制减小了全局内存加锁的开销。

* 当程序需要小对象（小于 32K）时，会直接从 `MCache` 中分配，对象被回收后控制内存返回给全局控制堆（`MHeap`）；
    从 `MCache` 中分配避免了在全局控制堆上频繁加锁。
* 当需要大对象时，会直接从全局控制堆上以页（4KB）为单位进行分配。因此大对象总是页对齐的。

### 垃圾回收
Go 语言使用标记清除算法来完成垃圾回收，整个回收过程会 stop the world.  
标记阶段从 root 区域出发，扫描所有直接或间接引用的对象；清除阶段直接扫描堆区，对未被标记的对象进行回收。  
由于标记过程是一个树形的操作，所以这个过程被并行化，以提升速度。

## 网络
Go 通过运行时层面对 epoll/kqueue 的封装来实现非阻塞 io.

封装层次：
* 平台相关的 API 封装
* 平台独立的 runtime 封装
* 用户级别的库封装（如 `net`）

在 `runtime.main` 启动时，会运行 `newm(sysmon, nil)`，而 sysmon 就会每隔一段时间执行 `runtime.epoll`。sysmon 的地位要比 gc 重要的多，而且会频繁执行，所以会单独为它分配一个系统线程（m）来运行。
