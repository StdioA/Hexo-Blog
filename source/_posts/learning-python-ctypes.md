title: Python 学习之 ctypes
date: 2016-12-30 11:27:00
categories:
- Python
tags:
- python
- ctypes
toc: true

---

课设写不完，要炸了… Windows UI 好难写，只好用 PyQt 解决… 怎么把 Python 和 C 模块连起来？就用 ctypes.

<!-- more -->

# 1. 简介
ctypes 是一个 Python 外部库。它提供 C 语言兼容的数据结构，而且我们可以通过它来调用 DLL 或共享库(shared library)的函数。

# 2. 快速上手
编写 `mod.c`:
```C
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
 
int add(int a, int b)
{
    return a+b;
}
 
char *hello(char *name)
{
    char *str;
    str = (char *)malloc(100);
    sprintf(str, "Hello, %s", name);
    return str;
}
```

将 C 程序编译为共享库：`gcc -fPIC -shared mod.c -o mod.so`

编写 Python 代码：
```Python
# coding: utf-8
from ctypes import *
 
mod = CDLL("./mod.so")
 
print(mod.add(1, 2))        # 3
 
hello = mod.hello
hello.restype = c_char_p    # Set response type,
                            # otherwise an address(integer) will be returned
 
world = create_string_buffer(b"World")
res = hello(byref(world))   # Pass by pointer
print(res)
```

完。

剩下的等写完课设再更。
