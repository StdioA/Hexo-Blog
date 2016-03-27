title: 点点技能点之Ajax
date: 2015-09-22 21:07:00
categories:
- 网络
tags:
- Ajax
- 前端开发

---

之前自己没怎么写过前端，一直以为Ajax写起来很麻烦，也没怎么接触Ajax，今天写了个小网页用到了Ajax，感觉它并不是想象的那么复杂，写完网页随手总结一下。（喂我这个技能点是不是点的太晚了啊！）

<!-- more -->

# 1. 引子
没啥可写的…

# 2. 科普
[W3School - Ajax教程](http://www.w3school.com.cn/ajax/index.asp)，就这个吧，科普完毕。

# 3. 笔记
## 3.1 Ajax的特点
Ajax为Asynchronous Javascript And XML的缩写，即异步JS与XML，因为它具有异步特性，所以我们可以在网页加载中或加载结束的任意时刻使用，通过与服务器进行少量数据交换来实现网页异步更新。利用Ajax可以直接更改网页的一部分来动态刷新网页，而不需要刷新整个网页。

## 3.2 使用方法——原生方式
```javascript
var ajax = new XMLHttpRequest();                  // IE6不能用这个方式，不管它了
ajax.onreadystatechange = function() {
    if(ajax.readyState==4 && ajax.status==200)
    {
        // do something with ajax.responseText
        // 还有一个ajax.responseXML, 也不管它了，json大法好
        alert(ajax.responseText);
    }
}
ajax.open("GET", "./get?var1=val1&var2=val2", true);     // 异步方式运行
ajax.send();
```
```javascript
var ajax = new XMLHttpRequest();
ajax.open("POST", "./post", false);                      // 同步方式运行
ajax.setRequestHeader("Content-type","application/x-www-form-urlencoded");
ajax.send("var1=val1&var2=val2");
// do something with ajax.responseText
alert(ajax.responseText);
```

## 3.3 原生的好麻烦！——jQuery Ajax
[jQuery](http://jquery.com/)，棒棒的前端库，不多说。
```javascript
// $.get(URL, callback)
$.get("./get?var1=val1&var2=val2", function(data, status) {
        // do something with data & status
    });
// $.post(URL,data,callback)
$.post("./post", {
        var1: "val1",
        var2: "val2"
    }, function(data, status) {
        // do something with data & status
    });
// $.getJSON(url, data, success(data,status,xhr))
// data和回调函数success可选
// success中，data参数必需，其它可选
$.getJSON("./get", function(data) {
        // do something with data
    });
```

# 4. 附带的小零碎
`window.setInterval(getData, time);`用于实现轮询。

# 5. 总结 & 后记
Ajax简单实用功能强大应用广泛，通过简单的步骤实现前后端的数据交换，来实现网页的动态加载。  
开学到现在感觉一直不在状态（实在是想吐槽自己的一群神一样的舍友），今天找了个安静的、没人打扰的地方认认真真写了写代码，总结了一下自己到现在所学的知识才发现自己学的东西是如此零散、不成体系，仿佛从大一开始就在乱点技能点。这一年要多深入学习一个特定方面的知识（暂定后端开发），多写代码积累经验，不能再东学一点西学一点了。
