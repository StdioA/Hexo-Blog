title: Python学习之unittest
date: 2015-11-12 18:53:02
excerpt: "This is my post excerpt"
categories:
- Python
tags:
- python
- unittest

---

单元测试，工程开发中重要的一环。

<!-- more -->

# 1. 简介
单元测试是用来对一个模块、一个函数或者一个类来进行正确性检验的测试工作。

# 2. 使用unittest进行测试
unittest是Python自带的单元测试模块，通过编写测试类，在测试类中编写测试函数的方式进行测试。
不知道该测试什么，就对python的list对象测试一下好了。

## 2.1 编写测试文件
```python
# coding: utf-8

import unittest                                 # 使用unittest模块

class TestStringMethods(unittest.TestCase):     # 测试类需继承自unittest.TestCase
    def setUp(self):
        # 在函数里进行测试之前需要进行的一些准备工作，setUp在每次测试之前都会运行一次
        self.d = []

    def test_empty(self):
        """\
        Empty testcase
        """
        self.assertEqual([], [])                # 若a和b相等则通过
        self.assertNotEqual([1], [])            # 若不相等则通过

    def test_bool(self):
        """\
        Bool transform testcase
        """
        self.assertTrue([1])                    # 若内部的表达式或对象为真则通过
        self.assertFalse([])                    # 为假则通过
        self.assertTrue(bool([1]))
        self.assertFalse(bool([]))

    def test_append(self):
        """\
        Append testcase
        """
        list_ = []
        list_.append(1)
        self.assertEqual(list_, [1])
        self.assertNotEqual(list_, [])
        list_.append(2)
        self.assertIn(1, list_)                 # 若a in b则通过
        self.assertNotIn(3, list_)              # 若a not in b则通过
        self.assertTrue(list_[0] == 1)
        self.assertEqual(list_[1], 2)

    def test_instance(self):
        list_ = []
        self.assertIsInstance(list_, list)      # 若isinstance(a, b)为真则通过
        self.assertNotIsInstance(list_, str)    # 若isinstance(a, b)为假则通过

    def test_index(self):
        """\
        Index testcase
        """
        list_ = []
        with self.assertRaises(IndexError):     # 检测是否有异常抛出，有指定异常抛出则通过
            a = list_[0]

        list_ = [1, 2, 3]
        self.assertEqual(list_[1], 2)
        self.assertEqual(list_[-2], 2)

        with self.assertRaises(IndexError):
            a = list_[4]

    def tearDown(self):
        # 在测试之后需要进行的一些处理事项，tearDown在每次测试之后都会运行一次
        del self.d

if __name__ == '__main__':
    __import__("sys").argv.append("-v")         # 采用verbose方式，输出测试信息
    unittest.main()
```

## 2.2 运行测试
* 直接运行测试程序，输出:
```bash
test_append (__main__.TestStringMethods)
Append testcase ... ok
test_bool (__main__.TestStringMethods)
Bool transform testcase ... ok
test_empty (__main__.TestStringMethods)
Empty testcase ... ok
test_index (__main__.TestStringMethods)
Index testcase ... ok
test_instance (__main__.TestStringMethods) ... ok

----------------------------------------------------------------------
Ran 5 tests in 0.001s

OK
```
* 运行`python -m unittest test.py`也可以进行测试；
* 运行`python -m unittest`, unittest会在当前文件夹中寻找测试类进行测试（真智能）；
* 若有测试未通过，unittest会在测试时报告FAIL, 并在测试结束后将所有未通过测试的项目列出。

# 3. 使用Nose进行单元测试
Nose是一个对unittest的扩展测试框架，能自动发现并运行测试。使用Nose，可以将单元测试代码的编写变得更简单，不用再构造测试类，只需要在以`test`开头的文件中建立以`test`开头的函数即可。

## 3.1 编写测试文件
```python
# coding: utf-8

import nose

def test_nose_installed_successfully():
    import nose                         # 运行测试代码
    assert True                         # assert True表示测试成功

def test_obviously_failed():
    assert False                        # assert False表示测试失败，测试时会报"FAIL"

def test_returns_an_exception():
    raise ValueError                    # 若抛出除AssertionError的异常，测试时会报"ERROR"

nose.main()
```

## 3.2 运行测试文件
nose自带可执行文件，所以只需要输入`nosetests [测试文件名]`即可，若测试文件名为空，则nosetest会在当前文件夹寻找所有测试。以下命令格式均可接受：
```bash
nosetests test.module
nosetests another.test:TestCase.test_method
nosetests a.test:TestCase
nosetests /path/to/test/file.py:test_function
```

运行`nosetests test_with_nose -v `，输出：
```bash
$ nosetests test_with_nose -v
test_with_nose.test_nose_installed_successfully ... ok
test_with_nose.test_obviously_failed ... FAIL
test_with_nose.test_returns_an_exception ... ERROR

======================================================================
ERROR: test_with_nose.test_returns_an_exception
----------------------------------------------------------------------
Traceback (most recent call last):
  File "c:\python35\lib\site-packages\nose\case.py", line 198, in runTest
    self.test(*self.arg)
  File "C:\Users\Stdio\Desktop\temp\utest\test_with_nose.py", line 13, in test_returns_an_exception
    raise ValueError
ValueError

======================================================================
FAIL: test_with_nose.test_obviously_failed
----------------------------------------------------------------------
Traceback (most recent call last):
 File "c:\python35\lib\site-packages\nose\case.py", line 198, in runTest
   self.test(*self.arg)
  File "C:\Users\Stdio\Desktop\temp\utest\test_with_nose.py", line 10, in test_obviously_failed
   assert False
AssertionError

----------------------------------------------------------------------
Ran 3 tests in 0.003s

FAILED (errors=1, failures=1)
```

# 4. 参考文档
1. [unitest - Python Doc](https://docs.python.org/3/library/unittest.html)
2. [单元测试 - 廖雪峰的python教程](http://www.liaoxuefeng.com/wiki/0014316089557264a6b348958f449949df42a6d3a2e542c000/00143191629979802b566644aa84656b50cd484ec4a7838000)
3. [Nose documentation](https://nose.readthedocs.org/en/latest/)
4. [Python的学习（十八）---- 单元测试工具nose](http://blog.csdn.net/linda1000/article/details/8533349)

# 5. 后记
拖了一个月，这个坑再不填日子没法过了(╯‵□′)╯︵┻━┻ 
到现在没给自己的代码写过测试，也是醉…
一直想转Py3一直没转，昨天改了环境变量，强迫自己用一阵子，多去看看Py3的特性，有时间整理一下。
Git的坑还没填完不过自己的Pro Git看的差不多了，貌似也满足自己的需求了，准备弃坑。
下一篇有可能是SQLAlchemy, MongoDB, Flask…哎，随它去吧。
