title: Python学习之logging
date: 2015-10-21 16:00:00
categories:
- Python
tags:
- python
- logging

---

以前写代码的时候，所有信息包括调试信息全部输出在屏幕上，有的时候会看起来乱糟糟的，这时就需要logging模块来记录信息或生成日志文件。

<!-- more -->

# 1. 简介
logging是一个用来记录日志信息的模块，它可以输出信息到stdout或者自定的日志文件中。

# 2. 基本
输出信息：
```python
logging.debug("A debug message")
logging.info("A info message")
logging.warn("A warning message")
logging.error("A error message")
logging.critical("A critical message")
```

日志级别：
logging默认的日志级别为WARNING，在当前级别下，只有warning及以上的日志可以被记录。logging的级别可通过`logging.basicConfig`或`logging.setLevel`来修改。级别大小关系为CRITICAL(50) > ERROR(40) > WARNING(30) > INFO(20) > DEBUG(10) > NOTSET(0), 需要注意的是，一旦记录信息，logging的日志级别就不可再更改。
```python
>>> logging.basicConfig(level=logging.INFO)
>>> logging.info("Info")
INFO:root:Info
>>> logging.basicConfig(level=logging.DEBUG)
>>> logging.debug("Debug")                      # 没有输出
```

# 3. 日志格式设置及日志文件操作

```python
import logging
import time

logging.basicConfig(format="%(asctime)s %(levelname)s %(message)s",
                        datefmt="%Y %b %d %H:%M:%S",
                        filename="./log.log",
                        filemode="w",               # default is "a"
                        level=logging.INFO)

while True:
    for i in range(6):
        logging.log(i*10, "a log")                  # logging.log(level, msg)
        time.sleep(1)
```

log.log输出(可以用tail -f命令实时查看):
```
2015 Oct 21 14:41:36 INFO a log
2015 Oct 21 14:41:37 WARNING a log
2015 Oct 21 14:41:38 ERROR a log
2015 Oct 21 14:41:41 CRITICAL a log
2015 Oct 21 14:41:42 INFO a log
2015 Oct 21 14:41:43 WARNING a log
2015 Oct 21 14:41:44 ERROR a log
2015 Oct 21 14:41:45 CRITICAL a log
```

# 4. 将日志输出到多个流中
```python
import logging
import sys

# 可以通过logging.basicConfig设置一个默认流

console = logging.StreamHandler(stream=sys.stdout)                    # 默认流为sys.stderr
console.setLevel(logging.INFO)
formatter = logging.Formatter('%(name)-12s: %(levelname)-8s %(message)s')
console.setFormatter(formatter)
logging.getLogger().addHandler(console)

files = logging.FileHandler("log2.log", mode="a", encoding="utf-8")   # 设置文件流
files.setLevel(logging.WARNING)
formatter = logging.Formatter("%(levelname)s %(message)s")
files.setFormatter(formatter)
logging.getLogger().addHandler(files)

for i in range(1, 6):
    logging.log(i*10, logging.getLevelName(i*10).lower())

```

# 5. 设置多个logger以记录不同信息
```python
# coding: utf-8

import logging
import sys

log1 = logging.Logger("0.0")

console = logging.StreamHandler(sys.stdout)                               # 默认流为sys.stderr
console.setLevel(logging.INFO)
formatter = logging.Formatter('%(name)-12s: %(levelname)-8s %(message)s')
console.setFormatter(formatter)
log1.addHandler(console)

log2 = logging.Logger("-.-")

files = logging.FileHandler("log2.log", mode="a", encoding="utf-8")     # 设置文件流
files.setLevel(logging.WARNING)
formatter = logging.Formatter("%(name)s %(levelname)s %(message)s")
files.setFormatter(formatter)
log2.addHandler(files)

for i in range(1, 6):
    log1.log(i*10, logging.getLevelName(i*10).lower())
    log2.log(i*10, logging.getLevelName(i*10).lower())

logging.critical("AHHH! I'm the root logger but you forget me!")        # 默认使用logging时logger name为"root"
logging.getLogger("root").info("Of course not!")

```

# 5. 参考文档
1. [Logging HOWTO](https://docs.python.org/2.7/howto/logging.html)
2. [python 的日志logging模块学习](http://www.cnblogs.com/dkblog/archive/2011/08/26/2155018.html)

# 6. 后记
又填完一个坑，最近要填的坑好多…
