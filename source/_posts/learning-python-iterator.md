title: Python学习之迭代器
date: 2015-10-26 15:43:00
categories:
- Python
tags:
- python
- 迭代器
toc: true

---

昨晚用python写了个简单的链表，突然想起了迭代器，就随手整理一下，顺便过一下itertools模块。

<!-- more -->

# 1. 简介
迭代器(iterator)是Python中一种用来进行惰性迭代的数据类型，迭代器可以惰性地在需要时生成数据并返回已进行迭代，而不需要在开始进行时生成所有数据（有的时候也不可能生成所有数据，比如斐波那契数列的无穷迭代等）然后一个一个返回。

# 2. 用法
## 2.1 最基础的例子: iter()函数生成迭代器
```python
list_ = [1,2]
it = iter(list_)                        # 返回一个迭代器

it.next()                               # 获取一个迭代器元素，返回1
next(it)                                # 使用next内建函数，等同于it.next(), 返回2
it.next()                               # 迭代结束，触发StopIteration异常

list_ = [1,2]
it = iter(list_)
for x in it:                            # 用for循环对迭代器进行循环迭代
    print x                             # 迭代结束时，StopIteration不会被触发
```
值得注意的是，迭代器只可向前迭代，不能向后迭代，获取之前已返回过的值（当然绝大多数的时候也没必要向后迭代）。

## 2.2 生成器
`[x*2 for x in range(10) if x%2==0]`, 这段代码为列表生成式（也称作列表解析），可以将一个列表进行转化；将方括号变为圆括号，`(x*2 for x in range(10) if x%2==0)`, 则该代码变为一个生成器(generator)。
生成器可以像迭代器一样进行迭代，还可以通过生成器的成员函数和生成器内部的代码进行数据交换。

# 3. 自定义迭代器及生成器
## 3.1 自定义迭代器
在编写类的时候，可以通过定义类的`__iter__`函数来使类可以转化为迭代器。在`__iter__`函数结束时，会自动引发`StopIteration`异常。例: 
```python
class Counter(object):
    def __init__(self):
        self.number = 0

    def __iter__(self):
        while True:
            yield self.number
            self.number += 1

a = Counter()
it = iter(a)
print it.next()                 # 0
print it.next()                 # 1

for num in it:
    print num                   # 因为该迭代器为无穷迭代，所以会导致死循环
```

## 3.2 自定义生成器
可以用函数自定义生成器。例：
```python
def fib():
    a, b = 1, 1
    while True:
        yield a                 # 用yield函数来在迭代时返回值，在下次迭代时，自动从yield语句的下一条语句开始执行
        a, b = b, a+b

k = fib()
for i in range(10):
    print next(k)               # 1 1 2 3 5 8 13 21 34 55

print type(k)                   # <generator object fib at 0x0268BDF0>
```

## 3.3 生成器的进阶用法
`generator.next()`用来进行迭代并获取返回值；`generator.close()`用来关闭生成器，并在下次迭代时出发StopIteration异常。
```python
a = fib()
print next(a)                   # 1
a.close()                       # 关闭迭代器
print next(a)                   # StopIteration
```
`generator.send(arg)`可以向生成器内部传递对象，`generator.throw(typ[,val[,tb]])`可以向生成器内部传递异常（包括类型，具体异常对象和Traceback）；在调用这两个函数后，生成器立刻进行迭代并返回值（或触发StopIteration）。具体操作：
```python
def counter():
    num = 0
    while True:
        try:
            ret = yield num
        except ValueError:
            print "ValueError caught"

        if ret is not None:
            num = ret
        num += 1

c = counter()
print c.next()                          # 0
print c.send(3)                         # 传输3，内部yield语句返回3，然后进行下次迭代，生成4
print c.next()                          # 5
c.throw(ValueError)                     # 传递异常，内部yield函数触发异常，然后进行异常处理，输出"ValueError caught"，若未处理，则异常会向上层抛出
print c.next()                          # 6
```

# 4. itertools模块
itertools是Python自带的一个模块，包含很多使用函数，用来对一个或多个可迭代对象进行操作后返回一个迭代器。具体函数列表：

1. 无限迭代器
    * `count(start, [step])` 
      从start开始，以后每个元素都加上step。step默认值为1。
      `count(5)` -> 5 6 7 …
    * `cycle(p)` 
      迭代至p的最后一个元素之后，从p的第一个元素重新迭代。
      `cycle('abc')` -> a b c a b c …
    * `repeat(elem [,n])`
      无限重复或重复n次返回elem。
      `repeat("Ah", 3) ` -> "Ah" "Ah" "Ah"

2. 在最短的序列结束迭代时停止迭代
    * chain(p, q, ...)
      迭代至序列p的最后一个元素后，从q的第一个元素开始，直到所有序列终止。
      `chain("ABC", "DEF")` -> A B C D E F
    * `compress(data, selectors)` 
      如果bool(selectors[n])为True，则next()返回data[n]，否则跳过data[n]。 
      `compress('ABCDEF', [1,0,1,0,1,1])` -> A C E F
    * `dropwhile(pred, seq)` 
      当pred对seq[n]的调用返回False时才开始迭代。 
      `dropwhile(lambda x: x<5, [1,4,6,4,1])` -> 6 4 1
    * `takewhile(pred, seq)` 
      dropwhile的相反版本。 
      `takewhile(lambda x: x<5, [1,4,6,4,1])` -> 1 4
    * `ifilter(pred, seq)`
      内建函数filter的迭代器版本。 
      `ifilter(lambda x: x%2, range(10))` -> 1 3 5 7 9
    * `ifilterfalse(pred, seq)` 
      ifilter的相反版本。 
      `ifilterfalse(lambda x: x%2, range(10))` -> 0 2 4 6 8
    * `imap(func, p, q, ...) `
      内建函数map的迭代器版本。 
      `imap(pow, (2,3,10), (5,2,3))` -> 32 9 1000
    * `starmap(func, seq) `
      将seq的每个元素以变长参数(*args)的形式调用func。 
      `starmap(pow, [(2,5), (3,2), (10,3)])` -> 32 9 1000
    * `izip(p, q, ...) `
      内建函数zip的迭代器版本。 
      `izip('ABCD', 'xy')` -> Ax By
    * `izip_longest(p, q, ..., fillvalue=None)`
      izip的取最长序列的版本，短序列将填入fillvalue。 
      `izip_longest('ABCD', 'xy', fillvalue='-')` -> Ax By C- D-
    * `tee(it, n) `
      返回n个迭代器it的复制迭代器。
    * `groupby(iterable[, keyfunc])`
      这个函数功能类似于SQL的分组。使用groupby前，首先需要使用相同的keyfunc对iterable进行排序，比如调用内建的sorted函数。然后，groupby返回迭代器，每次迭代的元素是元组(key值, iterable中具有相同key值的元素的集合的子迭代器)。或许看看Python的排序指南对理解这个函数有帮助。 
      `groupby([0, 0, 0, 1, 1, 1, 2, 2, 2])` -> (0, (0 0 0)) (1, (1 1 1)) (2, (2 2 2))

3. 组合迭代器
    * `product(p, q, ... [repeat=1])` 
    生成笛卡尔积。 
    `product('ABCD', repeat=2)` --> AA AB AC AD BA BB BC BD CA CB CC CD DA DB DC DD
    * `permutations(p[, r])` 
    生成全排列。 
    `permutations('ABCD', 2)` --> AB AC AD BA BC BD CA CB CD DA DB DC
    * `combinations(p, r)` 
    生成组合。 
    `combinations('ABCD', 2)` --> AB AC AD BC BD CD
    * `combinations_with_replacement()` 
    生成排列元素(p, q), 且p<q.
    `combinations_with_replacement('ABCD', 2)` --> AA AB AC AD BB BC BD CC CD DD

部分文字来源：<http://www.cnblogs.com/huxi/archive/2011/07/01/2095931.html>
        

# 5. 参考文档
1. [itertools - Python Doc](https://docs.python.org/2.7/library/itertools.html#module-itertools)
2. [生成器 - 廖雪峰的Python教程](http://www.liaoxuefeng.com/wiki/001374738125095c955c1e6d8bb493182103fac9270762a000/00138681965108490cb4c13182e472f8d87830f13be6e88000)
3. [Python 迭代器 & \_\_iter\_\_方法](http://blog.csdn.net/bluebird_237/article/details/38894617)
4. [Python函数式编程指南（三）：迭代器](http://www.cnblogs.com/huxi/archive/2011/07/01/2095931.html)
