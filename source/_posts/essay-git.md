title: 随手记之Git
date: 2015-10-14 18:28:00
categories:
- 随手记
tags:
- Git
- 编程工具
toc: true

---

从一年前开始使用Git, 一直没系统整理过Git的命令，前几天在部署代码的时候出了问题不知道该如何解决，于是决定整理一份Git的使用方法。本博文会持续更新。

<!-- more -->

# 1. 什么是Git
Linus大神写的分布式版本控制工具，具体请访问官网<http://git-scm.com>.
Wikipedia链接: [Git (software)](https://en.wikipedia.org/wiki/Git_(software))

# 2. 基本操作

初始化版本仓库: `git init`
从远程仓库克隆: `git clone [url] [repo_name]`  
例: `git clone https://github.com/user/repo my_repo`


![life-cycle](/pics/lifecycle.png)
<center><small>文件的状态变化周期</small></center>

检查文件状态: `git status`
状态简览: `git status -s`
```lll
$ git status -s
 M README
MM Rakefile
A  lib/git.rb
M  lib/simplegit.rb
?? LICENSE.txt
```
| 字母 | 所代表的状态 |
| :- | :-: |
| ?? | 未追踪 |
| &nbsp;&nbsp;&nbsp;M | 已修改，未暂存 |
| MM | 修改后暂存，然后又修改 |
| A | 新添加到暂存区 |
| M | 修改后添加到暂存区 |

忽略文件： 编辑`.gitignore`文件。  
文件 .gitignore 的格式规范如下：
* 所有空行或者以 ＃ 开头的行都会被 Git 忽略。
* 可以使用标准的 glob 模式匹配。
* 匹配模式可以以（/）开头防止递归。
* 匹配模式可以以（/）结尾指定目录。
* 要忽略指定模式以外的文件或目录，可以在模式前加上惊叹号（!）取反。
具体例子见[Git-基础](http://git-scm.com/book/zh/v2/Git-%E5%9F%BA%E7%A1%80-%E8%AE%B0%E5%BD%95%E6%AF%8F%E6%AC%A1%E6%9B%B4%E6%96%B0%E5%88%B0%E4%BB%93%E5%BA%93)。

查看尚未暂存的文件更新部分: `git diff`
查看已暂存的文件更新的部分: `git diff --staged`

提交更新: `git commit`
`git commit -a` = `git add --all; git commit`
提交时输入单行信息(Commit log): `git commit -m`


添加文件: `git add [filename]`, 使git跟踪文件
删除文件: `git rm [filename]`, 使git停止跟踪文件并将文件删除
停止跟踪文件: `git rm --cached`, 停止跟踪但不删除文件
移动文件: `git mv`, 规则和`mv`基本相同

查看提交历史: `git log`
显示每次提交的内容差异: `git log -p`
查看每次提交的简略统计信息: `git log --stat`
查看每次提交的代码更改详情：`git log --cc`
显示ASCII图形表示的分支合并历史: `git log --graph`
粗略显示: `git log --oneline --graph --decorate --all`

更多请输入`git log --help`查看man page.

更改上次提交: `git commit --amend`
取消暂存: `git reset HEAD [filename]`
取消暂存并丢弃现有的更改: `git reset HEAD --hard [filename]`, 未提交的更改会丢失
撤销对文件的更改: `git checkout -- [filename]`, 未提交的更改会丢失

`git push origin --delete [branch]`, 删除远程分支

# 3. 分支操作

`git branch [branch]`, 在当前引用上建立新分支
`git checkout -b [branch]`, 建立新分支并检出到该分支上

commit引用方式:  
* 直接hash引用: `d921970aadf03b3cf0e71becdaab3147ba71cdef`, `d921970`
* 分支引用: `HEAD`, `master`
* 相对引用: `HEAD^`, `HEAD^^`, `HEAD~2`
* 若`a`是一个合并提交，有两个父引用，则`a^`为a的第一父引用，`a^2`为a的第二父引用。
* 条件引用: `master@{yesterday}`, `master@{1.week.before}`
* 引用区间: 
    * `refA..refB`选择从refA和refB的共同祖先开始直到refB的所有提交。例:  
    ```
    若 1 - 2 - 3 - 4 ← refA
             \ 5 - 6 ← refB
    refA..refB 即为6 5, refB..refA 即为4 3.
    origin/master..master为master分支上还未提交到远端的所有引用
    ```
    * 区间筛选，例：以下三条命令等价：
        `$ git log refA..refB`
        `$ git log ^refA refB`
        `$ git log refB --not refA `
    * `refA...refB`选择出被两个引用中的一个包含但又不被两者同时包含的提交，即`refA..refB`+`refB..refA`
    ```
    $ git log --left-to-right refA...refB
    < 4
    < 3
    > 6
    > 5
    ```

# 4. 工作储藏
`git stash "comments"`, 储藏所有工作，包括已添加的和已更改未添加的
`git stash apply stash@{0}`, 恢复储藏，可能会产生冲突，解决冲突后`git add`添加
`git stash list`, 列出储藏栈
`git stash show [stashname]`, 查看储藏的更改
`git stash pop`, 应用栈顶储藏并弹出储藏
`git stash drop [stashname]`, 删除储藏
`git stash branch [branchname]`, 应用储藏到某分支并切换到该分支
`git stash --keep-index`, 只储藏已更改未添加的改动，不储藏 已添加的
`git stash --include-untracked`, 储藏未追踪的文件（未追踪文件）并将其从工作目录中删除
`git stash --all`, 贮藏所有文件


# 5. git reset
`git reset --soft HEAD^`, 将HEAD指针及当前分支指针指向HEAD^, 工作目录中所有的文件均添加暂存(staged), 类似“撤销`git commit`命令”，但如果重新提交，会创造一个hash不同的commit，即时提交内容完全相同。
`git reset --mixed HEAD^`, 将HEAD指针及当前分支指针指向HEAD^, 工作目录中的更改未暂存(modified but not staged), 类似“撤销`git add`和`git commit`命令”
`git reset --hard HEAD^`, 将HEAD指针及当前分支指针指向HEAD^, 工作目录中的更改完全丢失，相当于“撤销更改, `git add`和`git commit`”.
