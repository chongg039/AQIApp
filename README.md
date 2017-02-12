## AQI app
进入`src`目录下执行`make`编译，编译后目录结构如下：
```go
~/AQIApp
	|--build		// 可执行文件
	|--log			// 日志文件
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

#### 爬虫

爬虫引擎“每分钟”获取一次 [绿色呼吸](http://www.pm25.com/rank.html) 的空气质量数据存入mongodb，以“天”为`collection`，城市及其数据作为单个字段存放，数据格式为“年月日天时”，例如名为“2017021021”的`collection`中存放的是“2017年2月10日21时”的全国天气数据（不包含县及县以下行政区划）。

#### API服务器

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

#### 使用Supervisor守护进程

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