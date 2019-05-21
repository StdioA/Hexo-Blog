---
title: 用堆找出最小的 N 个数
author: David Dai
tags:
  - Python
  - 算法
  - 堆
categories:
  - 随笔
date: 2019-05-21 16:26:34
---

不知道为啥，突然想水一篇很水的算法文章。

<!--more-->

今天整理 MySQL 的笔记，看到了这样一句话：
> MySQL 在执行 `ORDER BY x LIMIT n` 这类语句，且 `LIMIT` 的数量有限时（比如只需要 3 条数据），MySQL 会尽量通过堆来构建优先队列，减少排序所需的时间。

这是堆的一个经典应用：从海量数据中找出最大（小）的 n 个数。  
之前只用堆写过堆排，没有用堆处理过在线算法，所以就写了写。

用一句话概括这个算法：要找最小的数，就要构建大顶堆。

在处理数据时，我们会构建一个大顶堆 `H`，那么 `H[0]` 的值也就是当前数据中**最小**的 N 个数中的**最大值**，也就是第 N 小的数。  
当处理新的数时，如果这个数小于堆顶的数，那么就把它变成堆顶，然后再对堆进行维护，以保证有序。

此算法的时间复杂度为 `O(MlogN)`，空间复杂度为 `O(N)`。其中，M 是海量数据的数量，而要求出的最小数的数量。

最后上代码，使用 Python 语言编写。

```python
import random


class HeapForMinN:
    """求最小的 N 个数"""

    def __init__(self, size):
        self.size = size
        self.heap = []
    
    def add(self, num):
        if len(self.heap) < self.size:
            # 数据太少
            self.heap.append(num)
        elif len(self.heap) + 1 == self.size:
            # 数据数量够了，构建堆
            self.heap.append(num)
            self.build_heap()
        elif num < self.heap[0]:
            # 使用新数替换堆顶
            self.heap[0] = num
            self.max_heapify(0)

    def build_heap(self):
        length = len(self.heap)
        for i in range((length + 1) // 2, -1, -1):
            self.max_heapify(i)

    def max_heapify(self, index):
        length = len(self.heap)
        while index < length:
            maxi = index
            i1, i2 = index * 2 + 1, index * 2 + 2
            if i1 < length and self.heap[maxi] < self.heap[i1]:
                maxi = i1
            if i2 < length and self.heap[maxi] < self.heap[i2]:
                maxi = i2
            if maxi == index:
                break
            self.heap[index], self.heap[maxi] = self.heap[maxi], self.heap[index]
            index = maxi

    def min_n(self):
        return self.heap


if __name__ == "__main__":
    list_ = list(range(1, 10001)) * 2
    random.shuffle(list_)
    N = 10
    hn = HeapForMinN(N)
    for i in list_:
        hn.add(i)
    min_n = hn.min_n()
    print(min_n)        # [5, 4, 5, 4, 3, 3, 1, 1, 2, 2]
    assert sorted(min_n) == sorted(list_)[:n]
```
