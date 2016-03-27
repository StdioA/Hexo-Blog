title: Python学习之collections
date: 2015-10-29 20:55:00
categories:
- Python
tags:
- python
- collections

---
    
collection是Python内建的一个实用工具包，提供一些使用的容器，用于对传统容器类型进行功能提升。

<!-- more -->

# 1. 简介
上面就是简介，不会扯了，贴两个链接好了。
[Collections - 廖雪峰的Python教程](http://www.liaoxuefeng.com/wiki/001374738125095c955c1e6d8bb493182103fac9270762a000/001411031239400f7181f65f33a4623bc42276a605debf6000)
[Collections - High-performance container datatypes - Python Doc](https://docs.python.org/2.7/library/collections.html#module-collections)

# 2. namedtuple
## 2.1 功能简介
`namedtuple`为一个函数，用于产生一个tuple类的子类。由该类实例化后的对象可以像tuple一样pack/unpack，也可以通过预先定义好的成员名称访问对象内部的数据。例：
```bash
>>> from collections import namedtuple
>>> Point = namedtuple("Point", ['x', 'y'])
>>> Point                                           # namedtuple函数返回一个类
<class '__main__.Point'>
>>> Point.__base__
<type 'tuple'>
>>> Point(x=1, y=2)
Point(x=1, y=2)
>>> Point(1,2)                                      # 两种定义方式均可，但是要注意的是，Point的参数长度是不可变的
Point(x=1, y=2)
>>> p = Point(1, 2)
>>> p.x
1
>>> x, y = p                                        # 可以像tuple一样进行unpack
>>> x
1
>>> Point._make((1,2))                              # pack的方式和tuple不同
Point(x=1, y=2)
```

## 2.2 成员函数
* `somenamedtuple._make(iterable)`, 将一个可迭代对象（list, tuple等）转化为一个namedtuple对象。
```python
>>> t = [11, 22]
>>> Point._make(t)
Point(x=11, y=22)
```
* `somenamedtuple._asdict()`, 将对象转化一个为键值和数据对应的OrderedDict.
```python
>>> p._asdict()
OrderedDict([('x', 11), ('y', 22)])
```
* `somenamedtuple._replace(kwargs)`, 返回一个新对象，该对象中的值按照_replace函数中的参数所改变。
```python
>>> p = Point(x=11, y=22)
>>> p._replace(x=33)
Point(x=33, y=22)
```
* `somenamedtuple._fields`, 返回一个所有键值构成的dict.
```python
>>> p._fields                                                       # 查看键值名字
('x', 'y')

>>> Color = namedtuple('Color', 'red green blue')
>>> Pixel = namedtuple('Pixel', Point._fields + Color._fields)      # 将Color和Point“合并”
>>> Pixel(11, 22, 128, 255, 0)
Pixel(x=11, y=22, red=128, green=255, blue=0)
```

## 2.3 应用场景
列举一个应用场景：在进行sql查询时，`cur.fetchone()`会返回一个tuple而不是dict，我们可以定义一个namedtuple, 然后将sql的查询数据转为namedtuple类，然后通过成员函数访问。举个栗子（私货）：
```python
Account = namedtuple("Account", ["username", "password"])   # 定义Account类
...
cur = self.db.cursor()
cur.execute("SELECT username, password FROM account\
                WHERE valid=1\
                ORDER BY RANDOM()\
                LIMIT 1")
account= Account._make(cur.fetchone())                      # 将数据tuple转为Account对象
self.postdata["username"] = account.username                # 直接访问成员数据
self.postdata["password"] = account.password
...
# self.login(username, password)
self.login(account)                                         # 可以将account整体作为参数，使程序看起来更简洁
...
```

# 3. deque - 双端队列
## 3.1 介绍
deque为Double-Ended Queue的简写，用于提供一个可以快速进行两端的插入或删除操作的队列，同时可以像list一样通过下标访问数据。由于list是储存在线性空间的，其插入/删除数据的时间复杂度为O(n)，所以当list比较长的时候，在其头部插入/删除数据的操作需要耗费大量时间。所以如果需要对list对象进行大量头和尾的插入删除操作时，使用deque会使程序的运行效率更高（当然，使用Queue模块的Queue和LifoQueue也是不错的选择）。需要注意的是，deque对象不支持在中间位置的插入操作。
deque对象也被包含在Queue模块中。

## 3.2 成员函数
* `append(x)`, `appendleft(x)`, 在deque的尾/头插入对象x, 与list类似；
* `extend(iterable)`, `extendleft(iterable)`, 与list的extend类似；
* `pop(x)`, `popleft(x)`, 与list的pop类似；
* `clear()`, 清空deque的所有元素；
* `count(x)`, 统计deque中x元素的个数；
* `remove(x)`, 删除deque中的x元素，若x元素不存在，则触发ValueError;
* `reverse()`, 将deque反转，返回None, 与list的reverse类似;
* `rotate(i)`, 将deque最末尾i个元素取出并插入头部，若i<0, 则将头部|i|个元素取出添加到尾部。

# 4. Counter - 计数器
## 4.1 介绍
Counter类是dict的子类，提供一个计数器，可对hashable的对象进行计数。需要注意的是，对象的计数可以小于0.

## 4.2 基本操作
直接上例子：
```python
>>> c = Counter()                           # 空计数器
>>> c = Counter('gallahad')                 # 对iterable对象中的元素进行计数
>>> c = Counter({'red': 4, 'blue': 2})      # 从map中读取元素和对应计数
>>> c = Counter(cats=4, dogs=8)             # 从关键字参数中获取
>>> c
Counter({'dogs': 8, 'cats': 4})
>>> c["dogs"]                               # 获取计数
8
>>> c["elephants"]                          # 对于不存在的元素，会返回0 
0
```

## 4.3 成员函数
* `elements()`, 返回一个迭代器，包含Counter中的所有元素，若元素x的计数为n，则x在迭代器中会出现n次。
```python
c = Counter(a=4, b=2, c=0, d=-2)
>>> list(c.elements())
['a', 'a', 'a', 'a', 'b', 'b']
```
* `most_common([n])`, 返回计数最多的n个元素。
```python
>>> Counter('abracadabra').most_common(3)
[('a', 5), ('r', 2), ('b', 2)]
```
* `update([iterable-or-mapping])`, 在现有计数上添加参数所对应的计数。
```python
>>> c = Counter("1233")
>>> c
Counter({'3': 2, '1': 1, '2': 1})
>>> c.update("2345")
>>> c
Counter({'3': 3, '2': 2, '1': 1, '5': 1, '4': 1})
```
* `substact([iterable-or-mapping])`, 在现有计数上减去参数所对应的计数。

可以用在Counter对象上的通用函数和常见用法有：
* `sum(c.values())`, 所有计数之和
* `c.clear()`, 将计数重置
* `list(c)`, 返回所有元素构成的列表, 其中元素不会重复
* `set(c)`, 转化为set
* `dict(c)`, 转化为dict
* `c.items()`, 转化为由(元素, 计数)构成的列表
* `Counter(dict(list_of_pairs))`, 从上述列表转化为Counter对象
* `c.most_common()[:-n-1:-1]`, 出现最少的n个元素
* `c += Counter()`, 移除所有元素计数为0的元素
注意: 将某元素的计数赋值为0, 并不能在计数器的元素列表中删除该元素。若要删除该元素，需要用`del`.
```python
>>> c = Counter("12")
>>> c
Counter({'1': 1, '2': 1})
>>> c["0"]
0
>>> c
Counter({'1': 1, '2': 1})
>>> c["0"] = 0
>>> c
Counter({'1': 1, '2': 1, '0': 0})
>>> del c["0"]
>>> c
Counter({'1': 1, '2': 1})
```

# 5. defaultdict
## 5.1 介绍
defaultdict类是dict类的子类，包含dict类的所有功能。与dict不同的是，在调用不存在的键值时，dict会抛出`KeyError`异常，而defaultdict会执行`__missing__`函数, 通过调用该函数来进行操作或返回值。
定义方式如下:
```python
d = defaultdict()                                   # 这种定义方式返回的对象跟普通dict无异
a = defaultdict(lambda: "N/A")                      # 定义__missing__函数

print a["0"]                                        # 返回"N/A", 同时a["0"]被赋值为"N/A"
```

## 5.2 应用场景
```python
>>> s = [('yellow', 1), ('blue', 2), ('yellow', 3), ('blue', 4), ('red', 1)]
>>> d = defaultdict(list)                                           # 若键值不存在，则赋值为[]
>>> for k, v in s:
...     d[k].append(v)
...
>>> d.items()
[('blue', [2, 4]), ('red', [1]), ('yellow', [1, 3])]
```

# 6. OrderdDict
## 6.1 简介
使用`dict`时，Key是无序的。在对`dict`做迭代时，我们无法确定Key的顺序。如果要保持Key的顺序，可以用`OrderedDict`.
`OrderedDict`的键值顺序按照键值被插入的顺序排列，而不是Key本身的顺序。
`OrderdDict`的键值可以通过调用`popitem`函数被弹出，弹出时弹出(key, item)的键值对。
```python
>>> a = OrderedDict({"one": 1})                     # 可通过dict定义，也可传入可迭代对象
>>> a["zero"] = 0
>>> a["two"] = 2
>>> a
OrderedDict([('one', 1), ('zero', 0), ('two', 2)])
>>> a.popitem()                                     # 弹出最后一个键值对
('two', 2)
>>> a.popitem(last=False)                           # 弹出第一个键值对
('one', 1)
>>> a
OrderedDict([('zero', 0)])
```

# 7. 参考文档
1. [8.3 collections — High-performance container datatypes](https://docs.python.org/2.7/library/collections.html#collections.OrderedDict)
2. [collections - 廖雪峰的官方教程](http://www.liaoxuefeng.com/wiki/001374738125095c955c1e6d8bb493182103fac9270762a000/001411031239400f7181f65f33a4623bc42276a605debf6000)

# 8. 后记
这篇博文写了两天，感觉在翻译文档，有点枯燥。一直觉得应该多去了解Python的内置常用模块，不要求完全熟练，至少要有个印象，这样在开发中需要使用这些模块的时候，才能够想起来用它。
所以，又填完一个坑。下一个应该是unittest吧。
