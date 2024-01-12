# 项目介绍

基于Vue3+ElementPlus+Gin+Gorm编写的一款管理平台

# 后端启动

进入main.go文件，执行main函数即可。

# 后端项目地址

github：https://github.com/Lvan826199/xz-go-frame

gitee：https://gitee.com/xiaozai-van-liu/xz-go-frame

# 前端项目地址


# 后端依赖介绍

可以配置相关国内代理

```shell
##### goproxy
https://goproxy.io/zh/
# 配置
GOPROXY=https://goproxy.io/zh/,direct


##### 七牛云
https://goproxy.cn
# 配置
GOPROXY=https://goproxy.cn,direct


##### 阿里云
https://mirrors.aliyun.com/goproxy/
# 配置
GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```



可以直接根据go.mod进行下载

```shell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
go mod download
```

## 用到的技术

**gin框架**

```shell
go get -u github.com/gin-gonic/gin
```

**session库**

```shell
go get -u github.com/gin-contrib/sessions
```

**viper配置管理库**

```shell
go get github.com/spf13/viper
```
**GORM框架+MySQL驱动**

```shell
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

**验证码模块**

官网 : https://github.com/mojocn/base64Captcha

```shell
go get github.com/mojocn/base64Captcha
```

**JWT模块**

官网：https://github.com/golang-jwt/jwt

```shell
go get -u github.com/golang-jwt/jwt/v5
```

**防缓存击穿singleflight**

```shell
go get golang.org/x/sync/singleflightGroup
```

# 前端依赖安装


# 其他依赖安装
