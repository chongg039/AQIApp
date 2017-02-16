## AQI app
进入`src`目录下执行`make`编译，编译后目录结构如下：
```go
~/AQIApp
	|--build		// 可执行文件
	|--log			// 日志文件
	|--old			// 旧版源码
	|--src			// 源码文件
	|--README.md	// README
```

源码包含两个主要程序：

```go
～/src
	|--crawler  // 爬虫引擎
	|--server   // API服务器
	|--Makefile // make脚本
```

## 爬虫引擎

使用的第三方包：

```go
"github.com/PuerkitoBio/goquery"
"github.com/go-sql-driver/mysql"
"github.com/robfig/cron"
```



~~爬虫引擎“每分钟”获取一次 [绿色呼吸](http://www.pm25.com/rank.html) 的空气质量数据存入mongoDB，以“小时”作为`collection`，城市及其数据作为单个字段存放，数据格式为“年月日天时”，例如名为“2017021021”的`collection`中存放的是“2017年2月10日21时”的全国天气数据（不包含县及县以下行政区划）。~~

2017年2月16日将服务从mongoDB迁移到Mysql，建立一个空气质量表统一存放：

```sql
/*
假设数据库为"testcity"，表名为"aqi"
建表语句中
	id：统一格式为"2017021622"，表示2017年2月16日22时
	city：城市名
	aqi：AQI数值
	time：统计时间
其中将"id"和"city"作为复合主键，允许单个字段的重复
*/
CREATE TABLE `aqi` (
  `id` int(11) NOT NULL,
  `city` varchar(50) NOT NULL DEFAULT '',
  `aqi` int(11) NOT NULL,
  `time` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`,`city`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

每小时更新数据为363条左右，若单条获取并写入数据耗时约9.85s，而同样单条写入mongoDB仅用时200ms（暂时不明白原因），非常影响性能，尤其是到百万级别数据写入更为如此。

解决方案是采用事务插入，获取数据后将数据写入事务统一插入，节省了每次发送连接和开关数据库的开销。优化后360余条数据插入耗时约70ms，从打开数据库开始计算时间不到500ms。

使用第三方定时模块[cron](https://github.com/robfig/cron)，检测当前时间的在每小时的第一分钟之内，定时每秒执行爬虫函数。使用[绿色呼吸](http://www.pm25.com/rank.html) 页面上的数据时间，检测数据库中是否有这个小时的数据，没有的话将数据存入，有的话不采取任何操作。

## API服务器

使用的第三方包：

```go
"github.com/gorilla/mux"
"github.com/go-sql-driver/mysql"
```



API服务器提供`JSON`格式的数据返回，请求示例：

**（！！以下所有城市名都需要汉字encoding转义）**

返回当前时间的`JSON`数据：

GET `http://localhost:8088/aqi/成都&now`

```json
{
    "City": "成都",
    "AQI": "117",
    "Time": "统计时间：2017-02-12 17:00"
}
```

返回特定时间的`JSON`数据：

GET `http://localhost:8088/aqi/成都&2017021104`

```json
// 2017021104:2017年2月11日4时
{
    "City": "成都",
    "AQI": "73",
    "Time": "统计时间：2017-02-11 04:00"
}
```

返回当天所有的`JSON`数据：

GET `http://localhost:8088/aqi/成都&today`：返回当天所有的`JSON`数据

```json
{
    "DataItems": [
        {
            "City": "成都",
            "AQI": "79",
            "Time": "统计时间：2017-02-12 00:00"
        },
        ......
        {
            "City": "成都",
            "AQI": "117",
            "Time": "统计时间：2017-02-12 16:00"
        },
        {
            "City": "成都",
            "AQI": "117",
            "Time": "统计时间：2017-02-12 17:00"
        }
    ]
}
```

出错，或数据库无数据时，返回错误示例：

GET `http://localhost:8088/aqi/成都&2017021001`

```json
{
    "code": 404,
    "text": "Not Found"
}
```

## 使用Supervisor守护进程

安装Supervisor，在`/etc/supervisord.conf`添加：

```sh
# 修改为自己的目录
[program:AQIServer]
command=/home/gc/AQIApp/build/server
autostart=true
autorestart=true
startsecs=10
redirect_stderr = true
stdout_logfile=/home/gc/AQIApp/log/server.log
stdout_logfile_maxbytes=20MB
stdout_logfile_backups=10
stdout_capture_maxbytes=1MB

[program:AQICrawler]
command=/home/gc/AQIApp/build/crawler
autostart=true
autorestart=true
startsecs=15
startretries = 3
redirect_stderr = true
stdout_logfile=/home/gc/AQIApp/log/crawler.log
stdout_logfile_maxbytes=20MB
stdout_logfile_backups=10
stdout_capture_maxbytes=1MB
```

启动`Supervisor`：

```shell
sudo supervisord -c /etc/supervisord.conf
```

使用`Supervisorctl`查看运行情况：

```shell
sudo supervisorctl -c /etc/supervisord.conf
# 更多 supervisorctl 操作请自行查阅
```

查看	`Supervisor`运行端口并关闭：

```shell
ps -ef | grep supervisor | grep -v grep
# 或通过`lsof -i :端口` 查询到pid
# kill -9 ${pid}
```