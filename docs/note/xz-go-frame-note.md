# 架构分析

 使用gin作为web服务，使用vue作为前端服务，来构建一个admin管理后台系统。

**模板文件**

Goland创建go文件添加作者信息：

- File -> Settings -> Editor -> File and Code Templates -> Go File

- 右侧填入模板

- ```go
  /*
  * @Author: 梦无矶小仔
  * @Date:   ${DATE} ${TIME}
  */
  package ${GO_PACKAGE_NAME}
  ```

## web框架Gin

安装gin的组件

```go
go get -u github.com/gin-gonic/gin
```

定义go的初始化

```go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	//  初始化 gin服务
	rootRouter := gin.Default()
	// 启动HTTP服务,可以修改端口
	address := fmt.Sprintf(":%d", 8088)
	rootRouter.Run(address)
}

```



## 初始化-1

main.go中进行初始化

```go
package main

import (
   "fmt"
   "xz-go-frame/initlization"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
   // 开始初始化gin路由服务
   initlization.WebRouterInit()
   fmt.Println("启动xz-go-frame后端成功")
}
```

## 路由-1

新建initlization包，创建router.go

```go
package initlization

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WebRouterInit() {
	// 初始化 gin 服务
	rootRouter := gin.Default()

	// 登录路由
	rootRouter.GET("/login", func(context *gin.Context) {
		context.JSON(http.StatusOK,"我是gin")
	})

	// 启动HTTP服务
	rootRouter.Run("127.0.0.1:8888")
}

```



启动之后可以直接访问`127.0.0.1:8888/login`，就能访问你的页面了。

![image-20231225141313328](images/image-20231225141313328.png)



**接下来开始简单的一些封装。**

## 封装api

新建api包

```go
api
	- admin
		-- index.go
	- bbs
		-- bbs.go
	- login
		-- login.go
	- video
		-- video.go
```

### index.go

```go
package admin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminApi struct {
}

// 登录处理的逻辑

func (e *AdminApi) Index(c *gin.Context) {
	// 获取会话
	session := sessions.Default(c)
	// 获取登录用户信息
	user := session.Get("user")
	username := user.(string)
	c.JSON(http.StatusOK, "我是gin"+username)

}

```

### bbs.go

```go
package bbs

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BbsApi struct {
}

// 首页处理
func (e *BbsApi) BbsIndex(c *gin.Context) {
	username, _ := c.Get("username")
	// 可以获取到的login放入session的数据
	c.JSON(http.StatusOK, "我是bbs的首页："+username.(string))
}

// 获取明细
func (e *BbsApi) GetBbsDetailById(c *gin.Context) {
	username, _ := c.Get("username")
	// 可以获取到login放入session的数据
	param := c.Param("id")
	c.JSON(http.StatusOK,"我是bbs的名字，参数："+param+" : "+username.(string))
}

```



### login.go

```go
package login

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginApi struct {
}

func (e *LoginApi) Login(c *gin.Context) {
	// session ---- 是一种所有请求之间的数据共享机制，为什么会出现session，是因为http请求是一种无状态。
	// 什么叫无状态：就是指，用户在浏览器输入方位地址的时候，地址请求到服务区，到响应服务，并不会存储任何数据在客户端或者服务端，
	// 也是就：一次request---response就意味着内存消亡，也就以为整个过程请求和响应过程结束。
	// 但是往往在开发中，我们可能要存存储一些信息，让各个请求之间进行共享。所有就出现了session会话机制
	// session会话机制其实是一种服务端存储技术，底层原理是一个map
	// 比如：我登录的时候，要把用户信息存储session中，然后给 map[key]any =
	// key = sdf365454klsdflsd --sessionid

	// 初始化session对象
	session := sessions.Default(c)
	// 存放用户信息到session中
	session.Set("user", "xiaozai") // map[sessionid] == map[user][xiaozai]
	// 记住一定要调用save方法，否则内存不会写入进去
	session.Save()
	username := session.Get("user")
	c.JSON(http.StatusOK, "我是gin:登录用户名："+username.(string))
}

```

### video.go

```go
package video

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Video struct {
}

// 首页处理
func (e *Video) VideoIndex(c *gin.Context) {
	username, _ := c.Get("username")
	// 可以获取到login放入session的数据
	c.JSON(http.StatusOK, "我是video的首页: "+username.(string))
}

// 获取明细

func (e *Video) GetVideoDetailById(c *gin.Context) {
	username, _ := c.Get("username")
	// 可以获取到login放入session的数据
	param := c.Param("id")
	c.JSON(http.StatusOK, "我是VideoApi的名字,参数:"+param+"  ： "+username.(string))
}

```

## 抽离router

新建router包

```go
router
	- bbs
		-- bbs.go
	- video
		-- video.go
```



### bbs.go

```go
package bbs

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/api/bbs"
)

type BbsRouter struct {
}

func (bBsRouter BbsRouter) InitBBSRouter(group *gin.RouterGroup) {
	// 帖子路由
	bbsApi := bbs.BbsApi{}
	bbsApiGroup := group.Group("bbs")
	{
		// 函数封装
		bbsApiGroup.GET("/indexx", bbsApi.BbsIndex)
		bbsApiGroup.GET("/get/:id", bbsApi.GetBbsDetailById)
	}
}

```

### video.go

```go
package video

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/api/video"
)

type VideoRouter struct {
}

func (videoRouter VideoRouter) InitVideoRouter(group *gin.RouterGroup) {
	// 帖子路由
	videoApi := video.Video{}
	videoGroup := group.Group("video")
	{
		videoGroup.GET("/index", videoApi.VideoIndex)
		videoGroup.GET("/get/:id", videoApi.GetVideoDetailById)
	}
}

```



## 封装中间件

新建middle包

```go
middle
	- loginfilter.go
```

### loginfilter.go

```go
package middle

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取会话
		session := sessions.Default(c)
		// 获取用户登录信息
		user := session.Get("user")
		// 如果用户没有登录，直接重定向返回登录
		if user == nil {
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort() // 拦截，到这里就不会往下执行了
		}
		// 取出会话信息
		username := user.(string)
		// 把session用户信息，放入到context中，给后续路由进行使用
		// 好处就是： router中方法不需要再次获取session再来拿会话中的信息
		c.Set("username", username)
		c.Next() //放行，默认就会放行
	}
}

```



## 封装路由

initlization -> router.go

```go
package initlization

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"xz-go-frame/api/login"
	"xz-go-frame/middle"
	"xz-go-frame/router/bbs"
	"xz-go-frame/router/video"
)

func WebRouterInit() {
	// 初始化 gin 服务
	rootRouter := gin.Default()
	
	// 创建cookie存储
	store := cookie.NewStore([]byte("secret"))
	// 路由上加入session中间件
	rootRouter.Use(sessions.Sessions("mysession",store))
	
	bbsRouter := bbs.BbsRouter{}
	videoRouter := video.VideoRouter{}
	loginApi := login.LoginApi{}
	

	// 登录路由
	rootRouter.GET("/login", loginApi.Login)
	
	// 首页路由
	AdminGroup := rootRouter.Group("admin")
	AdminGroup.Use(middle.LoginInterceptor())
	{
		bbsRouter.InitBBSRouter(AdminGroup)
		videoRouter.InitVideoRouter(AdminGroup)
	}

	// 启动HTTP服务,可以修改端口
	address := fmt.Sprintf(":%d", 8088)
	rootRouter.Run(address)
}

```



到此，就拥有了一个session登录拦截中间件请求的服务。

```shell
1、直接访问 http://127.0.0.1:8088/admin/video/index
会跳转至登录 http://127.0.0.1:8088/login
2、登录成功后再次访问才是正常 http://127.0.0.1:8088/admin/video/index
3、http://127.0.0.1:8088/admin/video/get/100
显示"我是VideoApi的名字,参数:100  ： xiaozai"
```

# 关于项目中的配置viper

详见知识点viper库

```go
configfile
	- application.yaml
```

## 1、下载viper
```go
go get github.com/spf13/viper
```
## 2、编写一个yaml的配置文件application.yaml
```yaml
# 数据配置的设定 tab对齐，shift+tab 回退对齐
mysql:
  database:
    host: 127.0.0.1
    user: root
    dbname: test
    pwd: 123456
# 服务组配置信息
server:
  port: 8088
  cookiname: mysession
```
## 3、基本读取

在initlization包下新建`init_viper.go`文件。

```go
/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 14:43
 */
package initlization

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func InitViper() {
	// 获取项目的执行路径
	path, err := os.Getwd() // D:\Z_Enviroment\GoWorks\src\xz-go-frame
	if err != nil {
		panic(err)
	}
	// 初始化一个viper解析配置对象
	config := viper.New()
	// 开始设置从哪个目录下去找yaml文件
	config.AddConfigPath(path + "/configfile") // 设置读取的文件路径
	// 设置配置文件的名字
	config.SetConfigName("application") // 设置读取的文件名
	// 设置配置文件的后缀
	config.SetConfigType("yaml") // 设置文件类型

	// 尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	// 监控配置文件的变化
	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件发生变化：", e.Name)
		//// 把改变的值重新放入到config配置中去
		//if err = config.Unmarshal(&config); err!= nil {
		//	fmt.Println(err)
		//}
	})

	// 打印文件读取出来的内容
	fmt.Println(len(config.Get("mysql.database.host").(string))) // 9
	fmt.Println(config.Get("mysql.database.host")) // 127.0.0.1
	fmt.Println(config.Get("mysql.database.user")) // root
	fmt.Println(config.Get("mysql.database.dbname")) // test
	fmt.Println(config.Get("mysql.database.pwd")) // 123456
	fmt.Println(config.Get("server.port")) // 8888
	fmt.Println(config.Get("server.cookiname")) // mysession

}

```

测试代码，修改为如下代码进行测试配置文件的读取。

```go
func main() {
	//  开始初始化配置文件
	initlization.InitViper()
	fmt.Println("初始化配置文件成功！")
	// 开始初始化gin路由服务
	//initlization.WebRouterInit()
	//fmt.Println("启动xz-go-frame后端成功")
}

```



进一步封装配置见代码

# 全局global属性配置

配置文件`application.yaml`

```yaml
# 服务端口的配置
server:
  port: 8989
  context: /
# 数据库配置
# "root:123456@tcp(127.0.0.1:3306)/xz-go-frame-db?charset=utf8&parseTime=True&loc=Local", // DSN data source name
database:
  mysql:
    host: 127.0.0.1
    port: 3306
    dbname: xz-go-frame-db
    username: root
    password: 123456
    config: charset=utf8&parseTime=True&loc=Local
# nosql数据的配置
nosql:
  redis:
    host: 127.0.0.1
    port: 3306
    password: 123
    db: 0
  es:
    host: 127.0.0.1
    port: 9300
    password: 456
ksd:
  alipay:
    appid: 12454545
    screat: 45587
    paths: 1,2,3
    detail:
      id: 100
      name: xiaozai
    map:
      name: xiaozai
      age: 18
      phone: 3245646
    urls:
      - 1
      - 2
      - 3
    routers:
      - id: 100
        url: /user/
        filter: 1
      - id: 200
        url: /video/
        filter: 2

```

1、新建global包，在下面新建global.go文件

```go
/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 16:29
 */
package global

import "gorm.io/gorm"

var (
	Yaml map[string]any
	Config *parse.Config
)
```



2、新建commons包，在下面新建相关配置包和文件进行属性管理。

```go
commons
	- parse
		- Config.go
		- Database.go
		- Ksd.go
		- NoSql.go
```

**Config.go**

```go
package parse

// 配置
type Config struct {
	// 数据库
	Database Database `mapstructure:"database" json:"database" yaml:"database"`
	Ksd      Ksd      `mapstructure:"ksd" json:"ksd" yaml:"ksd"`
	NoSQL    NoSQL    `mapstructure:"nosql" json:"nosql" yaml:"nosql"`
}

```

**Database.go**

```go
/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 17:02
 */
package parse

// Database 配置展示
// database:--------------------------------struct
//
//	mysql: --------------------------------struct
//		host: 127.0.0.1 -------------------field
//		port: 3306 -------------------field
//		dbname: xz-go-frame-db -------------------field
//		username: root -------------------field
//		password: 123456 -------------------field
//		config: charset=utf8&parseTime=True&loc=Local -------------------field
type Database struct {
	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
}

type Mysql struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string    `mapstructure:"port" json:"port" yaml:"port"`
	Dbname   string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Config   string `mapstructure:"config" json:"config" yaml:"config"`
}

```

**Ksd.go**

```go
/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 17:03
 */
package parse

// Ksd 数据格式展示
/*
ksd:
  alipay:
    appid: 12454545
    screat: 45587
    paths: 1,2,3
    detail:
      id: 100
      name: feige
    map:
      name: feige
      age: 30
      phone: 3245646
    urls:
      - 1
      - 2
      - 3
    routers:
      - id: 100
        url: /user/
        filter: 1
      - id: 200
        url: /video/
        filter: 2


*/
type Ksd struct {
	Alipay Alipay `mapstructure:"alipay" json:"alipay" yaml:"alipay"`
}

type Alipay struct {
	Appid   string         `mapstructure:"appid" json:"appid" yaml:"appid"`
	Screat  string         `mapstructure:"screat" json:"screat" yaml:"screat"`
	URLS    []string       `mapstructure:"urls" json:"urls" yaml:"urls"`
	Paths   []string       `mapstructure:"paths" json:"path" yaml:"paths"`
	Routers []Router       `mapstructure:"routers" json:"routers" yaml:"routers"`
	Detail  Detail         `mapstructure:"detail" json:"detail" yaml:"detail"`
	Map     map[string]any `mapstructure:"map" json:"map" yaml:"map"`
}

type Detail struct {
	Id   int    `mapstructure:"id" json:"id" yaml:"id"`
	Name string `mapstructure:"name" json:"name" yaml:"name"`
}

type Router struct {
	Id     int    `mapstructure:"id" json:"id" yaml:"id"`
	Url    string `mapstructure:"url" json:"url" yaml:"url"`
	Filter string `mapstructure:"filter" json:"filter" yaml:"filter"`
}

```

**NoSql.go**

```go
/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 17:03
 */
package parse

// NoSQL # nosql数据的配置
// nosql:
//
//	redis:
//		host: 127.0.0.1
//		port: 3306
//		password:
//		db: 0
type NoSQL struct {
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

```



修改`init_viper.go`文件的代码，当配置文件进行修改时，执行回调函数实时更改配置信息

```go
/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 14:43
 */
package initlization

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"xz-go-frame/global"
)

func InitViper() {
	// 获取项目的执行路径
	path, err := os.Getwd() // D:\Z_Enviroment\GoWorks\src\xz-go-frame
	if err != nil {
		panic(err)
	}
	// 初始化一个viper解析配置对象
	config := viper.New()
	// 开始设置从哪个目录下去找yaml文件
	config.AddConfigPath(path + "/configfile") // 设置读取的文件路径
	// 设置配置文件的名字
	config.SetConfigName("application") // 设置读取的文件名
	// 设置配置文件的后缀
	config.SetConfigType("yaml") // 设置文件类型

	// 尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	// 监控配置文件的变化
	config.WatchConfig()
	// 当配置文件发生变化时，执行回调函数
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件发生变化：", e.Name)
		// 把改变的值重新放入到config配置中去
		if err = config.Unmarshal(&global.Config); err != nil {
			fmt.Println(err)
		}
	})

	// 这里才是把yaml配置文件解析放入到Config对象的过程---map---config
	if err = config.Unmarshal(&global.Config); err != nil {
		fmt.Println(err)
	}

	// 遍历打印文件读取来的内容
	keys := config.AllKeys()
	dataMap := make(map[string]any)
	for _, key := range keys {
		fmt.Println("yaml存在的key是: "+key+" : ", config.Get(key))
		dataMap[key] = config.Get(key)
	}

	global.Yaml = dataMap


}

```



输出结果

```shell
yaml存在的key是: database.mysql.dbname :  xz-go-frame-db
yaml存在的key是: nosql.redis.host :  127.0.0.1
yaml存在的key是: ksd.alipay.paths :  1,2,3
yaml存在的key是: ksd.alipay.map.phone :  3245646
yaml存在的key是: database.mysql.host :  127.0.0.1
yaml存在的key是: database.mysql.username :  root
yaml存在的key是: nosql.redis.password :  123
yaml存在的key是: ksd.alipay.map.name :  xiaozai
yaml存在的key是: ksd.alipay.urls :  [1 2 3]
yaml存在的key是: nosql.es.port :  9300
yaml存在的key是: ksd.alipay.routers :  [map[filter:1 id:100 url:/user/] map[filter:2 id:200 url:/video/]]
yaml存在的key是: server.port :  8989
yaml存在的key是: server.context :  /
yaml存在的key是: database.mysql.password :  123456
yaml存在的key是: database.mysql.config :  charset=utf8&parseTime=True&loc=Local
yaml存在的key是: nosql.redis.db :  0
yaml存在的key是: nosql.redis.port :  3306
yaml存在的key是: ksd.alipay.screat :  45587
yaml存在的key是: ksd.alipay.detail.id :  100
yaml存在的key是: ksd.alipay.detail.name :  xiaozai
yaml存在的key是: ksd.alipay.map.age :  18
yaml存在的key是: database.mysql.port :  3306
yaml存在的key是: nosql.es.host :  127.0.0.1
yaml存在的key是: nosql.es.password :  456
yaml存在的key是: ksd.alipay.appid :  12454545
初始化配置文件成功！
```



# 关于项目中如何整合GORM框架

官方文档：[连接到数据库 | GORM - The fantastic ORM library for Golang, aims to be developer friendly.](https://gorm.io/zh_CN/docs/connecting_to_the_database.html)

https://gorm.io/zh_CN/docs/

## 1、整合gorm框架

### 安装

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

编写main.go来连接数据库

需要自己在本地安装好mysql新建数据库，这里不做演示。

### 初始化数据库

global.go文件代码修改

```go
var (
	Yaml   map[string]any
	Config *parse.Config
	XZ_DB  *gorm.DB
)
```

commons -> parse -> Database.go新增方法

```go
func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}
```



在initlization下新建`init_gorm.go`

```go

```









# 知识点

记录框架涉及到的知识点

## viper库

在 Go 语言中，Viper 是一个流行的配置管理库。它被设计用来处理应用程序配置需求，如读取、写入和管理配置文件。Viper 支持多种不同的配置格式，包括 JSON、TOML、YAML、HCL、envfile 和 Java properties 配置文件。

Viper 的一些关键特性包括：

1. **设置默认值**：你可以为不同的配置选项设置默认值。
2. **读取配置文件**：Viper 可以从文件系统中的文件、环境变量、命令行参数等多个来源读取配置。
3. **重载配置**：在运行时，Viper 可以让你重载配置文件，无需重新启动应用。
4. **环境变量**：Viper 可以绑定到环境变量，使得配置可以通过环境变量进行覆盖。
5. **远程配置系统**：Viper 可以从远程配置系统（如 etcd 或 Consul）读取并监视配置变化。
6. **命令行参数解析**：Viper 集成了 Cobra 库，可以与命令行参数解析库一起使用。

基本上，Viper 被用来处理应用程序中所有的配置需求，使得配置管理变得简单而统一。它特别适合用在12因素应用程序和云原生应用程序中，这些应用程序经常从环境变量和服务发现层面读取配置。

使用 Viper，开发者可以轻松地编写读取配置的代码，而不必担心配置数据来自哪里，或者是什么格式，这帮助开发者专注于构建应用程序的业务逻辑。

### 基本使用

Viper 是一个用于 Go 应用程序的配置解决方案，它允许你从多种来源加载配置，包括 JSON、YAML、TOML、环境变量和命令行参数。以下是 Viper 的一些基本用法：

首先，你需要安装 Viper。假设你已经设置了 Go 环境，你可以使用 `go get` 命令来安装 Viper：

```sh
go get github.com/spf13/viper
```

接下来，你可以在你的 Go 程序中使用 Viper。以下是一个简单的例子，展示了如何使用 Viper 来加载配置文件和读取配置项：

```go
package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func main() {
	// 设置配置文件的名称和类型
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")      // 查找配置文件所在的路径

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// 显示配置文件路径
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	// 获取并打印一个配置值
	if viper.IsSet("item.key") {
		value := viper.Get("item.key") // 返回 interface{}
		fmt.Printf("Value: %v\n", value)
	} else {
		fmt.Println("Key not found")
	}

	// 直接获取字符串类型的配置值
	strValue := viper.GetString("item.key") // 返回 string
	fmt.Printf("String Value: %s\n", strValue)

	// 从环境变量读取
	viper.AutomaticEnv() // 自动读取匹配的环境变量

	// 假设有一个环境变量名为 PREFIX_ITEM_KEY
	envValue := viper.Get("PREFIX_ITEM_KEY") // 返回 interface{}
	fmt.Printf("Environment Value: %v\n", envValue)
}
```

在上面的例子中，我们首先设置了配置文件的名称和类型，并告诉 Viper 在当前目录中查找配置文件。然后我们尝试读取配置文件，并在成功加载后，使用 `viper.Get` 和 `viper.GetString` 方法来获取配置项的值。我们还使用了 `viper.AutomaticEnv` 来自动读取环境变量。

你需要有一个名为 `config.yaml` 的配置文件在你的工作目录中，内容如下：

```yaml
item:
  key: value
```

此外，如果你有一个环境变量 `PREFIX_ITEM_KEY`，Viper 也可以读取它的值。

这只是 Viper 的一些基本功能。Viper 还支持更复杂的配置结构、配置文件热重载、远程配置等特性。



# OS库

