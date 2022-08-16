# TimuShow 微服务项目 WEB

本项目使用 NaCos 作为配置中心，Consul作为服务中心，Raids 为手机验证码存储，Mysql8为用户等数据存储，接入aliyun短信服务，实现简易的分布式微服务项目，各子项目配置信息如下:

## NaCos 存储的配置文件信息

user-web.json
```json
{
    "name": "user-web",
    "host":"192.168.0.162",
    "tags":["user","web"],
    "port": 8021,
    "user-srv": {
        "name": "user-srv"
    },
    "jwt": {
        "key": "加密字符串"
    },
    "sms": {
       "key": "阿里云的key",
       "secrect": "阿里云的secrect"
    },
    "redis": {
        "host": "127.0.0.1",
        "port": 6379,
        "expire": 300
    },
    "consul": {
        "host": "127.0.0.1",
        "port": 8500
    }
}
```

goods-web.json
```json
{
    "name": "goods-web",
    "port": 8022,
    "host":"192.168.0.162",
    "tags":["goods","web"],
    "goods-srv": {
        "name": "goods-srv"
    },
    "jwt": {
        "key": "加密字符串"
    },
    "consul": {
        "host": "127.0.0.1",
        "port": 8500
    }
}
```