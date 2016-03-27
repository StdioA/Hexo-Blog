title: Python学习之上下文管理
date: 2015-10-22 20:00:00
categories:
- Python
tags:
- python
- 上下文管理

---

最近突发奇想，想写一个能改变当前输出环境，输出彩色文字的上下文管理器，于是学习了一下上下文管理。

<!-- more -->

# 1. 简介
上下文管理器(context manager)是Python2.5开始支持的一种语法，用于规定某个对象的使用范围。一旦进入或者离开该使用范围，会有特殊操作被调用 (比如为对象分配或者释放内存)。一个最简单的例子:
```python
with file("a.txt", "r") as f:
    s = f.read()
    print f.closed                  # 此时文件是打开的

print f.closed                      # 变量s和f依然存在，但此时f已关闭
```

# 2. 编写上下文管理器
为一个类编写上下文管理器时，需要定义类的`__enter__`和`__exit__`函数。

## 2.1 \_\_enter\_\_函数
`__enter__`函数规定如下：
> contextmanager.\_\_enter\_\_()
>   Enter the runtime context and return either this object or another object related to the runtime context. The value returned by this method is bound to the identifier in the as clause of with statements using this context manager.
> 进入当前上下文，并返回当前对象或另一个与当前上下文相关的对象。被该函数返回的变量会通过该上下文管理器与with文法中as后的声明所绑定。

例：
```python
class BracketAdder(object):
    def __enter__(self):
        # do something with self and runtime context
        sys.stdout.write("(")           # 输出左括号
        return self

    # __exit__函数略

with BracketAdder() as ba:
    # do_something

```

## 2.2 \_\_exit\_\_函数
`__exit__`函数格式如下:
> contexmanager.\_\_exit\_\_(exc\_type, exc\_val, exc\_tb)

*文档太长懒得翻译了*\_(:зゝ∠)\_

如果代码块中出现异常，`exc_type`为异常类型，`exc_val`为该异常，`exc_tb`为traceback.`__exit__`函数应返回一个bool类型，若返回值为`True`, 则在代码块及`__exit__`函数运行结束后不会抛出任何异常，然后立即执行后面的代码；若返回值为`False`，则在`__exit__`函数运行结束后抛出异常。  
若代码块中未出现异常，则三个变量均为`None`.
例：
```python
class BracketAdder(object):
    # __enter__函数略
    def __exit__(self, exc_type, exc_val, exc_tb):
        # do something with self and runtime context
        sys.stdout.write(')')
        if exc_type == NameError:
            return True
        else:
            return False

with BracketAdder():
    print a                    # 若执行这一条引发NameError异常，则异常不会被抛出
    a = 1/0                    # 若执行这一条引发ZeroDivisionError, 则异常会被抛出
```

# 3. 综合示例
```python
# coding: utf-8

import sys

class CManager(object):
    def __init__(self):
        self.in_context = False

    def __enter__(self):
        self.in_context = True
        print "I'm entering the context"
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.in_context = False

        print "I'm leaving the context"
        
        if exc_type == NameError:       # 拦截NameError异常
            return True
        else:
            return False

    def show(self):
        print "I'm {status}in the context.".format(
                    status="" if self.in_context else "not ")


with CManager() as ba:
    ba.show()
    # blahblah
    raise NameError                     # 触发NameError异常，代码块后面的语句实际不会执行，而该异常会在__exit__函数中被拦截
    print 'k'

ba.show()

```
程序输出:
```text
I'm entering the context
I'm in the context.
I'm leaving the context
I'm not in the context.
```


# 4. 参考文档
1. [Context Manager Types - Python Doc](https://docs.python.org/2.7/library/stdtypes.html#index-37)
2. [Python深入02 上下文管理器](http://www.cnblogs.com/vamei/archive/2012/11/23/2772445.html)
3. [浅谈 Python 的 with 语句](http://python.jobbole.com/82494/)

# 5. 后记
这些小知识点平时只是粗略的看一下，只能达到会用的程度，有时候觉得只有自己整理一遍才能真正理解它们，才能熟练运用。  
所以…又填完一个坑。周末简单写写PyQt4吧…
