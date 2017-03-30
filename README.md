
> Golang的代理池服务，采集免费的代理资源为爬虫提供有效的代理。

### 1、代理池设计

　　代理池由四部分组成：

* Getter：

　　代理获取接口，目前有6个免费代理源，每调用一次就会抓取这个6个网站最新的100个代理放入Channel，可自行添加额外的代理获取接口；

* Channel：

　　临时存放采集来的代理，通过访问稳定的网站去验证代理的有效性，有效则并存入数据库；

* Schedule：

　　用定时的计划任务去检测数据库中代理IP的可用性，删除不可用的代理。同时也会主动通过Getter去获取最新代理；

* Api：

　　代理池的访问接口，提供get接口输出JSON，方便爬虫直接使用。
<!--#### 功能图纸-->
![设计](https://pic2.zhimg.com/v2-f2756da2986aa8a8cab1f9562a115b55_b.png)

### 2、代码实现

* Api：
　　api接口相关代码，提供`get`接口，输出JSON；

* Storage：
　　数据库相关代码，数据库采用Mongo；

* Getter：
　　代理获取的相关代码，目前抓取：[快代理](http://www.kuaidaili.com)、[代理66](http://www.66ip.cn)、[IP181](http://www.ip181.com)、[有代理](http://www.youdaili.net/Daili/http/)、[西刺代理](http://www.xicidaili.com/nn/)、[guobanjia](http://www.goubanjia.com/free/gngn/index)这个六个网站的免费代理，经测试这些网站每天更新的可用代理只有六七十个，当然也支持自己扩展代理接口；

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
        "table": "pool",
        "event": "event"
    },
    "host": ":8080"
}
```

### 3、安装及使用

下载代码：
```
go get -u github.com/henson/ProxyPool

```

启动：
```
配置好相应的config.json

go build

./ProxyPool
```

使用：
```
访问：http://localhost:8080/v1/ip
```