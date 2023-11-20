# socks5proxy

在本地监听端口，所有发送到该端口的socks5流量都会被通过代理转发，现适配51代理及通过fofa搜索代理池。

## 使用方法

首先在config.yaml中配置fofa的email和key以及51代理的地址

使用51代理

```
./socks5proxy -51
```

使用fofa代理

```
./socks5proxy -fofa
```

使用fofa进行代理前需要使用参数fofaup更新fofa代理池

```
./socks5proxy -fofaup
```

设置更换ip时间（秒）

```
./socks5proxy -51 -t 10
```

设置监听端口

```
./socks5proxy -51 -t 10 -p 8081
```

