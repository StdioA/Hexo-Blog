title: Javascript学习总结
date: 2016-02-24 19:31:42
categories:
- Javascript
tags:
- Javascript
- node.js
- 前端开发

---

从开始写Javascript到现在，已经有一个月了，这一个月学了不少新姿势，随便写写，简单整理一下。

<!-- more -->

# 1. Javascript

这个没什么好整理的…随便写几条。

```javascript
(function () {
    var a;                          // 防止变量作用域提升
    do_something();
}());
 
some_list.map(function (obj) {
        this.do_something(obj);
    }, this);                       // 把this传到map里面的匿名函数中，否则里面的this为undefined
 
setInterval(function () {}, 1000);
setTimeout(function () {}, 1000);
```

# 2. jQuery

## 2.1 选择器
```javascript
$("input[type=text]")   // 选择属性
$("ul>li:eq(3)")        // 选择第4个li元素
$("ul>li").eq(3)        // 选择所有"ul>li"中的第4个元素（注意与上面那个选择器的不同）
$("ul>li:even")         // 选择所有奇数li元素（odd同理）
```

## 2.2 杂项

获取表单内容时，要用`$("...").val()`而不是`$("...").text()`.

# 3. Node.js

不知道写什么，随便写几个好玩的库：
* `cheerio`, 在服务器端解析html，跟jQuery用法差不多
* `chalk`, 输出彩色文字
* `gulp`, 流式自动化构建工具，后面细写

---
May.2 2017 更新
## 3.1 Node.js 文档阅读笔记
### 3.1.1 console.timer
计时工具。

```javascript
console.time('100-elements');
for (let i = 0; i < 100; i++);
console.timeEnd('100-elements');
// 100-elements: 0.238ms
```

### 3.1.2 Buffer
Buffer 在处理文件或流时可能会用到，在处理文件时，跟 Python 的 bytes 有些相似。

创建 Buffer: `Buffer.alloc(10)` 或 `Buffer.from([1, 2, 3])`；  
长度： `buf.length`；  
切分： `buf.slice([start, [end]])`；  
字符串与 Buffer 转换：
```javascript
> b = Buffer.from("哦")
<Buffer e5 93 a6>
> b[0]
229
> b.toString('utf-8')
'哦'
> b.toString('base64')
'5ZOm'
```

### 3.1.3 child_process
生成子进程，执行文件：
```javascript
const spawn = require('child_process').spawn;
const ls = spawn('ls', ['-lh', '/usr']);
 
ls.stdout.on('data', (data) => {
  console.log(`stdout: ${data}`);
});
 
ls.stderr.on('data', (data) => {
  console.log(`stderr: ${data}`);
});
 
ls.on('close', (code) => {
  console.log(`child process exited with code ${code}`);
});
```

如果进程可以立即结束（比如 ls），或者不需要实时查看 stdout 的输出，可以使用 `child_process.exec`:

```javascript
const exec = require('child_process').exec;
exec('cat *.js bad_file | wc -l', (error, stdout, stderr) => {
  if (error) {
    console.error(`exec error: ${error}`);
    return;
  }
  console.log(`stdout: ${stdout}`);
  console.log(`stderr: ${stderr}`);
});
```

### 3.1.4 Path
用于处理文件路径，跟 os.path 类似。

```javascript
path.join('/foo', 'bar', 'baz/asdf', 'quux', '..')
// Returns: '/foo/bar/baz/asdf'
 
path.normalize('/a/b//../c/./d/')
'/a/c/d/'
 
path.relative('/data/orandea/test/aaa', '/data/orandea/impl/bbb')
// Returns: '../../impl/bbb'
```

# 4. React.js

## 4.1 JSX与Babel
JSX是一种语言，babel是一种预处理工具。
JSX可以在浏览器中转换为javascript并执行，有两个前提：
1. 包含了`browser.js`
2. script标签的类型为`text/babel`

看[这里](https://facebook.github.io/react/docs/getting-started.html)。

在很久以前，JSX代码转换库包含在React库里，名叫`JSXTransformer.js`, 那时script标签的类型为`text/jsx`, 但是后来JSX代码改用babel来转换了，所以script标签的类型也就变为了`text/babel`, 代码转换库也不再由React提供。

吐槽：最初看tutorial的时候看的[中文页面](http://reactjs.cn/react/docs/getting-started.html)，这个页面很久很久没有更新过了，很多内容都过时了，在离线转换的时候需要安装包，npm报出一堆`package deprecated`的信息，整个人一个大写的卧槽。后来看了上面的英文页面才知道`JSXTransformer`已经不再使用了，现在大家都用[Babel](https://babeljs.io/).

## 4.2 React

1. [阮一峰大大写的React入门教程](http://www.ruanyifeng.com/blog/2015/03/react.html)很不错，我基本上是看这个入门的。

2. 写JSX时，要注意DOM节点的`class`和`for`属性要写为`className`和`htmlFor`，因为`class`和`for`是javascript的关键字。

3. 组件的生命周期：看[官方文档](https://facebook.github.io/react/docs/working-with-the-browser.html#component-lifecycle)，`render`函数中**不要**改变组件的state，若组件的props改变而需要相应更改state, 则要在`componentWillReceiveProps`函数中完成state更改。函数执行完后，render函数会被执行，组件重新渲染。

# 5. gulp

## 5.1 简介
流式自动化构建工具，用于各类文件的转换（如jsx→js→min.js），监控文件变化，搭建静态文件服务器等，可以类比为makefile, 拥有种类繁多的插件。

[gulp官网](http://gulpjs.com/)，[中文网，包含中文文档](http://www.gulpjs.com.cn/)。

## 5.2 操作

* `gulp.task(name[, deps], fn)`, 注册一个任务
    * name: 任务名称
    * deps: 依赖的前置任务，string[]
    * fn: 任务函数，在里面写该任务需要完成的具体事项

* `gulp.src(globs[, options])`, 输出一个满足匹配模式的stream, stream可以用pipe连接起来，类比shell的管道`|`.

* `gulp.dest(path[, options])`, 将stream写到某个path当中。

* `gulp.watch(glob [, opts], tasks)` 或 `gulp.watch(glob [, opts, cb])`, 监控满足匹配模式的文件，若文件变化，则执行某些任务。

直接看一个例子：
```javascript
gulp.task('render', ['array', 'of', 'task', 'names'], function () {
    gulp.src('./client/templates/*.jade')                               // 找到原路径所有的jade文件
        .pipe(jade())                                                   // 渲染模板
        .pipe(gulp.dest('./build/templates'))                           // 输出到某目录
        .pipe(minify())                                                 // minify
        .pipe(gulp.dest('./build/minified_templates'));                 // 输出到另一目录
});
 
gulp.task('watch', ["compress"], function () {
    var watcher = gulp.watch('./public/src/*.jsx', ['compress']);       // 监控文件，若文件变化则执行compress任务
    watcher.on('change', function (event) {                             // 监听change事件
        console.log('File ' + event.path + ' was ' + event.type + ', running tasks...');
    });
});
```

# 6. React工具集成: React+Babel+gulp

用gulp执行任务，用babel转换JSX.

## 6.1 gulp
所需插件:
* gulp
* gulp-babel
* gulp-uglify (压缩js文件用，可选)

gulpfile:
```javascript
var gulp  = require('gulp');
var babel = require('gulp-babel');
var uglify = require('gulp-uglify');

gulp.task('transform', function () {
    return gulp.src('./public/src/*.jsx')
                .pipe(babel())
                .pipe(gulp.dest('./public/build'));
});

gulp.task("compress", ["transform"], function () {
    return gulp.src('./public/build/!(*.min).js')
            .pipe(uglify())
            .pipe(rename({ suffix: ".min" }))
            .pipe(gulp.dest('./public/build'))
});
```

## 6.2 Babel
所需插件：
* babel-preset-es2015
* babel-preset-react

`.babelrc`文件：
```json
{
    "presets": [
        "es2015",
        "react"
    ],
    "plugins": []
}
```

若在全局安装了`babel-cli`，则可以用babel命令转换文件：
1. 若当前目录存在上述babelrc文件：
    执行`babel public\src --out-dir public\build`
2. 若当前目录不存在bebelrc文件：
    执行`babel --presets react public\src --out-dir public\build`

# 7. Semantic UI

超级棒的一个前端组件库，去看[文档](http://www.semantic-ui.cn/)吧。

# 8. 后记

终于写完了。

又一个月没有写东西了，代码写了不少，所学到的知识却没有及时整理下来。真是太懒，懒得整理知识。到现在还欠着一篇Django的生产环境配置的文章没写（其实就是一条命令+nginx配置而已），改天补上。

寒假算是荒废过去了，本来能够写更多东西的，却被一些事情打乱了计划。

以前觉得javascript是一门很糟糕的语言，代码写起来很乱，四五层回调看起来头都大了，但真正写了一个月以后感觉舒服了很多，很多代码写起来得心应手。JS的生态也让我很喜欢，各种工具组件层出不穷，使用起来也及其方便。

这一个月所写的js项目：
* [carrez](https://github.com/StdioA/carrez)
    前后端都有，后端Express, 解析html, 前端使用ajax和后端进行信息交互
* [starwars](https://github.com/StdioA/starwars)
    纯前端项目，使用React, 用了很久的`browser.js`, 本地工具集成鼓捣了很久才鼓捣明白
* [git_modifier](https://github.com/StdioA/git_modifier)
    算是半个JS项目（Github告诉我我这是个JS项目，因为JS代码占比最大），后端Flask, 前端React, 说是Web App其实是个本地项目，用来读取及修改本地git repo的commit信息用的，奇怪的需求，写来自己用，方便伪造commit信息，帮某人作弊😂😂😂。

后面看看要不要写`Mocha`吧。

最后，郑重感谢phoenixe同学给了我一个比较系统地学习和使用javascript的机会。谢谢你。

以上。
