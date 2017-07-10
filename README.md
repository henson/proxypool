[![Travis Status for henson/ProxyPool](https://travis-ci.org/henson/ProxyPool.svg?branch=master)](https://travis-ci.org/henson/ProxyPool) [![Go Report Card](https://goreportcard.com/badge/github.com/henson/ProxyPool)](https://goreportcard.com/report/github.com/henson/ProxyPool)

# Golang实现的IP代理池

> 采集免费的代理资源为爬虫提供有效的代理


### 1、代理池设计

　　代理池由四部分组成：

* Getter：

　　代理获取接口，目前有**9**个免费代理源，每调用一次就会抓取这些网站最新的100个代理放入Channel，可自行[添加额外的代理获取接口](#4添加自定义代理采集接口)；

* Channel：

　　临时存放采集来的代理，通过访问稳定的网站去验证代理的有效性，有效则存入数据库；

* Schedule：

　　用定时的计划任务去检测数据库中代理IP的可用性，删除不可用的代理。同时也会主动通过Getter去获取最新代理；

* Api：

　　代理池的访问接口，提供get接口输出JSON，方便爬虫直接使用。

### 2、代码实现

* Api：

　　api接口相关代码，提供`get`接口，输出JSON；

* Storage：

　　数据库相关代码，数据库采用Mongo；

* Getter：

　　代理获取接口，目前抓取这九个网站的免费代理，当然也支持自己扩展代理接口；

1. [快代理](http://www.kuaidaili.com)
2. [代理66](http://www.66ip.cn)
3. ~~[IP181](http://www.ip181.com)~~
4. ~~[有代理](http://www.youdaili.net/Daili/http/)~~
5. [西刺代理](http://www.xicidaili.com/nn/)
6. [guobanjia](http://www.goubanjia.com/free/gngn/index)
7. [讯代理](http://www.xdaili.cn/freeproxy.html)
8. [无忧代理](http://www.data5u.com/free/index.shtml)
9. [Proxylist+](https://list.proxylistplus.com)

* Schedule：

　　定时任务，目前在main.go中以轮询方式实现，后期会改进；

* Util：

　　存放一些公共的模块、方法或函数，包含`Config`:读取配置文件config.json；

* 其他文件：

　　配置文件:config.json，数据库配置和代理获取接口配置；

```
{
    "mongo": {
        "addr": "mongodb://127.0.0.1:27017/",
        "db": "temp",
        "table": "pool"
    },
    "host": ":8080"
}
```

### 3、安装及使用

因为有些代理网站使用了加密页面、混淆代码等反爬技术，要正确采集到代理数据得用到 [PhantomJS](http://phantomjs.org/) ，必须提前先装好。

另外，本项目用到的依赖库有：
```
gopkg.in/mgo.v2
github.com/PuerkitoBio/goquery
github.com/parnurzeal/gorequest
github.com/nladuo/go-phantomjs-fetcher
```

下载本项目：
```
go get -u github.com/henson/ProxyPool
```

然后配置好相应的config.json并启动：
```
go build
./ProxyPool
```

随机输出可用的代理：
```
GET http://localhost:8080/v1/ip
```
![HTTP](pics/http.png)

随机输出HTTPS代理：
```
GET http://localhost:8080/v1/https
```
![HTTPS](pics/https.png)

### 4、添加自定义代理采集接口

其实很简单，只需要在getter包下新增一个采集函数（如例子的Data5u()），甚至可以不需要新建一个go文件（新建文件是为了方便归档采集函数，如例子5u.go）。

```
// 5u.go
// Data5u get ip from data5u.com
func Data5u() (result []*models.IP) {
	//处理逻辑
	...
	log.Println("Data5u done.")
	return
}
```

然后在main.go的run函数中添加、删除或注释掉该采集函数的调用即可。

```
func run(ipChan chan<- *models.IP) {
	var wg sync.WaitGroup
	funs := []func() []*models.IP{
		getter.Data5u,
		getter.IP66,
		getter.KDL,
		getter.GBJ,
		getter.Xici,
		getter.XDL,
		//getter.IP181,
		//getter.YDL,
		getter.PLP,
	}
    ...
}
```

### 5、异常恢复

之前，偶尔会有朋友跟我反映程序无法编译，经过检查发现都是代理网站发生了变化（或修改了页面或关闭了网站），以致于采集程序原先设计的爬虫不能正常工作而导致了错误的发生。为此，我修改了代码，加入了容错机制，即便爬虫出错了也不会影响到主体程序的运行。出错的采集进程会被主线程忽略，其它正常的采集进程仍将继续工作。

### 6、诚挚的感谢

- 首先感谢您的使用，如果觉得程序还不错也能帮助您解决实际问题，不妨添个赞以鼓励本人继续努力，谢谢！
- 如果您对程序有任何建议和意见，也欢迎提交issue。
- 当然，如果您愿意贡献代码和我一起改进本程序，那再好不过了。