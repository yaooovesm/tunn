# Tunn - NetworkTunnel

<br>

本项目旨在于实现高速、安全、高质量的跨网络通信。通过创建网络隧道实现为用户提供虚拟局域网环境。支持多平台多种传输协议以及加密方式。适用于简单的异地组网等场景。

<br>
<br>
<br>

### 特性

--------

#### 支持的操作系统

已测试的操作系统：

- Windows 7/10/11
- CentOS 7.x
- Ubuntu 20.x

理论支持：

- Windows 7+
- 支持tun设备的Linux发行版

#### 支持的传输协议

TCP / KCP / WS / WSS

#### 支持的加密方式

AES256 / AES192 / AES128 / XOR / SM4 / TEA / XTEA / Salsa20 / Blowfish

### 更新

------

2022/05/27 @ 1.0.0.220527

- WebUI
- 系统路由导入自动化
- 客户端配置操作

2022/05/10 @ 0.0.1

- 分离自项目 [Tunnel](https://gitee.com/jackrabbit872568318/tunnel)

2022/05/09 @ 历史更新

- 解决在NAT网路中的通信问题
- windows客户端支持
- 数据包CRC32校验
- ws/wss/tcp/kcp通信协议支持
- 路由暴露/导入
- 多通道传输
- 数据加密
- ...

### 编译

------

需要安装Go1.18.2或者更高版本 [下载](https://golang.google.cn/dl/)

准备

```shell
#拉取仓库
git clone https://gitee.com/jackrabbit872568318/tunn.git

#进入目录
cd ./tunn

#下载依赖
set GO111MODULE=on
go mod tidy

#进入cmd目录
cd cmd
```

编译

```shell
# @linux
go build -o tunn
```

或

```shell
# @windows
go build -o tunn.exe
```

### 使用

------

#### 客户端配置示例

[配置文件](../config/config_full.json)

```json lines
{
  "user": {
    //Hub用户
    "Account": "account",
    //Hub密码 (在设置密码时将会自动连接)
    "Password": "password"
  },
  "auth": {
    //认证服务器地址
    "Address": "aaa.bbb.ccc",
    //认证服务器端口
    "Port": 10241
  },
  "security": {
    //Hub证书
    "cert": "cert.pem"
  },
  "admin": {
    //控制台地址
    "address": "127.0.0.1",
    //控制台端口
    "port": 8080,
    //控制台用户
    "user": "admin",
    //控制台密码
    "password": "P@ssw0rd"
  }
}
```

#### 启动

启动参数

- -c 指定配置文件路径

示例：

```shell
# @linux
./tunn -c config.json
```

或

```shell
# @windows
tunn.exe -c config.json
```

启动成功如图

![img](./img/powershell_startup.png)

打开浏览器
![img](./img/admin_login.png)

进入控制台
![img](./img/admin_main.png)