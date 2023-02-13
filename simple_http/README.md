# simple-http

## 简介

基于 Hertz Web 框架 + GORM 对象关系映射框架实现的简单 Web 单体应用，实现了抖音项目方案中的所有接口功能。项目中使用 MySQL 作为关系型数据库管理系统、MinIO 作为 OSS 对象存储服务，分别实现对用户数据存储和管理、用户上传视频文件和封面图片等资源分布式存储。

接口说明详见：https://bytedance.feishu.cn/docs/doccnKrCsU5Iac6eftnFBdsXTof#0lXbNY

注：本项目作为后续 RPC 微服务项目实现的参考。

## 依赖服务

**FFmpeg**：

- [下载地址](https://www.ffmpeg.org/download.html#build-windows)。
- 主要在 /publish/action/ 视频上传接口服务中，用于对用户上传视频进行封面截取并输出 .jpeg 格式封面图片。

**MySQL (v8.0.32)**：

- MySQL DNS 在 pkg/configs/db.go 中：

  ```go
  package configs
  const MySQLDataBaseDSN = "gorm:gorm@tcp(localhost:3308)/douyin?charset=utf8&parseTime=True&loc=Local&clientFoundRows=true"
  ```

- 库为 douyin，建库后运行 resources/sql/init.sql 文件，即可创建相关表结构。

**MinIO**：

- [下载地址](http://www.minio.io)，找教程自建即可（[国内网址](http://www.minio.org.cn/download.shtml#/kubernetes)）。

- 相关配置 （包括运行地址、用户名和密码等）在 pkg/configs/minio.go 中。
  注：使用自建的 MinIO 服务，除基本的用户名、密码等参数配置，还需要修改 ServerAddr 配置：

  ```go
  package configs
  const ServerAddr = "http://192.168.0.107:9000"
  ```

  修改为运行 MinIO 服务的主机网络 IP 和设定的端口，并且保证 App 端可以通过网络访问到该服务地址（同局域网环境下或公网环境）。

- 需要手动创建桶 (名称为 douyin)，并设置访问权限：

  ![image-20230213233523322](doc/img/image-20230213233523322.png)

  ![image-20230213233713305](doc/img/image-20230213233713305.png)

*Docker (v20.10.22)*：

- [下载地址](https://www.docker.com)。
- 容器中暂时只运行了 MySQL，可以不安装，直接使用本地 MySQL 服务即可。

注：上述依赖服务也可以直接使用我个人公网环境（服务器） 上的，需要修改 MySQL DNS 配置：
```go
package configs
const MySQLDataBaseDSN = "lortee:^LuoYi0813@tcp(rm-bp14n02e97n12gio0so.mysql.rds.aliyuncs.com:3306)/douyin?charset=utf8&parseTime=True&loc=Local&clientFoundRows=true"
```

MinIO 配置：

```go
package configs
const Endpoint = "114.55.42.185:9000"
const ServerAddr = "http://114.55.42.185:9000"
```
PS: 上述配置修改后，尽量避免将修改部分推送到代码仓库中。

## 打包命令

``` shell
cd simple_http/cmd
# 打包 linux 程序 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o simple-http-linux main.go
# 打包 mac 程序 
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o simple-http-mac main.go
# 打包 windows 程序 
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o simple-http-windows.exe main.go
```

## 运行

直接执行打包程序：

```shell
cd simple_http/cmd
./simple-http-linux
# ./simple-http-mac
# ./simple-http-windows.exe
```

保证手机或模拟器在同一局域网中，安装并配置 App 端 BaseUrl 地址为 http://本机网络 IP + 端口 (:8080) 地址即可；也可以使用已经部署在个人服务器上的项目， BaseUrl 地址为 http://114.55.42.185:8080。

最新版 App 在[使用文档](https://bytedance.feishu.cn/docs/doccnM9KkBAdyDhg8qaeGlIz7S7)中获取。