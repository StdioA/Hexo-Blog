title: Atom 浅度体验感受
date: 2016-03-27 10:40:03
categories:
- 随笔
tags:
- Atom
- Sublime Text
- 文本编辑器
- Git
toc: true

---

昨天看到了 Atom 官方发布了一篇 [Atom Flight Manual](http://flight-manual.atom.io/), 看了两页觉得有些感兴趣所以下了一个体验一下。

<!-- more -->

很久以前就听说了 `Atom` 这款编辑器，Github 宣称它是"A hackable text editor
for the 21st Century", 有些感兴趣于是就搞了一个内测码下载下来来玩了玩，但是，当时 Atom 的用户体验并不好，加上可以使用的插件少之又少，就放弃了它。昨天从微博上看到了官方的 [Atom Flight Manual](http://flight-manual.atom.io/)，决定重新体验一下这款所谓“属于 21 世纪的编辑器”。

当然，作为一个使用 Sublime Text 将近 2 年的用户，我肯定会将 Atom 跟 Sublime Text（下称 ST）进行比较咯。

# 1. 第一印象
一打开Atom界面，整个 UI 还是浓浓的仿 ST 风格，快捷键都基本一样，可能是因为 ST 的 UI 确实很简洁漂亮吧。  
不过 ST 是收费的（无期限评估使用也是收费），而 Atom 是免费的，这点要赞一个~当然这根本无法打消暑假购买 ST 授权码的决心。

# 2. 启动速度
Atom 的启动速度显然跟 ST 没法比。实测在装了 15 个插件以后，Atom 的启动速度要比 ST3 慢 5 倍左右。当然，它的启动速度应该会比装了一堆插件的 ST2 快一点…ST2 的启动速度貌似是历史遗留问题了。

# 3. 用户体验
感觉 Atom 的用户体验要比上次好的多。

来说一说我看好 Atom 的地方吧。

1. Web-Based  
    Atom 编辑器是基于 V8，整个 Atom 编辑页面就是一个网页。不知道为什么，感觉网页对我来说更有亲切感:flushed:，当然这也使得 Atom 更加容易修改，方便了插件的开发。

2. Git Diff & Markdown Preview 集成  
    Atom 自带了两个插件：`Git Diff` 和 `Markdown Preview`.  
    `Git Diff` 可以实时显示当前文件的 Git Diff 信息。嗯那个 Git Diff 还是蛮漂亮的~   
    ![Git Diff](/pics/atom/git_diff.jpg)  
    `Markdown Preview` 可以提供 Github Markdown 实时预览。然而作为一个可以对 Markdown 进行人肉编译的人，根本不需要这样的插件嗯:joy:  
    同时，Atom 还集成了 Git 的常用功能，如当前目录所在 Git 分支。这个功能也要点个赞。  

3. Package/Theme 安装及配置  
    这点上 Atom 做的要比 ST 好得多。Atom 可以在线查看 Package 的 README 信息，每一个 Package 也有其独立的配置页面，不必像 ST 那样直接去修改配置文件。  
    ![Package Settings](/pics/atom/package_settings.jpg)

    当然，Atom 自带了一个叫做 `apm` 的包管理工具，这个就不错评价咯，人家只是浅度使用嘛（傲娇脸

# 4. Hackable 简单体验——编辑器字体修改
之前用 ST 的时候为了将 Consolas 和微软雅黑结合起来，着实头疼了好久，到了 Atom 里，打开 Stylesheet 配置文件，一行代码搞定:satisfied:
```css
atom-text-editor {
  font-family: Consolas, "Microsoft Yahei", sans-serif;
}
```

# 5. 总结
Atom 确实比以前好用了很多。我大概不会卸载它，而是把它留在电脑里，等到 ST 用腻了可以换换口味~  
然而尝试了 Atom 并更换主题、安装了几个实用的插件以后，我还是跑到 Sublime 那里安装了类似的插件:joy:

配置后的 Atom:  
![Atom 个性化](/pics/atom/atom_after.jpg)  
使用 Atom 前的 ST:  
![ST 个性化前](/pics/atom/st_before.jpg)  
使用 Atom 后，安装了插件`Color Highlighter`, `GitGutter` 和 `SublimeGit`的 ST:  
![ST 个性化后](/pics/atom/st_after.jpg)

好吧…我把 ST 配置的像 Atom 了（摊手

完。
