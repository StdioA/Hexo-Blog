title: 解决 MySQL 编码问题
date: 2016-05-09 09:40:46
categories:
- 开发
tags:
- MySQL
- Django
toc: true

---

很久以前写了一个 Django 项目，数据库用的 MySQL. 很久没用，后来发现 MySQL 编码设置有问题，导致中文全部变成了问号。

<!-- more -->

# 1. 设置 MySQL 服务端默认字符集
在 `my.cnf` 中设置默认字符集：
```
[mysqld]
character-set-server=utf8
```

同时，可以设置 MySQL 客户端的默认字符集：
```
[mysql]
character-set-default=utf-8
```

配置完成后，重启 MySQL 服务，输入 `SHOW VARIABLES LIKE '%CHAR';` 检查字符集是否正确。

# 2. 手动设置数据库编码
修改默认编码其实是没用的:joy:，还需要手动设置数据库编码，而数据库编码需要在建立时指定，所以…真是一个悲伤的故事:cry:（反正我的服务也没人用

重新建立数据库：
```
DROP DATABASE spaste;
CREATE DATABASE spaste DEFAULT CHARACTER SET utf8;
```

重建数据库后，输入 `python manager.py migrate` 重新迁移数据库。

# 3. 参考资料
1. [Configuring the Character Set and Collation for Applications - MySQL 5.7 Reference Manual](https://dev.mysql.com/doc/refman/5.7/en/charset-applications.html)
2. [django 解决 mysql 数据库输入中文乱码问题](http://shikezhi.com/html/2015/mysql_0827/204288.html)

# 4. 后记
好短的一篇文章…

最近在爬 B 站用户的公开用户数据，数据库用了 MongoDB, 爬完以后好好玩一玩:flushed:  
B 站基佬好多啊\_(:зゝ∠)\_
