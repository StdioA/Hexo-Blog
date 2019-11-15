---
title: Golang 学习记录
author: David Dai
tags:
  - Golang
categories:
  - Golang
date: 2019-06-19 15:27
toc: true
---

这几个月在考虑从 Python 转向 Golang，所以专门学习了 Golang.  
这里是 Golang 学习的一些记录。

<!--more-->
学习笔记稍后再整理（咕），先列一下我这几个月看过的各种教程吧。

## 阅读列表
1. [Go by Example](https://gobyexample.xgwang.me/)
    > _Go by Example_ 是对 Go 基于实践的介绍，包含一系列带有标注说明的示例程序。

    真·快速上手必备。

2. [《Go语言四十二章经》](https://github.com/ffhelicopter/Go42/)
    > 《Go语言四十二章经》详细讲述了Go语言规范与语法细节以及在开发中常见的误区；通过对标准库包和著名第三方包的实际运用，来启发读者深刻理解Go语言的核心思维，仔细琢磨经典代码设计模式，引领读者进入Go语言开发的更高阶段。

    讲解详细、信息量超大的 Go 语言教程。

3. [Go Tour](https://tour.go-zh.org/)
    经典的 Golang 官方教程。

    上大学的时候啃过英文的 Go Tour，但是没啃完。  
    看完《Go 语言四十二章经》以后，我只用了一个半小时就把整套教程刷完了。

4. [Go-Mega](https://go-mega.bonfy.im/)
    作者模仿 The Flask Mega-Tutorial 写的 Go 语言 MVC 开发教程，使用裸 `http` 包来进行 Web 开发。

5. [《Go 语言圣经》](https://github.com/gopl-zh/gopl-zh.github.com)
    也是经典教程了。  
    不知道为什么，书的内容让我感觉稍微有点点晦涩。

---
大概这个时候被推荐了一个叫 [Exercism](https://exercism.io/tracks/go) 的网站，上面有一些有趣（或鬼畜）的题目可以用来练习编码熟练度。  
两天完成了大概三分之一的题目，然后就搁置在一边了。  
代码在[这里](https://github.com/StdioA/exercism-go)。

---

6. [《Go 语言标准库》](https://github.com/polaris1119/The-Golang-Standard-Library-by-Example)
    > Golang标准库。对于程序员而言，标准库与语言本身同样重要，它好比一个百宝箱，能为各种常见的任务提供完美的解决方案。以示例驱动的方式讲解Golang的标准库。

    写的很棒，日后也可以做工具书来查询标准库用法。  
    只不过这本书没有写完，有点可惜。

7. [《Go 语言高级编程》](https://github.com/chai2010/advanced-go-programming-book)
    > 《Go语言高级编程》开源图书，涵盖CGO、Go汇编语言、RPC实现、Protobuf插件实现、Web框架实现、分布式系统等高阶主题(完稿)

    我看的时候跳过了 CGO 和汇编的部分。  
    内容…Emmmm…有点杂。比如分布式系统章节的大部分内容并不是在讲 Go，而是在讲后端的解决方案、技术选型，以及各种成熟产品的使用。不过还是值得一看。

---
看完这些以后，我用了两天时间，使用 [gin](https://gin-gonic.com) 重新实现了大三的时候写的迷你博客的后端部分，只不过此时 Go-Mega 的内容已经忘记了不少，所以 MVC 的模式可能实现的不够规范。代码在[这](https://github.com/StdioA/inside-go)。

---

8. [《Go Web 编程》](https://astaxie.gitbooks.io/build-web-application-with-golang/zh/)
    这本书有点老了，内容也比较简单，大概翻了翻，不到俩小时就看完了。

9. [《深入解析 Go》](http://github.com/tiancaiamao/go-internals/)
    深入讲了 Go 语言的内部实现，应该对理解 Go 底层会有很大帮助。

    这本书刚开始看。书中的 Go 语言版本较老，一些底层数据结构的逻辑是用 C 实现的；而 Go 语言从 1.5 版本以后就实现了自举，而我的手上的代码是 Go 1.12.5，所以在翻看代码时可能会发现一些实现上的细小差别。

---
此外，适合入门的书籍还有[《Go 入门指南》](https://github.com/Unknwon/the-way-to-go_ZH_CN)，不过个人更推荐使用《Go语言四十二章经》来入门。

再推荐一位大佬的博客，他对 Go 语言的一些底层细节做了深入的分析和易懂的讲解。<https://halfrost.com/tag/go/>

## 学习笔记
[《深入解析 Go》笔记](/2019/06/go-internal-note/)，书讲的很深入，我记得不够深入。

基础内容的笔记稍后奉上（大概吧）。
