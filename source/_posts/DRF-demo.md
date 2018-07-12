title: Django REST Framework 入门
date: 2017-03-09 14:52
categories:
- Python
tags:
- Python
- Django
- Django REST Framework
toc: true

---

拖了很久，终于用 Django REST Framework 写了个小 Demo.

<!-- more -->

# 1. Django REST Framework
[Django REST Framework](http://www.django-rest-framework.org/) 是一个用来构造 Web API 的、强大而灵活的工具包。

最早认识这个框架是在我翻译的[《5 个最受人喜爱的开源 Django 包》](https://linux.cn/article-7679-1.html)中，而真正接触它的时候已经是在五个月之后。那个时候在扇贝实习，第一次体会到了这个框架的简单易用。但那个时候也仅仅是做码农堆代码，并没有对这个框架产生特别深的了解。  
前两天写了一个小 Demo，尝试着用比较少的代码实现对某个模型的增删改查操作，今天简单把实现过程记录一下。完整的代码在[这里](https://github.com/StdioA/drf_demo)。

# 2. 一点准备工作
先安装依赖（推荐使用 virtualenv，虽然我从来不用:new_moon_with_face:）：`pip install Django djangorestframework`.

然后创建工程，创建 app.
```bash
django-admin startproject drf_demo
cd drf_demo
./manage.py startapp demo
```

简单写一个 model：
```python
# demo/models.py
from django.db import models
 
class Data(models.Model):
    content = models.CharField(max_length=128)
```


迁移数据库：
```bash
./manage.py makemigrations
./manage.py migrate
```

在项目设置中添加 APP，顺便把 DRF 也添加进去：
```python
# drf_demo/settings.py
INSTALLED_APPS = [
    # ...
    'rest_framework',
    'demo'
]
```

在 root URL conf 中添加 APP 的 URL：
```python
# drf_femo/urls.py
urlpatterns = [
    # ...
    url(r'^', include('demo.urls'))
]
```

# 3. 开始动手写 API
如果在 Django 中写一个页面，你需要：
1. 在 `urls.py` 中注册 view；
2. 在 `views.py` 中编写 view；
3. 在 `templates` 文件夹中编写模板。

相对地，如果使用 DRF 创建一组 API，你需要：
1. 在 `urls.py` 中定义并注册 router；
2. 在 `views.py` 中定义 ViewSet；
3. 在 `serializers.py` 中定义 serializer.

## 3.1 创建 serializer
serializer 是什么？  
简单而言，serializer 就是一种根据配置，用来把数据在 model instance 及 raw data 之间转换的对象。

比如：我们可以针对刚刚定义的 model `Data` 来创建一个 serializer:
```python
# demo/serializers.py
from .models import Data
from rest_framework.serializers import ModelSerializer
 
 
class DataSerializer(ModelSerializer):
    class Meta:
        model = Data
        fields = ("id", "content")
```

然后，我们就可以通过这个 serializer 将数据在 `Data` 对象及 `JSON` 数据之间转换：

```python
from demo.serializers import Data, DataSerializer
from rest_framework.renderers import JSONRenderer   # 用于 JSON 渲染及解析
from rest_framework.parsers import JSONParser
from django.utils.six import BytesIO
 
instance = Data(content="test")
instance.save()
 
serializer = DataSerializer(instance)
JSONRenderer().render(serializer.data)              # b'{"id":1,"content":"test"}'
 
raw = b'{"id":2,"content":"another"}'
stream = BytesIO(raw)                               # 将 JSON 数据变成 Python dict
data = JSONParser().parse(stream)
 
serializer = DataSerializer(data=data)
ins = serializer.save()                             # <Data: Data object>
ins.__dict__                                        # {'content': 'another', 'id': 2}
```

## 3.2 创建 ViewSet
非常简单，在 `views.py` 中定义一个 ViewSet 类，标明对应的 model 以及 serializer 即可：
```python
# demo/views.py
from rest_framework.viewsets import ModelViewSet
from .models import Data
from .serializers import DataSerializer
 
 
class DataViewSet(ModelViewSet):
    queryset = Data.objects.all()
    serializer_class = DataSerializer
```

## 3.3 定义、配置并注册 router
```python
# demo/urls.py
from django.conf.urls import url, include
from rest_framework.routers import DefaultRouter
from .views import DataViewSet
 
router = DefaultRouter()                        # 定义 router
router.register('data', DataViewSet)            # 注册 viewset
 
urlpatterns = [
    url(r'^', include(router.urls)),            # 在 urlpatterns 里包含 router
]
```

# 4. 大功告成！

启动服务器，访问 `/data/`，你就会看到 DRF 精美的调试界面，真是感动:sob: 

按照 RESTful API 规范，在列表界面，你可以通过 POST 表单来创建新对象：  
![DRF 列表界面](/pics/DRF/DRF-list-view.jpg)

在详情界面（`/data/1/`），你可以通过 PUT 表单修改对象，或通过 DELETE 按钮来删除对象:  
![DRF 详情界面](/pics/DRF/DRF-detail-view.jpg)

当然，如果你乐意，使用 CURL 来调试也是可以的，服务器会给你返回 JSON 而不是 HTML:

```bash
$  curl -XPUT --data "content=edit" http://localhost:8000/data/1/
{"id":1,"content":"edit"}
```

# 5. 等等…别忘了写测试！
DRF 提供了一系列工具来协助 RESTful API 测试，比如 `rest_framework.test.APIClient`. `APIClient` 在发送请求时会根据配置自动设定 `Content-Type`, 而不是用 Django 自带 `djanto.test.Client` 的 `application/octet-stream`，方便了 API 测试。  
测试代码不在这里贴出（有些无趣），感兴趣的同学可以戳戳戳[这里](https://github.com/StdioA/drf_demo/blob/master/demo/tests.py)。

教程结束。  
这两天我还会继续看 DRF，后面还会继续更新相关内容。
toc: true

---

更新…  
然而并没有更新。  
我看完了 DRF 的 tutorial，也在某个项目里碰到了坑/难用的地方，不过这些都没有太多可写的地方，所以就算了吧，遇到问题翻文档就好，DRF 的文档写的还是挺不错的。
