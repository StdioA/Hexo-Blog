title: LAMP环境搭建心得
date: 2015-07-09 20:05:00
categories:
- 网络
tags:
- Linux
- Apache
- MySQL
- PHP
toc: true

---

闲来无事，在虚拟机上搭了一个LAMP服务器环境，把安装及配置过程记了下来。

<!-- more -->

## 1. 引子
### 1.1 环境版本
此次搭建的LAMP环境版本：

* Ubuntu 14.04 LTS
* Apache 2.4.7
* mysql 5.6.19
* php 5.5.9
    
### 1.2 写(hu)在(che)开(yi)头(tong)

额…其实没什么好说的，自己一直想自己动手搭建、配置一个服务器，暑假之前师太（别问是谁）说如果要搞安全的话最好先自己从头搭一个服务器，把各种服务弄清楚，对整个架构有一个系统的理解，这样再深入搞安全的话接受一些观念也会更快更容易；但是因为自己太懒，再加上上学期忙成狗（其实还是太懒），一直没有去做这件事。暑假在一个小公司做软件测试，每天好像也没什么事干，有大把的时间做自己的事情，于是自己用了一中午加半个下午的时间照着一份[指南](http://segmentfault.com/q/1010000002397754)把它搭好了。 _不过话说回来，软件测试真的很无聊_\_(:зゝ∠)\_
toc: true

---
## 2. 科普
[LAMP](https://zh.wikipedia.org/wiki/LAMP): [Linux](https://zh.wikipedia.org/wiki/Linux) + [Apache](https://zh.wikipedia.org/wiki/Apache_HTTP_Server) + [Mysql](https://zh.wikipedia.org/wiki/MySQL) + [PHP](https://zh.wikipedia.org/wiki/PHP)     
科普结束。   
刚看到LAMP里面的P还能指Python 0.0
toc: true

---
## 3. 搭建过程
### 3.1 安装   
1. Linux    
我手上现在没实体机了，只有一个树莓派，我也不想每天带着它去上班，何况AMD架构上面软件好像少一点点，更何况树莓派性能挺差的（此处省略一坨借口），所以我只用Virtual Box装了一个虚拟机。
说到Linux，选一个用起来比较舒服的的发行版还是挺重要的。Linux发行版众多，一般用Red Hat或者CentOS（RH的社区版）或者Ubuntu Server来做服务器，不过…这学期用Debian系发行版用习惯了，再换到RH系的感觉有点不适应，于是我选择了Ubuntu Server. 当然，如果你想锻炼一下，推荐使用Arch Linux来搭建服务器。     
下载Linux镜像，搭虚拟机，配置虚拟网络&SSH，更改软件源，更新软件，配置自己需要的vim & tmux & vsftpd，不多说，想详细了解的可以去看某Linux虚拟机安装及配置指南（代号[PA0](http://cslab.nju.edu.cn/ics/index.php/Ics:2013/PA0)）。上学期装Linux装了绝不下10遍，再说下去自己都要吐了。       
不过值得一提的是，Ubuntu Server安装程序的用户体验简直棒，安装过程中有一步是设置键盘布局，以前都要自己去一个长长的列表里翻自己的键盘布局（通常是US），而Ubuntu Server提供了一个小脚本来进行自动检测：依照提示敲几个字母/符号，再回答一个问题，安装程序会自动检测出适合你的键盘布局。

2. Apache
输入`apt-get install apache2`命令安装apache.
安装过程中apache服务已经启动，如果未启动，则输入`service apache2 start`启动apache服务。
启动后访问服务器ip，会出现apache的测试页面。    
![](/pics/apache.jpg)

3. MySQL
输入`apt-get install mysql-server-5.6 mysql-client-5.6`进行安装。
安装过程中需要输入MySQL root密码。 

4. PHP
输入`apt-get install php5 libapache2-mod-php5`安装php, 安装过后需要输入`service apache2 restart`重启apache服务。

### 3.2 配置
1. 安装phpMyAdmin    
输入`apt-get install phpmyadmin`进行安装，安装的时候会提示输入mysql的root密码，并且提示新建一个数据库，当然也可以按需求不新建。  
安装好以后访问http://localhost/phpmyadmin/index.php ，登录之后页面下方会有警告“缺少 mcrypt 扩展。请检查 PHP 配置。”此时按照指南的方法做没有效果，经百度+Google后找到了解决方案：安装php5-mcrypt后，更改php.ini后问题未解决，根据[官方的mcrypt安装指南](http://php.net/manual/en/mcrypt.setup.php)，输入`php5enmod mcrypt`后，问题解决。

2. MySQL命令行无法启动     
注：这段是自己瞎折腾的，啥都没看就乱玩遇到的问题。
输入mysql后遇到问题：
  ```
  ➜  ~  mysql
  ERROR 1045 (28000): Access denied for user 'stdio'@'localhost' (using password: NO)
  ```

  尝试用root权限运行，得到同样的结果。
  ```
  ➜  ~  sudo mysql
  ERROR 1045 (28000): Access denied for user 'root'@'localhost' (using password: NO)
  ```

  一番百度+google+SegmentFault后，找到正确进入命令行的姿势：
  ```
  ➜  ~  sudo mysql -u root -p
  Enter password:
  Welcome to the MySQL monitor.  Commands end with ; or \g.
  Your MySQL connection id is 55
  Server version: 5.6.19-0ubuntu0.14.04.1 (Ubuntu)
  
  Copyright (c) 2000, 2014, Oracle and/or its affiliates. All rights reserved.
  
  Oracle is a registered trademark of Oracle Corporation and/or its
  affiliates. Other names may be trademarks of their respective
  owners.
  
  Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
  
  mysql>
  ```
3. 设置Apache虚拟目录     
把所有文件全部放在/var/www下真的好麻烦，何况普通账户没有/var/www的写权限，所以设置一个alias，将某个Apache虚拟目录映射到home目录下，以后操作起来就会方便很多。
修改/etc/apache2/mods-enabled/alias.conf文件，添加如下行，然后重启Apache服务：
    ```
    Alias /web "/home/stdio/websites/"
    
    <Directory "/home/stdio/websites/">
        Options None 
        AllowOverride None
        Order allow,deny
        Allow from all
    </Directory>
    ```
    然而在我访问http://localhost/web时，却得到了503 Forbidden的状态码，各种乱访问无果，于是在网上乱搜解决方案，有让改httpd的（httpd跟Apache有啥关系），有改alias配置的（我的alias配置的没有问题啊），最后看到了一个方案，查看apache2.conf的目录权限配置。
    修改/etc/apache2/apache2.conf文件，发现以下设置：
    ```
    <Directory />
        Options FollowSymLinks
        AllowOverride None
        Require all denied
    </Directory>

    <Directory /usr/share>
        AllowOverride None
        Require all granted
    </Directory>
    
    <Directory /var/www/>
        Options Indexes FollowSymLinks
        AllowOverride None
        Require all granted
    </Directory>
    ```
    因为Apache的默认配置是不能访问/的，所以我没有对~/websites的访问权限（这里逻辑好混乱）。添加配置：
    ```
    <Directory /home/stdio/websites>
        Options FollowSymLinks
        AllowOverride None
        Require all granted
    </Directory>
    ```
    重启Apache服务，http://localhost/web目录下的文件均可正常访问。
    toc: true

---
## 4. 乱搞
去年开学的时候用php写过一个小的文件浏览器（就像Apache自带的文件服务器那样的），闲得无聊想把它部署到自己刚搭好的服务器上，看看能不能正常运行，于是就把文件传到服务器上访问，**不出意外，失败了**。然后就找呀找呀找bug，找到一个小bug（请自动脑补背景音乐），找了半个点最后发现，在从配置文档读取根目录路径的时候，会在目录结尾加一个空格（现在想起来觉得应该是^M）导致路径拼接时出错，于是在$rootpath前面加了trim，然后就好了…我真是能作\_(:зゝ∠)\_
toc: true

---
## 5. 总结
1. 自己动手搭建LAMP环境还是一件挺有意思的事情，遇到问题自己去找答案自己解决，最后所有的服务全都正常运行时还是有一点点成就感的~
2. 半年没碰PHP，一共就写了不到10行代码，还写错了一半，比如把`phpinfo()`写成`php_info`，忘了在`<?`后面加php神马的…（我记得以前谁跟我说`<?`后面可以不加php的啊）
3. 一篇文章写了一晚上。好久没写过博文了，写这篇文章主要是把自己的经验记下来，如果这篇文章可以帮到谁的话，那当然更好~
4. Linux挺好玩的，比软件测试好玩多了！（果然到了最后还是要黑一把测试）
toc: true

---
## 6. 参考文档
1. [ubuntu下搭建LAMP](http://segmentfault.com/a/1190000000619342)
2. [PHP Mcrypt Installing/Configuring](http://php.net/manual/en/mcrypt.setup.php)
3. [apache服务出现Forbidden 403问题的解决方法总结](http://www.douban.com/note/410696698/)
4. [ERROR 1045 (28000): Access denied for user 'root'@'localhost' (using password: NO)](http://segmentfault.com/q/1010000000263069)
