title: 随手记之Vue.js
date: 2016-03-28 21:40:03
categories:
- Javascript
tags:
- Javascript
- node.js
- vue.js
- 前端开发

---

之前的那门 JS 和前端课结课了，最后花了一周时间写了一个 project，使用 Vue.js 的一套工具写了一个单页面应用，先将使用到及学到的知识来整理一下。

<!-- more -->

想了想，之前自己学习 `Vue.js` 的时候一直在翻文档，翻到哪些让人眼前一亮的功能时就将它加进自己的项目里，并没有什么系统地进行学习，所以想到哪写到哪好了。前方高乱预警。

# 1. Vue.js
之前用 `React` 用得有些不爽了所以想换换口味, 于是 用[`Vue.js`](http://vuejs.org.cn/) 写了一个小应用的前端部分。它的轻量以及MVVM架构让我在两天之内爱上了它，于是决定用它代替 React 去完成 final project 的前端部分。

不知道为什么，我觉得 `Vue.js` 比 `React` 更容易上手，所以很快就学会了它。想了想确实不知道该整理一些什么，因为 [Vue 的官方教程](http://vuejs.org.cn/guide/)写的确实很直白清楚，所以在此不再赘述，有时间的话可能会考虑写一个 Vue 的教程。

# 2. vue-cli
[vue-cli](https://github.com/vuejs/vue-cli) 是一个 Vue.js 官方提供的脚手架工具，你可以使用它来轻松地构建出一个应用 `Vue.js` 的工程。官方提供 5 种模板，当然你也可以在 `github` 上或者本地构建自己的模板然后使用 vue-cli 生成工程。

[具体使用方法](https://github.com/vuejs/vue-cli#usage)：生成一个工程极为简单，只需一两条命令即可。

```bash
$ vue init webpack my-project     // 使用 webpack 模板生成工程 my-project
$ cd my-project
$ npm install
$ npm run dev
```
然后，借助 `webpack` 与 `vue-loader`, 你可以将一个 `Vue` 组件的模板、核心js代码和CSS写在**一个文件里**，甚至还可以使用 `CSS` 与 `HTML` 的预处理器，像这样：  
![](/pics/vue/vue-component-with-pre-processors.png)  
<span style="font-size: 0.75em;">图片来自[ Vue.js 官方教程 - 构建大型应用](http://vuejs.org.cn/guide/application.html)。</span>

# 3. Vue-router
[`vue-router`](https://github.com/vuejs/vue-router/) 是 `vue` 官方提供的路由模块，可以实现 `SPA` 中的路由操作。具体文档看[这里](http://vuejs.github.io/vue-router/zh-cn/index.html)。顺便提一句，`vue-router` 的中文文档已经过时了。

## 3.1 初始化
首先使用 `npm` 安装 `vue-router`, 然后在程序入口点配置 `vue-router`.

```javascript
import Vue from 'vue'
import Router from 'vue-router'
import App from './components/App'                    // 程序的核心 Vue 应用
import HomePageView from './components/HomePageView'  // 导入所有的 View 组件
import ItemView from './components/ItemView'
 
Vue.use(Router)                                       // Vue 配置
var router = new Router()                             // 生成路由对象
router.map({                                          // 配置路由
  '/': {
    name: 'homepage',
    component: HomePageView
  },
  '/home': {
    name: 'homepage',
    component: HomePageView
  },
  '/item/:id': {                                    // 支持动态路径
    name: 'item',
    component: ItemView
  }
}
 
router.redirect({                                   // 设置重定向选项
  '*': '/home'
})
 
router.start(App, '[app]')                           // 挂载 Vue 主应用
```

然后在 `App.vue` 的 template 中设置 `router-view`.

```html
<template>
  <router-view
    keep-alive
    transition="fade"
    transition-mode="out-in">
  </router-view>
</template>
```

## 3.2 vue-router 的使用
`Vue Router` 对象被嵌入到每个 `vue` 组件中，可以在组件中调用 `this.$router` 来控制 router 对象，如进行页面跳转等。

此外，还可以在页面切换时在组件的 `route` 配置中使用路由切换钩子控制 `vue-router`，详情请看[文档](http://vuejs.github.io/vue-router/zh-cn/pipeline/index.html)

# 4. Vuex
[`Vuex`](https://github.com/vuejs/vuex/) 是一个借鉴于 `Flux`，但是专门为 `Vue.js` 所设计的状态管理方案。`Flux` 采用了 `Action → Dispatcher → Store → View` 的状态管理机制，而 `Vuex` 跟它差不多：Vue 组件调用 action，action dispatches mutation, mutation 改变 store 中的 state，state 改变 View. 下面是 `Vuex` 的数据流。

![Vuex Data Flow](/pics/vue/vuex.png)

## 4.1 使用方法
程序入口点：
```js
import Vuex from 'vuex'

Vue.use(Vuex)
```

主应用：
```js
import store from '../store'

export default {
  store,                      // 引用store
  ...
}
```

其它组件：
```javascript
export default {
  ...,
  vuex: {
    // 定义 getter 从 store 中获取 state 并注册至应用中
    getters: {
      logged_in: function (state) {
        return state.user.logged_in
      }
    },
    // 定义 action, 组件可在自己的函数中调用 action 来 dispatch mutations.
    actions: {
      login: ({ dispatch }, user) => {
        dispatch('LOGIN', user)
      }
    }
  },
  ...
}
```

关于 store 及 mutation 的定义方式，请参考[ Vuex 文档](http://vuejs.github.io/vuex/en/tutorial.html)。

# 5. 总结 & 后记
跟 React 相比，个人感觉 Vue 要更容易上手，易于使用，文档也很清楚（比 Hexo 的高到不知道那里去了，前两天被 Hexo 整疯了必须要黑一下它）；Vue 的一系列工具也很易于使用，与 Vue 整合度高，可以在组件中方便地进行操作。

前端课结课，准备退坑了，过几天可能会学习并整理一些后端的知识。
