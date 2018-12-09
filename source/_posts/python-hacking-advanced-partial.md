title: 'Python Hacking: “高级”偏函数'
author: David Dai
tags:
  - Python
categories:
  - Python
date: 2018-12-02 10:06:00
toc: true

---
本文讲解了一个需求的解决方案，而这个奇葩需求你在 99.93% 场景下都不会遇到，就算遇到了，也一定有其它更简单的解决方案。

<!--more-->
# 0. 引言
> `>>> print((lambda x:None).__code__.__doc__)`
> code(argcount, kwonlyargcount, nlocals, stacksize, flags, codestring,  
>     constants, names, varnames, filename, name, firstlineno,  
>     lnotab[, freevars[, cellvars]])  
> 
> Create a code object.  **Not for the faint of heart.**

# 1. 需求
[`mock.patch`](https://docs.python.org/3/library/unittest.mock.html#patch) 对象在用做装饰器时，会生成一个偏函数，来将原始函数的第一个位置参数覆盖为一个 mock 对象，如：
```python
@mock.patch('func')
def test_something(mock, a, b):
    pass
```
此处 `mock` 参数就是 patch 对应的 mock 对象。

而 [pytest](https://pytest.org/) 有一个有点厉害的功能：它会读取测试函数的参数，并在 `conftest.py` 中寻找每个参数对应的同名 [fixture](http://doc.pytest.org/en/latest/fixture.html) 并加载。

前两天就遇到了这样的问题：我要做一个类似 `patch` 的装饰器，在 pytest 测试函数执行时向内部注入一个参数，而且要保证这个参数在被 pytest 解析时**不能暴露在参数列表中**，否则 pytest 会因为找不到参数的 fixture 而报错。

```python
@my_patch('func')
def test_something(mock, fixa, fixb):
    # mock 参数由 @my_patch 注入
    pass

# 在外界看来，test_something 的参数只能有 fixa 和 fixb，而不能有 mock
import inspect
assert str(inspect.signature(test_something)) == "(fixa, fixb)"
```

# 2. 铺垫一下：简单的偏函数装饰器实现
如果想实现一个简单的偏函数装饰器，非常简单：

```python
def my_partial(*partial_args, **partial_kwargs):
    def decorator(func):
        def wrapped(*args, **kwargs):
            args = partial_args + args
            partial_kwargs.update(kwargs)
            return func(*(args), **partial_kwargs)
        return wrapped
    return decorator

@my_partial(1, 2, d=5)
def func(a, b, c, d=4):
    return sum([a, b, c, d])

assert func(3) == 11
```

但这样的话，`func` 的参数声明将会变成 `wrapped` 的参数声明，也就是 `(*args, **kwargs)`，而在这个场景中，`func` 的参数声明应该是 `(c, d=5)`，才能够真正满足 pytest 的要求。  
于是，我们就需要在定义 `wrapped` 函数的时候，对它的参数声明进行定制化。

# 3. 来点基础知识
## 3.1 获取函数参数声明
获取完整的函数参数声明，要涉及到函数的几个私有属性：
* `__code__`: 编译过的函数代码对象，类型为 `types.CodeType`
* `__defaults__`: 函数的序列参数默认值，类型为 `None` 或 `tuple`
* `__kwdefaults__`: 函数的关键字参数默认值，类型为 `None` 或 `dict`

其中 `code` 对象中的几个属性也需要用到：
* `co_varnames`: 函数声明中所有参数的变量名，其中 `*` 和 `**` 可变参数的变量名会放在最后
* `co_argcount`: 序列参数的数量
* `co_kwonlyargcount`: 严格关键词参数的数量
* `co_flags`: 函数性质标记，`0x04` 位声明这个函数是否用到了 `*args`，而 `0x08` 位声明这个函数是否用到了 `**kwargs`

有了这些属性，我们就可以获取到函数的参数声明信息。具体实现可以看[下面的代码汇总](#-21-%E8%8E%B7%E5%8F%96%E5%87%BD%E6%95%B0%E5%8F%82%E6%95%B0%E5%A3%B0%E6%98%8E%E7%9A%84-funcparser)。

当然，Python 3 还支持[函数注解](https://www.python.org/dev/peps/pep-3107/)，需要用到函数的 `__annotations__` 属性，但在我的代码中没有对这部分进行解析。

## 3.2 inspect 库
当然，这些轮子，Python 自带库 [`inspect`](https://docs.python.org/3/library/inspect.html) 都已经帮我们造好了。

`inspect.getfullargspec(f)` 可以获取到以上所有的参数信息：
```python
>>> def func(a, b=1, *args, c=2, **kwargs):
...     pass
...
>>> inspect.getfullargspec(func)
FullArgSpec(args=['a', 'b'], varargs='args', varkw='kwargs', defaults=(1,),
            kwonlyargs=['c'], kwonlydefaults={'c': 2}, annotations={})
```

而 `inspect.signature` 更加强大，它除了解析函数外，还支持“模拟调用”函数，对调用方式进行合法性验证，并展示调用之后函数中每个参数的值：
```python
>>> s = inspect.signature(func)     # func 定义见上
>>> s
<Signature (a, b=1, *args, c=2, **kwargs)>
>>> b = s.bind(0)
>>> b
<BoundArguments (a=0)>
>>> b.apply_defaults()
>>> b
<BoundArguments (a=0, b=1, args=(), c=2, kwargs={})>
```

实现了这些功能，我们就可以对给定函数进行参数解析并修改了。

~~这都算哪门子基础知识嘛…~~

# 4. 开始动手
先来写一个简单的、只支持序列参数，且不支持参数默认值的偏函数实现。

```python
def func(a, b, c):
    return sum([a, b, c])

if __name__ == '__main__':
    print(func(1, 2, 3,), func, signature(func))
    # 6 <function func at 0x02D23B28> (a, b, c)
    partial = nb_partial(3, 4)(func)
    # 我是个装饰器，相信你看得懂
    print(partial(5), partial, signature(partial, follow_wrapped=False))
    # 12 <function func at 0x02DA0DB0> (c)
```

注：`sigature(follow_wrapped=False)` 会改变 `signature` 的默认行为：只查看函数本身的定义，而不会通过 `__wrapped__` 链去查找原函数。

## 4.1 大体框架
搭个架子：

```python
from functools import wraps

def nb_partial(func, *partial_args):
    pre_arg_str = ', '.join(map(repr, partial_args))
    # 这里是所有预定义函数的值（"3, 4"）

    code = func.__code__
    if len(partial_args) > code.co_argcount:
        raise TypeError("Too many positional arguments")

    args = code.co_varnames[:code.co_argcount]
    post_args = args[len(partial_args):]
    post_args_str = ', '.join(post_args)
    # 这里是排除预定义函数后，剩下的参数名列表（"c"）

    def wrapped(c):             # To be implemented
        return func(3, 4, c)

    return wraps(func)(wrapped)
```

现在我们有了需要的变量值和变量名信息。可是，我们还需要根据这些变量信息**动态**定义 `wrapped` 函数，这个就有点麻烦了…

## 4.2 黑科技：code object & FunctionType

Python 毕竟是“万物皆对象”的动态语言。  
从字符串编译一段代码？没问题！用函数类定义一个函数对象？也没什么问题！  
牛（就）逼（是）了（慢）。

[`compile`](https://docs.python.org/3/library/functions.html#compile) 内置函数可以让我们动态地通过字符串来编译代码，来生成一个代码对象。这段代码可以是一个完整的程序，一两个定义，也可以是几个表达式。  
而通过 `types.FunctionType`，我们就可以使用代码对象，并向其中注入一些信息（如 `globals`），即可生成一个可用的函数。

```python
func_def = """
def f(a):
    return a + 1
"""
module_code = compile(func_def, '<>', 'exec')
function_code = next(code for code in module_code.co_consts
                     if isinstance(code, types.CodeType))
# 使用 exec 模式编译一段代码，并从中获取定义好的 code 常量（也就是 f 函数的代码）
func = types.FunctionType(function_code, {})
# 生成一个函数，第二个参数是 `globals`. 因为函数中不会调用到全局变量，所以这里传递空字典
print(func, func(2))
# <function f at 0x00E78AE0> 3
```

有了这个，我们就可以从一段字符串中动态生成一个 Python 函数了。

```python
def nb_partial(*partial_args):
    """
    生成一个偏函数，且新函数的参数列表中会去除
    """
    def decorator(func):
        pre_arg_str = ', '.join(map(repr, partial_args))

        code = func.__code__
        args = code.co_varnames[:code.co_argcount]
        post_args = args[len(partial_args):]
        post_args_str = ', '.join(post_args)

        func_def = (f"def wrapped({post_args_str}):\n"
                    f"    return func({pre_arg_str}, {post_args_str})")
        module_code = compile(func_def, '<>', 'exec')
        function_code = next(code for code in module_code.co_consts
                             if isinstance(code, types.CodeType))
        wrapped = types.FunctionType(function_code,
                                     {'func': func})
        # wrapped 函数中引用了一个外部变量，所以在定义函数时需要把 func 传进去
        return wraps(func)(wrapped)
    return decorator
```

这里有一个小缺陷没有解决掉：这个函数中的 `func` 是一个全局变量，而不是闭包中的自由变量。  
转念一想，装饰器中的原始 `func` 本身就不是自由变量啊…脑残了 Orz，继续继续。

所以，通过动态定义函数，我们就可以实现还原参数列表的的偏函数功能。大体步骤如下：  
1. 提取原函数参数
2. 处理新函数参数
3. 构造函数定义字符串
4. 动态定义新函数

完整实现见[下面](#-23-高级偏函数装饰器的进阶实现)。

# 5. 再改进一下？
我们的原函数目前只支持了最简单的函数声明方式，而无法支持关键词参数，参数默认值等。

结合上面写的 `FuncParser`，我们还可以对动态函数的定义做进一步改进。  
完整实现有点长，还是[下面见]()。

当然，这个实现还有一些可以改进之处：
1. 偏函数定义时可以支持关键字参数，这需要进行更多的判断，比如定义时传进去的关键字参数，有可能在原函数中是序列参数；  
2. 有很多对象并没有很好地实现 `__repr__` 魔术方法，导致 `repr(arg)` 后生成的字符串在函数定义中并不能做到完整还原。所以我们最好通过 `globals` 将它直接传进新的函数中，而不是使用字符串来进行参数传递。

突然犯懒，就不再做进一步的实现了。

# -2. 代码汇总
## -2.1 获取函数参数声明的 `FuncParser`
```python
from collections import namedtuple

FuncArgs = namedtuple('FuncArgs', ('args', 'defaults'))

class FuncParser:

    @classmethod
    def parse_vars(cls, func):
        code = func.__code__
        argc = code.co_argcount
        kwargc = code.co_kwonlyargcount
        varnames = code.co_varnames

        res = {
            "arg": varnames[:argc],
            "kwarg": varnames[argc:argc + kwargc],
            "*args": None,
            "**kwargs": None
        }
        flag = code.co_flags
        if flag & 0x04:
            # Use *args
            res['*args'] = varnames[argc + kwargc]
        if flag & 0x08:
            # Use **kwargs
            res['**kwargs'] = varnames[argc + kwargc + bool(flag & 0x04)]

        # parse defaults
        defaults = func.__defaults__ or ()
        defaults_map = dict(zip(reversed(res['arg']), reversed(defaults)))
        defaults_map.update(func.__kwdefaults__ or {})
        return FuncArgs(res, defaults_map)

    @classmethod
    def build_param_str(cls, func_args):
        func_vars, defaults_map = func_args
        params = []
        # arg
        for arg in func_vars['arg']:
            if arg in defaults_map:
                default_val = defaults_map[arg]
                params.append(f'{arg}={default_val}')
            else:
                params.append(arg)

        # *args
        if func_vars['*args']:
            params.append(f"*{func_vars['*args']}")
        else:
            if func_vars['kwarg']:
                params.append('*')

        # kwarg
        for kwarg in func_vars['kwarg']:
            if kwarg in defaults_map:
                default_val = defaults_map[kwarg]
                params.append(f'{kwarg}={default_val}')
            else:
                # 一种函数定义不会遇到的情况，但在后面会用到它
                params.append(f'{kwarg}={kwarg}')

        # **kwargs
        if func_vars['**kwargs']:
            params.append(f"**{func_vars['**kwargs']}")

        param = ', '.join(params)
        param_str = '({})'.format(param)
        return param_str

    @classmethod
    def analyse_func_param(cls, func):     
        func_args = cls.parse_vars(func)
        name = func.__name__
        return cls.build_param_str(func_args)

def test_func(a, *args, b=1, **kwargs):
    pass

if __name__ == '__main__':
    import inspect
    import types
    funcs =  [f for f in locals().values()
              if isinstance(f, types.FunctionType)]
    for f in funcs:
        fstr = FuncParser.analyse_func_param(f)
        assert fstr == str(inspect.signature(f))
        # "(a, *args, b=1, **kwargs)"

```

## -2.2 “高级”偏函数装饰器的初步实现
```python
import types
from functools import wraps
from inspect import signature

def nb_partial(*partial_args):
    """
    生成一个偏函数，且新函数的参数列表中会去除
    """
    def decorator(func):
        pre_arg_str = ', '.join(map(repr, partial_args))

        code = func.__code__
        if len(partial_args) > code.co_argcount:
            raise TypeError("Too many positional arguments")

        args = code.co_varnames[:code.co_argcount]
        post_args = args[len(partial_args):]
        post_args_str = ', '.join(post_args)

        func_def = (f"def wrapped({post_args_str}):\n"
                    f"    return func({pre_arg_str}, {post_args_str})")
        module_code = compile(func_def, '<>', 'exec')
        function_code = next(code for code in module_code.co_consts
                             if isinstance(code, types.CodeType))
        wrapped = types.FunctionType(function_code,
                                     {'func': func})
        return wraps(func)(wrapped)
    return decorator

def func(a, b, c):
    return sum([a, b, c])

if __name__ == '__main__':
    print(func(1, 2, 3,), func, signature(func))
    # 6 <function func at 0x03B43B28> (a, b, c)
    partial = nb_partial(3, 4)(func)
    # The decorator
    print(partial(5), partial, signature(partial, follow_wrapped=False))
    # 12 <function func at 0x03BC0DB0> (c)
```


## -2.3 “高级”偏函数装饰器的进阶实现
```python
import types
from functools import wraps
from inspect import signature
from funcparser import FuncParser, FuncArgs
# 上面写好的函数解析器

def nb_partial(*partial_args):
    """
    生成一个偏函数，且新函数的参数列表中会去除
    """
    def decorator(func):
        code = func.__code__
        if len(partial_args) > code.co_argcount:
            raise TypeError("Too many positional arguments")

        func_vars, defaults_map = FuncParser.parse_vars(func)

        # 构建函数定义
        def_vars = func_vars.copy()
        def_vars["arg"] = def_vars["arg"][len(partial_args):]
        def_str = FuncParser.build_param_str(FuncArgs(def_vars, defaults_map))

        # 构造函数调用
        call_vars = func_vars.copy()
        call_vars["arg"] = list(call_vars["arg"])
        call_vars["arg"][:len(partial_args)] = list(map(repr, partial_args))
        call_str = FuncParser.build_param_str(FuncArgs(call_vars, {}))
        call_str = call_str.replace('*, ', '')

        func_def = (f"def wrapped{def_str}:\n"
                    f"    return func{call_str}")
        module_code = compile(func_def, '<>', 'exec')
        function_code = next(code for code in module_code.co_consts
                             if isinstance(code, types.CodeType))

        # 设置序列参数默认值
        argdefs = []
        for arg in reversed(func_vars["arg"]):
            if arg not in defaults_map:
                break
            argdefs.append(defaults_map[arg])
        argdefs.reverse()
        wrapped = types.FunctionType(function_code,
                                     {'func': func},
                                     argdefs=tuple(argdefs))
        # 设置关键字参数默认值
        wrapped.__kwdefaults__ = func.__kwdefaults__
        return wraps(func)(wrapped)
    return decorator

def func(a, b=2, c=3, *, d=5):
    print(locals())
    return sum([a, b, c, d])

if __name__ == '__main__':
    print(func(1, 2, 3, d=6), func, signature(func))
    # locals: {'a': 1, 'b': 2, 'c': 3, 'd': 6}
    # 12 <function func at 0x034B3AE0> (a, b=2, c=3, *, d=5)
    # The decorator
    partial = nb_partial(3)(func)
    print(partial(6, d=7), partial, signature(partial, follow_wrapped=False))
    # locals: {'a': 3, 'b': 6, 'c': 3, 'd': 7}
    # 19 <function func at 0x00875390> (b=2, c=3, *, d=5)
```

# -1. Reference & 延伸阅读
* [Python 自定义函数的特殊属性（收藏专用）](https://segmentfault.com/a/1190000005685090)
* [Python 文档：Data Model](https://docs.python.org/3/reference/datamodel.html)
* [Python 创建动态函数](http://blog.soliloquize.org/2016/07/06/Python%E5%8A%A8%E6%80%81%E5%88%9B%E5%BB%BA%E5%87%BD%E6%95%B0/)
