# 架构分析

 使用gin作为web服务，使用vue作为前端服务，来构建一个admin管理后台系统。

**模板文件**

Goland创建go文件添加作者信息：

- File -> Settings -> Editor -> File and Code Templates -> Go File

- 右侧填入模板

- ```go
  /* 
  @Author: 梦无矶小仔
  @Date:   ${DATE} ${TIME}
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

# 配置viper

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



# 整合GORM框架

官方文档：[连接到数据库 | GORM - The fantastic ORM library for Golang, aims to be developer friendly.](https://gorm.io/zh_CN/docs/connecting_to_the_database.html)

https://gorm.io/zh_CN/docs/

## Gorm的安装

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

编写main.go来连接数据库

需要自己在本地安装好mysql新建数据库，这里不做演示。

## 初始化数据库&建表&简单封装

### 1、global.go文件代码修改

```go
var (
	Yaml   map[string]any
	Config *parse.Config
	XZ_DB  *gorm.DB
)
```

### 2、commons -> parse -> Database.go新增方法

```go
func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}
```

### 3、新建model文件夹

```shell
model
	- user
		- xz_user.go
```

```go
/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 19:03
 */
package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Account   string         `gorm:"unique" json:"account"`
	Password  string         `json:"password"`
	Email     *string        `json:"email"`
	Age       uint8          `json:"age"`
	Birthday  time.Time      `json:"birthday"`
	CreatedAt time.Time      `gorm:"autoUpdateTime" json:"createAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// 覆盖生成表
func (User) TableName() string {
	return "xz_admin_user"
}
```



### 4、commons ->新建orm->新建registertable.go注册表

```go
package orm

func RegisterTable() {
	db := global.XZ_DB
	// 注册和声明model
	db.AutoMigrate(&user.User{})
}
```



### 5、在initlization下新建`init_gorm.go`，初始化数据库

```go
/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 16:32
 */
package initlization

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"xz-go-frame/commons/orm"
	"xz-go-frame/global"
)

func InitMySQL() {
	m := global.Config.Database.Mysql
	fmt.Println(m.Dsn())
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               m.Dsn(), // DNS data source name
		DefaultStringSize: 191,     //string类型字段的默认长度
		//DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}))

	// 如果报错，请检查数据库配置
	if err != nil {
		panic("数据连接出错了" + err.Error()) // 把程序直接阻断，把数据连接好了在启动
	}

	global.XZ_DB = db // 数据库信息全局变量

	// 初始化数据库
	orm.RegisterTable()

	fmt.Println("数据库初始化完成,开始运行：", db)
}

```



### 6、修改main.go的代码

```go
func main() {
	//  开始初始化配置文件
	initlization.InitViper()
	fmt.Println("初始化配置文件成功！")
	// 初始化数据库
	initlization.InitMySQL()
	//开始初始化gin路由服务
	initlization.WebRouterInit()
	fmt.Println("启动xz-go-frame后端成功")

}
```



### 7、启动服务,控制台显示如下

```shell
root:123456@tcp(127.0.0.1:3306)/xz-go-frame-db?charset=utf8&parseTime=True&loc=Local
数据库初始化完成,开始运行： &{0xc0000a6090 <nil> 0 0xc0003ea000 1}
```



### 8、数据库中新增了`xz_admin_user表`

![image-20231227191043180](images/image-20231227191043180.png)

到此，你的gorm框架整合完毕，后续进行相关封装



# 层级封装

api层我新建了一个v1文件夹，把之前的代码放到了这个v1文件夹下，起到api版本管理的作用。

## service层

这一层主要是操作数据库，给api层提供数据

新建service文件夹在根目录

```shell
service
	- user
		- xz_user.go
```



xz_user.go代码如下

```go
/*
* @Author: 梦无矶小仔
* @Date:   2024/1/11 14:11
 */
package user

import (
	"xz-go-frame/global"
	"xz-go-frame/model/user"
)

// 对用户表的数据层处理
type UserService struct{}

// nil 是go空值处理，必须是指针类型
func (service *UserService) GetUserByAccount(account string) (user *user.User, err error) {
	// 根据account进行查询
	err = global.XZ_DB.Where("account = ?", account).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

```



路由router层和api层之后会加入enter.go文件，起到初始化聚合加载的作用。



# 关于Http.status响应状态码

### 1、 HTTP Status Code 1xx 请求信息

这一组[状态码](https://so.csdn.net/so/search?q=状态码&spm=1001.2101.3001.7020)表明这是一个临时性响应。此响应仅由状态行和可选的HTTP头组成，以一个空行结尾。由于HTTP／1.0未定义任何1xx状态码，所以不要向HTTP／1.0客户端发送1xx响应。

| Http状态码 | Http Status Code                                             | Http状态码含义中文说明 |
| :--------- | :----------------------------------------------------------- | :--------------------- |
| 100        | [100 Continue](https://seo.juziseo.com/doc/http_code/100)    | 请继续请求             |
| 101        | [101 Switching Protocols](https://seo.juziseo.com/doc/http_code/101) | 请切换协议             |
| 102        | [102 Processing](https://seo.juziseo.com/doc/http_code/102)  | 将继续执行请求         |



### 2、 HTTP Status Code 2xx 成功状态

这一组状态码表明客户端的请求已经被服务器端成功接收并正确解析。

| Http状态码    | Http Status Code                                             | Http状态码含义中文说明                  |
| :------------ | :----------------------------------------------------------- | :-------------------------------------- |
| ***\*200\**** | [200 OK](https://seo.juziseo.com/doc/http_code/200)          | 请求成功                                |
| ***\*201\**** | [201 Created](https://seo.juziseo.com/doc/http_code/201)     | 请求已被接受，等待资源响应              |
| ***\*202\**** | [202 Accepted](https://seo.juziseo.com/doc/http_code/202)    | 请求已被接受，但尚未处理                |
| ***\*203\**** | [203 Non-Authoritative Information](https://seo.juziseo.com/doc/http_code/203) | 请求已成功处理，结果来自第三方拷贝      |
| ***\*204\**** | [204 No Content](https://seo.juziseo.com/doc/http_code/204)  | 请求已成功处理，但无返回内容            |
| ***\*205\**** | [205 Reset Content](https://seo.juziseo.com/doc/http_code/205) | 请求已成功处理，但需重置内容            |
| ***\*206\**** | [206 Partial Content](https://seo.juziseo.com/doc/http_code/206) | 请求已成功处理，但仅返回了部分内容      |
| ***\*207\**** | [207 Multi-Status](https://seo.juziseo.com/doc/http_code/207) | 请求已成功处理，返回了多个状态的XML消息 |
| ***\*208\**** | [208 Already Reported](https://seo.juziseo.com/doc/http_code/208) | 响应已发送                              |
| ***\*226\**** | [226 IM Used](https://seo.juziseo.com/doc/http_code/226)     | 已完成响应                              |

### 3、 HTTP Status Code 3xx 重定向状态

这一组状态码表示客户端需要采取更进一步的行动来完成请求。通常，这些状态码用来重定向，后续的请求地址（重定向目标）在本次响应的Location域中指明。

| Http状态码    | Http Status Code                                             | Http状态码含义中文说明         |
| :------------ | :----------------------------------------------------------- | :----------------------------- |
| ***\*300\**** | [300 Multiple Choices](https://seo.juziseo.com/doc/http_code/300) | 返回多条重定向供选择           |
| ***\*301\**** | [301 Moved Permanently](https://seo.juziseo.com/doc/http_code/301) | 永久重定向                     |
| ***\*302\**** | [302 Found](https://seo.juziseo.com/doc/http_code/302)       | 临时重定向                     |
| ***\*303\**** | [303 See Other](https://seo.juziseo.com/doc/http_code/303)   | 当前请求的资源在其它地址       |
| ***\*304\**** | [304 Not Modified](https://seo.juziseo.com/doc/http_code/304) | 请求资源与本地缓存相同，未修改 |
| ***\*305\**** | [305 Use Proxy](https://seo.juziseo.com/doc/http_code/305)   | 必须通过代理访问               |
| ***\*306\**** | [306 (已废弃)Switch Proxy](https://seo.juziseo.com/doc/http_code/306) | (已废弃)请切换代理             |
| ***\*307\**** | [307 Temporary Redirect](https://seo.juziseo.com/doc/http_code/307) | 临时重定向，同302              |
| ***\*308\**** | [308 Permanent Redirect](https://seo.juziseo.com/doc/http_code/308) | 永久重定向，且禁止改变http方法 |

### 4、 HTTP Status Code 4xx 客户端错误

这一组状态码表示客户端的请求存在错误，导致服务器无法处理。除非响应的是一个HEAD请求，否则服务器就应该返回一个解释当前错误状况的实体，以及这是临时的还是永久性的状况。这些状态码适用于任何请求方法。浏览器应当向用户显示任何包含在此类错误响应中的实体内容。

| Http状态码    | Http Status Code                                             | Http状态码含义中文说明               |
| :------------ | :----------------------------------------------------------- | :----------------------------------- |
| ***\*400\**** | [400 Bad Request](https://seo.juziseo.com/doc/http_code/400) | 请求错误，通常是访问的域名未绑定引起 |
| ***\*401\**** | [401 Unauthorized](https://seo.juziseo.com/doc/http_code/401) | 需要身份认证验证                     |
| ***\*402\**** | [402 Payment Required](https://seo.juziseo.com/doc/http_code/402) | -                                    |
| ***\*403\**** | [403 Forbidden](https://seo.juziseo.com/doc/http_code/403)   | 禁止访问                             |
| ***\*404\**** | [404 Not Found](https://seo.juziseo.com/doc/http_code/404)   | 请求的内容未找到或已删除             |
| ***\*405\**** | [405 Method Not Allowed](https://seo.juziseo.com/doc/http_code/405) | 不允许的请求方法                     |
| ***\*406\**** | [406 Not Acceptable](https://seo.juziseo.com/doc/http_code/406) | 无法响应，因资源无法满足客户端条件   |
| ***\*407\**** | [407 Proxy Authentication Required](https://seo.juziseo.com/doc/http_code/407) | 要求通过代理的身份认证               |
| ***\*408\**** | [408 Request Timeout](https://seo.juziseo.com/doc/http_code/408) | 请求超时                             |
| ***\*409\**** | [409 Conflict](https://seo.juziseo.com/doc/http_code/409)    | 存在冲突                             |
| ***\*410\**** | [410 Gone](https://seo.juziseo.com/doc/http_code/410)        | 资源已经不存在(过去存在)             |
| ***\*411\**** | [411 Length Required](https://seo.juziseo.com/doc/http_code/411) | 无法处理该请求                       |
| ***\*412\**** | [412 Precondition Failed](https://seo.juziseo.com/doc/http_code/412) | 请求条件错误                         |
| ***\*413\**** | [413 Payload Too Large](https://seo.juziseo.com/doc/http_code/413) | 请求的实体过大                       |
| ***\*414\**** | [414 Request-URI Too Long](https://seo.juziseo.com/doc/http_code/414) | 请求的URI过长                        |
| ***\*415\**** | [415 Unsupported Media Type](https://seo.juziseo.com/doc/http_code/415) | 无法处理的媒体格式                   |
| ***\*416\**** | [416 Range Not Satisfiable](https://seo.juziseo.com/doc/http_code/416) | 请求的范围无效                       |
| ***\*417\**** | [417 Expectation Failed](https://seo.juziseo.com/doc/http_code/417) | 无法满足的Expect                     |
| ***\*418\**** | [418 I'm a teapot](https://seo.juziseo.com/doc/http_code/418) | 愚人节笑话                           |
| ***\*421\**** | [421 There are too many connections from your internet address](https://seo.juziseo.com/doc/http_code/421) | 连接数超限                           |
| ***\*422\**** | [422 Unprocessable Entity](https://seo.juziseo.com/doc/http_code/422) | 请求的语义错误                       |
| ***\*423\**** | [423 Locked](https://seo.juziseo.com/doc/http_code/423)      | 当前资源被锁定                       |
| ***\*424\**** | [424 Failed Dependency](https://seo.juziseo.com/doc/http_code/424) | 当前请求失败                         |
| ***\*425\**** | [425 Unordered Collection](https://seo.juziseo.com/doc/http_code/425) | 未知                                 |
| ***\*426\**** | [426 Upgrade Required](https://seo.juziseo.com/doc/http_code/426) | 请切换到TLS/1.0                      |
| ***\*428\**** | [428 Precondition Required](https://seo.juziseo.com/doc/http_code/428) | 请求未带条件                         |
| ***\*429\**** | [429 Too Many Requests](https://seo.juziseo.com/doc/http_code/429) | 并发请求过多                         |
| ***\*431\**** | [431 Request Header Fields Too Large](https://seo.juziseo.com/doc/http_code/431) | 请求头过大                           |
| ***\*449\**** | [449 Retry With](https://seo.juziseo.com/doc/http_code/449)  | 请重试                               |
| ***\*451\**** | [451 Unavailable For Legal Reasons](https://seo.juziseo.com/doc/http_code/451) | 访问被拒绝（法律的要求）             |
| ***\*499\**** | [499 Client Closed Request](https://seo.juziseo.com/doc/http_code/499) | 客户端主动关闭了连接                 |

### 5、 HTTP Status Code 5xx 服务器错误状态

这一组状态码说明服务器在处理请求的过程中有错误或者异常状态发生，也有可能是服务器意识到以当前的软硬件资源无法完成对请求的处理。除非这是一个HEAD请求，否则服务器应当包含一个解释当前错误状态以及这个状况是临时的还是永久的解释信息实体。浏览器应当向用户展示任何在当前响应中被包含的实体。

| Http状态码    | Http Status Code                                             | Http状态码含义中文说明   |
| :------------ | :----------------------------------------------------------- | :----------------------- |
| ***\*500\**** | [500 Internal Server Error](https://seo.juziseo.com/doc/http_code/500) | 服务器端程序错误         |
| ***\*501\**** | [501 Not Implemented](https://seo.juziseo.com/doc/http_code/501) | 服务器不支持的请求方法   |
| ***\*502\**** | [502 Bad Gateway](https://seo.juziseo.com/doc/http_code/502) | 网关无响应               |
| ***\*503\**** | [503 Service Unavailable](https://seo.juziseo.com/doc/http_code/503) | 服务器端临时错误         |
| ***\*504\**** | [504 Gateway Timeout](https://seo.juziseo.com/doc/http_code/504) | 网关超时                 |
| ***\*505\**** | [505 HTTP Version Not Supported](https://seo.juziseo.com/doc/http_code/505) | 服务器不支持的HTTP版本   |
| ***\*506\**** | [506 Variant Also Negotiates](https://seo.juziseo.com/doc/http_code/506) | 服务器内部配置错误       |
| ***\*507\**** | [507 Insufficient Storage](https://seo.juziseo.com/doc/http_code/507) | 服务器无法存储请求       |
| ***\*508\**** | [508 Loop Detected](https://seo.juziseo.com/doc/http_code/508) | 服务器因死循环而终止操作 |
| ***\*509\**** | [509 Bandwidth Limit Exceeded](https://seo.juziseo.com/doc/http_code/509) | 服务器带宽限制           |
| ***\*510\**** | [510 Not Extended](https://seo.juziseo.com/doc/http_code/510) | 获取资源策略未被满足     |
| ***\*511\**** | [511 Network Authentication Required](https://seo.juziseo.com/doc/http_code/511) | 需验证以许可连接         |
| ***\*599\**** | [599 Network Connect Timeout Error](https://seo.juziseo.com/doc/http_code/599) | 网络连接超时             |

上面的这些状态全部都是描述请求和响应的整个过程的状态。不包含业务的状态。我举例例子

```go
c.JSON(http.StatusOK, "你输入的账号和密码有误!!!")
```

在开发中一个接口：成功只有一种情况，但是失败和错误就有N种情况。那么这些N情况的返回到底选择什么样子状态就变得非常的重要。你接下来就困难选择症（你不知道到底选择什么？），而且你在这些状态找不到适合的，所以业务的状态处理不应该选择web框架中提供的。因为这些状态根本就不是让你来做业务错误的状态监控，别人专门去web请求和响应的状态的监听。

那么怎么处理。通过自己的方式来定义状态。但是所有的返回都用：http.StatusOK

- 只要请求和响应是正常的，无论正确和错误，我们都用http.StatusOK来返回，但是区分和界定用自己定义的状态来定义业务、

这也就是为什么我们要做自定义返回的意义和价值了.

# 集成验证码功能

- 短信验证码的重要性一：作为身份证明。
  - 在生活中，短信验证码随处可见，网络产品在开发过程中，几乎都会加入短信验证模块，如网站、app用户注册、安全登录、找回密码、绑定手机、手机银行转账等等，这就是短信验证码的重要性。
- 短信验证码的重要性二：提高注册信息的真实性，防止恶意注册。
  - 与以往的数字验证、图片验证相比，短信验证码更能防止恶意用户注册。一些朋友可能知道，市场上一些不法之徒使用作弊器等工具恶意注册、攻击企业网站，导致网站服务器无法承载而瘫痪，严重时会影响企业网站的运作。而且短信验证码的应用，能很好地识别用户身份的真实性，一个用户只能注册一个账户，有效地避免了恶意注册。

### 1:  下载验证码的组件

- 官网：https://github.com/mojocn/base64Captcha

- 下载模块：

  ```
  go get github.com/mojocn/base64Captcha
  ```

### 2: 定义生成验证码的接口

```shell
api -> v1 -> code -> code.go
```

```go
/*
@Author: 梦无矶小仔
@Date:   2024/1/11 14:32
*/

package code

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"xz-go-frame/commons/response"
)

// 验证码生成
type CodeApi struct{}

// 1、定义验证的store -- 默认是存储在go的内存中
var store = base64Captcha.DefaultMemStore

// 2、创建验证码
func (api *CodeApi) CreateCaptcha(c *gin.Context) {
	// 2、生成验证码的类型，这个默认类型是一个数字的验证码
	driver := &base64Captcha.DriverDigit{Height: 70, Width: 240, Length: 6, MaxSkew: 0.8, DotCount: 120}
	// 3、调用NewCaptcha方法生成具体的验证码对象
	captcha := base64Captcha.NewCaptcha(driver, store)
	// 4、调用Generate()生成具体base64验证码的图片地址和id
	// id 是用于后续校验使用，后续根据id和用户输入的验证码去调用store的get方法，就可以得到你输入的验证码是否正确，正确的true，错误的false
	id, baseURL, _, err := captcha.Generate()
	if err != nil {
		response.Fail(40001, "验证生成错误!", c)
		return
	}

	response.Ok(map[string]any{"id": id, "baseURL": baseURL}, c)
}

//func (api *CodeApi) CreateCaptcha(c *gin.Context) {
//	// 2：生成验证码的类型,这个默认类型是一个数字的验证码
//	driver := &base64Captcha.DriverString{
//		Height:          40,
//		Width:           240,
//		NoiseCount:      0,
//		ShowLineOptions: 2 | 2,
//		Length:          6,
//		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
//		BgColor: &color.RGBA{
//			R: 3,
//			G: 102,
//			B: 214,
//			A: 125,
//		},
//		Fonts: []string{"wqy-microhei.ttc"},
//	}
//	// 3：调用NewCaptcha方法生成具体验证码对象
//	captcha := base64Captcha.NewCaptcha(driver, store)
//	// 4: 调用Generate()生成具体base64验证码的图片地址，和id
//	// id 是用于后续校验使用，后续根据id和用户输入的验证码去调用store的get方法，就可以得到你输入的验证码是否正确，正确true,错误false
//	id, baseURL, err := captcha.Generate()
//
//	if err != nil {
//		response.Fail(40001, "验证生成错误", c)
//		return
//	}
//
//	response.Ok(map[string]any{"id": id, "baseURL": baseURL}, c)
//}

// 3、开始校验用户输入的验证码是否是正确的
func (api *CodeApi) VerifyCaptcha(c *gin.Context) {

	type BaseCaptcha struct {
		Id   string `json:"id"`
		Code string `json:"code"`
	}
	baseCaptcha := BaseCaptcha{}
	// 开始把用户输入的id和code进行绑定
	err2 := c.ShouldBindQuery(&baseCaptcha)

	if err2 != nil {
		response.Fail(402, "参数绑定失败", c)
		return
	}

	// 开始校验验证码是否正确
	verify := store.Verify(baseCaptcha.Id, baseCaptcha.Code,true)
	
	if verify {
		response.Ok("success", c)
	} else {
		response.Fail(403, "您输入的验证码有误！", c)
	}
	

}

```

封装了一个响应方法，`commons->response->response.go`

```go
/*
@Author: 梦无矶小仔
@Date:   2024/1/11 14:53
*/
package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Result(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

var (
	CODE = 20000
	MSG  = "success"
)

/*
Ok
请求响应成功
*/
func Ok(data any, c *gin.Context) {
	Result(CODE, MSG, data, c)
}

/*
Fail
请求响应失败（无响应数据）
*/
func Fail(code int,msg string,c *gin.Context) {
	Result(code, msg, map[string]any{}, c)
}

/*
FailWithData
请求响应失败（有响应数据）
*/
func FailWithData(code int,msg string,data any,c *gin.Context) {
	Result(code, msg, data, c)
}

```

### 3:定义生成验证码的router

router -> code -> code.go

```go
/*
@Author: 梦无矶小仔
@Date:   2024/1/11 16:23
*/
package code

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/api/v1/code"
)

type CodeRouter struct{}

func (e *CodeRouter) InitCodeRouter(Router *gin.Engine) {
	codeApi := code.CodeApi{}
	// 这个路由多了一个对post\put请求的中间件处理，而这个中间件做了一些对post和put参数的处理和一些公共信息的处理
	coureseRouter := Router.Group("code")
	{
		coureseRouter.GET("get", codeApi.CreateCaptcha)
		coureseRouter.GET("verify", codeApi.VerifyCaptcha)
	}
	
}

```

### 4:注册路由

```go
// 验证码接口
codeRouter := code.CodeRouter{}
codeRouter.InitCodeRouter(rootRouter)
```

到这一步，就可以去测试验证码功能了。

### 5：跨域请求处理

在之后和前端联调，需要设置跨域，所以还要新增跨域请求的代码。

`common -> filter -> cors.go`

```go
/*
@Author: 梦无矶小仔
@Date:   2024/1/11 17:28
*/
package filter

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// 注意这一行，不能配置为通配符“*”号 比如未来写域名或者你想授权的域名都可以
		//c.Header("Access-Control-Allow-Origin", "http://localhost:8088")
		c.Header("Access-Control-Allow-Origin", "*")
		// 响应头表示是否可以将对请求的响应暴露给页面。返回true则可以，其他值均不可以。
		c.Header("Access-Control-Allow-Credentials", "true")
		// 表示此次请求中可以使用那些header字段
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Cookie, Content-Length,Origin,cache-control,X-Requested-With, Content-Type, Accept, Authorization, Token, Timestamp, UserId") // 我们自定义的header字段都需要在这里声明
		// 表示此次请求中可以使用那些请求方法 GET/POST(多个使用逗号隔开)
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			//c.AbortWithStatus(http.StatusNoContent)
			c.AbortWithStatus(http.StatusOK)
		}
		// 处理请求
		c.Next()
	}
}

```



# 构建前端

## 1、 打开 vite的官网

https://vitejs.cn/

## 2、开始构建项目

> ​	兼容性注意
>
> Vite 需要 [Node.js](https://nodejs.org/en/) 版本 14.18+，16+。然而，有些模板需要依赖更高的 Node 版本才能正常运行，当你的包管理器发出警告时，请注意升级你的 Node 版本。

使用nodejs 下载一些工具，由于国内访问很多外网的限制，会出现下载失败的问题,这个时候，需要配置路径为阿里的免费 registry，如下：

官方链接：[npmmirror 镜像站](https://www.npmmirror.com/)

```shell
1、临时使用
npm --registry https://registry.npmmirror.com install express
2、永久使用
npm config set registry https://registry.npmmirror.com

-- 配置后可通过下面方式来验证是否成功
npm config get registry
-- 或npm info express

# 如果需要恢复成原来的官方地址只需要执行如下命令
npm config set registry https://registry.npmjs.org

#使用cnpm
#安装阿里的cnpm，然后在使用时直接将npm命令替换成cnpm命令即可（淘宝的已经停止解析）
npm install -g cnpm --registry=https://registry.npmmirror.com
```



## 3、进入目标文件夹下执行如下命令

```js
npm create vite@latest
```

当前下载的是vite@5.1.0

## 4、勾选内容如下

![image-20240110155723715](images/image-20240110155723715.png)

```shell
D:\Y_WebProject>npm create vite@latest
√ Project name: ... xz-vue-admin
√ Select a framework: » Vue
√ Select a variant: » Customize with create-vue ↗
Need to install the following packages:
  create-vue@3.9.1
Ok to proceed? (y) y
[##################] - reify:create-vue: http fetch GET 200 https://cdn.npmmirror.com/packages/create-vue/3.9.1/create- 

```

vue.js相关配置勾选，蓝色的就是我选择的，左右箭头可以进行选择，回车表示确认。

![image-20240110160254149](images/image-20240110160254149.png)



## 5、整体命令如下

```bash
D:\Y_WebProject>npm create vite@latest
√ Project name: ... xz-vue-admin
√ Select a framework: » Vue
√ Select a variant: » Customize with create-vue ↗
Need to install the following packages:
  create-vue@3.9.1
Ok to proceed? (y) y

Vue.js - The Progressive JavaScript Framework

√ 是否使用 TypeScript 语法？ ... 否 / 是
√ 是否启用 JSX 支持？ ... 否 / 是
√ 是否引入 Vue Router 进行单页面应用开发？ ... 否 / 是
√ 是否引入 Pinia 用于状态管理？ ... 否 / 是
√ 是否引入 Vitest 用于单元测试？ ... 否 / 是
√ 是否要引入一款端到端（End to End）测试工具？ » 不需要
√ 是否引入 ESLint 用于代码质量检测？ ... 否 / 是
√ 是否引入 Prettier 用于代码格式化？ ... 否 / 是

正在构建项目 D:\Y_WebProject\xz-vue-admin...

项目构建完成，可执行以下命令：

  cd xz-vue-admin
  npm install
  npm run format
  npm run dev
```

此时可以看到我们的项目创建好啦

![image-20240110160601617](images/image-20240110160601617.png)

## 6、下载依赖

cd到我们创建好的项目中，执行命令下载相关依赖

```bash
cd xz-vue-admin
npm install
```

显示如下表示下载依赖完成

```shell
D:\Y_WebProject>cd xz-vue-admin

D:\Y_WebProject\xz-vue-admin>npm install

added 152 packages in 2m
```

## 7、启动项目

```shell
npm run dev
```

显示如下

```shell
  VITE v5.0.11  ready in 438 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
  ➜  press h + enter to show help
```

## 8、访问项目

访问本地的5173端口，显示页面如下表示你的项目构建成功

![image-20240110162618137](images/image-20240110162618137.png)



## 9、同步到自己的仓库进行版本管理即可。

修改项目根目录下的`.git/config`文件内容，把需要同步的仓库链接都填进去，这样可以多个仓库同时同步。

```yaml
[core]
	repositoryformatversion = 0
	filemode = false
	bare = false
	logallrefupdates = true
	symlinks = false
	ignorecase = true
[remote "origin"]
	# github
	url = https://github.com/Lvan826199/xz-vue-admin.git
	# gitee
	url = https://gitee.com/xiaozai-van-liu/xz-vue-admin.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
	remote = origin
	merge = refs/heads/master

```



## 03、实战案例：验证码使用到登录和注册

- 1：下载验证码模块组件

- 2：通过验证码返回的captchaid和用户输入的code和手机号码调用短信发送接口
- 3：通过captchaid和code进行验证比较看是否输入的正确的，如果正确返回success,否则返回错误
- 4：思考

具体查看视频和代码。核心js和文件

login.js

```js
var vue = new Vue({
    el:"#app",
    data:{
        // 控制登录按钮是否可以登录.
        btndisabled:true,
        // 60s倒计时
        sendcount:60,
        // 默认显示发送短信
        sendflag:true,
        // 用于清楚倒计时
        sendTimer:null,
        // 倒计时状态
        sendmsg:"发送短信",
        // 验证码
        codeimg:"",
        // 登录的数据
        user:{
            phone:"15074816437",
            phonecode:"",
            captchaId:""
        }
    },
    created(){
      this.handleCaptchaCode();
    },
    methods:{
        toLogin(){

            if(!this.user.phone){
                alert("请输入手机号码")
                this.$refs.phoneRef.focus();
                return;
            }

            if(!this.user.phonecode){
                alert("请输入手机短信码")
                this.$refs.phonecodeRef.focus();
                return;
            }

            // 正则校验手机号码合法性
            // if(!/phonerege/.test(phone)){
            //     alert("请输入合法的手机号码!")
            //     this.$refs.phoneRef.focus();
            //     return;
            // }

            axios.post("/api/logined",this.user).then(res=>{
                if(res.data.code == 200){
                    window.location.href = "/"
                }else{
                    if(res.data.code == 601){  // 601: 短信验证码输入有误
                        alert(res.data.message)
                        this.$refs.phonecodeRef.focus();
                        this.user.phonecode = "";
                    }else if(res.data.code == 602){ // 602: 手机号码不存在
                        alert(res.data.message)
                        this.$refs.phoneRef.focus();
                        this.user.phone = "";
                    }
                }
            })
        },

        // 发送短信
        handleSendPhone(){
            var phone = this.user.phone;
            if(!phone){
                alert("请正确的输入手机号码!")
                this.$refs.phoneRef.focus();
                return;
            }
            // 正则校验手机号码合法性 phonerege = ?
            // if(!/phonerege/.test(phone)){
            //     alert("请输入合法的手机号码!")
            //     this.$refs.phoneRef.focus();
            //     return;
            // }

            // 更改前端的状态
            this.handleChangeSendMsg();
            // 发送短信-------发送短信接口
            axios.post("/api/sendsms",{"phone":phone}).then(res=>{
                if(res.data == "success"){
                    alert("短信发送成功!!!!");
                    // 打开登录按钮。--同时禁止发送短信按钮
                    this.btndisabled = false;
                    // 恢复短信发送的状态
                    this.sendmsg = "发送短信";
                    this.sendcount = 60;
                    // 关闭定时任务
                    if(this.sendTimer)clearInterval(this.sendTimer);
                }else{
                    alert(res.data)
                }
            })
        },

        // 更改文案和倒计时
        handleChangeSendMsg(){
            this.sendmsg = this.sendcount+"s";
            if(this.sendTimer)clearInterval(this.sendTimer);
            // 开始倒计时状态
            this.sendflag = false;
            // 开启倒计时
            this.sendTimer = setInterval(()=>{
                if(this.sendcount<=1){
                    // 关闭倒计时
                    this.sendflag = true;
                    clearInterval(this.sendTimer);
                    this.sendmsg = "发送短信"
                    this.sendcount = 60;
                    return;
                }
                this.sendcount--;
                this.sendmsg = this.sendcount+"s";
            },1000)
        },

        handleCaptchaCode(){
            axios.post("/code/captcha").then(res=>{
                this.codeimg = res.data.img;
                this.user.captchaId = res.data.captchaId;
            })
        },
    }
})
```



# 实现登录和验证码的校验

### 1: 配置登录路由和首页

如果报错记得安装一下：sass

```sh
pnpm install sass sass-loader
```

```vue
<template>
    <div class="login-box">
        <div class="loginbox">
            <div class="login-wrap">
                <h1 class="header">{{ title }}</h1>
                <form action="#">
                    <div class="ksd-el-items"><input type="text" class="ksd-login-input"  placeholder="请输入账号"></div>
                    <div class="ksd-el-items"><input type="text" class="ksd-login-input" placeholder="请输入密码"></div>
                    <div class="ksd-el-items"><input type="text" class="ksd-login-input" placeholder="请输入验证码"></div>
                    <div class="ksd-el-items"><input type="button" class="ksd-login-btn" value="登录"></div>            
                </form>
            </div>
        </div>
        <div class="imgbox">
            <img src="../assets/imgs/login_left.svg" alt="">    
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
const title = ref("我是一个登录页面")
</script>

<style scoped lang="scss">
    .ksd-el-items{margin: 15px 0;}
    .ksd-login-input{border:1px solid #eee;padding:12px 8px;width: 100%;box-sizing: border-box;outline: none;border-radius: 4px;}

    .ksd-login-btn{border:1px solid #eee;padding:12px 8px;width: 100%;box-sizing: border-box;
        background:#2196F3;color:#fff;border-radius:4px;cursor: pointer;}
        .ksd-login-btn:hover{background:#1789e7;}
    .login-box{
        display: flex;
        flex-wrap: wrap;
        background: url("../assets/imgs/login_background.jpg");
        background-size:cover;
        .loginbox{
            width: 50%;height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            .header{margin-bottom: 10px;}
            .login-wrap{
                width: 500px;
                height: 444px;
                padding:20px 100px;
                box-sizing: border-box;
                border-radius: 8px;
                box-shadow: 0 0 10px #fafafa;
                background: rgba(255,255,255,0.6);
                text-align: center;
                display: flex;
                flex-direction: column;
                justify-content: center;
            }
        }
        .imgbox{
            width: 50%;
            height: 100vh;
            display: flex;
            align-items: center;
        }
    }
</style>
```

2：配置路由

```js
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'index',
      component: () => import('../views/index.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/login.vue')
    }
  ]
})

export default router

```

3: 访问login.vue

http://localhost:8777/login

![image-20230714220145203](images/image-20230714220145203.png)



4： 调用验证码的接口和登录接口

安装axios的异步请求组件

```js
pnpm install axios
```

开始调用接口?

```js
import { onMounted, ref } from 'vue';
import axios from 'axios'
const title = ref("我是一个登录页面")

// 根据axios官方文档开始调用生成验证码的接口
const handleGetCapatcha = async () => {
    const resp = await axios.get("http://localhost:8989/code/get")
    console.log('resp', resp)
}

// 用生命周期去加载生成验证码
onMounted(() => {
    handleGetCapatcha()
})

```

![image-20230714220815155](images/image-20230714220815155.png)

上面的错误很明显的告诉，你浏览器http://localhost:8777/login，服务的接口地址：http://localhost:8989/code/get 

电脑 localhost 代表本机 127.0.0.1

- 8777 你服务端口
- /login  其实就是服务里某个请求资源

—A酒店 广州127.0.01

- 8989  就是你酒店某个房子号
- /code/get  可能



# JWT

`JSON Web Token(JWT)`是一个常用语HTTP的客户端和服务端间进行身份认证和鉴权的标准规范，使用JWT可以允许我们在用户和服务器之间传递安全可靠的信息。

在开始学习**[JWT](https://link.zhihu.com/?target=https%3A//jwt.io/)**之前，我们可以先了解下早期的几种方案。

### **token、cookie、session的区别**

**Cookie**

Cookie总是保存在客户端中，按在客户端中的存储位置，可分为`内存Cookie`和`硬盘Cookie`。

内存Cookie由浏览器维护，保存在内存中，浏览器关闭后就消失了，其存在时间是短暂的。硬盘Cookie保存在硬盘里，有一个过期时间，除非用户手工清理或到了过期时间，硬盘Cookie不会被删除，其存在时间是长期的。所以，按存在时间，可分为`非持久Cookie和持久Cookie`。

cookie 是一个非常具体的东西，指的就是浏览器里面能永久存储的一种数据，仅仅是浏览器实现的一种数据存储功能。

`cookie由服务器生成，发送给浏览器`，浏览器把cookie以key-value形式保存到某个目录下的文本文件内，下一次请求同一网站时会把该cookie发送给服务器。由于cookie是存在客户端上的，所以浏览器加入了一些限制确保cookie不会被恶意使用，同时不会占据太多磁盘空间，所以每个域的cookie数量是有限的。

**Session**

Session字面意思是会话，主要用来标识自己的身份。比如在无状态的api服务在多次请求数据库时，如何知道是同一个用户，这个就可以通过session的机制，服务器要知道当前发请求给自己的是谁

为了区分客户端请求，`服务端会给具体的客户端生成身份标识session`，然后客户端每次向服务器发请求的时候，都带上这个“身份标识”，服务器就知道这个请求来自于谁了。

至于客户端如何保存该标识，可以有很多方式，对于浏览器而言，一般都是使用`cookie`的方式

服务器使用session把用户信息临时保存了服务器上，用户离开网站就会销毁，这种凭证存储方式相对于cookie来说更加安全，但是session会有一个缺陷: 如果web服务器做了负载均衡，那么下一个操作请求到了另一台服务器的时候session会丢失。

因此，通常企业里会使用`redis,memcached`缓存中间件来实现session的共享，此时web服务器就是一个完全无状态的存在，所有的用户凭证可以通过共享session的方式存取，当前session的过期和销毁机制需要用户做控制。

**Token**

token的意思是“令牌”，是用户身份的验证方式，最简单的token组成: `uid(用户唯一标识)`+`time(当前时间戳)`+`sign(签名,由token的前几位+盐以哈希算法压缩成一定长度的十六进制字符串)`，同时还可以将不变的参数也放进token

这里我们主要想讲的就是`Json Web Token`，也就是本篇的主题:JWT

### **Json-Web-Token(JWT)介绍**

一般而言，用户注册登陆后会生成一个jwt token返回给浏览器，浏览器向服务端请求数据时携带`token`，服务器端使用`signature`中定义的方式进行解码，进而对token进行解析和验证。

### **JWT Token组成部分**

![img](images/v2-dea40372d962a09ac050a6d17e9dd2b2_720w.jpeg)

- header: 用来指定使用的算法(HMAC SHA256 RSA)和token类型(如JWT)
- payload: 包含声明(要求)，声明通常是用户信息或其他数据的声明，比如用户id，名称，邮箱等. 声明可分为三种: registered,public,private
- signature: 用来保证JWT的真实性，可以使用不同的算法

**header**

```text
{
    "alg": "HS256",
    "typ": "JWT"
}
```

对上面的json进行base64编码即可得到JWT的第一个部分

**payload**

- registered claims: 预定义的声明，通常会放置一些预定义字段，比如过期时间，主题等(iss:issuer,exp:expiration time,sub:subject,aud:audience)
- public claims: 可以设置公开定义的字段
- private claims: 用于统一使用他们的各方之间的共享信息

```text
{
    "sub": "xxx-api",
    "name": "bgbiao.top",
    "admin": true
}
```

对payload部分的json进行base64编码后即可得到JWT的第二个部分

`注意:` 不要在header和payload中放置敏感信息，除非信息本身已经做过脱敏处理

**signature**

为了得到签名部分，必须有编码过的header和payload，以及一个秘钥，签名算法使用header中指定的那个，然后对其进行签名即可

```
HMACSHA256(base64UrlEncode(header)+"."+base64UrlEncode(payload),secret)
```

签名是`用于验证消息在传递过程中有没有被更改`，并且，对于使用私钥签名的token，它还可以验证JWT的发送方是否为它所称的发送方。

在**[jwt.io](https://link.zhihu.com/?target=https%3A//jwt.io/)**网站中，提供了一些JWT token的编码，验证以及生成jwt的工具。

下图就是一个典型的jwt-token的组成部分。

![img](images/v2-becbc64838d787ca7683b257958f1d21_720w.webp)

### **什么时候用JWT**

- Authorization(授权): 典型场景，用户请求的token中包含了该令牌允许的路由，服务和资源。单点登录其实就是现在广泛使用JWT的一个特性
- Information Exchange(信息交换): 对于安全的在各方之间传输信息而言，JSON Web Tokens无疑是一种很好的方式.因为JWTs可以被签名，例如，用公钥/私钥对，你可以确定发送人就是它们所说的那个人。另外，由于签名是使用头和有效负载计算的，您还可以验证内容没有被篡改

### **JWT(Json Web Tokens)是如何工作的**

![img](images/v2-82c5f75466da70b96bfd238e0f2924b3_720w.jpeg)

所以，基本上整个过程分为两个阶段：

第一个阶段，**客户端向服务端获取token。**

第二阶段，**客户端带着该token去请求相关的资源。**

通常比较重要的是，服务端如何根据指定的规则进行token的生成。

在认证的时候，当用户用他们的凭证成功登录以后，一个JSON Web Token将会被返回。

此后，token就是用户凭证了，你必须非常小心以防止出现安全问题。

一般而言，你保存令牌的时候不应该超过你所需要它的时间。

无论何时用户想要访问受保护的路由或者资源的时候，用户代理（通常是浏览器）都应该带上JWT，典型的，通常放在Authorization header中，用Bearer schema: `Authorization: Bearer <token>`

服务器上的受保护的路由将会检查Authorization header中的JWT是否有效，如果有效，则用户可以访问受保护的资源。如果JWT包含足够多的必需的数据，那么就可以减少对某些操作的数据库查询的需要，尽管可能并不总是如此。

如果token是在授权头（Authorization header）中发送的，那么跨源资源共享(CORS)将不会成为问题，因为它不使用cookie.

![img](images/v2-50bb47d56a98d247ab5909a0fc4ddcc1_720w.jpeg)

- 客户端向授权接口请求授权
- 服务端授权后返回一个access token给客户端
- 客户端使用access token访问受保护的资源

### **基于Token的身份认证和基于服务器的身份认证**

**1.基于服务器的认证**

前面说到过session，cookie以及token的区别，在之前传统的做法就是基于存储在服务器上的session来做用户的身份认证，但是通常会有如下问题:

- Sessions: 认证通过后需要将用户的session数据保存在内存中，随着认证用户的增加，内存开销会大
- 扩展性: 由于session存储在内存中，扩展性会受限，虽然后期可以使用redis,memcached来缓存数据
- CORS: 当多个终端访问同一份数据时，可能会遇到禁止请求的问题
- CSRF: 用户容易受到CSRF攻击

**2.Session和JWT Token的异同**

都可以存储用户相关信息，但是session存储在服务端，JWT存储在客户端

![img](images/v2-1f9a28101bc26d90fdc057ba04310caa_720w.jpeg)

**3.基于Token的身份认证如何工作**

基于Token的身份认证是无状态的，服务器或者session中不会存储任何用户信息.(很好的解决了共享session的问题)

- 用户携带用户名和密码请求获取token(接口数据中可使用appId,appKey)
- 服务端校验用户凭证，并返回用户或客户端一个Token
- 客户端存储token,并在请求头中携带Token
- 服务端校验token并返回数据

**注意:**

- 随后客户端的每次请求都需要使用token
- token应该放在header中
- 需要将服务器设置为接收所有域的请求: `Access-Control-Allow-Origin: *`

**4.用Token的好处**

- 无状态和可扩展性
- 安全: 防止CSRF攻击;token过期重新认证

**5.JWT和OAuth的区别**

- 1.OAuth2是一种授权框架 ，JWT是一种认证协议
- 2.无论使用哪种方式切记用HTTPS来保证数据的安全性
- 3.OAuth2用在`使用第三方账号登录的情况`(比如使用weibo, qq, github登录某个app)，而`JWT是用在前后端分离`, 需要简单的对后台API进行保护时使用

# 使用Gin框架集成JWT

在Golang语言中，**[jwt-go](https://link.zhihu.com/?target=https%3A//github.com/dgrijalva/jwt-go)**库提供了一些jwt编码和验证的工具，因此我们很容易使用该库来实现token认证。

另外，我们也知道**[gin](https://link.zhihu.com/?target=https%3A//github.com/gin-gonic/gin)**框架中支持用户自定义middleware，我们可以很好的将jwt相关的逻辑封装在middleware中，然后对具体的接口进行认证。

## JWT的使用

下载组件

官网：  https://github.com/golang-jwt/jwt

```go
go get -u github.com/golang-jwt/jwt/v5
```

参考文章：https://blog.csdn.net/qq_39463535/article/details/133061861

防止缓存击穿神器-singleflight,参考资料见知识点。

```go
go get golang.org/x/sync/singleflight
```

### jwt的基本使用封装

在commons下新建jwtgo -> jwt_go.go

```go
/*
@Author: 梦无矶小仔
@Date:   2024/1/12 17:14
*/
package jwtgo

import (
	"errors"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/sync/singleflight"
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

// 定义一个JWT对象
type JWT struct {
	// 声明签名信息
	SigningKey []byte
}

// 初始化jwt对象
func NewJWT() *JWT {
	return &JWT{
		[]byte("www.mengwuji.com"),
	}
}

// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
	// 续期使用
	BufferTime int64
	// RegisteredClaims 内嵌标准的声明
	jwt.RegisteredClaims
}

// 创建一个token
// 指定编码的算法为jwt.SigningMethodHS256
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// 返回一个token的结构体指针
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error) {
	singleflightGroup := &singleflight.Group{}
	v, err, _ := singleflightGroup.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	// 将token中的claims信息解析出来并断言成用户自定义的有效载荷结构
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
```

### 定义jwt的model

model -> jwt -> jwt.go

```go
/*
@Author: 梦无矶小仔
@Date:   2024/1/12 18:23
*/
package jwt

import (
	"gorm.io/gorm"
	"time"
)

type JwtBlacklist struct {
	ID        uint           `gorm:"primarykey;comment:主键ID"` // 主键ID
	CreatedAt time.Time      `gorm:"type:datetime(0);comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"type:datetime(0);comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"` // 删除时间
	Jwt       string         `gorm:"type:text;comment:jwt"`
}

```

### 注册jwt表

commons -> orm -> registertable.go

```go
package orm

import (
	"xz-go-frame/global"
	"xz-go-frame/model/jwt"
	"xz-go-frame/model/user"
)

func RegisterTable() {
	db := global.XZ_DB
	// 注册和声明model
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(jwt.JwtBlacklist{})
}

```

### 定义JWT的中间件

middle -> jwt.go

```go
/*
@Author: 梦无矶小仔
@Date:   2024/1/15 11:36
*/
package middle

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/commons/response"
)


// 定义一个JWTAuth的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过http header中的token解析来认证
		// 获取token
		// 我们这里jwt鉴权取头部信息 Authorization 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Fail(701, "请求未携带token，无权限访问", c)
			c.Abort()
			return
		}
		// 生成jwt的对象
		myJwt := jwtgo.NewJWT()
		// parseToken 解析token包含的信息
		customClaims, err := myJwt.ParseToken(token)
		// 如果解析失败就出现异常
		if err != nil {
			response.Fail(60001, "token失效了", c)
			c.Abort()
			return
			}
		// 让后续的路由方法可以直接通过c.Get("claims")
		c.Set("claims", customClaims)
		c.Next()
		}
		
}

```

## 更新login相关Api内容

api -> v1 -> login -> login.go

```go
package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mojocn/base64Captcha"
	"time"
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/commons/response"
	service "xz-go-frame/service/user"
)

// 登录业务
type LoginApi struct {
}

// 1、定义验证码的store -- 默认是存储在go的内存中
var store = base64Captcha.DefaultMemStore

// 登录接口的处理
func (api *LoginApi) ToLogined(c *gin.Context) {
	type LoginParam struct {
		Account  string
		Code     string
		CodeId   string
		Password string
	}
	// 1、获取用户在页面上输入账号和密码，开始在数据库里对数据进行校验
	userService := service.UserService{}
	param := LoginParam{}
	err2 := c.ShouldBindJSON(&param)
	if err2 != nil {
		response.Fail(60002, "参数绑定有误", c)
		return
	}

	//if len(param.Code) == 0 {
	//	response.Fail(60002, "请输入验证码", c)
	//	return
	//}
	//
	//if len(param.CodeId) == 0 {
	//	response.Fail(60002, "验证码获取失败", c)
	//	return
	//}
	//
	//// 开始校验验证码是否正确
	//verify := store.Verify(param.CodeId, param.Code, true)
	//if !verify {
	//	response.Fail(60002, "你输入的验证码有误!!", c)
	//	return
	//}
	inputAccount := param.Account
	inputPassword := param.Password

	if len(inputAccount) == 0 {
		response.Fail(60002, "请输入账号", c)
		return
	}

	if len(inputPassword) == 0 {
		response.Fail(60002, "请输入密码", c)
		return
	}

	dbUser, err := userService.GetUserByAccount(inputAccount)
	if err != nil {
		response.Fail(60002, "您输入的账号或密码错误", c)
		return
	}

	// 校验用户的账号密码输入是否和数据库一致
	if dbUser != nil && dbUser.Password == inputPassword {
		// 1、jwt 生成token
		myJwt := jwtgo.NewJWT()
		// 2、生成token
		token, err2 := myJwt.CreateToken(jwtgo.CustomClaims{
			dbUser.ID,
			dbUser.Name,
			int64(1545),
			jwt.RegisteredClaims{
				Audience:  jwt.ClaimStrings{"XZ-USER"},                                            // 受众
				Issuer:    "MWJ-ADMIN",                                                            // 签发者
				IssuedAt:  jwt.NewNumericDate(time.Now()),                                         // 签发时间
				NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),                              // 生效时间
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Second * 60 * 60 * 24 * 7)), // 过期时间 7天
			},
		})
		fmt.Println("当前时间是：", time.Now().Unix())
		fmt.Println("签发时间：" + time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println("生效时间：" + time.Now().Add(-1000*time.Second).Format("2006-01-02 15:04:05"))
		fmt.Println("过期时间：" + time.Now().Add(1*time.Second*60*60*24*7).Format("2006-01-02 15:04:05"))
		if err2 != nil {
			response.Fail(60002, "登录失败,token颁发不成功！", c)
			return
		}
		response.Ok(map[string]any{"user": dbUser, "token": token}, c)

	} else {
		response.Fail(60002, "你输入的账号和密码有误", c)
	}
}

```

## 更新login的路由

router-> login -> login.go

```go
package login

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/api/v1/login"
)

// 登录路由
type LoginRouter struct{}

func (router *LoginRouter) InitLoginRouter(Router *gin.Engine) {
	loginApi := login.LoginApi{}
	// 单个定义
	//Router.GET("/login/toLogin", loginApi.ToLogined)
	//Router.GET("/login/toReg", loginApi.ToLogined)
	//Router.GET("/login/forget", loginApi.ToLogined)

	// 用组定义 ---》 推荐
	loginRouter := Router.Group("/login")
	{
		loginRouter.POST("/toLogin", loginApi.ToLogined)
	}
}
```

## 新建启动服务初始化

initlization -> init_server.go

```go
/*
@Author: 梦无矶小仔
@Date:   2024/1/15 13:48
*/
package initlization

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type server interface {
	ListenAndServe() error
}

func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

```

## 路由初始化时加入JWT

initlization -> init_router.go

```go
package initlization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"xz-go-frame/commons/filter"
	"xz-go-frame/global"
	"xz-go-frame/middle"
	"xz-go-frame/router"
	"xz-go-frame/router/code"
	"xz-go-frame/router/login"
)

func InitGinRouter() *gin.Engine {
	// 初始化 gin 服务
	ginServer := gin.Default()

	// 提供服务组
	videoRouter := router.RouterWebGroupApp.Video.VideoRouter

	// 解决接口的跨域问题
	ginServer.Use(filter.Cors())
	// 登录接口
	loginRouter := login.LoginRouter{}
	// 验证码接口
	codeRouter := code.CodeRouter{}

	// 接口隔离，比如登录，健康检查都不需要拦截和做任何处理
	loginRouter.InitLoginRouter(ginServer)
	codeRouter.InitCodeRouter(ginServer)

	// 业务模块接口
	publicGroup := ginServer.Group("/api")
	// 只要是api接口都使用jwt拦截
	publicGroup.Use(middle.JWTAuth())
	{
		videoRouter.InitVideoRouter(publicGroup)
	}

	fmt.Println("router register success")
	return ginServer
}

func RunServer() {
	// 初始化路由
	Router := InitGinRouter()
	// 为用户头像和文件提供静态地址
	Router.StaticFS("/static", http.Dir("/static"))
	address := fmt.Sprintf(":%d", global.Yaml["server.port"])
	// 启动HTTP服务，courseController
	s := initServer(address, Router)
	// 保证文本顺序输出
	time.Sleep(10*time.Microsecond)
	
	s2 := s.ListenAndServe().Error()
	fmt.Println("服务启动完毕",s2)

}

```

## 修改main.go的启动函数

```go
func main() {
	//  开始初始化配置文件
	initlization.InitViper()
	fmt.Println("初始化配置文件成功！")
	// 初始化数据库
	initlization.InitMySQL()
	//开始初始化gin路由服务
	initlization.RunServer() // 修改的地方
	fmt.Println("启动xz-go-frame后端成功")

}
```



## token的时限多长才合适？

- 面对极度敏感的信息，如钱或银行数据，那就根本不要在本地存放Token，只存放在内存中。这样，随着App关闭，Token也就没有了。（一次性token）
- 将Token的时限设置成较短的时间（如1小时）。
- 对于那些虽然敏感但跟钱没关系，如健身App的进度，这个时间可以设置得长一点，如1个月。
- 对于像游戏或社交类App，时间可以更长些，半年或1年。

并且，文章还建议增加一个“Token吊销”过程来应对Token被盗的情形，类似于当发现银行卡或电话卡丢失，用户主动挂失的过程。

```go
github.com/songzhibin97/gkit
```

## 修改video接口的jwt验证

api -> v1 -> video -> video.go

```go
package video

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/commons/response"
)

type Video struct {
}

// 查询video
func (videoController *Video) FindVideos(c *gin.Context) {
	claims, _ := c.Get("claims")
	customClaims := claims.(*jwtgo.CustomClaims)
	fmt.Println(customClaims.UserId)
	fmt.Println(customClaims.Username)
	response.Ok("success FindVideos", c)
}

// 获取视频明细
func (videoController *Video) GetByID(c *gin.Context) {
	// 绑定参数用来获取/:id 这个方式
	// id := c.Param("id")
	// 绑定参数 ?ids = 1111
	claims, _ := c.Get("claims")
	customClaims := claims.(*jwtgo.CustomClaims)
	fmt.Println(customClaims.UserId)
	fmt.Println(customClaims.Username)
	response.Ok("success GetByID", c)
}
```

### 修改video的router

router -> video -> video.go

```go
package video

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/api/v1/video"
)
type VideoRouter struct {
}
func (videoRouter *VideoRouter) InitVideoRouter(group *gin.RouterGroup) {
	// 帖子路由
	videoApi := video.Video{}
	videoGroup := group.Group("video")
	{
		videoGroup.GET("find", videoApi.FindVideos)
		videoGroup.GET("get", videoApi.GetByID)
	}
}
```



## 校验jwt测试

校验，在Navicat连接数据库，创建几个用户账号密码用于测试。

![image-20240115141303997](images/image-20240115141303997.png)

后端打印了jwt的时间

```shell
当前时间是： 1705299158
签发时间：2024-01-15 14:12:38
生效时间：2024-01-15 13:55:58
过期时间：2024-01-22 14:12:38
```



请求头携带Autorization的token访问 `http://localhost:8088/api/video/find`

![image-20240115145820747](images/image-20240115145820747.png)

后端打印出来了解析出来的用户id和用户名

```go
1
xiaozai
```

## 实现JWT的续期

续期的原理：你只要在有效内的一个时间点，进行续期接口。续期其实就指：重新生成一个新的token。用新的token来替换旧的token。旧的token必须拉入黑名单。设置过期时间 

一句话：新老更替。

- 首先，创建一个token , 有效期：1分钟 60
- 在有效时间点上比如：50秒的时候，我就重新创建一个新的token .来替换旧的token即可。
- 然后在旧的token拉入黑名单，然后写入可以定时删除的内存中（redis）

**方案如下：**

如果：过期时间 - 当前时间 < 缓冲时间

- 重新创建一个新的token ，把新的token返回给浏览器，替换旧token
- 把旧token放入黑名单，然后继续查看是否是用新的token来进行请求了
- 为什么要旧的token放黑名单：有了新欢旧要忘记旧爱。（redis）
  - 思考题：那么这个旧token会不很多，答案是的，所以你要定时或者人工的去处理和清除token表
  - 解决方案：可以考虑使用redis。因为redis有自动删除和设定时间的能力。

### JWT黑名单处理

commons -> jwtgo -> jwt_black_list.go

```go
package jwtgo

import (
	"errors"
	"gorm.io/gorm"
	"xz-go-frame/global"
	"xz-go-frame/model/jwt"
)

type JwtService struct{}

func (jwtService *JwtService) JsonInBlacklist(jwtList jwt.JwtBlacklist) (err error) {
	err = global.XZ_DB.Create(&jwtList).Error
	if err != nil {
		return
	}
	return
}

func (jwtService *JwtService) IsBlacklist(jwttoken string) bool {
	//_, ok := global.BlackCache.Get(jwt)
	//return ok
	err := global.XZ_DB.Where("jwt = ?", jwttoken).First(&jwt.JwtBlacklist{}).Error
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !isNotFound
}

```

### 修改jwt中间件内容

middle -> jwt.go

```go
/*
@Author: 梦无矶小仔
@Date:   2024/1/15 11:36
*/
package middle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/commons/response"
	"xz-go-frame/utils"
	jwtdb "xz-go-frame/model/jwt"
)



var jwtService = jwtgo.JwtService{}

// 定义一个JWTAuth的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
        
		...

		// 增加jwt-token的续期功能
		// 判断过期时间 - now  < buffertime 就开始续期 ep 1d -- no
		fmt.Println("customClaims.ExpiresAt", customClaims.ExpiresAt)
		fmt.Println("time.Now().Unix()", time.Now().Unix())
		fmt.Println("customClaims.ExpiresAt - time.Now().Unix()", customClaims.ExpiresAt.Unix()-time.Now().Unix())
		fmt.Println("customClaims.BufferTime", customClaims.BufferTime)

		if customClaims.ExpiresAt.Unix()-time.Now().Unix() < customClaims.BufferTime {
			// 1、生成一个新的token
			// 2、用c把新的token返回页面
			fmt.Println("开始续期.....")
			// 获取7天的过期时间
			eptime, _ := utils.ParseDuration("7d")
			// 用当前时间 + eptime 就是新的token过期时间
			customClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(eptime))
			// 生成新的token
			newToken, _ := myJwt.CreateTokenByOldToken(token, *customClaims)
			// 输出给浏览器 --- request --- header --- 给服务端
			// 输出给浏览器 --- reponse --- header --- 给浏览器
			c.Header("new-authorization", newToken)
			c.Header("new-expires-at", strconv.FormatInt(customClaims.ExpiresAt.Unix(), 10))
			// 如果生成新的token了，旧的token怎么办？jwt没有提供一个机制让旧token失效。
			_ = jwtService.JsonInBlacklist(jwtdb.JwtBlacklist{Jwt: token})
		}

		// 让后续的路由方法可以直接通过c.Get("claims")
		c.Set("claims", customClaims)
		c.Next()
	}

}

```



## 简单封装login接口（含jwt-token）

api -> v1 -> login - > login.go

```go
package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mojocn/base64Captcha"
	"time"
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/commons/response"
	"xz-go-frame/model/user"
	service "xz-go-frame/service/user"
	"xz-go-frame/utils"
)

// 登录业务
type LoginApi struct {
}

// 1、定义验证码的store -- 默认是存储在go的内存中
var store = base64Captcha.DefaultMemStore

// 登录接口的处理
func (api *LoginApi) ToLogined(c *gin.Context) {
	type LoginParam struct {
		Account  string
		Code     string
		CodeId   string
		Password string
	}
	// 1、获取用户在页面上输入账号和密码，开始在数据库里对数据进行校验
	userService := service.UserService{}
	param := LoginParam{}
	err2 := c.ShouldBindJSON(&param)
	if err2 != nil {
		response.Fail(60002, "参数绑定有误", c)
		return
	}

	//if len(param.Code) == 0 {
	//	response.Fail(60002, "请输入验证码", c)
	//	return
	//}
	//
	//if len(param.CodeId) == 0 {
	//	response.Fail(60002, "验证码获取失败", c)
	//	return
	//}
	//
	//// 开始校验验证码是否正确
	//verify := store.Verify(param.CodeId, param.Code, true)
	//if !verify {
	//	response.Fail(60002, "你输入的验证码有误!!", c)
	//	return
	//}
	inputAccount := param.Account
	inputPassword := param.Password

	if len(inputAccount) == 0 {
		response.Fail(60002, "请输入账号", c)
		return
	}

	if len(inputPassword) == 0 {
		response.Fail(60002, "请输入密码", c)
		return
	}

	dbUser, err := userService.GetUserByAccount(inputAccount)
	if err != nil {
		response.Fail(60002, "您输入的账号或密码错误", c)
		return
	}

	// 校验用户的账号密码输入是否和数据库一致
	if dbUser != nil && dbUser.Password == inputPassword {
		token := api.generaterToken(c, dbUser)
		// 根据用户id查询用户的角色
		roles := [2]map[string]any{}
		m1 := map[string]any{"id": 1, "name": "超级管理员"}
		m2 := map[string]any{"id": 2, "name": "财务"}
		roles[0] = m1
		roles[1] = m2
		// 根据用户id查询用户的角色的权限
		permissions := [2]map[string]any{}
		pm1 := map[string]any{"code": 10001, "name": "保存用户"}
		pm2 := map[string]any{"code": 20001, "name": "删除用户"}
		permissions[0] = pm1
		permissions[1] = pm2

		response.Ok(map[string]any{"user": dbUser, "token": token, "roles": roles, "permissions": permissions}, c)

	} else {
		response.Fail(60002, "你输入的账号和密码有误", c)
	}
}

/*
根据用户信息创建一个token
*/
func (api *LoginApi) generaterToken(c *gin.Context, dbUser *user.User) string {
	// 设置token续期的缓冲时间
	bf, _ := utils.ParseDuration("1d")
	ep, _ := utils.ParseDuration("7d")

	// 1、jwt 生成token
	myJwt := jwtgo.NewJWT()
	// 2、生成token
	token, err2 := myJwt.CreateToken(jwtgo.CustomClaims{
		dbUser.ID,
		dbUser.Name,
		int64(bf / time.Second),
		jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"XZ-USER"},               // 受众
			Issuer:    "MWJ-ADMIN",                               // 签发者
			IssuedAt:  jwt.NewNumericDate(time.Now()),            // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 过期时间 7天
		},
	})
	fmt.Println("当前时间是：", time.Now().Unix())
	fmt.Println("缓冲时间是：", int64(bf/time.Second))
	fmt.Println("签发时间：" + time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("生效时间：" + time.Now().Add(-1000).Format("2006-01-02 15:04:05"))
	fmt.Println("过期时间：" + time.Now().Add(ep).Format("2006-01-02 15:04:05"))
	if err2 != nil {
		response.Fail(60002, "登录失败,token颁发不成功！", c)
	}
	return token
}

```

# 11 、 关于逻辑删除得问题和处理

## 01、gorm默认机制

### gorm.Model

GORM 定义一个 `gorm.Model` 结构体，其包括字段 `ID`、`CreatedAt`、`UpdatedAt`、`DeletedAt ` 。

- 其中这里得deletedAt就是用于逻辑删除控制得字段，如果null 就代表没有删除，如果有时间就说你执行过delete from 才会把删除时写入到数据库表中
- 如果有字段，未来做任何得查询都自动跟上条件deletedAt is null 



### 修改逻辑删除的默认规则

如果你先修改默认的规则，从时间变成0/1这种方式，你必须如下执行

1: 先安装组件

```go
gorm.io/plugin/soft_delete
```

2: 把deletedAt删掉，修改如下：

IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt;default:0" json:"isDeleted"`

```go
package global

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type GVA_MODEL struct {
	ID        uint      `gorm:"primarykey;comment:主键ID" json:"id"` // 主键ID
	CreatedAt time.Time `gorm:"type:datetime(0);autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime(0);autoUpdateTime;comment:更新时间" json:"updatedAt"`
	//DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"` // 删除时间
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt;default:0" json:"isDeleted"`
}

```

前面的那些步骤和流程一个都不能省去，比如：注册表

```go
package orm

import (
	"xkginweb/global"
	bbs2 "xkginweb/model/entity/bbs"
	"xkginweb/model/entity/jwt"
	sys2 "xkginweb/model/entity/sys"
	user2 "xkginweb/model/entity/user"
	video2 "xkginweb/model/entity/video"
)

func RegisterTable() {
	db := global.KSD_DB
	// 注册和声明model
	db.AutoMigrate(user2.XkUser{})
	db.AutoMigrate(user2.XkUserAuthor{})
	// 系统用户，角色，权限表
	db.AutoMigrate(sys2.SysApis{})
	db.AutoMigrate(sys2.SysMenus{})
	db.AutoMigrate(sys2.SysRoleApis{})
	db.AutoMigrate(sys2.SysRoleMenus{})
	db.AutoMigrate(sys2.SysRoles{})
	db.AutoMigrate(sys2.SysUserRoles{})
	db.AutoMigrate(sys2.SysUser{})
	// 视频表
	db.AutoMigrate(video2.XkVideo{})
	db.AutoMigrate(video2.XkVideoCategory{})
	db.AutoMigrate(video2.XkVideoChapterLesson{})
	// 社区
	db.AutoMigrate(bbs2.XkBbs{})
	db.AutoMigrate(bbs2.BbsCategory{})

	// 声明一下jwt模型
	db.AutoMigrate(jwt.JwtBlacklist{})
}

```



3：然后重启查看效果即可。

- 其中这里得isDeleted就是用于逻辑删除控制得字段，如果0就代表没有删除，如果是1就是删除
- 未来你执行任何的删除操作就变成update table set is_deleted = 1,update_time = now() where id = 1
- 未来做任何得查询都自动跟上条件is_deleted = 0



4: 我要把删除和未删除全部查询出来？

往往在做后台管理系统的时候，你就必须要把删除和未删除全部查询出来。那么你就必须加上：.Unscoped() 来进行处理这样会把默认机制打破。不在跟已删除过滤。如下：

```go
// 查询分页
func (service *SysUserService) LoadSysUserPage(info request.PageInfo) (list interface{}, total int64, err error) {
	// 获取分页的参数信息
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 准备查询那个数据库表
	db := global.KSD_DB.Model(&sys.SysUser{})

	// 准备切片帖子数组
	var sysUserList []sys.SysUser

	// 加条件
	if info.Keyword != "" {
		db = db.Where("(username like ? or account like ? )", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}

	// 排序默时间降序降序
	db = db.Order("created_at desc")

	// 查询中枢
	err = db.Unscoped().Count(&total).Error
	if err != nil {
		return sysUserList, total, err
	} else {
		// 执行查询
		err = db.Unscoped().Limit(limit).Offset(offset).Find(&sysUserList).Error
	}

	// 结果返回
	return sysUserList, total, err
}
```





# 目录

国际化

32 - 1-2 静态菜单顶部侧边 - 动态菜单 - 页面布局处理

- layout/Index.vue中进行布局
- layout的component里面就是头部、菜单、内容区域的各部分组件
- 侧边菜单颜色修改，菜单折叠
- 路由管理pinia
- 难点 - 动态路由

32 -3- 菜单表的设计

- elementplus中的dialog、message、messagebox（自己封装，像之前的kva一样）

- 后台系统中用户权限角色的设计
- 表格滚动展示，固定表头和跳转页面，只滚动中间表格内容，over-flow；auto，计算属性控制页面放大缩小表格显示问题
- kva.js这个文件需要自己写一遍封装

33 - 

- 面包屑、顶部页头 - homepageHeader.vue
- src -> components -> index.js 全自动化注册全局组件
- 后端增加用户退出的接口，清除已登录账户的token
- 头像组件Avatar - 分割线 divided属性
- 登录状态通过路由访问登录页面，需要进行处理

34 - 

- 控制面板前端内容整理
- 隐藏和折叠
- 菜单导航 - 标签页Tab
- 状态管理，stores -> menuTab.js
- 后端 -- 建表 - model - sys下面的所有，注册表

35-1

- 后端zap日志处理，增加log模块，了解gin的自带日志
- gorm日志和zap日志，init_gorm.go中写
- 用户表了名字，字段看下有没有变化
- 实现系统用户表的管理（.md文件中有步骤）
- 前端settings.js设置菜单的宽度适配

35-2

- 修改了登录页面
- 骨架屏 - elementplus - skeleton
- 数字动画（重复内容）
- 后端用户增删改查
- 密码加密

36

- 读取表的信息
- modle增加了vo和entity，数据载体
- 校验验证 - validate

37

- 逻辑删除-物理删除，0值处理，unscope

- 角色权限设置，表设计，前端展示设计

38

- 数据库 - 事务 - 原子性、隔离性、持久性、一致性
- 新增用户前后端

39

- gorm更新0值失败的问题
- throttle限制连续请求-防抖-前端

40

- 泛型约束

41

- 后端（少）
  - 菜单管理递归
- 前端(多)
  - 菜单管理 - 添加子菜单 - 删除 - 添加- 编辑

42

- 后端
  - 查询角色对应的菜单
- 前端
  - MenuTab的动态变化，关闭新增，前一个后一个关闭
  - Tab切换保持数据维持 - keepAlive
  - 页面缓存加载 Suspense default fallback
  - form响应该值问题，数据格式改成form.value
  - 角色切换

43

- 后端
  - 角色菜单,角色授予sys_role_menus.go
- 前端
  - PageSIdebar.vue 子菜单排序问题 sort((a,b)=>a.sort-b.sort)
  - 全局根据角色权限响应式变换页面

44

- 后端
  - sql批量插入
- 前端
  - menu授权
  - 角色授权管理右侧树状RoleMenuApis.vue

45-1

- 后端
  - 又讲了一下泛型
  - 权限分配，权限管理
  - sys_casbin权限拦截
- 前端
  - 权限管理页面

45-2

- 后端
  - 认识casbin
  - https://github.com/casbin/gorm-adapter
  - 数据库需要建一个表casbin_rule
- 前端 - 无

46

- 后端
  - 在login模块里面加
- 前端
  - 在user.js添加uuid

47

- 前半部分讲了下后端打包，项目部署，看笔记文档就行，没有难的地方
- 后半部分讲了下手动前端部署



# 前端部分

# 页面的布局处理



## 官网参考

https://element-plus.gitee.io/zh-CN/component/table.html#%E5%B8%A6%E8%BE%B9%E6%A1%86%E8%A1%A8%E6%A0%BC

## 后台管理系统

- 登录
- 首页
  - 控制面板
  - 应用管理
  - 系统管理
  - 订单管理
  - 设计管理
  - 用户管理
  - 角色权限
  - 评论管理
  - 等等等
- 404错误页面



## 前端路由设计

- 延后到beego

## 大布局

- 登录
- 首页
  - 上–头部信息
  - 下
    - 左边（菜单）
    - 右边 （点击菜单具体内容）
- 404错误页面



实现的方式就是：vue-router

### 1：安装

```js
npm install vue-router@next
pnpm install vue-router@next
yarn add vue-router@next
```

### 2: 在src下定义router目录并新建index.js

```js
import { createRouter, createWebHashHistory } from 'vue-router'
import NProgress from 'nprogress'
// 获取状态管理的token
import { useUserStore } from '@/stores/user.js'
import Layout from "@/layout/index.vue";
import PageFrame from "@/layout/components/PageFrame.vue";
// 显示右上角螺旋加载提示
NProgress.configure({ showSpinner: true })

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "Home",
      component: Layout
  ]
})

export default router

```

### 3：  注册路由，在main.js中如下：

```js
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
const app = createApp(App)
app.use(router)
app.mount('#app')
```

### 4： 开始定义路由页面SPA

这页面都定义在views目录下, 除非是大页面或者公共页面，我们可以直接放入到views目录下。比如：404，login.vue，index.vue。如果是页面的模块建议建设模块的目录然后在里面定义的路由页面。比如：/logs/Operation.vue, /sys/User.vue,/sys/Role.vue

![image-20230729202421316](images/image-20230729202421316.png)



**==而且注意：页面路由的命名尽量用大写开头。==**

### 5:  静态路由注册

给每个SPA页面进行注册路由。找到router/index.js开始进行静态配置如下：

如果报错请安装：

```js
pnpm install nprogress
npm install nprogress
yarn add nprogress
```

```js
import { createRouter, createWebHashHistory } from 'vue-router'
import NProgress from 'nprogress'
// 获取状态管理的token
import { useUserStore } from '@/stores/user.js'
// 显示右上角螺旋加载提示
NProgress.configure({ showSpinner: true })

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "Home",
      component: () => import('@/views/Index.vue')
    },
    {
      path: '/newindex',
      name: 'Newindex',
      meta: { title: "newindex" },
      component: () => import('@/views/NewIndex.vue')
    },
    {
      path: '/login',
      name: 'Login',
      meta: { title: "login" },
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})



router.beforeEach(async (to) => {
  //开启进度条
  NProgress.start()
  const useStore = useUserStore();
  // 判断是否登录
  if (!useStore.isLogin && to.name !== 'Login') {
    // 这里的query就是为了记录用户最后一次访问的路径，这个路径是通过to的参数获取
    // 后续在登录成功以后，就可以根据这个path的参数，然后调整到你最后一次访问的路径
    return { name: 'Login', query: { "path": to.path } }
  } 
  return true
})

router.afterEach(() => {
  //完成进度条
  NProgress.done()
})


export default router

```

###  设置页面的落脚点router-view的位置

这个落脚点就是：App.vue 其实是入口页面。我准备准备一个空白页面

```vue
<template>
    <router-view></router-view>
</template>
```

关键代码：`   <router-view></router-view>`

router-路由，view-视图。你通过vue-router访问所有的路由，都会在路由视图位置进行渲染和加载

路由的加载又两种方式：

- 异步加载：component: () => import('@/views/NotFound.vue')

- 静态加载：

  - import Home from "@/views/Index.vue";

  - ```JS
    const router = createRouter({
      history: createWebHashHistory(import.meta.env.BASE_URL),
      routes: [
        {
          path: "/",
          name: "Home",
          component: Home
        }
      ]
    })
    ```

### 6: 启动项目

- http://127.0.0.1:8777/#/login
- http://127.0.0.1:8777/#/
- http://127.0.0.1:8777/#/newindex
- http://127.0.0.1:8777/#/loginxxxx—-进入404错误页面





## 大布局设计Layout

- 登录—设计完毕
- 首页
  - 上–头部信息
  - 下
    - 左边（菜单）
    - 右边 （点击菜单具体内容）
- 404错误页面——设计完毕

首页如何设计，其实在很多后台系统，都使用布局来命名layout. 在目录下新建一个layout目录。然后新建一个index.vue如下：

![image-20230729205050140](images/image-20230729205050140.png)

layout/index.vue如下：

```vue
<template>
  <div>我是首页</div>
</template>
```

开始在router/index.js去进行配置首页路由：

```js
import { createRouter, createWebHashHistory } from 'vue-router'
import NProgress from 'nprogress'
// 获取状态管理的token
import { useUserStore } from '@/stores/user.js'
// 显示右上角螺旋加载提示
NProgress.configure({ showSpinner: true })
import Layout from "@/layout/Index.vue";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "Home",
      component: Layout
    },
    {
      path: '/newindex',
      name: 'Newindex',
      meta: { title: "newindex" },
      component: () => import('@/views/NewIndex.vue')
    },

    {
      path: '/login',
      name: 'Login',
      meta: { title: "login" },
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})



router.beforeEach(async (to) => {
  //开启进度条
  NProgress.start()
  const useStore = useUserStore();
  // 判断是否登录
  if (!useStore.isLogin && to.name !== 'Login') {
    // 这里的query就是为了记录用户最后一次访问的路径，这个路径是通过to的参数获取
    // 后续在登录成功以后，就可以根据这个path的参数，然后调整到你最后一次访问的路径
    return { name: 'Login', query: { "path": to.path } }
  } 
  return true
})

router.afterEach(() => {
  //完成进度条
  NProgress.done()
})


export default router

```

重新启动：`npm run dev` 关闭上次启按多次：ctrl+c 即可。

启动访问：http://localhost:8777/#/

![image-20230729205542754](images/image-20230729205542754.png)

在layout/Index.vue页面中开始进行布局，因为如果把头部，菜单，内容的展示区域都放在一起的话，其实是没问题。但是后续的维护是变得困难和臃肿，所以要进行组件化，（其实就是把布局首页进行分割成若干子组件，然后在Index.vue汇合）如下：

1: 头部组件 PageHeader.vue

```vue
<template>
  <div class="header-cont">
   头部组件
  </div>
</template>
<script setup>


</script>
<style lang="scss">
.header-cont {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding:0 20px;
  a {
    color: inherit;
    text-decoration: none;
  }
  h1 {
    margin: 0;
    font-size: 20px;
  }
  .gap {
    margin-right: 20px;
  }
  .right {
    .lang {
      font-size: 14px;
      .item {
        cursor: pointer;
        &.active {
          font-size: 18px;
          font-weight: bold;
        }
      }
    }
  }
  .el-dropdown {
    color: inherit;
  }
}
</style>
```

2：菜单组件 PageSilder.vue

```vue
<template>
<div class="page-sidebar">
   左侧菜单
</div>
</template>

<script  setup>
</script>

<style lang="scss">
$side-width: 200px;
.page-sidebar {
  height: 100vh;
  background: #000;
  .sidemenu.el-menu,
  .sidemenu .el-sub-menu > .el-menu {
    --el-menu-text-color: #ccc;
    --el-menu-hover-bg-color: #060251;
    --el-menu-border-color: transparent;
    --el-menu-bg-color: #000;
    .el-menu-item {
      &.is-active {
        background-color: var(--el-menu-hover-bg-color);
      }
    }
  }
  .sidemenu.el-menu:not(.el-menu--collapse) {
    width: $side-width;
  }
  .collape-bar {
    color: #fff;
    font-size: 16px;
    line-height: 36px;
    text-align: center;

    .c-icon {
      cursor: pointer;
    }
  }
}
</style>
```

3：内容区域组件 PageMain.vue

```vue
<template>
  <router-view></router-view>
</template>
<script>
export default {
  
}
</script>
<style lang="">
  
</style>

```

然后访问首页如下：http://localhost:8777/#/

![image-20230729211755648](images/image-20230729211755648.png)

假设在开发中我们有一个控制面板菜单，而这个菜单的内容必须渲染到layout的右侧内容区域也就是pagemain.vue的位置，如何实现呢？

1： 在page-main组件设定一个router-view标记

2： 然后把 dashboard路由设置成为layout的在子路由如下：

```js
import { createRouter, createWebHashHistory } from 'vue-router'
import NProgress from 'nprogress'
// 获取状态管理的token
import { useUserStore } from '@/stores/user.js'
// 显示右上角螺旋加载提示
NProgress.configure({ showSpinner: true })
import Layout from "@/layout/Index.vue";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "Home",
      component: Layout,
      children:[
        {
          path: '/dashboard',
          name: 'Dashboard',
          meta: { title: "dashboard" },
          component: () => import('@/views/Dashboard.vue')
        }    
      ]
    },    
    {
      path: '/login',
      name: 'Login',
      meta: { title: "login" },
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})



router.beforeEach(async (to) => {
  //开启进度条
  NProgress.start()
  const useStore = useUserStore();
  // 判断是否登录
  if (!useStore.isLogin && to.name !== 'Login') {
    // 这里的query就是为了记录用户最后一次访问的路径，这个路径是通过to的参数获取
    // 后续在登录成功以后，就可以根据这个path的参数，然后调整到你最后一次访问的路径
    return { name: 'Login', query: { "path": to.path } }
  } 
  return true
})

router.afterEach(() => {
  //完成进度条
  NProgress.done()
})


export default router

```

当访问http://localhost:8777/#/dashboard 的时候就会把dashborad路由对应的spa页面/views/DashBoard.vue的页面内容渲染到layout/page-main.vue的router-view的位置。从而实现点击菜单展示到右侧，后续的业务菜单原理是一样的。如下：

![image-20230729213120041](images/image-20230729213120041.png)





## 静态菜单–不查数据库

其实通过上面的案例其实就已经很清晰了。你需要自己去定义一个列表或者使用组件来形成一个菜单组件。从而给每个菜单项绑定一个路由地址，你就可以实现了所有的菜单结构。

### 自定义菜单

router/index.js把所有的模块都注册进来

```js
import { createRouter, createWebHashHistory } from 'vue-router'
import NProgress from 'nprogress'
// 获取状态管理的token
import { useUserStore } from '@/stores/user.js'
// 显示右上角螺旋加载提示
NProgress.configure({ showSpinner: true })
import Layout from "@/layout/Index.vue";
import PageMain from "@/layout/components/PageMain.vue";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "Home",
      component: Layout,
      children:[
        {
          path: 'dashboard',
          name: 'Dashboard',
          meta: { title: "dashboard" },
          component: () => import('@/views/Dashboard.vue')
        },
        {
          path: "app",
          name: "App",
          meta: {title:"app"},
          redirect: '/app/user',
          component: PageMain,
          children: [
            {
              path: "user",
              name: "AppUser",
              meta: {title:"AppUser"},
              component: () => import("@/views/app/User.vue"),
            },
            {
              path: "dept",
              name: "AppDept",
              meta: {title:"AppDept"},
              component: () => import("@/views/app/Dept.vue"),
            },
            {
              path: "role",
              name: "AppRole",
              meta: {title:"AppRole"},
              component: () => import("@/views/app/Role.vue"),
            },
            {
              path: "resource",
              name: "AppResource",
              meta: {title:"AppResource"},
              component: () => import("@/views/app/Resource.vue"),
            }
          ],
        },
        {
          path: 'sys',
          meta: {title:"sys"},
          redirect: '/sys/user',
          component: PageMain,
          children: [
            {
              path: "user",
              name: "SysUser",
              meta: {title:"SysUser"},
              component: () => import("@/views/sys/User.vue"),
            },
            {
              path: "notice",
              name: "SysNotice",
              meta: {title:"SysNotice"},
              component: () => import("@/views/sys/Notice.vue"),
            }
          ],
        },
        {
          path: "logs",
          name: "LogsManagement",
          meta: {title:"logs"},
          component: PageMain,
          redirect: '/logs/visit',
          children: [
            {
              path: "visit",
              name: "VisitsLog",
              meta: {title:"VisitsLog"},
              component: () => import("@/views/logs/Visit.vue"),
            },
            {
              path: "operation",
              name: "OprationsLog",
              meta: {title:"OprationsLog"},
              component: () => import("@/views/logs/Operation.vue"),
            }
          ],
        },       
        ]
      },    
    {
      path: '/login',
      name: 'Login',
      meta: { title: "login" },
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})



router.beforeEach(async (to) => {
  //开启进度条
  NProgress.start()
  const useStore = useUserStore();
  // 判断是否登录
  if (!useStore.isLogin && to.name !== 'Login') {
    // 这里的query就是为了记录用户最后一次访问的路径，这个路径是通过to的参数获取
    // 后续在登录成功以后，就可以根据这个path的参数，然后调整到你最后一次访问的路径
    return { name: 'Login', query: { "path": to.path } }
  } 
  return true
})

router.afterEach(() => {
  //完成进度条
  NProgress.done()
})


export default router

```

自定义菜单-pagesilder.vue

```vue
<template>
<div class="page-sidebar">
   <ul>
      <li><a style="color:#fff" href="/#/dashboard">控制面板</a></li>
      <li><a style="color:#fff" href="/#/app/user">用户</a></li>
      <li><a style="color:#fff" href="/#/app/role">角色</a></li>
      <li><a style="color:#fff" href="/#/app/dept">部门</a></li>
      <li><a style="color:#fff" href="/#/app/resource">资源</a></li>
   </ul>
   <ul>
      <li><a style="color:#fff" href="/#/logs/visit">日志访问</a></li>
      <li><a style="color:#fff" href="/#/logs/operation">日志operation</a></li>
   </ul>
   <ul>
      <li><a style="color:#fff" href="/#/sys/user">系统用户</a></li>
      <li><a style="color:#fff" href="/#/sys/notice">系统通知</a></li>
   </ul>
</div>
</template>

<script  setup>
</script>

<style lang="scss">
$side-width: 200px;
.page-sidebar {
  height: 100vh;
  background: #000;
  .sidemenu.el-menu,
  .sidemenu .el-sub-menu > .el-menu {
    --el-menu-text-color: #ccc;
    --el-menu-hover-bg-color: #060251;
    --el-menu-border-color: transparent;
    --el-menu-bg-color: #000;
    .el-menu-item {
      &.is-active {
        background-color: var(--el-menu-hover-bg-color);
      }
    }
  }
  .sidemenu.el-menu:not(.el-menu--collapse) {
    width: $side-width;
  }
  .collape-bar {
    color: #fff;
    font-size: 16px;
    line-height: 36px;
    text-align: center;

    .c-icon {
      cursor: pointer;
    }
  }
}
</style>
```

然后点击访问各个模块就可以在右侧进行展示和显示了。如下：

![image-20230729215035476](images/image-20230729215035476.png)

未来如果你要添加系统角色如下：步骤如下：

1： 在views/sys/Role.vue

2:  然后在sys路由下继续添加一个路由子项:

```js
import { createRouter, createWebHashHistory } from 'vue-router'
import NProgress from 'nprogress'
// 获取状态管理的token
import { useUserStore } from '@/stores/user.js'
// 显示右上角螺旋加载提示
NProgress.configure({ showSpinner: true })
import Layout from "@/layout/Index.vue";
import PageMain from "@/layout/components/PageMain.vue";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "Home",
      component: Layout,
      children:[
        {
          path: 'dashboard',
          name: 'Dashboard',
          meta: { title: "dashboard" },
          component: () => import('@/views/Dashboard.vue')
        },
        {
          path: "app",
          name: "App",
          meta: {title:"app"},
          redirect: '/app/user',
          component: PageMain,
          children: [
            {
              path: "user",
              name: "AppUser",
              meta: {title:"AppUser"},
              component: () => import("@/views/app/User.vue"),
            },
            {
              path: "dept",
              name: "AppDept",
              meta: {title:"AppDept"},
              component: () => import("@/views/app/Dept.vue"),
            },
            {
              path: "role",
              name: "AppRole",
              meta: {title:"AppRole"},
              component: () => import("@/views/app/Role.vue"),
            },
            {
              path: "resource",
              name: "AppResource",
              meta: {title:"AppResource"},
              component: () => import("@/views/app/Resource.vue"),
            }
          ],
        },
        {
          path: 'sys',
          meta: {title:"sys"},
          redirect: '/sys/user',
          component: PageMain,
          children: [
            {
              path: "user",
              name: "SysUser",
              meta: {title:"SysUser"},
              component: () => import("@/views/sys/User.vue"),
            },
            {
              path: "notice",
              name: "SysNotice",
              meta: {title:"SysNotice"},
              component: () => import("@/views/sys/Notice.vue"),
            },
            {
              path: "role",
              name: "SysRole",
              meta: {title:"SysRole"},
              component: () => import("@/views/sys/Role.vue"),
            }
          ],
        },
        {
          path: "logs",
          name: "LogsManagement",
          meta: {title:"logs"},
          component: PageMain,
          redirect: '/logs/visit',
          children: [
            {
              path: "visit",
              name: "VisitsLog",
              meta: {title:"VisitsLog"},
              component: () => import("@/views/logs/Visit.vue"),
            },
            {
              path: "operation",
              name: "OprationsLog",
              meta: {title:"OprationsLog"},
              component: () => import("@/views/logs/Operation.vue"),
            }
          ],
        },       
        ]
      },    
    {
      path: '/login',
      name: 'Login',
      meta: { title: "login" },
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})



router.beforeEach(async (to) => {
  //开启进度条
  NProgress.start()
  const useStore = useUserStore();
  // 判断是否登录
  if (!useStore.isLogin && to.name !== 'Login') {
    // 这里的query就是为了记录用户最后一次访问的路径，这个路径是通过to的参数获取
    // 后续在登录成功以后，就可以根据这个path的参数，然后调整到你最后一次访问的路径
    return { name: 'Login', query: { "path": to.path } }
  } 
  return true
})

router.afterEach(() => {
  //完成进度条
  NProgress.done()
})


export default router

```

然后在page-silder.vue去增加一个菜单

```vue
<template>
<div class="page-sidebar">
   <ul>
      <li><a style="color:#fff" href="/#/dashboard">控制面板</a></li>
      <li><a style="color:#fff" href="/#/app/user">用户</a></li>
      <li><a style="color:#fff" href="/#/app/role">角色</a></li>
      <li><a style="color:#fff" href="/#/app/dept">部门</a></li>
      <li><a style="color:#fff" href="/#/app/resource">资源</a></li>
   </ul>
   <ul>
      <li><a style="color:#fff" href="/#/logs/visit">日志访问</a></li>
      <li><a style="color:#fff" href="/#/logs/operation">日志operation</a></li>
   </ul>
   <ul>
      <li><a style="color:#fff" href="/#/sys/user">系统用户</a></li>
      <li><a style="color:#fff" href="/#/sys/notice">系统通知</a></li>
      <li><a style="color:#fff" href="/#/sys/role">系统角色</a></li>
   </ul>
</div>
</template>

<script  setup>
</script>

<style lang="scss">
$side-width: 200px;
.page-sidebar {
  height: 100vh;
  background: #000;
  .sidemenu.el-menu,
  .sidemenu .el-sub-menu > .el-menu {
    --el-menu-text-color: #ccc;
    --el-menu-hover-bg-color: #060251;
    --el-menu-border-color: transparent;
    --el-menu-bg-color: #000;
    .el-menu-item {
      &.is-active {
        background-color: var(--el-menu-hover-bg-color);
      }
    }
  }
  .sidemenu.el-menu:not(.el-menu--collapse) {
    width: $side-width;
  }
  .collape-bar {
    color: #fff;
    font-size: 16px;
    line-height: 36px;
    text-align: center;

    .c-icon {
      cursor: pointer;
    }
  }
}
</style>
```

http://localhost:8777/#/sys/role

![image-20230729215404923](images/image-20230729215404923.png)

这种静态配置有利也有弊端，利：比较清晰和简单，一个页面一个路由一个菜单。然后访问，刷新也会自动定位。

### 配置思维

新建一个js

```js
export const menuTreeData = [
    {
        id:1,
        parentId:0,
        name:'App',
        path:'/app',
        icon:'menu',
        children:[
            {
                id:11,
                parentId:1,
                name:'AppUser',
                path:'/app/user',
                icon:'user',
            },
            {
                id:12,
                parentId:1,
                name:'AppDept',
                path:'/app/dept',
                icon:'office-building',
            },
            {
                id:13,
                parentId:1,
                name:'AppRole',
                path:'/app/role',
                icon:'avatar',
            },
            {
                id:14,
                parentId:1,
                name:'AppResource',
                path:'/app/resource',
                icon:'manage',
            }
        ]
    },
    {
        id:2,
        parentId:0,
        name:'Sys',
        path:'/sys',
        icon:'setting',
        children:[
            {
                id:21,
                parentId:1,
                name:'SysUser',
                path:'/sys/user',
                icon:'user-filled',
            },
            {
                id:22,
                parentId:1,
                name:'SysNotice',
                path:'/sys/notice',
                icon:'',
            }
        ]
    },
    {
        id:3,
        parentId:0,
        name:'Logs',
        path:'/logs',
        icon:'document',
        children:[
            {
                id:31,
                parentId:1,
                name:'LogsVisit',
                path:'/logs/visit',
                icon:'tickets',
            },
            {
                id:32,
                parentId:1,
                name:'LogsOperation',
                path:'/logs/operation',
                icon:'operation',
            }
        ]
    }
]


```

然后读取js的数据信息，进行处理，同时还可以进行国际化处理。你只要菜单名字使用国际化配置的key的名字然后使用{{ t(key)}} 就可以读取国际化的内容。如下：

```VUE
<template>
  <div class="page-sidebar">
    <div class="collape-bar">
      <el-icon class="cursor" @click="isCollapse = !isCollapse">
        <expand v-if="isCollapse" />
        <fold v-else />
      </el-icon>
    </div>
    <el-menu
      active-text-color="#ffd04b"
      background-color="#000000"
      text-color="#fff"
      router
      :default-active="defaultActive"
      class="sidemenu"
      :collapse="isCollapse"
      @open="handleOpen"
      @close="handleClose"
    >
      <el-sub-menu v-for="(item, i) in menuTree" :key="i" :index="item.path">
        <template #title>
          <el-icon v-if="item.icon"><component :is="item.icon"></component></el-icon>
          <span>{{ t(`menu.${item.name}`) }}</span>
        </template>
        <template v-for="(child, ci) in item.children" :key="ci">
            <el-menu-item :index="child.path">
              <el-icon><component :is="child.icon"></component></el-icon>
              {{ t(`menu.${child.name}`) }}
            </el-menu-item>
          </template>
      </el-sub-menu>
    </el-menu>
  </div>
  </template>
  
  <script  setup>
  const route = useRoute();
  const { t } = useI18n();
  const isCollapse = ref(false)
  const menuTree = ref([
    {
      id: 1,
      parentId: 0,
      name: 'App',
      path: "/app",
      icon: "location",
      children: [
        {
          id: 11,
          parentId: 1,
          name: 'AppUser',
          path: "/app/user",
          icon: "user",
        },
        {
          id: 12,
          parentId: 1,
          name: 'AppDept',
          path: "/app/dept",
          icon: "office-building",
        },
        {
          id: 13,
          parentId: 1,
          name: 'AppRole',
          path: "/app/role",
          icon: "avatar",
        },
        {
          id: 14,
          parentId: 1,
          name: 'AppResource',
          path: "/app/resource",
          icon: "management",
        },
      ],
    },
    {
      id: 2,
      parentId: 0,
      name: 'Sys',
      path: "/sys",
      icon: "setting",
      children: [
        {
          id: 21,
          parentId: 2,
          name: 'SysUser',
          path: "/sys/user",
          icon: "user-filled",
        },
        {
          id: 22,
          parentId: 2,
          name: 'SysNotice',
          path: "/sys/notice",
          icon: "chat-dot-round",
        },
      ],
    },
    {
      id: 3,
      parentId: 0,
      name: 'Logs',
      path: "/logs",
      icon: "document",
      children: [
        {
          id: 31,
          parentId: 3,
          name: 'LogsVisit',
          path: "/logs/visit",
          icon: "tickets",
        },
        {
          id: 32,
          parentId: 3,
          name: 'LogsOperation',
          path: "/logs/operation",
          icon: "operation",
        },
      ],
    },
  ])
  
  const defaultActive = computed(() => route.path || menuTree.value[0].path)
  const handleOpen = (key, keyPath) => {
    console.log(key, keyPath)
  }
  const handleClose = (key, keyPath) => {
    console.log(key, keyPath)
  }
  </script>
  
  <style lang="scss">
  $side-width: 200px;
  .page-sidebar {
    .sidemenu.el-menu,
    .sidemenu .el-sub-menu > .el-menu {
      --el-menu-text-color: #ccc;
      --el-menu-hover-bg-color: #060251;
      --el-menu-border-color: transparent;
      --el-menu-bg-color: #000;
      .el-menu-item {
        &.is-active {
          background-color: var(--el-menu-hover-bg-color);
        }
      }
    }
    .sidemenu.el-menu:not(.el-menu--collapse) {
      width: $side-width;
    }

    .collape-bar {
      color: #fff;
      font-size: 16px;
      line-height: 36px;
      text-align: center;
  
      .c-icon {
        cursor: pointer;
      }
    }
  }
  </style>
```



## 菜单的路由访问和选中的问题

通过上节课的处理我们使用el-menu的菜单组件，那么菜单组件你发现并没有编写的。如何完成el-menu菜单定位路由同事又可以选择和激活（选中）呢？方式如下：

- index 

  - 把index设置成路由地址，

  - 同时增加router属性true

  - ![image-20230730201617056](images/image-20230730201617056.png)

    

- default-active

  这个激活的当前当前，也就说如果你设定的default-active和某个index的值的一致，那么就自动激活这个菜单，并且选中。其实在开发中我们更多希望达到的效果当前访问路径是什么。那么就选择什么。并且刷新以后也会自动定位到当前访问路由的菜单处。

如何获取到当前访问的路由路径呢？

```js
import {useRoute} from 'vue-router'
// 这个是用来获取当前访问的路由信息,
const route = useRoute();
// 根据当前路由来激活菜单
const defaultActive = computed(()=>(route.path))

```



## 选中菜单颜色改变

使用属性来激活你选中的菜单的文本颜色

```vue
 <el-menu 
      active-text-color="#ffd04b" 
   >
```

也通过css来覆盖它

```html
<div class="page-siderbar">
 <el-menu 
      active-text-color="#ffd04b" 
      class="sidemenu"
   >
```

```css
.page-sidebar {
  .sidemenu.el-menu,
  .sidemenu .el-sub-menu>.el-menu {
    --el-menu-text-color: #ccc;
    --el-menu-hover-bg-color: #060251;
    --el-menu-border-color: transparent;
    --el-menu-bg-color: #000;

    .el-menu-item {
      &.is-active {
        background-color: var(--el-menu-hover-bg-color);
        color: #42d51d
      }
    }
  }
```

## 菜单的折叠的问题



```vue
<template>
	<el-menu 
      active-text-color="#ffd04b" 
      background-color="#000000" 
      text-color="#fff" 
      router
      :default-active="defaultActive" 
      class="sidemenu" 
      :collapse="isCollapse" 
      @open="handleOpen" 
      @close="handleClose">
....
    
</template>
 <script>
     // 默认情况下不折叠
    const isCollapse = ref(false)
 </script>
<style>
	/* elmenu菜单的折叠效果是通过属性：
  :collapse="isCollapse"  原理就在控制在不停切换elmenu="el-menu--collapse"样式信息
  1：ture 就折叠，就会使用图标宽度+padding作为菜单宽度
  2: false 就不折叠，那么就使用默认宽度：200px

  下面这行css是什么意思：
  如果菜单上存在el-menu--collapse样式就说明是折叠状态，就使用图标宽度+padding作为菜单宽度
  否则：就用我的width:200作为菜单宽度
*/
.sidemenu.el-menu:not(.el-menu--collapse) {
width: 200px;
}

</style>
```

## 菜单的国际化问题

1: 准备国际化的组件安装，看前面课程

2: 先准备国际化的配置信息

menu.js是中文：

```js
export default {
  menu: {
    DashBoard:"控制面板",
    App: "应用管理",
    AppUser: "用户管理",
    AppDept: "机构管理",
    AppRole: "角色管理",
    AppResource: "资源管理",
    AppPermission: "授权管理",
    Sys: "系统管理",
    SysUser: "用户管理",
    SysNotice: "公告管理",
    Logs: "审计管理",
    LogsVisit: "访问日志",
    LogsOperation: "操作日志",
  }
}

```

menuEn.js

```js
export default {
  menu: {
    DashBoard:"DashBoard",
    App: "Website",
    AppUser: "User",
    AppDept: "Department",
    AppRole: "Role",
    AppResource: "Resource",
    AppPermission: "Permission",
    Sys: "System",
    SysUser: "User",
    SysNotice: "Notice",
    Logs: "Logs",
    LogsVisit: "Visits",
    LogsOperation: "Operations",
  }
}

```

3：定义菜单的数据

```js
const menuTree = ref([
  {
    id: 4,
    parentId: 0,
    name: 'DashBoard',
    path: "/dashboard",
    icon: "home",
    children:[]
  },
  {
    id: 1,
    parentId: 0,
    name: 'App',
    path: "/app",
    icon: "location",
    children: [
      {
        id: 11,
        parentId: 1,
        name: 'AppUser',
        path: "/app/user",
        icon: "user",
      },
      {
        id: 12,
        parentId: 1,
        name: 'AppDept',
        path: "/app/dept",
        icon: "office-building",
      },
      {
        id: 13,
        parentId: 1,
        name: 'AppRole',
        path: "/app/role",
        icon: "avatar",
      },
      {
        id: 14,
        parentId: 1,
        name: 'AppResource',
        path: "/app/resource",
        icon: "management",
      },
    ],
  },
  {
    id: 2,
    parentId: 0,
    name: 'Sys',
    path: "/sys",
    icon: "setting",
    children: [
      {
        id: 21,
        parentId: 2,
        name: 'SysUser',
        path: "/sys/user",
        icon: "user-filled",
      },
      {
        id: 22,
        parentId: 2,
        name: 'SysNotice',
        path: "/sys/notice",
        icon: "chat-dot-round",
      },
    ],
  },
  {
    id: 3,
    parentId: 0,
    name: 'Logs',
    path: "/logs",
    icon: "document",
    children: [
      {
        id: 31,
        parentId: 3,
        name: 'LogsVisit',
        path: "/logs/visit",
        icon: "tickets",
      },
      {
        id: 32,
        parentId: 3,
        name: 'LogsOperation',
        path: "/logs/operation",
        icon: "operation",
      },
    ],
  }
])
```

注意这里菜单的name并不是明文，而是国际化的key的名字，你只要保持一致，然后在菜单显示的时候使用国际化的api方法t方法就可以把国际化中与之对应的语言的菜单信息显示出来。

```js
<template>  
    <span>{{ t('menu.AppUser') }}</span>
</template>
<script setup>
const { t } = useI18n();
</script>

```

## 关于项目中vue.vue-router,pinia,vuex自动导入的问题

在项目的vite.config.js文件中配置自动导入插件如下：

```js
import { fileURLToPath, URL } from 'node:url'
import { defineConfig,loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  return {
    // vite 配置
    plugins: [
      vue(),
      AutoImport({
        imports: ['vue', 'vue-router', 'pinia', 'vue-i18n'],
        resolvers: [ElementPlusResolver()],
      }),
      Components({
        resolvers: [ElementPlusResolver()],
      })
    ]
  }
})



```

这样配置以后，在后续的SPA或者SFC页面中，你可以省去vue,vue-router vuex ,vue-i18n的导入过程

```js
import { useI18n } from "vue-i18n";
import { ref,computed } from "vue";
import { useRoute } from "vue-router";

// 这个是用来获取当前访问的路由信息,
const route = useRoute();
const { t } = useI18n();
const isCollapse = ref(false)
// 根据当前路由来激活菜单
const defaultActive = computed(()=>(route.path))
```

配置以后你可以简化如下：

```js
// 这个是用来获取当前访问的路由信息,
const route = useRoute();
const { t } = useI18n();
const isCollapse = ref(false)
// 根据当前路由来激活菜单
const defaultActive = computed(()=>(route.path))
```

## 动态菜单—查数据库

关于动态菜单的数据获取的问题。一般有如下方式

1： 在登录的时候 （推荐）

2：路由的beforeEach 前置守卫中获取

在登录时候，我们会直接把用户对应角色，角色对应菜单和权限会全部查询出来。如下：

因为暂时还没有和真正数据库打交道。可以使用mock数据进行测试如下：

```js
export const users = [
    {
        name:"visitor",
        roleId:"visitor",
        password:"visitor"
    },
    {
        name:"master",
        roleId:"master",
        password:"master"
    },
    {
        name:"admin",
        roleId:"admin",
        password:"admin"
    }
]

// 模拟服务端角色对应菜单信息--超级管理员
export const menuTreeData = [
    {
        id:4,
        parentId:0,
        name:'DashBoard',
        path:'/dashboard',
        icon:'home',
        children:[]
    },
    {
        id:1,
        parentId:0,
        name:'App',
        path:'/app',
        icon:'menu',
        children:[
            {
                id:11,
                parentId:1,
                name:'AppUser',
                path:'/app/user',
                icon:'user',
            },
            {
                id:12,
                parentId:1,
                name:'AppDept',
                path:'/app/dept',
                icon:'office-building',
            },
            {
                id:13,
                parentId:1,
                name:'AppRole',
                path:'/app/role',
                icon:'avatar',
            },
            {
                id:14,
                parentId:1,
                name:'AppResource',
                path:'/app/resource',
                icon:'avatar',
            }
        ]
    },
    {
        id:2,
        parentId:0,
        name:'Sys',
        path:'/sys',
        icon:'setting',
        children:[
            {
                id:21,
                parentId:1,
                name:'SysUser',
                path:'/sys/user',
                icon:'user-filled',
            },
            {
                id:22,
                parentId:1,
                name:'SysNotice',
                path:'/sys/notice',
                icon:'user-filled',
            }
        ]
    },
    {
        id:3,
        parentId:0,
        name:'Logs',
        path:'/logs',
        icon:'document',
        children:[
            {
                id:31,
                parentId:1,
                name:'LogsVisit',
                path:'/logs/visit',
                icon:'tickets',
            },
            {
                id:32,
                parentId:1,
                name:'LogsOperation',
                path:'/logs/operation',
                icon:'operation',
            }
        ]
    }
]
```

然后在store/user.js中开始引入mock数据、同时定义menuTree的数据用于接收服务端或者mock的测试数据如下：

```js
import { defineStore } from 'pinia'
import request from '@/request'
import router from '@/router'
import { menuTreeData } from '@/mock/data.js'

//https://blog.csdn.net/weixin_62897746/article/details/129124364
//https://prazdevs.github.io/pinia-plugin-persistedstate/guide/
export const useUserStore = defineStore('user', {
  // 定义状态
  state: () => ({
    user: {},
    username: '',
    userId: '',
    token: '',
    age:10,
    male:1,
    role:[],
    permissions:[],
    // 路由菜单，用来接收服务端传递过来的菜单数据
    menuTree:[]
  }),

  // 就是一种计算属性的机制，定义的是函数，使用的是属性就相当于computed
  getters:{

    malestr(state){
      if(state.male==1)return "男"
      if(state.male==0)return "女"
      if(state.male==1)return "保密"
    },

    isLogin(state){
      return state.token ? true : false
    },

    roleName(state){
      return state.roles && state.roles.map(r=>r.name).join(",")
    },

    permissionCode(state){
      return state.permissions &&  state.permissions.map(r=>r.code).join(",")
    }
  },

  // 定义动作
  actions: {
   setToken(newtoken){
      this.token = newtoken
   },

   getToken(){
    return this.token
   },
   
   /* 登出*/
   async LoginOut (){
      this.token = ''
      this.user = {}
      this.username = ''
      this.userId = ''
      sessionStorage.clear()
      localStorage.clear()
      router.push({ name: 'Login', replace: true })
      window.location.reload()
  },
  
   async toLogin(loginUser){

      // 查询用户信息，角色，权限，角色对应菜单
      const resp = await request.post("login/toLogin", loginUser,{noToken:true})
      // 这个会回退，回退登录页
      var { user ,token,roles,permissions } = resp.data
      // 登录成功以后获取到菜单信息, 这里要调用一
      this.menuTree = menuTreeData
      // 把数据放入到状态管理中
      this.user = user
      this.userId = user.id
      this.username = user.name
      this.token = token
      this.roles = roles
      this.permissions = permissions
      return Promise.resolve(resp)
    }
  },
  persist: {
    key: 'kva-pinia-userstore',
    storage: localStorage,//sessionStorage
  }
})
```

然后你在pageSider.vue中可以读取到菜单信息如下：

```vue
<template>
  <div class="page-sidebar">
    <div class="collape-bar">
      <el-icon class="cursor" @click="isCollapse = !isCollapse">
        <expand v-if="isCollapse" />
        <fold v-else />
      </el-icon>
    </div>
    <el-menu 
      active-text-color="#ffd04b" 
      background-color="#000000" 
      text-color="#fff" 
      router
      :default-active="defaultActive" 
      class="sidemenu" 
      :collapse="isCollapse">
      <template v-for="(item, i) in menuTree" :key="i">
        <template v-if="item.children && item.children.length">
          <el-sub-menu :index="item.path">
            <template #title>
              <el-icon v-if="item.icon">
                <component :is="item.icon"></component>
              </el-icon>
              <span>{{ t(`menu.${item.name}`) }}</span>
            </template>
              <template v-for="(child, ci) in item.children" :key="ci">
                <el-menu-item :index="child.path">
                  <el-icon>
                    <component :is="child.icon"></component>
                  </el-icon>
                  {{ t(`menu.${child.name}`) }}
                </el-menu-item>
            </template>
          </el-sub-menu>
        </template>
        <template v-else>
          <el-menu-item :index="item.path">
            <el-icon v-if="item.icon">
              <component :is="item.icon"></component>
            </el-icon>
            <span>{{ t(`menu.${item.name}`) }}</span>
          </el-menu-item>
        </template>
      </template>
    </el-menu>
  </div>
</template>
  
<script  setup>
import { useUserStore } from '@/stores/user.js'
// 这个是用来获取当前访问的路由信息,
const route = useRoute();
const { t } = useI18n();
// 默认情况下不折叠
const isCollapse = ref(false)
// 根据当前路由来激活菜单
const defaultActive = computed(()=>(route.path))
// 获取状态管理的菜单信息
const userStore = useUserStore();
// 如何获取菜单数据呢？
const menuTree = computed(()=>userStore.menuTree)

</script>
<style lang="scss">
.page-sidebar {
  .sidemenu.el-menu,
  .sidemenu .el-sub-menu>.el-menu {
    --el-menu-text-color: #ccc;
    --el-menu-hover-bg-color: #060251;
    --el-menu-border-color: transparent;
    --el-menu-bg-color: #000;

    .el-menu-item {
      &.is-active {
        background-color: var(--el-menu-hover-bg-color);
        color: #42d51d
      }
    }
  }

   /* elmenu菜单的折叠效果是通过属性：
      :collapse="isCollapse"  原理就在控制在不停切换elmenu="el-menu--collapse"样式信息
      1：ture 就折叠，就会使用图标宽度+padding作为菜单宽度
      2: false 就不折叠，那么就使用默认宽度：200px

      下面这行css是什么意思：
      如果菜单上存在el-menu--collapse样式就说明是折叠状态，就使用图标宽度+padding作为菜单宽度
      否则：就用我的width:200作为菜单宽度
   */
  .sidemenu.el-menu:not(.el-menu--collapse) {
    width: 200px;
  }

  .collape-bar {
    color: #fff;
    font-size: 16px;
    line-height: 36px;
    text-align: center;

    .c-icon {
      cursor: pointer;
    }
  }
}
</style>
```

然后刷新访问，如果无效可以关闭服务器重新启动把缓存全部清空在登录在尝试如下：

![image-20230730211551695](images/image-20230730211551695.png)



## 动态菜单的思考

你觉得上面之所以可以，是因为现在数据库里的菜单和vue-router/index.js配置的路由菜单以及在views的定义菜单对应路由的spa都存在。所以你能够一个非常正常效果。

- views/mode/spa
- vue-router/index.js也配置spa的路由对应访问路径
- 数据库里刚好和他们一致。

但是往往开发中，你觉得会出现什么问题。比如：

- 一个添加了菜单，但是没有添加spa也没有添加vue-router肯定不可以访问。

如果你按照正常流程你应该是怎么做：

- 第一步：先定义spa
- 第二步：配置一个路由来调整spa页面 
- 第三步：在数据库在配置一个菜单



所谓的动态菜单：其实就简化后续模块的新增和变化，不需要在vue-router/index.js配置路由，然后在数据库中添加一个菜单，然后在views添加一个页面，你就访问到对应的菜单。

具体实现步骤：

1: 把所有home首页子路由全部去掉如下：

```js
import { createRouter, createWebHashHistory } from 'vue-router'
import NProgress from 'nprogress'
// 获取状态管理的token
import { useUserStore } from '@/stores/user.js'
// 显示右上角螺旋加载提示
NProgress.configure({ showSpinner: true })
import Layout from "@/layout/Index.vue";
import PageMain from "@/layout/components/PageMain.vue";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "Home",
      component: Layout,
      redirect:"/dashboard"
    },    
    {
      path: '/login',
      name: 'Login',
      meta: { title: "login" },
      component: () => import('@/views/Login.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})



router.beforeEach(async (to) => {
  //开启进度条
  NProgress.start()
  const useStore = useUserStore();
  // 判断是否登录
  if (!useStore.isLogin && to.name !== 'Login') {
    // 这里的query就是为了记录用户最后一次访问的路径，这个路径是通过to的参数获取
    // 后续在登录成功以后，就可以根据这个path的参数，然后调整到你最后一次访问的路径
    return { name: 'Login', query: { "path": to.path } }
  } 
  return true
})

router.afterEach(() => {
  //完成进度条
  NProgress.done()
})


export default router

```

2: 然后通过vue-router的js方法完成动态路由的创建，

一句话：就使用js的方法来完成首页children的拼接和处理。

而这个children的数据是通过：服务端返回的菜单数据和views定义路由页面组合而成的。

```js
children = [
    {
       // 这些从复服务端返回的菜单数据来
      path: 'dashboard',
      name: 'Dashboard',
      meta: { title: "dashboard" },
      // 这个从views里而来 ，这里命名就非常关键
      component: () => import('@/views/Dashboard.vue')
    }
]
```

```js
router.addRoute(“Home”,children)
```





## vue-router

https://router.vuejs.org/zh/guide/advanced/dynamic-routing.html

### 为什么要配置动态路由

- 因为菜单是由数据库来管理的
  - 如果在数据库里添加一个菜单，你必须要做两个事情你添加到数据库路由才会由意义
  - 第一个条件：必须在router/index.js的routes进行定义和配置
  - 第二个条件：你必须还为路由配置路由页面/views/xxx/xxxx.vue
- 因为后续开发我们的菜单是根据用户的角色来决定和显示的。我们如果去写死的话，就不灵活
- 你思考如果你角色由10个，那么你必须配置10不同的菜单，而10个菜单数据，其实就可能有些多了一些菜单，有些少了一些菜单，那么为什么不把他们存在数据库，来进行一种关联映射，然后定义接口。根据用户查询对应角色，然后根据角色查询对应菜单。然后在返回菜单数据。那么不完美了么。
- 我们唯独在后台中要去开发的就如果还把菜单添加，然后动态绑定给角色。然后又如何把角色绑定给用户。



## 动态路由又什么好处呢？

-  简化的前端路由的注册的过程，其实就完成了手动静态注册的过程。如果实在不喜欢，你可以考虑还原成静态注册方式。



## 关于elementplus中的dialog,message,messagebox的认识和使用



### messageBox

https://element-plus.gitee.io/zh-CN/component/message-box.html

从设计上来说，MessageBox 的作用是美化系统自带的 `alert`、`confirm` 和 `prompt`，因此适合展示较为简单的内容。 如果需要弹出较为复杂的内容，请使用 Dialog。

```js
import { ElMessage, ElMessageBox } from 'element-plus'
```

alert

```js
ElMessageBox.alert('This is a message', 'Title', {
    // if you want to disable its autofocus
    // autofocus: false,
    confirmButtonText: 'OK',
    callback: (action: Action) => {
      ElMessage({
        type: 'info',
        message: `action: ${action}`,
      })
    },
  })
```

confirm

```js
ElMessageBox.confirm(
    'proxy will permanently delete the file. Continue?',
    'Warning',
    {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  )
    .then(() => {
      ElMessage({
        type: 'success',
        message: 'Delete completed',
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: 'Delete canceled',
      })
    })
```

prompt

```js
ElMessageBox.prompt('Please input your e-mail', 'Tip', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    inputPattern:
      /[\w!#$%&'*+/=?^_`{|}~-]+(?:\.[\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[\w](?:[\w-]*[\w])?/,
    inputErrorMessage: 'Invalid Email',
  })
    .then(({ value }) => {
      ElMessage({
        type: 'success',
        message: `Your email is:${value}`,
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: 'Input canceled',
      })
    })
```

使用 HTML 片段

```js
 ElMessageBox.alert(
    '<strong>proxy is <i>HTML</i> string</strong>',
    'HTML String',
    {
      dangerouslyUseHTMLString: true,
    }
  )
```

自己封装的处理

```js
const KVA = {
    alert(title,content,options){
        // 默认值
        var defaultOptions = {icon:"warning",confirmButtonText:"确定",cancelButtonText:"取消"}
        // 用户传递和默认值就行覆盖处理
        var opts = {...defaultOptions,...options}
        return ElMessageBox.alert(content, title,{
            //确定按钮文本
            confirmButtonText: opts.confirmButtonText,
            // 内容支持html
            dangerouslyUseHTMLString: true,
            // 是否支持拖拽
            draggable: true,
            // 修改图标
            type: opts.icon
        })
    },
    confirm(title,content,options){
        // 默认值
        var defaultOptions = {icon:"warning",confirmButtonText:"确定",cancelButtonText:"取消"}
        // 用户传递和默认值就行覆盖处理
        var opts = {...defaultOptions,...options}
        // 然后提示
        return ElMessageBox.confirm(content, title, {
           //确定按钮文本
           confirmButtonText: opts.confirmButtonText,
           //取消按钮文本
           cancelButtonText: opts.cancelButtonText,
           // 内容支持html
           dangerouslyUseHTMLString: true,
           // 是否支持拖拽
           draggable: true,
           // 修改图标
           type: opts.icon,
        })
    },
    prompt(title,content,options){
        // 默认值
        var defaultOptions = {confirmButtonText:"确定",cancelButtonText:"取消"}
        // 用户传递和默认值就行覆盖处理
        var opts = {...defaultOptions,...options}
        return ElMessageBox.prompt(content, title, {
            //确定按钮文本
            confirmButtonText: opts.confirmButtonText,
            //取消按钮文本
            cancelButtonText: opts.cancelButtonText,
            // 内容支持html
            dangerouslyUseHTMLString: true,
            // 是否支持拖拽
            draggable: true,
            // 输入框的正则验证
            inputPattern: opts.pattern,
            // 验证的提示内容
            inputErrorMessage: opts.message||'请输入正确的内容',
          })
    }
}

export default  KVA
```

其实在elementplus注册到vue的时候如下：

```js
import ElementPlus from 'element-plus'
app.use(ElementPlus)
```

就会常用message,messagebox,dialog等都注册一份挂载到`app.config.globalProperties` 下。添加如下全局方法：`$msgbox`、 `$alert`、 `$confirm` 和 `$prompt`。 因此在 Vue 实例中可以采用本页面中的方式来调用`MessageBox`。 参数如下：

- `$msgbox(options)`
- `$alert(message, title, options)` 或 `$alert(message, options)`
- `$confirm(message, title, options)` 或 `$confirm(message, options)`
- `$prompt(message, title, options)` 或 `$prompt(message, options)`



那么如何拿到这个app.config.globalProperties对象。步骤如下：

```js
import {getCurrentInstance} from 'vue'
const {proxy} = getCurrentInstance();

proxy.$alert(message, title, options)
proxy.$confirm(message, title, options)
proxy.$prompt(message, title, options)
```

如果是vue2的话

```js
this.$alert(message, title, options)
this.$confirm(message, title, options)
this.$prompt(message, title, options)
```

## 样式调整问题

base.css

```css
html,
body {
  height: 100%;
}
#app {
  height: 100%;
  overflow: hidden;
}
.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}
.cursor {
  cursor: pointer;
}
.txt-c {
  text-align: center;
}
.w100p {
  width: 100%;
}

::-webkit-scrollbar {
  width: 6px;
  height: 6px;
  background-color: #f5f5f5;
}

::-webkit-scrollbar-button {
  background-color: rgba(0, 0, 0, 0.1);
}

::-webkit-scrollbar-track {
  background-color: #f5f5f5;
  border-radius: 5px;
}

::-webkit-scrollbar-thumb {
  background-color: #c9c9c9;
  border-radius: 5px;
}

.kva-container{padding:5px;overflow: hidden;}
.kva-contentbox{background:#fff;padding:15px;}
.kva-pagination-box{margin-top: 15px;display: flex;}
.kva-pagination-box.left{justify-content: flex-start;}
.kva-pagination-box.center{justify-content: center;}
.kva-pagination-box.right{justify-content: flex-end;}
.kva-form-search{border-bottom: 1px solid #e7e7e7;}
```

pagesilder.vue

```vue
<template>
  <div class="page-sidebar">
    <div class="collape-bar">
      <el-icon class="cursor" @click="isCollapse = !isCollapse">
        <expand v-if="isCollapse" />
        <fold v-else />
      </el-icon>
    </div>
    <el-menu 
      active-text-color="#333" 
      background-color="#ffffff" 
      text-color="#333" 
      router
      :default-active="defaultActive" 
      class="sidemenu" 
      :collapse="isCollapse">
      <template v-for="(item, i) in menuTree" :key="i">
        <template v-if="item.children && item.children.length">
          <el-sub-menu :index="item.path">
            <template #title>
              <el-icon v-if="item.icon">
                <component :is="item.icon"></component>
              </el-icon>
              <span>{{ t(`menu.${item.name}`) }}</span>
            </template>
              <template v-for="(child, ci) in item.children" :key="ci">
                <el-menu-item :index="child.path">
                  <el-icon>
                    <component :is="child.icon"></component>
                  </el-icon>
                  {{ t(`menu.${child.name}`) }}
                </el-menu-item>
            </template>
          </el-sub-menu>
        </template>
        <template v-else>
          <el-menu-item :index="item.path">
            <el-icon v-if="item.icon">
              <component :is="item.icon"></component>
            </el-icon>
            <span>{{ t(`menu.${item.name}`) }}</span>
          </el-menu-item>
        </template>
      </template>
    </el-menu>
  </div>
</template>
  
<script  setup>
import { useUserStore } from '@/stores/user.js'
// 这个是用来获取当前访问的路由信息,
const route = useRoute();
const { t } = useI18n();
// 默认情况下不折叠
const isCollapse = ref(false)
// 根据当前路由来激活菜单
const defaultActive = computed(()=>(route.path))
// 获取状态管理的菜单信息
const userStore = useUserStore();
// 如何获取菜单数据呢？
const menuTree = computed(()=>userStore.menuTree)

</script>
<style lang="scss">
$slider-width: 180px;
.page-sidebar {
  height: calc(100vh - 90px);
  overflow: hidden auto;
  .sidemenu.el-menu,
  .sidemenu .el-sub-menu>.el-menu {
    --el-menu-text-color: #ccc;
    --el-menu-hover-bg-color: #060251;
    --el-menu-border-color: transparent;
    --el-menu-bg-color: #001529;

    .el-menu-item {
      &.is-active {
        background-color: #e6f7ff;
        color: #1890ff
      }
    }
  }

   /* elmenu菜单的折叠效果是通过属性：
      :collapse="isCollapse"  原理就在控制在不停切换elmenu="el-menu--collapse"样式信息
      1：ture 就折叠，就会使用图标宽度+padding作为菜单宽度
      2: false 就不折叠，那么就使用默认宽度：200px

      下面这行css是什么意思：
      如果菜单上存在el-menu--collapse样式就说明是折叠状态，就使用图标宽度+padding作为菜单宽度
      否则：就用我的width:200作为菜单宽度
   */
  .sidemenu.el-menu:not(.el-menu--collapse) {
    width: $slider-width;
  }

  .collape-bar {
    color: #333;
    font-size: 16px;
    line-height: 36px;
    position: fixed;
    z-index: 2;
    width: 100%;
    left:20px;
    bottom: 0;
    .c-icon {
      cursor: pointer;
    }
  }
}
</style>
```

## 面包屑的处理和思考

https://element-plus.gitee.io/zh-CN/component/breadcrumb.html#%E5%9F%BA%E7%A1%80%E7%94%A8%E6%B3%95

有时候位了让我们路径更加的清晰，和让操作者知道你所在菜单的位置，一般会在右侧的页面增加面包屑

==作用：显示当前页面的路径，快速返回之前的任意页面。==

```vue
<el-breadcrumb separator="/">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>---root
  <el-breadcrumb-item>系统管理</el-breadcrumb-item>-----parent
  <el-breadcrumb-item>用户管理</el-breadcrumb-item>-----children
</el-breadcrumb>
```

如果使用手动的方式，每个页面去增加，这样维护起来是非常麻烦的。如果后续发生变化和变动，以及国际化的处理都会变得非常的麻烦。怎么办？其实你可以这样思考。我们访问路径是不可以拿到。`route.path` —–`/sys/user` 

- children === path = /sys/user—–然后开始遍历获取到信息—-SysUser—-然后使用国际化进行处理即可
- parent==== path =/sys 然后开始遍历获取到信息—-Sys—-然后使用国际化进行处理即可



实现步骤：

1: 定义组件 HomePageHeader.vue

```vue
<template>
    <div>
        <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">
                首页
            </el-breadcrumb-item>
            <el-breadcrumb-item v-if="parentName">
                <a href="javascript:void(0);">{{ t('menu.'+parentName) }}</a>
            </el-breadcrumb-item>
            <el-breadcrumb-item v-if="isChildren">
                <a href="javascript:void(0);">{{ t('menu.'+route.meta.name) }}</a>
            </el-breadcrumb-item>
        </el-breadcrumb>
        <div style="padding:15px 0">
            <slot></slot>
        </div>
    </div>
</template>
<script setup>
// 获取到当前路由
const route = useRoute()
// 获取国际化
const { t } = useI18n();
// 判断是不是有子元素, 因为在菜单中存在一种没有子的情况，这个时候就没有第二级。
const isChildren = ref(true)
// 获取菜单数据
import { menuTreeData } from '@/mock/data.js'

console.log('route',route)

// 开始截取当前的访问路径，比如：/sys/user
let parentPath = route.path.substring(0,route.path.indexOf('/',2))//得到的是：/sys
if(!parentPath){ 
    parentPath  = route.path
    // 代表你没有子元素
    isChildren.value = false;
}
// 如果有子元素，可以把去查找菜单信息
const parentName = menuTreeData.find(obj=>obj.path==parentPath).name

</script>
<style>
</style>
```

2: 组件必须要进行注册

建议使用全局注册，这样就不需要每个spa页面进行引入以后才能使用如下：

vue3插件机制如下：

```js
import HomePageHeader from './HomePageHeader.vue'
export default {
    install(app){
        // const modules = import.meta.glob('../components/**/*.vue');
        // for(let key in modules){
        //     var componentName = key.substring(key.lastIndexOf('/')+1,key.lastIndexOf("."))
        //     app.component(componentName,modules[key])
        // }

        // 全局注册组件
        app.component("HomePageHeader",HomePageHeader)
    }
}
```

你也可以使用全自动注册

```js
export default {
    install(app){
        // 全自动化过程注册全局组件，就不需要在引入在注册
        // 把src/components目录下的以.vue结尾的文件全部匹配出来。包括子孙目录下的.vue结尾的文件
         const modules = import.meta.glob('../components/**/*.vue');
         for(let key in modules){
             var componentName = key.substring(key.lastIndexOf('/')+1,key.lastIndexOf("."))
             app.component(componentName,modules[key])
         }
    }
}
```

这样的好处就是，不需要你增加组件，又来到处和注册一次，省去了这个步骤。

然后在main.js进行插件生效注册。

```js
import { createApp } from 'vue'
import KVAComponents from '@/components'

const app = createApp(App)
app.use(KVAComponents)
```

3: 然后在每个需要使用面包屑地方进行使用` <home-page-header>` 进行包裹即可，如下：

```vue
<template>
  <div class="kva-container">
    <div class="kva-contentbox">
      <home-page-header>
        <div class="kva-form-search">
          <el-form :inline="true" :model="queryParams">
            <el-form-item>
              <el-button type="primary" v-permission="[10001]" icon="Plus" @click="handleAdd">添加</el-button>
              <el-button type="danger"  v-permission="[20001]" icon="Delete" @click="handleDel">删除</el-button>
            </el-form-item>
            <el-form-item label="关键词：">
              <el-input v-model="queryParams.keyword" placeholder="请输入搜索关键词..." maxlength="10" clearable />
            </el-form-item>
            <el-form-item label="关键词：">
              <el-input v-model="queryParams.name" placeholder="请输入名字..." maxlength="10" clearable />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="Search" @click="onSubmit">搜索</el-button>
            </el-form-item>
          </el-form>
        </div>
        <el-table :data="tableData" style="width: 100%" height="calc(100vh - 218px)">
          <el-table-column fixed prop="date" label="Date" width="150" />
          <el-table-column prop="name" label="Name" width="120" />
          <el-table-column prop="state" label="State" width="120" />
          <el-table-column prop="city" label="City" width="320" />
          <el-table-column prop="address" label="Address" />
          <el-table-column fixed="right" prop="zip" label="Zip" width="120" />
        </el-table>
        <div class="kva-pagination-box">
          <el-pagination
            v-model:current-page="currentPage4"
            v-model:page-size="pageSize4"
            :page-sizes="[100, 200, 300, 400]"
            :small="small"
            :disabled="disabled"
            :background="background"
            layout="total, sizes, prev, pager, next, jumper"
            :total="400"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </home-page-header>
    </div>
  </div>
</template>

<script  setup>
import KVA from '@/utils/kva.js'
import { ElMessage } from 'element-plus';
const {proxy} = getCurrentInstance();

// 搜索属性定义
const queryParams = reactive({
name:"",
keyword:""
})

// 添加事件
const handleAdd = ()=>{
KVA.notify("注册提示","感谢你注册平台,<a href=''>点击此处进入查看</a>",3000,{type:"success",position:"br"})
}

// 删除事件
// const handleDel = async ()=>{
//   try{
//     const response =  await KVA.confirm("警告","你确定要抛弃我么？",{icon:"info"})
//     alert("去请求你要删除的异步请求的方法把")
//   }catch(e){
//     alert("你点击的是关闭或者取消按钮")
//   }
// }

const handleDel =  ()=>{
  KVA.confirm("警告","<strong>你确定要抛弃我么？</strong>",{icon:"success"}).then(()=>{
    KVA.message("去请求你要删除的异步请求的方法把")
  }).catch(err=>{
    KVA.error("你点击的是关闭或者取消按钮")
    //proxy.$message({message:"你点击的是关闭或者取消按钮",type:"success",showClose:true})
    //proxy.$message({message:"你点击的是关闭或者取消按钮",type:"warining",showClose:true})
    //proxy.$message({message:"你点击的是关闭或者取消按钮",type:"error",showClose:true})
  })

  // proxy.$confirm("<strong>你确定要抛弃我么？</strong>","警告",{type:"success",dangerouslyUseHTMLString:true}).then(()=>{
  //   alert("去请求你要删除的异步请求的方法把")
  // }).catch(err=>{
  //   alert("你点击的是关闭或者取消按钮")
  // })
}


const tableData = [
{
  date: '2016-05-03',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-02',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-04',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-01',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-08',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-06',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-07',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-03',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-02',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-04',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-01',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-08',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-06',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-07',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-03',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-02',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-04',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-01',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-08',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-06',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-07',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-03',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-02',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-04',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-01',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-08',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-06',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-07',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-03',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-02',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-04',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-01',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-08',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-06',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
{
  date: '2016-05-07',
  name: 'Tom',
  state: 'California',
  city: 'Los Angeles',
  address: 'No. 189, Grove St, Los Angeles',
  zip: 'CA 90036',
},
]
</script>



```

## 头部处理

PageHeader.vue页面

```vue
<template>
  <div class="header-cont">
    <div class="left">
      <h1>
        <router-link to="/">{{ t('uniLiveMangeSystem') }}</router-link>
      </h1>
    </div>
    <div class="right flex-center">
      <div class="lang gap">
        <span
          class="item"
          :class="{ active: locale === 'zh-cn' }"
          @click="changeLanguage('zh-cn')"
        >简体中文</span>
        /
        <span
          class="item"
          :class="{ active: locale === 'en' }"
          @click="changeLanguage('en')"
        >EN</span>
      </div>
      <template v-if="isLogin">
        <div class="gap">
          <router-link to="/personal/message">
            <el-badge :is-dot="!!unReadCount">
              <el-icon>
                <message />
              </el-icon>
            </el-badge>
          </router-link>
        </div>
        <el-dropdown trigger="click" @command="handleCommand">
          <div class="flex-center cursor">
            {{ username }}
            <el-icon>
              <caret-bottom />
            </el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="toPersonal">{{ t('personalCenter') }}</el-dropdown-item>
              <el-dropdown-item command="toLogout">{{ t('logout') }}</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>
      <template v-else-if="$route.name !== 'Login'">
        <router-link to="/login">{{ t('login') }}</router-link>
      </template>
    </div>
  </div>
</template>
<script setup>
import { useUserStore } from '@/stores/user.js'
const store = useUserStore()
const router = useRouter();
const { locale, t } = useI18n();
const isLogin = computed(() => store.token);
const userInfo = computed(() => store.user);
const username = computed(() => store.username)
const unReadCount = computed(() => 100);

const commands = ({
  toPersonal: () => {
    router.push('/personal')
  },
  toLogout: () => {
    store.LoginOut();
  }
});

// 语言切换
function changeLanguage(lang) {
  locale.value = lang
  localStorage.setItem('ksd-kva-language', lang)
}

function handleCommand(command) {
  commands[command] && commands[command]();
}

onMounted(()=>{
  locale.value = localStorage.getItem("ksd-kva-language") || 'zh-cn'
})

</script>
<style lang="scss">
.header-cont {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding:0 20px;
  a {
    color: inherit;
    text-decoration: none;
  }
  h1 {
    margin: 0;
    font-size: 20px;
  }
  .gap {
    margin-right: 20px;
  }
  .right {
    .lang {
      font-size: 14px;
      .item {
        cursor: pointer;
        &.active {
          font-size: 18px;
          font-weight: bold;
        }
      }
    }
  }
  .el-dropdown {
    color: inherit;
  }
}
</style>
```



### 国际化

```vue
<template>
  <div class="header-cont">
    <div class="left">
      <h1>
        <router-link to="/">{{ t('KvaAdminHome') }}</router-link>
      </h1>
    </div>
    
    <div class="right flex-center">
      <div class="lang gap">
        <span
          class="item"
          :class="{ active: locale === 'zh-cn' }"
          @click="changeLanguage('zh-cn')"
        >简体中文</span>
        /
        <span
          class="item"
          :class="{ active: locale === 'en' }"
          @click="changeLanguage('en')"
        >EN</span>
      </div>
      <template v-if="isLogin">
        <div class="gap">
          <router-link to="/personal/message">
            <el-badge :is-dot="!!unReadCount">
              <el-icon>
                <message />
              </el-icon>
            </el-badge>
          </router-link>
        </div>
        <el-dropdown trigger="click" @command="handleCommand">
          <div class="flex-center cursor">
            {{ username }}
            <el-icon>
              <caret-bottom />
            </el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="toPersonal">{{ t('personalCenter') }}</el-dropdown-item>
              <el-dropdown-item command="toLogout">{{ t('logout') }}</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>
      <template v-else-if="$route.name !== 'Login'">
        <router-link to="/login">{{ t('login') }}</router-link>
      </template>
    </div>
  </div>
</template>
<script setup>
import { useUserStore } from '@/stores/user.js'
const store = useUserStore()
const router = useRouter();
const { locale, t } = useI18n();
const isLogin = computed(() => store.token);
const username = computed(() => store.username)
const unReadCount = computed(() => 100);

const commands = ({
  toPersonal: () => {
    router.push('/personal')
  },
  toLogout: () => {
    store.LoginOut();
  }
});


function handleCommand(command) {
  commands[command] && commands[command]();
}



// 语言切换
function changeLanguage(lang) {
  // 把选择的语言进行切换
  locale.value = lang
  // 切换以后记得把本地缓存进行修改，否则只会生效当前，刷新就还原。
  localStorage.setItem('ksd-kva-language', lang)
}

// 用于读取本地缓存存储的语言是什么？
function initReadLocale(){
  locale.value = localStorage.getItem("ksd-kva-language") || 'zh-cn'
}

onMounted(()=>{
  initReadLocale();
})

</script>
<style lang="scss">
.header-cont {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding:0 20px;
  a {
    color: inherit;
    text-decoration: none;
  }
  h1 {
    margin: 0;
    font-size: 20px;
  }
  .gap {
    margin-right: 20px;
  }
  .right {
    .lang {
      font-size: 14px;
      .item {
        cursor: pointer;
        &.active {
          font-size: 16px;
          font-weight: bold;
        }
      }
    }
  }
  .el-dropdown {
    color: inherit;
  }
}
</style>
```

但是这里国际化仅限于我们自己控制的，那么element-plus组件那些国际化如下处理呢？如下：

在按照elememtplus的时候其实默认情况它语言包其实已经下载下来了。你只需要引入即可如下：找到i18n目录下的/index.js导入elementplus的国际化js文件如下：

```js
import { createI18n } from 'vue-i18n'
import zhLocale from './lang/zh'
import enLocale from './lang/en'
import zhCn from 'element-plus/es/locale/lang/zh-cn' //--------------------------------------这里是新增
import en from 'element-plus/es/locale/lang/en' //--------------------------------------这里是新增

const i18n = createI18n({
  legacy:false,
  fallbackLocale:'zh',
  locale:  localStorage.getItem("ksd-kva-language") || 'zh-cn', // 设置地区
  messages: {
    en: {
      ...enLocale,
      ...zhCn//--------------------------------------这里是新增
    },
    'zh-cn': {
      ...zhLocale,
      ...en //--------------------------------------这里是新增
    }
  }
})

export default i18n

export const elementLocales = { //--------------------------------------这里是新增
  'zh-cn': zhCn,
  en
}

```

然后找到App.vue下面的router-view包裹一一个标签组件如下：

```vue	
<template>
    <el-config-provider :locale="elementLocales[locale]">
      <router-view></router-view>
    </el-config-provider>
</template>

<script setup>
import { elementLocales } from '@/i18n'
const { locale } = useI18n();
locale.value = localStorage.getItem('locale') || 'zh-cn';
</script>
```

## 全屏处理

### 登出

```js
// 状态管理获取登录信息
import { useUserStore } from '@/stores/user.js'
const userStore = useUserStore()
// 下拉事件处理
const commands = ({
  //个人中心跳转
  toPersonal: () => {
    router.push('/personal')
  },
  // 退出方法
  toLogout: () => {
    userStore.logout();
  }
});
```

1： 找到状态管理定义退出方法  stores/user.js 在actions中增加logout退出方法如下：

```js
import { defineStore } from 'pinia'
import request from '@/request'
import router from '@/router'
import { menuTreeData } from '@/mock/data.js'

//https://blog.csdn.net/weixin_62897746/article/details/129124364
//https://prazdevs.github.io/pinia-plugin-persistedstate/guide/
export const useUserStore = defineStore('user', {
  // 定义状态
  state: () => ({
    routerLoaded:false,
    user: {},
    username: '',
    userId: '',
    token: '',
    age:10,
    male:1,
    role:[],
    permissions:[],
    // 路由菜单，用来接收服务端传递过来的菜单数据
    menuTree:[]
  }),

  // 就是一种计算属性的机制，定义的是函数，使用的是属性就相当于computed
  getters:{

    malestr(state){
      if(state.male==1)return "男"
      if(state.male==0)return "女"
      if(state.male==1)return "保密"
    },

    isLogin(state){
      return state.token ? true : false
    },

    roleName(state){
      return state.roles && state.roles.map(r=>r.name).join(",")
    },

    permissionCode(state){
      return state.permissions &&  state.permissions.map(r=>r.code).join(",")
    }
  },

  // 定义动作
  actions: {
   setToken(newtoken){
      this.token = newtoken
   },

   getToken(){
    return this.token
   },
   
   /* 登出*/
   async logout (){
      // 清除状态信息
      this.token = ''
      this.user = {}
      this.username = ''
      this.userId = ''
      this.role = []
      this.permissions = []
      this.menuTree = []
      // 清除自身的本地存储
      localStorage.removeItem("ksd-kva-language")
      localStorage.removeItem("kva-pinia-userstore")
      localStorage.removeItem("isWhitelist")
      // 然后跳转到登录
      router.push({ name: 'Login', replace: true })
  },
  
     
   async toLogin(loginUser){

      // 查询用户信息，角色，权限，角色对应菜单
      const resp = await request.post("login/toLogin", loginUser,{noToken:true})
      // 这个会回退，回退登录页
      var { user ,token,roles,permissions } = resp.data
      // 登录成功以后获取到菜单信息, 这里要调用一
      this.menuTree = menuTreeData
      // 把数据放入到状态管理中
      this.user = user
      this.userId = user.id
      this.username = user.name
      this.token = token
      this.roles = roles
      this.permissions = permissions
      return Promise.resolve(resp)
    }
  },
  persist: {
    key: 'kva-pinia-userstore',
    storage: localStorage,//sessionStorage
  }
})
```

退出在本地确实没问题，退出以后我们要明白一个逻辑，一个用户既然确定要退出了。那么就token就应该立即失效。不应该还有时效意义。那么怎么办。所以我们必须在服务端定义一个接口把当前的token拉入黑名单中，这才是最保险的做法如下：

1: 定义服务端的路由退出接口方法

```go
package login

import (
	"github.com/gin-gonic/gin"
	"xkginweb/api/v1/login"
)

// 登录路由
type LogoutRouter struct{}

func (router *LogoutRouter) InitLogoutRouter(Router *gin.RouterGroup) {
	logoutApi := login.LogOutApi{}
	// 用组定义--（推荐）
	loginRouter := Router.Group("/login")
	{
		loginRouter.POST("/logout", logoutApi.ToLogout)
	}
}

```

ToLogout方法

```go
package login

import (
	"github.com/gin-gonic/gin"
	"xkginweb/commons/jwtgo"
	"xkginweb/commons/response"
	"xkginweb/model/jwt"
)

// 登录业务
type LogOutApi struct{}

var jwtService = jwtgo.JwtService{}

// 退出接口
func (api *LogOutApi) ToLogout(c *gin.Context) {
	// 获取头部的token信息
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Fail(401, "请求未携带token，无权限访问", c)
		return
	}
	// 退出的token,加入到黑名单中
	err := jwtService.JsonInBlacklist(jwt.JwtBlacklist{Jwt: token})
	// 保存失败会进到到错误
	if err != nil {
		response.Fail(401, "token作废失败", c)
		return
	}
	// 如果保存到黑名单中说明,已经可以告知前端可以进行执行清理动作了
	response.Ok("token作废成功!", c)
}

```



2: 进行注册路由

```go
package initilization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"xkginweb/commons/filter"
	"xkginweb/commons/middle"
	"xkginweb/global"
	"xkginweb/router"
	"xkginweb/router/code"
	"xkginweb/router/login"
)

func InitGinRouter() *gin.Engine {
	// 创建gin服务
	ginServer := gin.Default()
	// 提供服务组
	courseRouter := router.RouterWebGroupApp.Course.CourseRouter
	videoRouter := router.RouterWebGroupApp.Video.VideoRouter

	// 解决接口的跨域问题
	ginServer.Use(filter.Cors())

	loginRouter := login.LoginRouter{}
	logoutRouter := login.LogoutRouter{}
	codeRouter := code.CodeRouter{}
	// 接口隔离，比如登录，健康检查都不需要拦截和做任何的处理
	// 业务模块接口，
	privateGroup := ginServer.Group("/api")
	// 不需要拦截就放注册中间间的前面,需要拦截的就放后面
	loginRouter.InitLoginRouter(privateGroup)
	codeRouter.InitCodeRouter(privateGroup)
	// 只要接口全部使用jwt拦截
	privateGroup.Use(middle.JWTAuth())
	{
		logoutRouter.InitLogoutRouter(privateGroup)
		videoRouter.InitVideoRouter(privateGroup)
		courseRouter.InitCourseRouter(privateGroup)
	}

	fmt.Println("router register success")
	return ginServer
}

func RunServer() {
	// 初始化路由
	Router := InitGinRouter()
	// 为用户头像和文件提供静态地址
	Router.StaticFS("/static", http.Dir("/static"))
	address := fmt.Sprintf(":%d", global.Yaml["server.port"])
	// 启动HTTP服务,courseController
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)

	s2 := s.ListenAndServe().Error()
	fmt.Println("服务启动完毕 ", s2)
}

```



3: 在前端定义退出的方法

```js
import request from '@/request/index.js'

/**
 * 退出登录
 */
export const handleLogout = ()=>{
    request.post("/login/logout")
}
```

4: 执行退出

找到stores/user.js的actions中的logout方法增加服务端的退出请求如下：

```js
import { defineStore } from 'pinia'
import request from '@/request'
import router from '@/router'
import { menuTreeData } from '@/mock/data.js'
import { handleLogout } from '../api/logout.js'

//https://blog.csdn.net/weixin_62897746/article/details/129124364
//https://prazdevs.github.io/pinia-plugin-persistedstate/guide/
export const useUserStore = defineStore('user', {
  // 定义状态
  state: () => ({
    routerLoaded:false,
    user: {},
    username: '',
    userId: '',
    token: '',
    role:[],
    permissions:[],
    // 路由菜单，用来接收服务端传递过来的菜单数据
    menuTree:[]
  }),

  // 就是一种计算属性的机制，定义的是函数，使用的是属性就相当于computed
  getters:{
    isLogin(state){
      return state.token ? true : false
    },

    roleName(state){
      return state.roles && state.roles.map(r=>r.name).join(",")
    },

    permissionCode(state){
      return state.permissions &&  state.permissions.map(r=>r.code).join(",")
    }
  },

  // 定义动作
  actions: {

   /* 设置token */ 
   setToken(newtoken){
      this.token = newtoken
   },

   /* 获取token*/
   getToken(){
    return this.token
   },
   
   /* 登出*/
   async logout (){
      // 执行服务端退出
      await handleLogout()
      // 清除状态信息
      this.token = ''
      this.user = {}
      this.username = ''
      this.userId = ''
      this.role = []
      this.permissions = []
      this.menuTree = []
      // 清除自身的本地存储
      localStorage.removeItem("ksd-kva-language")
      localStorage.removeItem("kva-pinia-userstore")
      localStorage.removeItem("isWhitelist")
      // 然后跳转到登录
      router.push({ name: 'Login', replace: true })
  },

  /* 登录*/
  async toLogin(loginUser){
      // 查询用户信息，角色，权限，角色对应菜单
      const resp = await request.post("login/toLogin", loginUser,{noToken:true})
      // 这个会回退，回退登录页
      var { user ,token,roles,permissions } = resp.data
      // 登录成功以后获取到菜单信息, 这里要调用一
      this.menuTree = menuTreeData
      // 把数据放入到状态管理中
      this.user = user
      this.userId = user.id
      this.username = user.name
      this.token = token
      this.roles = roles
      this.permissions = permissions
      return Promise.resolve(resp)
    }
  },
    
  persist: {
    key: 'kva-pinia-userstore',
    storage: localStorage,//sessionStorage
  }
})
```

然后查看：jwt_blacklists 中是否增加一条token记录。如果增加了说明就已经拉入到了黑名单中。就没有没有问题了。



### 用户和角色和头像展示

```vue
<template>
  <div class="header-cont">
    <div class="left">
      <h1>
        <router-link to="/">{{ t('KvaAdminHome') }}</router-link>
      </h1>
    </div>
    
    <div class="right flex-center">
      <!--全屏处理-->
      <div class="fullbox">
        <el-icon @click="handleFullChange(true)" v-if="!screenfullFlag" color="#fff"><FullScreen /></el-icon>
        <el-icon @click="handleFullChange(false)" v-else color="#fff"><Aim /></el-icon>
      </div>
      <!--国际化-->
      <div class="lang gap">
        <span
          class="item"
          :class="{ active: locale === 'zh-cn' }"
          @click="changeLanguage('zh-cn')"
        >简体中文</span>
        /
        <span
          class="item"
          :class="{ active: locale === 'en' }"
          @click="changeLanguage('en')"
        >EN</span>
      </div>
      <template v-if="isLogin">
        <div class="gap">
          <router-link to="/personal/message">
            <el-badge :is-dot="!!unReadCount">
              <el-icon>
                <message />
              </el-icon>
            </el-badge>
          </router-link>
        </div>
        <el-dropdown trigger="click" @command="handleCommand">
          <div class="flex-center cursor">
            <el-avatar size="small" :src="userStore.user.avatar" />
            <span class="uname"> {{ username }}</span> 
            <el-icon>
              <caret-bottom />
            </el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>当前角色：{{ currentRole.name}}</el-dropdown-item>
              <el-dropdown-item v-for="(item,index) in otherRoles" :key="index">切换角色：{{ item.name }}</el-dropdown-item>
              <el-dropdown-item divided command="toPersonal"><el-icon><User /></el-icon>{{ t('personalCenter') }}</el-dropdown-item>
              <el-dropdown-item divided command="toLogout"><el-icon><Pointer /></el-icon>{{ t('logout') }}</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>
      <template v-else-if="$route.name !== 'Login'">
        <router-link to="/login">{{ t('login') }}</router-link>
      </template>
    </div>
  </div>
</template>
<script setup>
// 状态管理获取登录信息
import KVA from '@/utils/kva.js'
import { useUserStore } from '@/stores/user.js'
const userStore = useUserStore()
// 路由跳转
const router = useRouter();
// 国际化处理
const { locale, t } = useI18n();
// 获取登录的信息
const isLogin = computed(() => userStore.token);
const username = computed(() => userStore.username)
// 消息未读取的数量
const unReadCount = computed(() => 100);
// 全屏处理
import screenfull from 'screenfull'
// 状态管理全屏按钮切换
const screenfullFlag = ref(false)
// 获取第一个以后角色方便进行切换
const currentRole = computed(()=>userStore.roles && userStore.roles.length && userStore.roles[0])
const otherRoles = computed(()=>userStore.roles && userStore.roles.length>1 && userStore.roles.filter((c,index)=>index > 0))

// 全屏事件处理
const handleFullChange = (flag) => {
  screenfull.toggle()
  screenfullFlag.value = flag
}

// 下拉事件处理
const commands = ({
  //个人中心跳转
  toPersonal: () => {
    router.push('/personal')
  },
  // 退出方法
  toLogout: () => {
    KVA.confirm("退出提示","您确定要离开吗?",{icon:"error"}).then(res=>{
      userStore.logout();
    })
  }
});

function handleCommand(command) {
  commands[command] && commands[command]();
}

// 语言切换
function changeLanguage(lang) {
  // 把选择的语言进行切换
  locale.value = lang
  // 切换以后记得把本地缓存进行修改，否则只会生效当前，刷新就还原。
  localStorage.setItem('ksd-kva-language', lang)
}

// 用于读取本地缓存存储的语言是什么？
function initReadLocale(){
  locale.value = localStorage.getItem("ksd-kva-language") || 'zh-cn'
}

onMounted(()=>{
  initReadLocale();
})

</script>
<style lang="scss">
.header-cont {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding:0 20px;
  a {
    color: inherit;
    text-decoration: none;
  }
  h1 {
    margin: 0;
    font-size: 20px;
  }
  .gap {
    margin-right: 20px;
  }
  .right {
    .uname{margin-left: 10px;}
    .fullbox{margin-right: 20px;cursor: pointer;}
    .lang {
      font-size: 14px;
      .item {
        cursor: pointer;
        &.active {
          font-size: 16px;
          font-weight: bold;
        }
      }
    }
  }
  .el-dropdown {
    color: inherit;
  }
}
</style>
```



## 如果登录了，立即跳转到后台首页

如果已经登录，如果我们又去访问登录，其实这属于无用操作。应该要处理掉。如果登录状态又去访问登录页面就应该直接让他跳转首页。

```js
router.beforeEach(async (to) => {
  //开启进度条
  NProgress.start()
  const userStore = useUserStore();

  // 如果当前是登录状态，用户访问又是登录，属于无用操作，应该跳转到首页去
  if(to.path === '/login'){
    if(userStore.isLogin){
        return {name:"Home"}
    }
    return true;
  }

  // 判断是否登录
  if (!userStore.isLogin && to.name !== 'Login') {
    // 这里的query就是为了记录用户最后一次访问的路径，这个路径是通过to的参数获取
    // 后续在登录成功以后，就可以根据这个path的参数，然后调整到你最后一次访问的路径
    return { name: 'Login', query: { "path": to.path } }
  } 

  // 动态加载路由
  addDynamic()

  // 如果刷新出现空白的问题，那么就使用下面这行代码
  if (!to.name && hasRoute(to)) {
    return { ...to };
  }
  // 查询是否注册
  return true
})
```





## 关于访问首页跳转到默认页面

在访问首页，我们不能够直接把新页面写死，一般来有两种做法：

- 1： 在前端进行配置 在src的目录下，新建一个整个系统的配置js文件如下。

  - src/setting.js 内容下：

    ```js
    export default {
      // 配置首页访问的时候，自动跳转到你指定defaultPage
      defaultPage: {name:"DashBoard",replace:true},
      // 指定菜单导航的个数
      menuCount: 10,
    }
    
    ```

  - 然后在需要的地方进行导入使用即可。

    比如：router/index.js 修改如下：

    ```js
    import { createRouter, createWebHashHistory } from 'vue-router'
    import NProgress from 'nprogress'
    // 获取状态管理的token
    import { useUserStore } from '@/stores/user.js'
    // 显示右上角螺旋加载提示
    NProgress.configure({ showSpinner: true })
    import Layout from "@/layout/Index.vue";
    import PageMain from "@/layout/components/PageMain.vue";
    import { menuTreeData } from '@/mock/data.js'
    import settings from '@/settings.js'
    
    
    const router = createRouter({
      history: createWebHashHistory(import.meta.env.BASE_URL),
      routes: [
        {
          path: "/",
          name: "Home",
          component: Layout
        },
        {
          path: '/login',
          name: 'Login',
          meta: { title: "login" },
          component: () => import('@/views/Login.vue')
        }
      ]
    })
    
    const router404 = {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/error/NotFound.vue')
    }
    
    router.beforeEach(async (to) => {
      //开启进度条
      NProgress.start()
      const userStore = useUserStore();
    
      // 如果当前是登录状态，用户访问又是登录，属于无用操作，应该跳转到首页去
      if(to.path === '/login'){
        if(userStore.isLogin){
            return {name:"Home"}
        }
        return true;
      }
    
      // 判断是否登录
      if (!userStore.isLogin && to.name !== 'Login') {
        // 这里的query就是为了记录用户最后一次访问的路径，这个路径是通过to的参数获取
        // 后续在登录成功以后，就可以根据这个path的参数，然后调整到你最后一次访问的路径
        return { name: 'Login', query: { "path": to.path } }
      } 
    
      // 动态加载路由---这里需要耗时---db--ajax-
      await addDynamic()
    
      // 如果刷新出现空白的问题，那么就使用下面这行代码
      if (!to.name && hasRoute(to)) {
        return { ...to };
      }
    
      // 如果访问的是首页，就跳转到/dashboard页面-------------------------------------新增代码
      if(to.path === "/"){
        // 读取默认菜单的默认页面,需要从数据库的菜单表中去读取
        return settings.defaultPage;
      }
    
      // 查询是否注册
      return true
    })
    
    // 动态路由
    function addDynamic(){
      const userStore = useUserStore();
      // 404可以这样处理
      router.addRoute(router404)
      // 必须服务器返回的菜单和views去碰撞形成一个完整的route信息，然后注册到home下
      addDynamicRoutes(menuTreeData)
      // 同时同步到状态管理中
      userStore.menuTree = menuTreeData;
    }
    
    // 这里是获取工程目录下的views下的所以的.vue结尾的SPA页面
    const modules = import.meta.glob('../views/**/*.vue');
    function addDynamicRoutes(menuTreeData,parent){
      // 开始循环遍历菜单信息
      menuTreeData.forEach((item,index) => {
        // 准备路由数据格式
        const route = {
          path: item.path,
          name: item.name,
          // 增加访问路径的元数据信息
          meta: {name: item.name,icon:item.icon},
          children:[]
        }
        // 如果存在parent,就说明有children
        if(parent){
          if(item.parentId!==0){
            // 这里就开始给子菜单匹配views下面的页面spa
            const compParr = item.path.replace("/", "").split("/");
            const l = compParr.length - 1; 
            const compPath = compParr
              .map((v, i) => {
                return i === l ? v.replace(/\w/, (L) => L.toUpperCase()) + ".vue" : v;
              })
              .join("/");
            route.path = compParr[l];
            // 设置动态组件
            route.component = modules[`../views/${compPath}`];
            parent.children.push(route);
          }
        }else{
          // 判断你是否有children
          if (item.children && item.children.length > 0) {
            // 这里的含义是：把匹配到菜单数据第一项作为首页的入口页面
            // /order-----redirect-----/order/list
            route.redirect = item.children[0].path;
            route.component = PageMain;
            // 递归
            addDynamicRoutes(item.children, route)
          }else{
            //route.component = modules[`../views/${item.name}.vue`] 
            route.component = modules[`../views/${item.name.toLowerCase()}/Index.vue`] 
          }
          router.addRoute("Home", route);
        }
      })
    }
    
    // 判断当前路由是否存在动态添加的路由数据中
    function hasRoute(to) {
      const item = router.getRoutes().find((item) => item.path === to.path);
      return !!item;
    }
    
    router.afterEach(() => {
      //完成进度条
      NProgress.done()
    })
    
    export default router
    
    ```

    

2：在数据库进行配置,其实就给菜单增加一个默认菜单。并且是唯一的



## 关于菜单伸缩处理自适应问题

控制右侧菜单的折叠和隐藏，然后小屏幕也可以进行一些适配。主要实现的方式是如下：

1： 监听浏览器屏幕的宽度 ，并且 如果浏览器屏幕的宽度<992px 就会进行折叠

```js
// 获取屏幕宽度
const screenWidth = ref(window.innerWidth)
onMounted(() => {
  // 然后监听浏览器的窗口resize事件，只要浏览器发生大小的变化就会触发。
  window.addEventListener("resize", () => {
    screenWidth.value = window.innerWidth
  })
})

//watch监听屏幕宽度的变化，进行侧边栏的收缩和展开
watch(screenWidth, (newValue, oldValue) => {
    // 如果浏览器的宽度小于992px时候就会菜单的处理折叠状态
    isCollapse.value = newValue < 992
})
```

2： 如果要控制隐藏那么你可以进行隐藏状态的控制

可配置性：src/settings.js

```js
export default {
  // 配置首页访问的时候，自动跳转到你指定defaultPage
  defaultPage: {path:"/dashboard",replace:true},
  // 指定菜单导航的个数
  menuCount: 10,
  // 菜单折叠屏幕宽度的大小
  collapseWidth: 992,
  // 菜单隐藏屏幕宽度的大小
  hiddenWidth: 640,
}

```

然后修改layout/PageSilder.vue如下：

```js
<template>
  <div class="page-sidebar" v-show="isHidden">
    <div class="collape-bar">
      <el-icon class="cursor" @click="isCollapse = !isCollapse">
        <expand v-if="isCollapse" />
        <fold v-else />
      </el-icon>
    </div>
    <el-menu 
      active-text-color="#333" 
      background-color="#ffffff" 
      text-color="#333" 
      router
      :default-active="defaultActive" 
      class="sidemenu" 
      :collapse="isCollapse">
      <template v-for="(item, i) in menuTree" :key="i">
        <template v-if="item.children && item.children.length">
          <el-sub-menu :index="item.path">
            <template #title>
              <el-icon v-if="item.icon">
                <component :is="item.icon"></component>
              </el-icon>
              <span>{{ t(`menu.${item.name}`) }}</span>
            </template>
              <template v-for="(child, ci) in item.children" :key="ci">
                <el-menu-item :index="child.path">
                  <el-icon>
                    <component :is="child.icon"></component>
                  </el-icon>
                  {{ t(`menu.${child.name}`) }}
                </el-menu-item>
            </template>
          </el-sub-menu>
        </template>
        <template v-else>
          <el-menu-item :index="item.path">
            <el-icon v-if="item.icon">
              <component :is="item.icon"></component>
            </el-icon>
            <span>{{ t(`menu.${item.name}`) }}</span>
          </el-menu-item>
        </template>
      </template>
    </el-menu>
  </div>
</template>
  
<script  setup>
import { useUserStore } from '@/stores/user.js'
import settings from '@/settings.js'
// 这个是用来获取当前访问的路由信息,
const route = useRoute();
const { t } = useI18n();
// 默认情况下不折叠
const isCollapse = ref(false)
const isHidden = ref(true)
// 根据当前路由来激活菜单
const defaultActive = computed(()=>(route.path))
// 获取状态管理的菜单信息
const userStore = useUserStore();
// 如何获取菜单数据呢？
const menuTree = computed(()=>userStore.menuTree)

// 获取屏幕宽度
const screenWidth = ref(window.innerWidth)
onMounted(() => {
  // 然后监听浏览器的窗口resize事件，只要浏览器发生大小的变化就会触发。
  window.addEventListener("resize", () => {
    screenWidth.value = window.innerWidth
  })
})

//watch监听屏幕宽度的变化，进行侧边栏的收缩和展开
watch(screenWidth, (newValue, oldValue) => {
    // 如果浏览器的宽度小于992px时候就会菜单的处理折叠状态
    isCollapse.value = newValue < settings.collapseWidth
    isHidden.value = !(newValue < settings.hiddenWidth)
})

</script>
<style lang="scss">
$slider-width: 180px;
.page-sidebar {
  height: calc(100vh - 90px);
  overflow: hidden auto;
  .sidemenu.el-menu,
  .sidemenu .el-sub-menu>.el-menu {
    --el-menu-text-color: #ccc;
    --el-menu-hover-bg-color: #060251;
    --el-menu-border-color: transparent;
    --el-menu-bg-color: #001529;

    .el-menu-item {
      &.is-active {
        background-color: #e6f7ff;
        color: #1890ff
      }
    }
  }

   /* elmenu菜单的折叠效果是通过属性：
      :collapse="isCollapse"  原理就在控制在不停切换elmenu="el-menu--collapse"样式信息
      1：ture 就折叠，就会使用图标宽度+padding作为菜单宽度
      2: false 就不折叠，那么就使用默认宽度：200px

      下面这行css是什么意思：
      如果菜单上存在el-menu--collapse样式就说明是折叠状态，就使用图标宽度+padding作为菜单宽度
      否则：就用我的width:200作为菜单宽度
   */
  .sidemenu.el-menu:not(.el-menu--collapse) {
    width: $slider-width;
  }

  .collape-bar {
    color: #333;
    font-size: 16px;
    line-height: 36px;
    position: fixed;
    z-index: 2;
    width: 100%;
    left:20px;
    bottom: 0;
    .c-icon {
      cursor: pointer;
    }
  }
}
</style>
```



## 控制面板

views/dashboard/Index.vue

```vue
<template>
    <div class="page admin-box" element-loading-text="正在加载中">
        <div class="gva-card-box">
            <div class="gva-card gva-top-card">
                <div class="gva-top-card-left">
                    <div class="gva-top-card-left-title">早安，管理员，请开始一天的工作吧</div>
                    <div class="gva-top-card-left-dot">今日晴，0℃ - 10℃，天气寒冷，注意添加衣物。</div>
                    <div class="gva-top-card-left-rows">
                        <div class="el-row">
                            <div class="el-col el-col-8 el-col-xs-24 el-col-sm-8">
                                <div class="flex-center"><i 
                                        class="el-icon dashboard-icon"><svg 
                                            xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024">
                                            <path fill="currentColor"
                                                d="M384 96a32 32 0 0 1 64 0v786.752a32 32 0 0 1-54.592 22.656L95.936 608a32 32 0 0 1 0-45.312h.128a32 32 0 0 1 45.184 0L384 805.632V96zm192 45.248a32 32 0 0 1 54.592-22.592L928.064 416a32 32 0 0 1 0 45.312h-.128a32 32 0 0 1-45.184 0L640 218.496V928a32 32 0 1 1-64 0V141.248z">
                                            </path>
                                        </svg></i> 今日流量 (1231231) </div>
                            </div>
                            <div class="el-col el-col-8 el-col-xs-24 el-col-sm-8">
                                <div class="flex-center"><i 
                                        class="el-icon dashboard-icon"><svg 
                                            xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024">
                                            <path fill="currentColor"
                                                d="M628.736 528.896A416 416 0 0 1 928 928H96a415.872 415.872 0 0 1 299.264-399.104L512 704l116.736-175.104zM720 304a208 208 0 1 1-416 0 208 208 0 0 1 416 0z">
                                            </path>
                                        </svg></i> 总用户数 (24001) </div>
                            </div>
                            <div class="el-col el-col-8 el-col-xs-24 el-col-sm-8">
                                <div class="flex-center"><i 
                                        class="el-icon dashboard-icon"><svg 
                                            xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024">
                                            <path fill="currentColor"
                                                d="M736 504a56 56 0 1 1 0-112 56 56 0 0 1 0 112zm-224 0a56 56 0 1 1 0-112 56 56 0 0 1 0 112zm-224 0a56 56 0 1 1 0-112 56 56 0 0 1 0 112zM128 128v640h192v160l224-160h352V128H128z">
                                            </path>
                                        </svg></i> 好评率 (99%) </div>
                            </div>
                        </div>
                    </div>
                    <div >
                        <div class="gva-top-card-left-item"> 使用教学： <a 
                                target="view_window" href="https://www.bilibili.com/video/BV1Rg411u7xH/"
                                style="color: rgb(64, 158, 255);">https://www.bilibili.com/video/BV1Rg411u7xH</a></div>
                        <div class="gva-top-card-left-item"> 插件仓库： <a 
                                target="view_window" href="https://plugin.gin-vue-admin.com/#/layout/home"
                                style="color: rgb(64, 158, 255);">https://plugin.gin-vue-admin.com</a></div>
                    </div>
                </div><img src="https://demo.gin-vue-admin.com/assets/dashboard-70e55b71.png"
                    class="gva-top-card-right" alt="">
            </div>
        </div>
        <div class="gva-card-box">
            <div class="el-card is-always-shadow gva-card quick-entrance">
                <div class="el-card__header">
                    <div class="card-header"><span >快捷入口</span></div>
                </div>
                <div class="el-card__body" style="">
                    <div class="el-row" style="margin-left: -10px; margin-right: -10px;">
                        <div class="el-col el-col-4 el-col-xs-8 is-guttered quick-entrance-items"
                            style="padding-right: 10px; padding-left: 10px;">
                            <div class="quick-entrance-item">
                                <div class="quick-entrance-item-icon"
                                    style="background-color: rgba(255, 156, 110, 0.3);"><i 
                                        class="el-icon"><svg xmlns="http://www.w3.org/2000/svg"
                                            viewBox="0 0 1024 1024" style="color: rgb(255, 156, 110);">
                                            <path fill="currentColor"
                                                d="M544 768v128h192a32 32 0 1 1 0 64H288a32 32 0 1 1 0-64h192V768H192A128 128 0 0 1 64 640V256a128 128 0 0 1 128-128h640a128 128 0 0 1 128 128v384a128 128 0 0 1-128 128H544zM192 192a64 64 0 0 0-64 64v384a64 64 0 0 0 64 64h640a64 64 0 0 0 64-64V256a64 64 0 0 0-64-64H192z">
                                            </path>
                                        </svg></i></div>
                                <p >用户管理</p>
                            </div>
                        </div>
                        <div class="el-col el-col-4 el-col-xs-8 is-guttered quick-entrance-items"
                            style="padding-right: 10px; padding-left: 10px;">
                            <div class="quick-entrance-item">
                                <div class="quick-entrance-item-icon"
                                    style="background-color: rgba(105, 192, 255, 0.3);"><i 
                                        class="el-icon"><svg xmlns="http://www.w3.org/2000/svg"
                                            viewBox="0 0 1024 1024" style="color: rgb(105, 192, 255);">
                                            <path fill="currentColor"
                                                d="M600.704 64a32 32 0 0 1 30.464 22.208l35.2 109.376c14.784 7.232 28.928 15.36 42.432 24.512l112.384-24.192a32 32 0 0 1 34.432 15.36L944.32 364.8a32 32 0 0 1-4.032 37.504l-77.12 85.12a357.12 357.12 0 0 1 0 49.024l77.12 85.248a32 32 0 0 1 4.032 37.504l-88.704 153.6a32 32 0 0 1-34.432 15.296L708.8 803.904c-13.44 9.088-27.648 17.28-42.368 24.512l-35.264 109.376A32 32 0 0 1 600.704 960H423.296a32 32 0 0 1-30.464-22.208L357.696 828.48a351.616 351.616 0 0 1-42.56-24.64l-112.32 24.256a32 32 0 0 1-34.432-15.36L79.68 659.2a32 32 0 0 1 4.032-37.504l77.12-85.248a357.12 357.12 0 0 1 0-48.896l-77.12-85.248A32 32 0 0 1 79.68 364.8l88.704-153.6a32 32 0 0 1 34.432-15.296l112.32 24.256c13.568-9.152 27.776-17.408 42.56-24.64l35.2-109.312A32 32 0 0 1 423.232 64H600.64zm-23.424 64H446.72l-36.352 113.088-24.512 11.968a294.113 294.113 0 0 0-34.816 20.096l-22.656 15.36-116.224-25.088-65.28 113.152 79.68 88.192-1.92 27.136a293.12 293.12 0 0 0 0 40.192l1.92 27.136-79.808 88.192 65.344 113.152 116.224-25.024 22.656 15.296a294.113 294.113 0 0 0 34.816 20.096l24.512 11.968L446.72 896h130.688l36.48-113.152 24.448-11.904a288.282 288.282 0 0 0 34.752-20.096l22.592-15.296 116.288 25.024 65.28-113.152-79.744-88.192 1.92-27.136a293.12 293.12 0 0 0 0-40.256l-1.92-27.136 79.808-88.128-65.344-113.152-116.288 24.96-22.592-15.232a287.616 287.616 0 0 0-34.752-20.096l-24.448-11.904L577.344 128zM512 320a192 192 0 1 1 0 384 192 192 0 0 1 0-384zm0 64a128 128 0 1 0 0 256 128 128 0 0 0 0-256z">
                                            </path>
                                        </svg></i></div>
                                <p >角色管理</p>
                            </div>
                        </div>
                        <div class="el-col el-col-4 el-col-xs-8 is-guttered quick-entrance-items"
                            style="padding-right: 10px; padding-left: 10px;">
                            <div class="quick-entrance-item">
                                <div class="quick-entrance-item-icon"
                                    style="background-color: rgba(179, 127, 235, 0.3);"><i 
                                        class="el-icon"><svg xmlns="http://www.w3.org/2000/svg"
                                            viewBox="0 0 1024 1024" style="color: rgb(179, 127, 235);">
                                            <path fill="currentColor"
                                                d="M160 448a32 32 0 0 1-32-32V160.064a32 32 0 0 1 32-32h256a32 32 0 0 1 32 32V416a32 32 0 0 1-32 32H160zm448 0a32 32 0 0 1-32-32V160.064a32 32 0 0 1 32-32h255.936a32 32 0 0 1 32 32V416a32 32 0 0 1-32 32H608zM160 896a32 32 0 0 1-32-32V608a32 32 0 0 1 32-32h256a32 32 0 0 1 32 32v256a32 32 0 0 1-32 32H160zm448 0a32 32 0 0 1-32-32V608a32 32 0 0 1 32-32h255.936a32 32 0 0 1 32 32v256a32 32 0 0 1-32 32H608z">
                                            </path>
                                        </svg></i></div>
                                <p >菜单管理</p>
                            </div>
                        </div>
                        <div class="el-col el-col-4 el-col-xs-8 is-guttered quick-entrance-items"
                            style="padding-right: 10px; padding-left: 10px;">
                            <div class="quick-entrance-item">
                                <div class="quick-entrance-item-icon"
                                    style="background-color: rgba(255, 214, 102, 0.3);"><i 
                                        class="el-icon"><svg xmlns="http://www.w3.org/2000/svg"
                                            viewBox="0 0 1024 1024" style="color: rgb(255, 214, 102);">
                                            <path fill="currentColor"
                                                d="M320 256a64 64 0 0 0-64 64v384a64 64 0 0 0 64 64h384a64 64 0 0 0 64-64V320a64 64 0 0 0-64-64H320zm0-64h384a128 128 0 0 1 128 128v384a128 128 0 0 1-128 128H320a128 128 0 0 1-128-128V320a128 128 0 0 1 128-128z">
                                            </path>
                                            <path fill="currentColor"
                                                d="M512 64a32 32 0 0 1 32 32v128h-64V96a32 32 0 0 1 32-32zm160 0a32 32 0 0 1 32 32v128h-64V96a32 32 0 0 1 32-32zm-320 0a32 32 0 0 1 32 32v128h-64V96a32 32 0 0 1 32-32zm160 896a32 32 0 0 1-32-32V800h64v128a32 32 0 0 1-32 32zm160 0a32 32 0 0 1-32-32V800h64v128a32 32 0 0 1-32 32zm-320 0a32 32 0 0 1-32-32V800h64v128a32 32 0 0 1-32 32zM64 512a32 32 0 0 1 32-32h128v64H96a32 32 0 0 1-32-32zm0-160a32 32 0 0 1 32-32h128v64H96a32 32 0 0 1-32-32zm0 320a32 32 0 0 1 32-32h128v64H96a32 32 0 0 1-32-32zm896-160a32 32 0 0 1-32 32H800v-64h128a32 32 0 0 1 32 32zm0-160a32 32 0 0 1-32 32H800v-64h128a32 32 0 0 1 32 32zm0 320a32 32 0 0 1-32 32H800v-64h128a32 32 0 0 1 32 32z">
                                            </path>
                                        </svg></i></div>
                                <p >代码生成器</p>
                            </div>
                        </div>
                        <div class="el-col el-col-4 el-col-xs-8 is-guttered quick-entrance-items"
                            style="padding-right: 10px; padding-left: 10px;">
                            <div class="quick-entrance-item">
                                <div class="quick-entrance-item-icon"
                                    style="background-color: rgba(255, 133, 192, 0.3);"><i 
                                        class="el-icon"><svg xmlns="http://www.w3.org/2000/svg"
                                            viewBox="0 0 1024 1024" style="color: rgb(255, 133, 192);">
                                            <path fill="currentColor"
                                                d="M805.504 320 640 154.496V320h165.504zM832 384H576V128H192v768h640V384zM160 64h480l256 256v608a32 32 0 0 1-32 32H160a32 32 0 0 1-32-32V96a32 32 0 0 1 32-32zm318.4 582.144 180.992-180.992L704.64 510.4 478.4 736.64 320 578.304l45.248-45.312L478.4 646.144z">
                                            </path>
                                        </svg></i></div>
                                <p >表单生成器</p>
                            </div>
                        </div>
                        <div class="el-col el-col-4 el-col-xs-8 is-guttered quick-entrance-items"
                            style="padding-right: 10px; padding-left: 10px;">
                            <div class="quick-entrance-item">
                                <div class="quick-entrance-item-icon"
                                    style="background-color: rgba(92, 219, 211, 0.3);"><i 
                                        class="el-icon"><svg xmlns="http://www.w3.org/2000/svg"
                                            viewBox="0 0 1024 1024" style="color: rgb(92, 219, 211);">
                                            <path fill="currentColor"
                                                d="M512 512a192 192 0 1 0 0-384 192 192 0 0 0 0 384zm0 64a256 256 0 1 1 0-512 256 256 0 0 1 0 512zm320 320v-96a96 96 0 0 0-96-96H288a96 96 0 0 0-96 96v96a32 32 0 1 1-64 0v-96a160 160 0 0 1 160-160h448a160 160 0 0 1 160 160v96a32 32 0 1 1-64 0z">
                                            </path>
                                        </svg></i></div>
                                <p >关于我们</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="gva-card-box">
            <div class="gva-card">
                <div class="card-header"><span >数据统计</span></div>
                <div class="echart-box">
                    <div class="el-row" style="margin-left: -10px; margin-right: -10px;">
                        <div class="el-col el-col-24 el-col-xs-24 el-col-sm-18 is-guttered"
                            style="padding-right: 10px; padding-left: 10px;">
                            <div class="dashboard-line-box">
                                <div class="dashboard-line-title"> 访问趋势 </div>
                                <div class="dashboard-line" _echarts_instance_="ec_1691233547121"
                                    style="user-select: none; -webkit-tap-highlight-color: rgba(0, 0, 0, 0);">
                                    <div
                                        style="position: relative; width: 1177px; height: 360px; padding: 0px; margin: 0px; border-width: 0px; cursor: default;">
                                        <canvas data-zr-dom-id="zr_0" width="1177" height="360"
                                            style="position: absolute; left: 0px; top: 0px; width: 1177px; height: 360px; user-select: none; -webkit-tap-highlight-color: rgba(0, 0, 0, 0); padding: 0px; margin: 0px; border-width: 0px;"></canvas>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="el-col el-col-24 el-col-xs-24 el-col-sm-6 is-guttered"
                            style="padding-right: 10px; padding-left: 10px;">
                            <div class="commit-table">
                                <div class="commit-table-title"> 更新日志 </div>
                                <div class="log">
                                    <div class="log-item">
                                        <div class="flex-1 flex key-box"><span data-v-144ac47f=""
                                                class="key top">1</span></div>
                                        <div class="flex-5 flex message">feat: update ci add
                                            CGO_ENABLED=0 (#1497)</div>
                                        <div class="flex-3 flex form">2023-08-04</div>
                                    </div>
                                    <div class="log-item">
                                        <div class="flex-1 flex key-box"><span data-v-144ac47f=""
                                                class="key top">2</span></div>
                                        <div class="flex-5 flex message">集成tailwindcss 替换新的登陆页面和版权信息组件
                                            (#1499)

                                            * 集成tailwindcss

                                            * fix: 修改登录，init 页面

                                            * Update package.json

                                            * 升级可升级的第三方库到新版本

                                            * 细节调整

                                            ---------

                                            Co-authored-by: bypanghu &lt;bypanghu@163.com&gt;
                                            Co-authored-by: task &lt;121913992@qq.com&gt;</div>
                                        <div class="flex-3 flex form">2023-08-04</div>
                                    </div>
                                    <div class="log-item">
                                        <div class="flex-1 flex key-box"><span data-v-144ac47f=""
                                                class="key top">3</span></div>
                                        <div class="flex-5 flex message">fix: macos和linux下代码定位命令不对 #1495
                                            (#1496)</div>
                                        <div class="flex-3 flex form">2023-08-02</div>
                                    </div>
                                    <div class="log-item">
                                        <div class="flex-1 flex key-box"><span data-v-144ac47f=""
                                                class="key">4</span></div>
                                        <div class="flex-5 flex message">Rename outer.go to other.go
                                            (#1490)</div>
                                        <div class="flex-3 flex form">2023-08-02</div>
                                    </div>
                                    <div class="log-item">
                                        <div class="flex-1 flex key-box"><span data-v-144ac47f=""
                                                class="key">5</span></div>
                                        <div class="flex-5 flex message">feat: use description to name
                                            apiGroup (#1494)</div>
                                        <div class="flex-3 flex form">2023-08-01</div>
                                    </div>
                                    <div class="log-item">
                                        <div class="flex-1 flex key-box"><span data-v-144ac47f=""
                                                class="key">6</span></div>
                                        <div class="flex-5 flex message">Merge pull request #1483 from
                                            ChengDaqi2023/oscs_fix_cis9at8au51vj78hfkt0

                                            fix(sec): upgrade golang.org/x/image to 0.5.0</div>
                                        <div class="flex-3 flex form">2023-07-21</div>
                                    </div>
                                    <div class="log-item">
                                        <div class="flex-1 flex key-box"><span data-v-144ac47f=""
                                                class="key">7</span></div>
                                        <div class="flex-5 flex message">update golang.org/x/image
                                            v0.0.0-20210220032944-ac19c3e999fb to 0.5.0</div>
                                        <div class="flex-3 flex form">2023-07-20</div>
                                    </div>
                                    <div class="log-item">
                                        <div class="flex-1 flex key-box"><span data-v-144ac47f=""
                                                class="key">8</span></div>
                                        <div class="flex-5 flex message">修复marked问题</div>
                                        <div class="flex-3 flex form">2023-07-19</div>
                                    </div>
                                    <div class="log-item">
                                        <div class="flex-1 flex key-box"><span data-v-144ac47f=""
                                                class="key">9</span></div>
                                        <div class="flex-5 flex message">增加多图片自动生成功能</div>
                                        <div class="flex-3 flex form">2023-07-18</div>
                                    </div>
                                    <div class="log-item">
                                        <div class="flex-1 flex key-box"><span data-v-144ac47f=""
                                                class="key">10</span></div>
                                        <div class="flex-5 flex message">修复element2.3.8以上i18n的问题</div>
                                        <div class="flex-3 flex form">2023-07-17</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { handleLoadMenus } from '@/api/sysmenu.js'


const handleLoadMenusData = async () => {
    const menuDatas = await handleLoadMenus()
    console.log('menuDatas', menuDatas)
}

onMounted(() => {
    handleLoadMenusData()
})
</script> 

<style lang="scss" scoped>

.page {
    background: #f0f2f5;
}

.page .gva-card-box {
    padding: 12px 16px
}

.page .gva-card-box+.gva-card-box {
    padding-top: 0
}

.page .gva-card {
    box-sizing: border-box;
    background-color: #fff;
    border-radius: 2px;
    height: auto;
    padding: 26px 30px;
    overflow: hidden;
    box-shadow: 0 0 7px 1px rgba(0, 0, 0, .03)
}

.page .gva-top-card {
    height: 260px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: #777
}

.page .gva-top-card-left {
    height: 100%;
    display: flex;
    flex-direction: column
}

.page .gva-top-card-left-title {
    font-size: 22px;
    color: #343844
}

.page .gva-top-card-left-dot {
    font-size: 16px;
    color: #6b7687;
    margin-top: 24px
}

.page .gva-top-card-left-rows {
    margin-top: 18px;
    color: #6b7687;
    width: 600px;
    align-items: center
}

.page .gva-top-card-left-item {
    margin-top: 14px
}

.page .gva-top-card-left-item+.gva-top-card-left-item {
    margin-top: 24px
}

.page .gva-top-card-right {
    height: 600px;
    width: 600px;
    margin-top: 28px
}

.page .el-card__header {
    padding: 0;
    border-bottom: none
}

.page .card-header {
    padding-bottom: 20px;
    border-bottom: 1px solid #e8e8e8
}

.page .quick-entrance-title {
    height: 30px;
    font-size: 22px;
    color: #333;
    width: 100%;
    border-bottom: 1px solid #eee
}

.page .quick-entrance-items {
    display: flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    color: #333
}

.page .quick-entrance-items .quick-entrance-item {
    padding: 16px 28px;
    margin-top: -16px;
    margin-bottom: -16px;
    border-radius: 4px;
    transition: all .2s;
    cursor: pointer;
    height: auto;
    text-align: center
}

.page .quick-entrance-items .quick-entrance-item:hover {
    box-shadow: 0 0 7px rgba(217, 217, 217, .55)
}

.page .quick-entrance-items .quick-entrance-item-icon {
    width: 50px;
    height: 50px !important;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto
}

.page .quick-entrance-items .quick-entrance-item-icon i {
    font-size: 24px
}

.page .quick-entrance-items .quick-entrance-item p {
    margin-top: 10px
}

.page .echart-box {
    padding: 14px
}

.dashboard-icon {
    font-size: 20px;
    color: #55a0f8;
    width: 30px;
    height: 30px;
    margin-right: 10px;
    display: flex;
    align-items: center
}

.flex-center {
    display: flex;
    align-items: center
}

@media (max-width: 750px) {
    .gva-card {
        padding: 20px 10px !important
    }

    .gva-card .gva-top-card {
        height: auto
    }

    .gva-card .gva-top-card-left-title {
        font-size: 20px !important
    }

    .gva-card .gva-top-card-left-rows {
        margin-top: 15px;
        align-items: center
    }

    .gva-card .gva-top-card-right {
        display: none
    }

    .gva-card .gva-middle-card-item {
        line-height: 20px
    }

    .gva-card .dashboard-icon {
        font-size: 18px
    }
}

.commit-table {
    background-color: #fff;
    height: 400px
}

.commit-table-title {
    font-weight: 600;
    margin-bottom: 12px
}

.commit-table .log-item {
    display: flex;
    justify-content: space-between;
    margin-top: 14px
}

.commit-table .log-item .key-box {
    justify-content: center
}

.commit-table .log-item .key {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: #F0F2F5;
    text-align: center;
    color: rgba(0, 0, 0, .65)
}

.commit-table .log-item .key.top {
    background: #314659;
    color: #fff
}

.commit-table .log-item .message {
    color: rgba(0, 0, 0, .65)
}

.commit-table .log-item .form {
    color: rgba(0, 0, 0, .65);
    margin-left: 12px
}

.commit-table .log-item .flex {
    line-height: 20px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap
}

.commit-table .log-item .flex-1 {
    flex: 1
}

.commit-table .log-item .flex-2 {
    flex: 2
}

.commit-table .log-item .flex-3 {
    flex: 3
}

.commit-table .log-item .flex-4 {
    flex: 4
}

.commit-table .log-item .flex-5 {
    flex: 5
}</style>
```

关于统计报表的处理和安装

1: 安装echarts

```
pnpm install echarts
```

2: 定义组件

```vue
<template>
    <div id="orderCharts" style="height:400px;width: 100%;">111111</div>
</template>
<script setup>
// 1: 引入echarts
import * as echarts from 'echarts';
// 2: 开始定义统计报表 
const handleLoadCharts = () => {
    // 基于准备好的dom，初始化echarts实例
    var myChart = echarts.init(document.getElementById('orderCharts'));
    // 绘制图表
    myChart.setOption({
        color: ['#80FFA5', '#00DDFF', '#37A2FF', '#FF0087', '#FFBF00'],
        title: {
            text: ''
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985'
            }
            }
        },
        legend: {
            data: ['Line 1', 'Line 2', 'Line 3', 'Line 4', 'Line 5']
        },
        toolbox: {
            feature: {
            saveAsImage: {}
            }
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        },
        xAxis: [
            {
            type: 'category',
            boundaryGap: false,
            data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
            }
        ],
        yAxis: [
            {
            type: 'value'
            }
        ],
        series: [
            {
            name: 'Line 1',
            type: 'line',
            stack: 'Total',
            smooth: true,
            lineStyle: {
                width: 0
            },
            showSymbol: false,
            areaStyle: {
                opacity: 0.8,
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                {
                    offset: 0,
                    color: 'rgb(128, 255, 165)'
                },
                {
                    offset: 1,
                    color: 'rgb(1, 191, 236)'
                }
                ])
            },
            emphasis: {
                focus: 'series'
            },
            data: [140, 232, 101, 264, 90, 340, 250]
            },
            {
            name: 'Line 2',
            type: 'line',
            stack: 'Total',
            smooth: true,
            lineStyle: {
                width: 0
            },
            showSymbol: false,
            areaStyle: {
                opacity: 0.8,
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                {
                    offset: 0,
                    color: 'rgb(0, 221, 255)'
                },
                {
                    offset: 1,
                    color: 'rgb(77, 119, 255)'
                }
                ])
            },
            emphasis: {
                focus: 'series'
            },
            data: [120, 282, 111, 234, 220, 340, 310]
            },
            {
            name: 'Line 3',
            type: 'line',
            stack: 'Total',
            smooth: true,
            lineStyle: {
                width: 0
            },
            showSymbol: false,
            areaStyle: {
                opacity: 0.8,
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                {
                    offset: 0,
                    color: 'rgb(55, 162, 255)'
                },
                {
                    offset: 1,
                    color: 'rgb(116, 21, 219)'
                }
                ])
            },
            emphasis: {
                focus: 'series'
            },
            data: [320, 132, 201, 334, 190, 130, 220]
            },
            {
            name: 'Line 4',
            type: 'line',
            stack: 'Total',
            smooth: true,
            lineStyle: {
                width: 0
            },
            showSymbol: false,
            areaStyle: {
                opacity: 0.8,
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                {
                    offset: 0,
                    color: 'rgb(255, 0, 135)'
                },
                {
                    offset: 1,
                    color: 'rgb(135, 0, 157)'
                }
                ])
            },
            emphasis: {
                focus: 'series'
            },
            data: [220, 402, 231, 134, 190, 230, 120]
            },
            {
            name: 'Line 5',
            type: 'line',
            stack: 'Total',
            smooth: true,
            lineStyle: {
                width: 0
            },
            showSymbol: false,
            label: {
                show: true,
                position: 'top'
            },
            areaStyle: {
                opacity: 0.8,
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                {
                    offset: 0,
                    color: 'rgb(255, 191, 0)'
                },
                {
                    offset: 1,
                    color: 'rgb(224, 62, 76)'
                }
                ])
            },
            emphasis: {
                focus: 'series'
            },
            data: [220, 302, 181, 234, 210, 290, 150]
            }
        ]
        });
}

onMounted(() => {
    handleLoadCharts()
})
</script>
<style lang="scss" scoped></style>
```

3: 使用组件

vue3如下：

```vue
<template>
	 <div class="dashboard-line-box">
       <order-stat-charts></order-stat-charts>
    </div>
</template>
<script setup>
import OrderStatCharts from './compoments/OrderStatCharts.vue';
</script>

```

vue2的如下：

```vue
<template>
	 <div class="dashboard-line-box">
       <order-stat-charts></order-stat-charts>
    </div>
</template>
<script > 
import OrderStatCharts from './compoments/OrderStatCharts.vue';
export default ({
    components:{
        OrderStatCharts
    },
    data(){
        return {
            
        }
    }
})
</script>

```

## 关于菜单表的设计



## 关于后台系统中用户，角色，权限的设计



## 用户授权



## 角色绑定菜单和绑定权限





## 菜单定位导航





## 骨架屏幕加载



## 踢下线



## 把gva项目中的业务搬家自己架构中



## 开始实现课程章节管理



## 订单管理



## 帖子管理



## 文章管理



## 定时器



## 日志

- zap

## elk



## 异步编程



## 自动构建

- model
- web
- router
- page









# Go 整合zap和日志文件分割处理



## 01、下载和安装

```go
# zap核心组件
go get -u go.uber.org/zap
# 日志文件分割,日志文件的保留
go get gopkg.in/natefinch/lumberjack.v2
```

## 02、日志文件对象的初始化

### 1： 定义全局的日志对象

这个日志对象就未来在你代码去使用，但是前提必须要初始化，如何完成初始化呢，在global.go中定义个日志对象即可。然后吧日志进行初始化。找到项目的global.go文件修改如下：

```go
package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"xkginweb/commons/parse"
)

var (
	Log        *zap.Logger // -------------------------- 新增代码
	SugarLog   *zap.SugaredLogger // -------------------------- 新增代码
	Lock       sync.RWMutex
	Yaml       map[string]interface{}
	Config     *parse.Config
	KSD_DB     *gorm.DB
	BlackCache local_cache.Cache
	REDIS      *redis.Client
)

```

### 2：在 [initilization](C:\Users\zxc\go\xkginweb\initilization)*中定义*init_zaplog.go的文件来初始化日志对象信息

```go
package initilization

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
	"xkginweb/global"
)

func InitLogger(mode string) {
	var (
		allCore []zapcore.Core
		core    zapcore.Core
	)
	encoder := getEncoder()
	writeSyncerInfo := getLumberJackWriterInfo()
	writeSyncerError := getLumberJackWriterError()
	// 日志是输出终端
	if mode == "debug" || mode == "info" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	}

	if mode == "error" {
		allCore = append(allCore, zapcore.NewCore(encoder, writeSyncerError, zapcore.ErrorLevel))
	}

	if mode == "info" {
		allCore = append(allCore, zapcore.NewCore(encoder, writeSyncerInfo, zapcore.InfoLevel))
	}

	core = zapcore.NewTee(allCore...)
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
	global.Log = logger
	global.SugarLog = logger.Sugar()
}

func getLumberJackWriterError() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./zap_error.log", // 日志文件位置
		MaxSize:    5,                 // 进行切割之前，日志文件最大值(单位：MB)，默认100MB
		MaxBackups: 5,                 // 保留旧文件的最大个数
		MaxAge:     1,                 // 保留旧文件的最大天数
		Compress:   false,             // 是否压缩/归档旧文件
	}

	// 输入文件和控制台
	//return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	// 只输出文件
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger))
}

func getLumberJackWriterInfo() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./zap_info.log", // 日志文件位置
		MaxSize:    5,                // 进行切割之前，日志文件最大值(单位：MB)，默认100MB
		MaxBackups: 5,                // 保留旧文件的最大个数
		MaxAge:     1,                // 保留旧文件的最大天数
		Compress:   false,            // 是否压缩/归档旧文件
	}

	// 输入文件和控制台
	//return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	// 只输出文件
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger))
}

// json的方式输出
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 空格的方式输出
//func getEncoder() zapcore.Encoder {
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
//	return zapcore.NewConsoleEncoder(encoderConfig)
//}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

```

### 3：在main.go中然后初始化 InitLogger的方法即可。

```go
package main

import (
	"xkginweb/initilization"
)

func main() {
	// 解析配置文件
	initilization.InitViper()
    // 初始化日志 开发的时候建议设置成：debug ，发布的时候建议设置成：info/error
	initilization.InitLogger("error")
	// 初始化中间 redis/mysql/mongodb
	initilization.InitMySQL()
	// 初始化缓存
	initilization.InitRedis()
	// 定时器
	// 并发问题解决方案
	// 异步编程
	// 初始化路由
	initilization.RunServer()
}

```

### 4：使用

#### sugar用法

```go
global.SugarLog.Infow("failed to fetch URL",
  // Structured context as loosely typed key-value pairs.
  "url", url,
  "attempt", 3,
  "backoff", time.Second,
)
global.SugarLog.Infof("Failed to fetch URL: %s", url)
```

#### 非sugar

```go
global.Log.Info("failed to fetch URL",
  // Structured context as strongly typed Field values.
  zap.String("url", url),
  zap.Int("attempt", 3),
  zap.Duration("backoff", time.Second),
)
```



## 日志对象是如何初始化呢？

### 日志格式输出的风格

- 默认情况是用空格来定义的，如下：

  ```go
  	// 日志级别是：debug 或者 info .那么就默认encoder输出格式化改成正常输出
  	if mode == "debug" || mode == "info" {
  		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
  		allCore = append(allCore, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
  	}
  ```

  上面代码告诉我们。在debug和info级别下。使用 `NewConsoleEncoder` 来进行对日志格式进行输出到控制台，而这种日志格式如下：

  ```go
  2023-08-07T20:33:02.366+0800    DEBUG   initilization/init_gorm.go:53   数据库连接成功。开始运行
  ```

- json格式输出

  ```go
  // 如果error错误级别
  if mode == "error" {
      allCore = append(allCore, zapcore.NewCore(encoder, writeSyncerError, zapcore.ErrorLevel))
  }
  
  // 如果是info级别，也写入日志文件
  if mode == "info" {
      allCore = append(allCore, zapcore.NewCore(encoder, writeSyncerInfo, zapcore.InfoLevel))
  }
  ```

  上面的代码告诉，如果在error或者info的日志级别下，会使用  `encoder == getEncoder()` , 代码如下：

  ```go
  // json的方式输出
  func getEncoder() zapcore.Encoder {
  	encoderConfig := zap.NewProductionEncoderConfig()
  	encoderConfig.EncodeTime = customTimeEncoder
  	return zapcore.NewJSONEncoder(encoderConfig)
  }
  ```

  上面代码就告诉，如果error和info的日志级别，采用的是json的方式来格式化你的日志内容。同时把格式化好json日志内容写入日志文件中

  debug(在开发阶段设置，因为你在开发阶段都已经错误或者问题都解决了才上生存。) > info > warn>error >Fatal

- 那么在开发中我们一般使用什么？info / error 

  





## 日志级别

debug(在开发阶段设置，因为你在开发阶段都已经错误或者问题都解决了才上生存。) > info > warn>error >Fatal

## 日志如何写到日志文件

日志写入到文件使用： `gopkg.in/natefinch/lumberjack.v2`

```go
// 错误日志写入到文件为止
func getLumberJackWriterError() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./kva_error.log", // 日志文件位置
		MaxSize:    5,                 // 进行切割之前，日志文件最大值(单位：MB)，默认100MB
		MaxBackups: 5,                 // 保留旧文件的最大个数
		MaxAge:     1,                 // 保留旧文件的最大天数
		Compress:   false,             // 是否压缩/归档旧文件
	}

	// 输入文件和控制台
	//return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	// 只输出文件
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger))
}

// info日志文件的指定
func getLumberJackWriterInfo() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./kva_info.log", // 日志文件位置
		MaxSize:    5,                // 进行切割之前，日志文件最大值(单位：MB)，默认100MB
		MaxBackups: 5,                // 保留旧文件的最大个数
		MaxAge:     1,                // 保留旧文件的最大天数
		Compress:   false,            // 是否压缩/归档旧文件
	}

	// 输入文件和控制台
	//return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	// 只输出文件
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger))
}

```





## 开始思考，如果扩展日志字段呢？

```
global.Log.Error("我是是一个error的日志级别")
global.Log.Fatal("我是fatal日志级别")
```

写入日志文件

```
{"level":"error","ts":"2023-08-07 20:57:50.955","caller":"xkginweb/main.go:19","msg":"我是是一个error的日志级别"}
{"level":"fatal","ts":"2023-08-07 20:57:50.968","caller":"xkginweb/main.go:20","msg":"我是fatal日志级别"}
{"level":"error","ts":"2023-08-07 20:59:15.305","caller":"xkginweb/main.go:16","msg":"我是是一个error的日志级别"}
{"level":"fatal","ts":"2023-08-07 20:59:15.321","caller":"xkginweb/main.go:20","msg":"我是fatal日志级别"}

```

默认情况下：

- level : 日志基本
- ts : 时间
- caller : 文件和行
- msg : 日志信息

如果扩展字段

```go
global.Log.Error("SugarLog 我是一个错误日志", zap.String("ip", "127.0.0.5"))
global.SugarLog.Errorw("SugarLog 我是一个错误日志", "ip", "127.0.0.1", "port", 8888, "url", "http://www.baidu.com")
global.SugarLog.Errorf("SugarLog 你访问的地址是：%s, 端口是：%d", "127.0.0.1", 8088)
```

## 日志又如何写到kafka中



```go
package initilization

import (
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

type LogKafka struct {
	Topic     string
	Producer  sarama.SyncProducer
	Partition int32
}

func (lk *LogKafka) Write(p []byte) (n int, err error) {
	// 构建消息
	msg := &sarama.ProducerMessage{
		Topic:     lk.Topic,
		Value:     sarama.ByteEncoder(p),
		Partition: lk.Partition,
	}
	// 发现消息
	_, _, err = lk.Producer.SendMessage(msg)
	if err != nil {
		return
	}

	return
}

func main() {
	// mode == debug 日志console输出，其他不输出；kafkaSwitch == false 默认输出到文件，kafkaSwitch == true 输出到kafka
	InitLoggerKafka("debug", true)
	// 输出日志
	sugar.Debugf("查询用户信息开始 id:%d", 1)
	sugar.Infof("查询用户信息成功 name:%s age:%d", "zhangSan", 20)
	sugar.Errorf("查询用户信息失败 error:%v", "未该查询到该用户信息")

	time.Sleep(time.Second * 1)
}

func InitLoggerKafka(mode string, kafkaSwitch bool) {
	var (
		err     error
		allCore []zapcore.Core
		core    zapcore.Core
	)

	// 日志是输出终端
	if mode == "debug" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	}

	if kafkaSwitch { // 日志输出kafka
		// kafka配置
		config := sarama.NewConfig()                     // 设置日志输入到Kafka的配置
		config.Producer.RequiredAcks = sarama.WaitForAll // 等待服务器所有副本都保存成功后的响应
		//config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机的分区类型
		config.Producer.Return.Successes = true // 是否等待成功后的响应,只有上面的RequiredAcks设置不是NoReponse这里才有用.
		config.Producer.Return.Errors = true    // 是否等待失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.

		// kafka连接
		var kl LogKafka
		kl.Topic = "LogTopic" // Topic(话题)：Kafka中用于区分不同类别信息的类别名称。由producer指定
		kl.Partition = 1      // Partition(分区)：Topic物理上的分组，一个topic可以分为多个partition，每个partition是一个有序的队列。partition中的每条消息都会被分配一个有序的id（offset）
		kl.Producer, err = sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
		if err != nil {
			panic(fmt.Sprintf("connect kafka failed: %+v\n", err))
		}
		encoder := getEncoderKafka()
		writeSyncer := zapcore.AddSync(&kl)
		allCore = append(allCore, zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel))
	} else { // 日志输出file
		encoder := getEncoder()
		writeSyncer := getLumberJackWriter()
		allCore = append(allCore, zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel))
	}

	core = zapcore.NewTee(allCore...)
	logger = zap.New(core, zap.AddCaller())
	defer logger.Sync()
	sugar = logger.Sugar()
}

func getLumberJackWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log", // 日志文件位置
		MaxSize:    1,            // 进行切割之前，日志文件最大值(单位：MB)，默认100MB
		MaxBackups: 5,            // 保留旧文件的最大个数
		MaxAge:     1,            // 保留旧文件的最大天数
		Compress:   false,        // 是否压缩/归档旧文件
	}

	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}

func getEncoderKafka() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoderKafka
	return zapcore.NewJSONEncoder(encoderConfig)
}

func customTimeEncoderKafka(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

```



## 封装web请求日志



```go

```

## gorm数据库日志的配置

```go
package initilization

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"xkginweb/commons/orm"
	"xkginweb/global"
)

func InitMySQL() {

	// 初始化gorm的日志
	newLogger := logger2.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger2.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger2.Info, // Log level
			IgnoreRecordNotFoundError: false,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,        // Don't include params in the SQL log
			Colorful:                  true,         // Disable color
		},
	)

	m := global.Config.Database.Mysql
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		// GORM 定义了这些日志级别：Silent、Error、Warn、Info
		//Logger: logger.Default.LogMode(logger.Info),
		Logger: newLogger,
	})

	// 如果报错
	if err != nil {
		global.Log.Error("数据连接出错了", zap.String("error", err.Error()))
		panic("数据连接出错了" + err.Error()) // 把程序直接阻断，把数据连接好了在启动
	}

	global.KSD_DB = db
	// 初始化数据库表
	orm.RegisterTable()

	// 日志输出
	global.Log.Debug("数据库连接成功。开始运行", zap.Any("db", db))
}

```

## gin的日志配置

```go
package initilization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"xkginweb/commons/filter"
	"xkginweb/commons/middle"
	"xkginweb/global"
	"xkginweb/router"
	"xkginweb/router/code"
	"xkginweb/router/login"
)

func InitGinRouter() *gin.Engine {
	// 打印gin的时候日志是否用颜色标出
	//gin.ForceConsoleColor()
	//gin.DisableConsoleColor()
	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 创建gin服务
	ginServer := gin.Default()
	// 提供服务组
	courseRouter := router.RouterWebGroupApp.Course.CourseRouter
	videoRouter := router.RouterWebGroupApp.Video.VideoRouter
	menusRouter := router.RouterWebGroupApp.SysMenu.SysMenusRouter

	// 解决接口的跨域问题
	ginServer.Use(filter.Cors())

	loginRouter := login.LoginRouter{}
	logoutRouter := login.LogoutRouter{}
	codeRouter := code.CodeRouter{}
	// 接口隔离，比如登录，健康检查都不需要拦截和做任何的处理
	// 业务模块接口，
	privateGroup := ginServer.Group("/api")
	// 不需要拦截就放注册中间间的前面,需要拦截的就放后面
	loginRouter.InitLoginRouter(privateGroup)
	codeRouter.InitCodeRouter(privateGroup)
	// 只要接口全部使用jwt拦截
	privateGroup.Use(middle.JWTAuth())
	{
		logoutRouter.InitLogoutRouter(privateGroup)
		videoRouter.InitVideoRouter(privateGroup)
		courseRouter.InitCourseRouter(privateGroup)
		menusRouter.InitSysMenusRouter(privateGroup)
	}

	fmt.Println("router register success")
	return ginServer
}

func RunServer() {
	// 初始化路由
	Router := InitGinRouter()
	// 为用户头像和文件提供静态地址
	Router.StaticFS("/static", http.Dir("/static"))
	address := fmt.Sprintf(":%d", global.Yaml["server.port"])
	// 启动HTTP服务,courseController
	s := initServer(address, Router)
	global.Log.Debug("服务启动成功：端口是：", zap.String("port", "8088"))
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)

	s2 := s.ListenAndServe().Error()
	global.Log.Info("服务启动完毕", zap.Any("s2", s2))
}

```

这里就告诉我们一个道理。你可以自己使用zap来完成gin的日志的事情，或者完成gorm日志的事情。



# Zap日志输出kafka、文件、console

## 前言

日志对于项目的重要性不言而喻，之前项目线上的日志都是zap输出到文件，再由filebeat读取输出到kafka，文件服务器保留了大量的日志文件，而且有时filebeat服务重启，可能会导致日志消费重复的问题。所以后面就考虑直接输出到kafka，这样可以减少filebeat的处理过程，且不会出现日志重复消费的问题。

## 一、Kafka服务

部署服务这里采用docker部署，毕竟部署过程不是重点；这里一共部署3个服务，zookeeper、kafka、kafdrop（kafka管理）。

```yaml
version: "3"
services:
  zookeeper:
    image: 'bitnami/zookeeper:3.7'
    ports:
      - '2181:2181'
    environment:
      - ALLOW_ANONYMOUS_LOGIN=yes
  kafka:
    image: 'bitnami/kafka:3.2.0'
    ports:
      - '9092:9092'
    environment:
      - KAFKA_BROKER_ID=1
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://127.0.0.1:9092
      - KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181
      - ALLOW_PLAINTEXT_LISTENER=yes
    depends_on:
      - zookeeper
  kafdrop:
    image: 'obsidiandynamics/kafdrop:3.30.0'
    restart: always
    ports:
      - "9000:9000"
    environment:
      - KAFKA_BROKERCONNECT=kafka:9092
    depends_on:
      - zookeeper
      - kafka
```

将上述代码放入"docker-compose.yml"文件，然后在该文件下执行"docker-compose up -d zookeeper kafka kafdrop"，等待镜像下载和服务器，这个该过程可能需要点时间。若安装成功，访问"http://localhost:59000/"，就可以进入kafka管理界面。

## 二Zap日志输出

### 1.输出console

```go
package main
 
import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)
 
var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)
 
func main() {
	// mode == debug 日志console数据，其他不输出
	InitLogger("debug")
 
	// 输出日志
	sugar.Debugf("查询用户信息开始 id:%d", 1)
	sugar.Infof("查询用户信息成功 name:%s age:%d", "zhangSan", 20)
	sugar.Errorf("查询用户信息失败 error:%v", "未该查询到该用户信息")
}
 
func InitLogger(mode string) {
	var (
		allCore []zapcore.Core
		core    zapcore.Core
	)
 
	// 进入debug模式，日志输出到终端
	if mode == "debug" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	}
 
	core = zapcore.NewTee(allCore...)
	logger = zap.New(core, zap.AddCaller())
	defer logger.Sync()
	sugar = logger.Sugar()
}
 
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
```

运行上述代码输出结果如下：
![img.png](images/img.png)

### 2.输出file

```go
func InitLogger(mode string) {
	var (
		allCore []zapcore.Core
		core    zapcore.Core
	)
 
	// 日志是输出终端
	if mode == "debug" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	}
 
	// 日志输出文件
	encoder := getEncoder()
	writeSyncer := getLumberJackWriter()
	allCore = append(allCore, zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel))
 
	core = zapcore.NewTee(allCore...)
	logger = zap.New(core, zap.AddCaller())
	defer logger.Sync()
	sugar = logger.Sugar()
}
 
func getLumberJackWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log", // 日志文件位置
		MaxSize:    1,            // 进行切割之前，日志文件最大值(单位：MB)，默认100MB
		MaxBackups: 5,            // 保留旧文件的最大个数
		MaxAge:     1,            // 保留旧文件的最大天数
		Compress:   false,        // 是否压缩/归档旧文件
	}
 
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}
 
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
```

运行上述代码，会生成test.log， 如果日志大于1MB，会进行自动分割，大家可以自己尝试；test.log内容如下：
![img_1.png](images/img_1.png)

###  3.输出kafka

```go
package main
 
import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)
 
var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)
 
type LogKafka struct {
	Topic     string
	Producer  sarama.SyncProducer
	Partition int32
}
 
func (lk *LogKafka) Write(p []byte) (n int, err error) {
	// 构建消息
	msg := &sarama.ProducerMessage{
		Topic:     lk.Topic,
		Value:     sarama.ByteEncoder(p),
		Partition: lk.Partition,
	}
	// 发现消息
	_, _, err = lk.Producer.SendMessage(msg)
	if err != nil {
		return
	}
 
	return
}
 
func main() {
	// mode == debug 日志console输出，其他不输出；kafkaSwitch == false 默认输出到文件，kafkaSwitch == true 输出到kafka
	InitLogger("debug", true)
 
	// 输出日志
	sugar.Debugf("查询用户信息开始 id:%d", 1)
	sugar.Infof("查询用户信息成功 name:%s age:%d", "zhangSan", 20)
	sugar.Errorf("查询用户信息失败 error:%v", "未该查询到该用户信息")
 
	time.Sleep(time.Second * 1)
}
 
func InitLogger(mode string, kafkaSwitch bool) {
	var (
		err     error
		allCore []zapcore.Core
		core    zapcore.Core
	)
 
	// 日志是输出终端
	if mode == "debug" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	}
 
	if kafkaSwitch { // 日志输出kafka
		// kafka配置
		config := sarama.NewConfig()                     // 设置日志输入到Kafka的配置
		config.Producer.RequiredAcks = sarama.WaitForAll // 等待服务器所有副本都保存成功后的响应
		//config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机的分区类型
		config.Producer.Return.Successes = true // 是否等待成功后的响应,只有上面的RequiredAcks设置不是NoReponse这里才有用.
		config.Producer.Return.Errors = true    // 是否等待失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
 
		// kafka连接
		var kl LogKafka
		kl.Topic = "LogTopic" // Topic(话题)：Kafka中用于区分不同类别信息的类别名称。由producer指定
		kl.Partition = 1      // Partition(分区)：Topic物理上的分组，一个topic可以分为多个partition，每个partition是一个有序的队列。partition中的每条消息都会被分配一个有序的id（offset）
		kl.Producer, err = sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
		if err != nil {
			panic(fmt.Sprintf("connect kafka failed: %+v\n", err))
		}
		encoder := getEncoder()
		writeSyncer := zapcore.AddSync(&kl)
		allCore = append(allCore, zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel))
	} else { // 日志输出file
		encoder := getEncoder()
		writeSyncer := getLumberJackWriter()
		allCore = append(allCore, zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel))
	}
 
	core = zapcore.NewTee(allCore...)
	logger = zap.New(core, zap.AddCaller())
	defer logger.Sync()
	sugar = logger.Sugar()
}
 
func getLumberJackWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log", // 日志文件位置
		MaxSize:    1,            // 进行切割之前，日志文件最大值(单位：MB)，默认100MB
		MaxBackups: 5,            // 保留旧文件的最大个数
		MaxAge:     1,            // 保留旧文件的最大天数
		Compress:   false,        // 是否压缩/归档旧文件
	}
 
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
}
 
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
 
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
```

上述代码通过"kafkaSwitch"控制日志输出到file还是kafka，这里输出到kafka，结果如下：
![img_2.png](images/img_2.png)
![img_3.png](images/img_3.png)

### 总结

通过上述列子，我们可以轻松在项目中实现日志输出，自由选择cosole、file、kafka，方便项目开发和问题排查。



# 实现系统用户表的管理



## 01、系统后台相关表的业务对照

- sys_user —后台用户表

  ![image-20230807214954498](images/image-20230807214954498.png)

- sys_apis — 权限表 

  ![image-20230807214945444](images/image-20230807214945444.png)

- sys_roles—角色表

  ![image-20230807214936137](images/image-20230807214936137.png)

- sys_menus – 菜单表

  ![image-20230807215002879](../../../../../../L_Learning/%25E6%25B5%258B%25E5%25BC%2580%25E8%25AF%25BE%25E7%25A8%258B/%25E7%258B%2582%25E7%25A5%259E/3-%25E9%25A1%25B9%25E7%259B%25AE%25E5%25AE%259E%25E6%2588%2598%2520-%2520GVA%25E5%2590%258E%25E5%258F%25B0%25E9%25A1%25B9%25E7%259B%25AE%25E7%25AE%25A1%25E7%2590%2586%25E5%25BC%2580%25E5%258F%2591/20230817%25EF%25BC%259A%25E7%25AC%25AC%25E4%25B8%2589%25E5%258D%2581%25E4%25B9%259D%25E8%25AF%25BE%25EF%25BC%259A%25E8%2587%25AA%25E5%25BB%25BA%25E9%25A1%25B9%25E7%259B%25AE-%25E8%25A7%2592%25E8%2589%25B2%25E3%2580%2581%25E8%258F%259C%25E5%258D%2595%25E3%2580%2581%25E6%259D%2583%25E9%2599%2590%25E6%25B7%25BB%25E5%258A%25A0%25EF%25BC%258C%25E8%258A%2582%25E6%25B5%2581%25E5%2592%258C%25E9%2598%25B2%25E6%258A%2596%25E7%259A%2584%25E5%25BA%2594%25E7%2594%25A8%25E5%2592%258C%25E5%25A4%2584%25E7%2590%2586(1)/%25E9%25A1%25B9%25E7%259B%25AE%25E7%25AC%2594%25E8%25AE%25B0/assets/image-20230807215002879.png)

- sys_role_apis—角色权限 

  ![image-20230807215102175](images/image-20230807215102175.png)

- sys_role_menus — 角色菜单

  ![image-20230807215114855](images/image-20230807215114855.png)

- sys_user_roles —- 用户角色

  ![image-20230807215039489](images/image-20230807215039489.png)



# 后台用户管理

## 01、 概述

系统用户，是不注册，一般都是 由超级管理员在后台添加和分配，然后给账号和初始密码给予小伙伴，然后在登录。

## 02、对应的表

```sql
CREATE TABLE `sys_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `uuid` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户UUID',
  `account` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户登录名',
  `password` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户登录密码',
  `username` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT '系统用户' COMMENT '用户昵称',
  `avatar` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'https://qmplusimg.henrongyi.top/gva_header.jpg' COMMENT '用户头像',
  `phone` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户手机号',
  `email` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户邮箱',
  `enable` bigint(20) DEFAULT '1' COMMENT '用户是否被冻结 1正常 2冻结',
  `is_deleted` bigint(20) unsigned DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_sys_users_uuid` (`uuid`) USING BTREE,
  KEY `idx_sys_users_username` (`account`) USING BTREE,
  KEY `idx_sys_users_account` (`account`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;
```

## 03、实现步骤

- model

  ```go
  package sys
  
  import (
  	uuid "github.com/satori/go.uuid"
  	"xkginweb/global"
  )
  
  type SysUser struct {
  	global.GVA_MODEL
  	UUID     uuid.UUID `json:"uuid" gorm:"index;comment:用户UUID"`                                                  // 用户UUID
  	Account  string    `json:"account" gorm:"index;comment:用户登录名"`                                                // 用户登录名
  	Password string    `json:"-"  gorm:"comment:用户登录密码"`                                                          // 用户登录密码
  	Username string    `json:"username" gorm:"default:系统用户;comment:用户昵称"`                                         // 用户昵称
  	Avatar   string    `json:"avatar" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
  	Phone    string    `json:"phone"  gorm:"comment:用户手机号"`                                                       // 用户手机号
  	Email    string    `json:"email"  gorm:"comment:用户邮箱"`                                                        // 用户邮箱
  	Enable   int       `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`                                   //用户是否被冻结 1正常 2冻结
  }
  
  func (s *SysUser) TableName() string {
  	return "sys_users"
  }
  
  ```

- model init

  ```
  package orm
  
  import (
  	"xkginweb/global"
  	"xkginweb/model/bbs"
  	"xkginweb/model/jwt"
  	"xkginweb/model/sys"
  	"xkginweb/model/user"
  	"xkginweb/model/video"
  )
  
  func RegisterTable() {
  	db := global.KSD_DB
  	// 注册和声明model
  	db.AutoMigrate(user.XkUser{})
  	db.AutoMigrate(user.XkUserAuthor{})
  	// 系统用户，角色，权限表
  	db.AutoMigrate(sys.SysApi{})
  	db.AutoMigrate(sys.SysMenus{})
  	db.AutoMigrate(sys.SysRoleApis{})
  	db.AutoMigrate(sys.SysRoleMenus{})
  	db.AutoMigrate(sys.SysRoles{})
  	db.AutoMigrate(sys.SysUserRoles{})
  	db.AutoMigrate(sys.SysUser{}) //-----------------新增代码
  	// 视频表
  	db.AutoMigrate(video.XkVideo{})
  	db.AutoMigrate(video.XkVideoCategory{})
  	db.AutoMigrate(video.XkVideoChapterLesson{})
  	// 社区
  	db.AutoMigrate(bbs.XkBbs{})
  	db.AutoMigrate(bbs.BbsCategory{})
  
  	// 声明一下jwt模型
  	db.AutoMigrate(jwt.JwtBlacklist{})
  }
  
  ```

- service

  ```go
  package sys
  
  import (
  	"xkginweb/global"
  	"xkginweb/model/comms/request"
  	"xkginweb/model/sys"
  )
  
  // 对用户表的数据层处理
  type SysUserService struct{}
  
  // 用于登录
  func (service *SysUserService) GetUserByAccount(account string) (sysUser *sys.SysUser, err error) {
  	// 根据account进行查询
  	err = global.KSD_DB.Where("account = ?", account).First(&sysUser).Error
  	if err != nil {
  		return nil, err
  	}
  	return sysUser, nil
  }
  
  // 添加
  func (service *SysUserService) SaveSysUser(sysUser *sys.SysUser) (err error) {
  	err = global.KSD_DB.Create(sysUser).Error
  	return err
  }
  
  // 修改
  func (menu *SysMenusService) UpdateSysUser(sysUser *sys.SysUser) (err error) {
  	err = global.KSD_DB.Model(sysUser).Updates(sysUser).Error
  	return err
  }
  
  // 删除
  func (menu *SysMenusService) DelSysUserById(id uint) (err error) {
  	var sysUser sys.SysUser
  	err = global.KSD_DB.Where("id = ?", id).Delete(&sysUser).Error
  	return err
  }
  
  // 批量删除
  func (menu *SysMenusService) DeleteSysUsersByIds(sysUsers []sys.SysUser) (err error) {
  	err = global.KSD_DB.Delete(&sysUsers).Error
  	return err
  }
  
  // 根据id查询信息
  func (menu *SysMenusService) GetSysUserByID(id uint) (sysUsers *sys.SysUser, err error) {
  	err = global.KSD_DB.Where("id = ?", id).First(&sysUsers).Error
  	return
  }
  
  // 查询分页
  func (menu *SysMenusService) LoadSysUserPage(info request.PageInfo) (list interface{}, total int64, err error) {
  	// 获取分页的参数信息
  	limit := info.PageSize
  	offset := info.PageSize * (info.Page - 1)
  
  	// 准备查询那个数据库表
  	db := global.KSD_DB.Model(&sys.SysUser{})
  
  	// 准备切片帖子数组
  	var sysUserList []sys.SysUser
  
  	// 加条件
  	if info.Keyword != "" {
  		db = db.Where("(username like ? or account like ? )", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
  	}
  
  	// 排序默时间降序降序
  	db = db.Order("create_at desc")
  
  	// 查询中枢
  	err = db.Count(&total).Error
  	if err != nil {
  		return sysUserList, total, err
  	} else {
  		// 执行查询
  		err = db.Limit(limit).Offset(offset).Find(&sysUserList).Error
  	}
  
  	// 结果返回
  	return sysUserList, total, err
  }
  
  ```

  

- api

  ```go
  package sys
  
  import (
  	"github.com/gin-gonic/gin"
  	"strconv"
  	"xkginweb/commons/response"
  	"xkginweb/global"
  	req "xkginweb/model/comms/request"
  	resp "xkginweb/model/comms/response"
  )
  
  type SysUsersApi struct{}
  
  /* 根据id查询用户信息 */
  func (api *SysUsersApi) GetById(c *gin.Context) {
  	// 根据id查询方法
  	id := c.Param("id")
  	// 根据id查询方法
  	parseUint, _ := strconv.ParseUint(id, 10, 64)
  	sysUser, err := sysUserService.GetSysUserByID(uint(parseUint))
  	if err != nil {
  		global.SugarLog.Errorf("查询用户: %s 失败", id)
  		response.FailWithMessage("查询用户失败", c)
  		return
  	}
  
  	response.Ok(sysUser, c)
  }
  
  /* 分页查询用户信息*/
  func (api *SysUsersApi) LoadSysUserPage(c *gin.Context) {
  	// 创建一个分页对象
  	var pageInfo req.PageInfo
  	// 把前端json的参数传入给PageInfo
  	err := c.ShouldBindJSON(&pageInfo)
  	if err != nil {
  		response.FailWithMessage(err.Error(), c)
  		return
  	}
  
  	xkBbsPage, total, err := sysUserService.LoadSysUserPage(pageInfo)
  	if err != nil {
  		response.FailWithMessage("获取失败"+err.Error(), c)
  		return
  	}
  	response.Ok(resp.PageResult{
  		List:     xkBbsPage,
  		Total:    total,
  		Page:     pageInfo.Page,
  		PageSize: pageInfo.PageSize,
  	}, c)
  }
  
  ```

  

- router

  ```go
  package sys
  
  import (
  	"github.com/gin-gonic/gin"
  	"xkginweb/api/v1/sys"
  )
  
  // 登录路由
  type SysUsersRouter struct{}
  
  func (router *SysUsersRouter) InitSysUsersRouter(Router *gin.RouterGroup) {
  	sysUsersApi := sys.SysUsersApi{}
  	// 用组定义--（推荐）
  	sysMenusRouter := Router.Group("/sys")
  	{
  		sysMenusRouter.GET("/user/get/:id", sysUsersApi.GetById)
  		sysMenusRouter.GET("/user/load", sysUsersApi.LoadSysUserPage)
  	}
  }
  
  ```

  

- router init

  ```go
  package initilization
  
  import (
  	"fmt"
  	"github.com/gin-gonic/gin"
  	"go.uber.org/zap"
  	"net/http"
  	"time"
  	"xkginweb/commons/filter"
  	"xkginweb/commons/middle"
  	"xkginweb/global"
  	"xkginweb/router"
  	"xkginweb/router/code"
  	"xkginweb/router/login"
  	"xkginweb/router/sys"
  )
  
  func InitGinRouter() *gin.Engine {
  	// 打印gin的时候日志是否用颜色标出
  	//gin.ForceConsoleColor()
  	//gin.DisableConsoleColor()
  	//f, _ := os.Create("gin.log")
  	//gin.DefaultWriter = io.MultiWriter(f)
  	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
  	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
  
  	// 创建gin服务
  	ginServer := gin.Default()
  	// 提供服务组
  	courseRouter := router.RouterWebGroupApp.Course.CourseRouter
  	videoRouter := router.RouterWebGroupApp.Video.VideoRouter
  	menusRouter := router.RouterWebGroupApp.SysMenu.SysMenusRouter
  
  	// 解决接口的跨域问题
  	ginServer.Use(filter.Cors())
  
  	loginRouter := login.LoginRouter{}
  	logoutRouter := login.LogoutRouter{}
  	codeRouter := code.CodeRouter{}
  	sysUserRouter := sys.SysUsersRouter{}
  	// 接口隔离，比如登录，健康检查都不需要拦截和做任何的处理
  	// 业务模块接口，
  	privateGroup := ginServer.Group("/api")
  	// 不需要拦截就放注册中间间的前面,需要拦截的就放后面
  	loginRouter.InitLoginRouter(privateGroup)
  	codeRouter.InitCodeRouter(privateGroup)
  	// 只要接口全部使用jwt拦截
  	privateGroup.Use(middle.JWTAuth())
  	{
  		logoutRouter.InitLogoutRouter(privateGroup)
  		videoRouter.InitVideoRouter(privateGroup)
  		courseRouter.InitCourseRouter(privateGroup)
  		menusRouter.InitSysMenusRouter(privateGroup)
  		sysUserRouter.InitSysUsersRouter(privateGroup)
  	}
  
  	fmt.Println("router register success")
  	return ginServer
  }
  
  func RunServer() {
  
  	// 初始化路由
  	Router := InitGinRouter()
  	// 为用户头像和文件提供静态地址
  	Router.StaticFS("/static", http.Dir("/static"))
  	address := fmt.Sprintf(":%d", global.Yaml["server.port"])
  	// 启动HTTP服务,courseController
  	s := initServer(address, Router)
  	global.Log.Debug("服务启动成功：端口是：", zap.String("port", "8088"))
  	// 保证文本顺序输出
  	// In order to ensure that the text order output can be deleted
  	time.Sleep(10 * time.Microsecond)
  
  	s2 := s.ListenAndServe().Error()
  	global.Log.Info("服务启动完毕", zap.Any("s2", s2))
  }
  
  ```

  

- 定义路由接口 

  ```js
  import request from '@/request/index.js'
  
  /**
   * 查询系统用户信息
   */
  export const LoadSysUser = ()=>{
     return request.get(`/sys/user/load`)
  }
  
  /**
   * 根据id查询系统用户信息
   */
  export const GetSysUserById = ( id )=>{
     return request.get(`/sys/user/${id}`)
  }
  
  
  ```





# 03、登录页面调整

你只需要把页面直接覆盖即可。包括背景图放置在src/assets/imgs目录下，直接在这里去找即可。

```vue
<template>
    <div class="login-box">
        <div class="imgbox">
            <div class="bgblue"></div> 
                <ul class="circles">
                    <li></li>
                    <li></li>
                    <li></li>
                    <li></li>
                    <li></li>
                    <li></li>
                    <li></li>
                    <li></li>
                    <li></li>
                    <li></li>
                </ul>
            <div class="w-full max-w-md">
                <div class="fz48 cof fw">欢迎光临</div>
                <div class="fz14 cof" style="margin-top:10px;line-height:24px;">欢迎来到好玩俱乐部，在这里和志同道合的朋友一起分享有趣的故事，一起组织有趣的活动...</div>
            </div>
        </div>
        <div class="loginbox">
            <div class="login-wrap">
                <h1 class="header fz32">{{ ctitle }}</h1>
                <form action="#">
                    <div class="ksd-el-items"><input type="text" v-model="loginUser.account" class="ksd-login-input"  placeholder="请输入账号"></div>
                    <div class="ksd-el-items"><input type="password" v-model="loginUser.password" class="ksd-login-input" placeholder="请输入密码" @keydown.enter="handleSubmit"></div>
                    <div class="ksd-el-items pr">
                        <input type="text" class="ksd-login-input" maxlength="6" v-model="loginUser.code" placeholder="请输入验证码"  @keydown.enter="handleSubmit">
                        <img v-if="codeURL" class="codeurl"  :src="codeURL" @click="handleGetCapatcha">
                    </div>
                    <div class="ksd-el-items"><input type="button" @click.prevent="handleSubmit" class="ksd-login-btn" value="登录"></div>            
                </form>
            </div>
        </div>
    </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';
import {useRouter,useRoute} from 'vue-router'
import {useUserStore} from '@/stores/user.js'
import {getCapatcha} from '@/api/code.js'
import KVA from "@/utils/kva.js";
const userStore = useUserStore();
const ctitle = ref(import.meta.env.VITE_APP_TITLE)

// 定义一个路由对象
const router = useRouter()
// 获取当前路由信息，用于获取当前路径参数，路径，HASH等
const route = useRoute();
// 准备接受图像验证码
const codeURL = ref("");
// 获取用户输入账号和验证码信息
const loginUser = reactive({
    code:"",
    account:"admin",
    password:"123456",
    codeId:""
})


// 根据axios官方文档开始调用生成验证码的接口
const handleGetCapatcha = async () => {
    const resp = await getCapatcha()
    const {baseURL,id} = resp.data
    codeURL.value = baseURL
    loginUser.codeId = id
}


// 提交表单
const  handleSubmit = async () => {
    // axios.post ---application/json---gin-request.body
    if(!loginUser.code){
        KVA.notifyError("请输入验证码")
        return;
    }
    if(!loginUser.account){
        KVA.notifyError("请输入账号")
        return;
    }
    if(!loginUser.password){
        KVA.notifyError("请输入密码")
        return;
    }

    // 把数据放入到状态管理中
    try{
        await userStore.toLogin(loginUser)
        var path = route.query.path || "/"
        router.push(path)
    }catch(e){
        if(e.code === 60002){
            loginUser.code = ""
            handleGetCapatcha()
        }
    }
}
// 用生命周期去加载生成验证码
onMounted(() => {
    handleGetCapatcha()
})

</script>

<style scoped lang="scss">
    .pr{position: relative;}
    .codeurl{position: absolute;top:5px;right:5px;width: 140px;}
    .ksd-el-items{margin: 15px 0;}
    .ksd-login-input{border:1px solid #eee;padding:16px 8px;width: 100%;box-sizing: border-box;outline: none;border-radius: 4px;}
    .ksd-login-btn{border:1px solid #eee;padding:16px 8px;width: 100%;box-sizing: border-box;
        background:#2196F3;color:#fff;border-radius:6px;cursor: pointer;}
    .ksd-login-btn:hover{background:#1789e7;}
    .login-box{
        display: flex;
        flex-wrap: wrap;
        background:#fff;
        .loginbox{
            width: 35%;height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            .header{margin-bottom: 30px;}
            .login-wrap{
                width: 560px;
                height: 444px;
                padding:20px 100px;
                box-sizing: border-box;
                border-radius: 8px;
                box-shadow: 0 0 10px #fafafa;
                background: rgba(255,255,255,0.6);
                text-align: center;
                display: flex;
                flex-direction: column;
                justify-content: center;
            }
        }
        .imgbox{
            width: 65%;
            height: 100vh;
            display: flex;
            align-items: center;
            justify-content:center;
            position:relative;
            background: url("../assets/imgs/login_background.jpg");
            background-size:cover;
            background-repeat:no-repeat;

            .bgblue{
                background-image:linear-gradient(to bottom,#4f46e5,#3b82f6);
                position:absolute;
                top:0;
                left:0;
                bottom:0;
                right:0;
                opacity:0.75;
            }
        }
    }

    .circles {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    overflow: hidden;
}
.circles li {
    position: absolute;
    display: block;
    list-style: none;
    width: 20px;
    height: 20px;
    background: rgba(255, 255, 255, 0.2);
    animation: animate 25s linear infinite;
    bottom: -150px;
}
.circles li:nth-child(1) {
    left: 25%;
    width: 80px;
    height: 80px;
    animation-delay: 0s;
}
.circles li:nth-child(2) {
    left: 10%;
    width: 20px;
    height: 20px;
    animation-delay: 2s;
    animation-duration: 12s;
}
.circles li:nth-child(3) {
    left: 70%;
    width: 20px;
    height: 20px;
    animation-delay: 4s;
}
.circles li:nth-child(4) {
    left: 40%;
    width: 60px;
    height: 60px;
    animation-delay: 0s;
    animation-duration: 18s;
}
.circles li:nth-child(5) {
    left: 65%;
    width: 20px;
    height: 20px;
    animation-delay: 0s;
}
.circles li:nth-child(6) {
    left: 75%;
    width: 110px;
    height: 110px;
    animation-delay: 3s;
}
.circles li:nth-child(7) {
    left: 35%;
    width: 150px;
    height: 150px;
    animation-delay: 7s;
}
.circles li:nth-child(8) {
    left: 50%;
    width: 25px;
    height: 25px;
    animation-delay: 15s;
    animation-duration: 45s;
}
.circles li:nth-child(9) {
    left: 20%;
    width: 15px;
    height: 15px;
    animation-delay: 2s;
    animation-duration: 35s;
}
.circles li:nth-child(10) {
    left: 85%;
    width: 150px;
    height: 150px;
    animation-delay: 0s;
    animation-duration: 11s;
}
@keyframes animate {
    0% {
        transform: translateY(0) rotate(0deg);
        opacity: 1;
        border-radius: 0;
    }

    100% {
        transform: translateY(-1000px) rotate(720deg);
        opacity: 0;
        border-radius: 50%;
    }
}
.max-w-md {
    max-width: 28rem;
    position:relative;
    z-index:10;
}
</style>
```

# 04、关于message框修改notify框

封装了notify关于error和success封装如下：

```js
const KVA = {
    alert(title,content,options){
        // 默认值
        var defaultOptions = {icon:"warning",confirmButtonText:"确定",cancelButtonText:"取消"}
        // 用户传递和默认值就行覆盖处理
        var opts = {...defaultOptions,...options}
        return ElMessageBox.alert(content, title,{
            //确定按钮文本
            confirmButtonText: opts.confirmButtonText,
            // 内容支持html
            dangerouslyUseHTMLString: true,
            // 是否支持拖拽
            draggable: true,
            // 修改图标
            type: opts.icon
        })
    },
    confirm(title,content,options){
        // 默认值
        var defaultOptions = {icon:"warning",confirmButtonText:"确定",cancelButtonText:"取消"}
        // 用户传递和默认值就行覆盖处理
        var opts = {...defaultOptions,...options}
        // 然后提示
        return ElMessageBox.confirm(content, title, {
           //确定按钮文本
           confirmButtonText: opts.confirmButtonText,
           //取消按钮文本
           cancelButtonText: opts.cancelButtonText,
           // 内容支持html
           dangerouslyUseHTMLString: true,
           // 是否支持拖拽
           draggable: true,
           // 修改图标
           type: opts.icon,
        })
    },
    prompt(title,content,options){
        // 默认值
        var defaultOptions = {confirmButtonText:"确定",cancelButtonText:"取消"}
        // 用户传递和默认值就行覆盖处理
        var opts = {...defaultOptions,...options}
        return ElMessageBox.prompt(content, title, {
            //确定按钮文本
            confirmButtonText: opts.confirmButtonText,
            //取消按钮文本
            cancelButtonText: opts.cancelButtonText,
            // 内容支持html
            dangerouslyUseHTMLString: true,
            // 是否支持拖拽
            draggable: true,
            // 输入框的正则验证
            inputPattern: opts.pattern,
            // 验证的提示内容
            inputErrorMessage: opts.message||'请输入正确的内容',
          })
    },
    message(message,type,duration=3000){
        //永远保持只有一个打开状态
        ElMessage.closeAll()
        return ElMessage({
            showClose: true,
            dangerouslyUseHTMLString: true,
            message,
            duration,
            type,
        })
    },
    success(message){
        return this.message(message,"success")
    },
    warning(message){
        return this.message(message,"warning")
    },
    error(message){
        return this.message(message,"error")
    },
    notifyError(message){//-------------------------------------新增代码
        return this.notify("提示",message,3000,{type:"error",position:"tr"})
    },
    notifySuccess(message){//-------------------------------------新增代码
        return this.notify("提示",message,3000,{type:"success",position:"tr"})
    },
    notify(title,message,duration=3000,options){
        // 默认值
        var defaultOptions = {type:"info",position:"tr"}
        // 用户传递和默认值就行覆盖处理
        var opts = {...defaultOptions,...options}
        //永远保持只有一个打开状态
        ElNotification.closeAll()
        var positionMap = {
            "tr":"top-right",
            "tl":"top-left",
            "br":"bottom-right",
            "bl":"bottom-left"
        }
        return ElNotification({
            title,
            message,
            duration: duration,
            type:opts.type,
            position: positionMap[opts.position||"tr"],
            dangerouslyUseHTMLString:true
        })
    }
}

export default  KVA
```

现在已经把项目中所以的KVA.message替换成了KVA.notifyError。改动的位置：Login.vue 



# 05、Skeleton 骨架屏开发

## 1：认识

在需要等待加载内容的位置设置一个骨架屏，某些场景下比 Loading 的视觉效果更好。

## 2：应用

### 整体控制

```vue
 <el-skeleton :loading="loading" animated>
 	 <template #template>----------------animated=true
		这里就是骨架屏元素的设置
     </template>
     <template #default> ----------------animated=false
		这里就是原本的内容
     </template>>
 </el-skeleton>
```

- animated : 是否需要动画效果
- loading : true 显示骨架的效果，false，隐藏骨架的效果。

### 单个控制

```vue
<el-skeleton/>
```

默认情况：rows=“3” 其实会显示4个子项出来。

案例如下：

```vue
<div style="background: #fff;padding:10px;">
    <el-skeleton :loading="true" animated> 
        <template #template>
            <el-skeleton-item variant="h1" />
        </template>
        <template #default>
            <h1>我是一个标题</h1>
        </template>
    </el-skeleton> 
</div>
```



## 如何把控制面板的整体页面做成骨架屏幕

1: 局部控制

```vue
 <el-skeleton :loading="loading" animated>
 	 <template #template>----------------animated=true
		这里就是骨架屏元素的设置
     </template>
     <template #default> ----------------animated=false
		这里就是原本的内容
     </template>>
 </el-skeleton>


 <el-skeleton :loading="loading" animated>
 	 <template #template>----------------animated=true
		这里就是骨架屏元素的设置
     </template>
     <template #default> ----------------animated=false
		这里就是原本的内容
     </template>>
 </el-skeleton>


 <el-skeleton :loading="loading" animated>
 	 <template #template>----------------animated=true
		这里就是骨架屏元素的设置
     </template>
     <template #default> ----------------animated=false
		这里就是原本的内容
     </template>>
 </el-skeleton>

 <el-skeleton :loading="loading" animated>
 	 <template #template>----------------animated=true
		这里就是骨架屏元素的设置
     </template>
     <template #default> ----------------animated=false
		这里就是原本的内容
     </template>>
 </el-skeleton>
```



2：学会整体控制（推荐）

```
<el-skeleton :loading="loading" animated>
 	 <template #template>----------------animated=true
		<div></div>
        <div></div>
        <div></div>
        <div></div>
        <div></div>
     </template>
     <template #default> ----------------animated=false
		<div></div>
        <div></div>
        <div></div>
        <div></div>
        <div></div>
     </template>>
 </el-skeleton>
```

案例如下：

```vue
<template>
    <el-skeleton animated :loading="true">
        <template #template>
            <!--card的骨架屏-->
            <div class="gva-card-box statebox">
                <el-row :gutter="12">
                    <el-col :span="6" :xs="12" v-for="i in 4" :key="i">
                        <el-card shadow="hover">
                            <div class="tit"><el-skeleton-item variant="div" style="width: 40%;height: 20px" /></div>
                            <div class="num">
                                <el-skeleton-item variant="div" style="width: 60%;height:30px"/>
                            </div>
                            <div class="info">
                                <el-skeleton-item variant="div" style="width: 80%;height:20px"/>
                            </div>
                        </el-card>
                    </el-col>
                </el-row>
            </div> 
            
            <!--统计报表的骨架-->
            <div class="gva-card-box">
                <div class="gva-card">
                    <div class="card-header"><span ><el-skeleton-item variant="h3" style="width:6%;height:30px;" /></span></div>
                    <div class="echart-box">
                        <div class="el-row" style="margin-left: -10px; margin-right: -10px;">
                            <div class="el-col el-col-24 el-col-xs-24 el-col-sm-18 is-guttered">
                                <div class="dashboard-line-box">
                                    <el-skeleton-item variant="h3" style="height:400px;" />
                                </div>
                            </div>
                            <div class="el-col el-col-24 el-col-xs-24 el-col-sm-6 is-guttered"
                                style="padding-right: 10px; padding-left: 10px;">
                                <div class="commit-table">
                                    <div class="commit-table-title"><el-skeleton-item variant="h3" style="width:20%;height:30px;" /> </div>
                                    <div class="log">
                                        <div class="log-item" v-for="i in 10" :key="i">
                                            <el-skeleton-item variant="h3" style="height:22px;" />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </template>
        <template #default>
            <slot></slot>
        </template>
    </el-skeleton>
</template>
<script setup>
import {useMenuTabStore} from '@/stores/menuTab.js'
const menuTabStore = useMenuTabStore()
const skLoading  = computed(()=>menuTabStore.skLoading)

onMounted(() => {
    setTimeout(() => {
        menuTabStore.skLoading = false;
    }, 1000)
})
</script>
```

如果未来你的骨架屏，不仅仅是一个地方进行设置，比如：控制面板，导航栏，菜单栏，头部栏，表格区域等等。那么如果你想用一个状态来控制他们所有的骨架屏的显示和隐藏。那么你就必须使用状态管理pinia来进行处理如下：

```js
import { defineStore } from 'pinia'

export const useSkeletonStore = defineStore('skeleton', {
  // 定义状态
  state: () => ({
    skLoading:true
  }),

  // 定义动作
  actions: {
   /* 设置loading */ 
   setLoading(loading){
      this.skLoading = loading
   }
  },
  persist: {
    key: 'kva-pinia-skeleton',
    storage: sessionStorage
  }
})
```



# 06、数字动画

1： 安装

```js
pnpm install animated-number-vue3
npm install animated-number-vue3
yarn add animated-number-vue3
```

2:在main.js中引入

```js
import AnimatedNumber from 'animated-number-vue3'
app.use(AnimatedNumber)
```

3: 使用

```js
<AnimatedNumber :from="0" :to="1000"></AnimatedNumber>
```

### 3.具体使用

> 3.1简单模式



```ruby
<AnimatedNumber :from="0" :to="1000"></AnimatedNumber>
```

> 3.2 slot模式



```jsx
//也可使用插槽来自定义界面等操作，option为整个动效包含内容的对象，
//item里面包含变动的数字，当from和to传的是一个数字时，为单数字动效，此时值必须为number
<AnimatedNumber :from="0" :to="1000">
  <template #default="{ option, item }">
    <h1>{{ item.number }}</h1>
  </template>
</AnimatedNumber>
```

> 3.3 复杂模式 如果想一次性为多个数字做动效，此插件也提供插槽来自定义，from和to的对象key必须为一样，一个开始一个结束的值



```xml
<AnimatedNumber :from="{
  number1:0,
  number2:0
}" :to="{
  number1:100,
  number2:100
}">
  <template #default="{ option, item }">
    <h1>{{ item.number1 }}</h1>
    <h1>{{ item.number2 }}</h1>
  </template>
</AnimatedNumber>
```

具体应用如下：

```vue
<template>
    <div class="gva-card-box statebox">
        <el-row :gutter="12">
            <el-col :span="6" :xs="12">
                <el-card shadow="hover">
                    <div class="tit">总销售额</div>
                    <div class="num">
                        <AnimatedNumber :from="0" :to="12560" duration="3000">
                            <template #default="{ option, item }">
                                <span>￥{{ item.number }}</span>
                            </template>
                        </AnimatedNumber>
                    </div>
                    <div class="info">
                        <AnimatedNumber :from="0" :to="8869" duration="3000">
                            <template #default="{ option, item }">
                                <span>日销售额：￥{{ item.number }}</span>
                            </template>
                        </AnimatedNumber>
                    </div>
                </el-card>
            </el-col>
            <el-col :span="6" :xs="12">
                <el-card shadow="hover">
                    <div class="tit">用户注册</div>
                    <div class="num">
                        <AnimatedNumber :from="0" :to="8846" duration="3000" />
                    </div>
                    <div class="info">
                        <AnimatedNumber :from="0" :to="1423" duration="3000">
                            <template #default="{ option, item }">
                                <span>日注册量：{{ item.number }}</span>
                            </template>
                        </AnimatedNumber>
                    </div>
                </el-card>
            </el-col>
            <el-col :span="6" :xs="12">
                <el-card shadow="hover">
                    <div class="tit">访问量</div>
                    <div class="num">
                        <AnimatedNumber :from="0" :to="6560" duration="3000" />
                    </div>
                    <div class="info">
                        <AnimatedNumber :from="0" :to="2423" duration="3000">
                            <template #default="{ option, item }">
                                <span>日访问量：{{ item.number }}</span>
                            </template>
                        </AnimatedNumber>
                    </div>
                </el-card>
            </el-col>
            <el-col :span="6" :xs="12">
                <el-card shadow="hover">
                    <div class="tit">支付笔数</div>
                    <div class="num">
                        <AnimatedNumber :from="0" :to="12589" duration="3000" />
                    </div>
                    <div class="info">
                        <AnimatedNumber :from="0" :to="80" duration="3000">
                            <template #default="{ option, item }">
                                <span>转化率：{{ item.number }} %</span>
                            </template>
                        </AnimatedNumber>
                    </div>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>
<script setup>
</script>
<style lang="scss" scoped>
.statebox {
    .el-card {
        padding: 10px;
        margin-bottom: 5px;
    }

    .tit {
        height: 22px;
        color: rgba(0, 0, 0, .45);
        font-size: 14px;
        line-height: 22px;
    }

    .num {
        height: 38px;
        margin-top: 10px;
        color: rgba(0, 0, 0, .75);
        font-size: 32px;
        line-height: 38px;
        white-space: nowrap;
        font-weight: bold;
        text-overflow: ellipsis;
        word-break: break-all;
    }

    .info {
        position: relative;
        width: 100%;
        margin-top: 20px;
    }
}</style>
```



# 07、如何开发保存

- 1： 先要确定你的id是不是数字，并且是不是自增的，如果不是：那就说你代码必须要手动方式去添加和维护。如果是自增，那么你程序是不需要去任何关于id的处理。
- 2：一定考虑那些字段是不为空，如果不为空的那么这列必须在代码用参数传递或者在项目写默认值
- 3：如果有些是又默认值。那么可以考虑不用传递。（status）
- 4： 创建和更新时间看看是否数据区维护还是框架维护了。



# 08、MD5加盐加密

```go
package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// 参数：需要加密的字符串
func getMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// md5加密
func Md5(str string) string {
	return getMd5(getMd5(PWD_SALT + str + PWD_SALT))
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

```

注意加密以后，那么就必须考虑到login的时候要加密比对

```go
// 登录的接口处理
func (api *LoginApi) ToLogined(c *gin.Context) {
	type LoginParam struct {
		Account  string
		Code     string
		CodeId   string
		Password string
	}

	// 1：获取用户在页面上输入的账号和密码开始和数据库里数据进行校验
	param := LoginParam{}
	err2 := c.ShouldBindJSON(&param)
	if err2 != nil {
		response.Fail(60002, "参数绑定有误", c)
		return
	}

	if len(param.Code) == 0 {
		response.Fail(60002, "请输入验证码", c)
		return
	}

	if len(param.CodeId) == 0 {
		response.Fail(60002, "验证码获取失败", c)
		return
	}

	// 开始校验验证码是否正确
	verify := store.Verify(param.CodeId, param.Code, true)
	if !verify {
		response.Fail(60002, "你输入的验证码有误!!", c)
		return
	}

	inputAccount := param.Account
	inputPassword := param.Password

	if len(inputAccount) == 0 {
		response.Fail(60002, "请输入账号", c)
		return
	}

	if len(inputPassword) == 0 {
		response.Fail(60002, "请输入密码", c)
		return
	}

	dbUser, err := sysUserService.GetUserByAccount(inputAccount)
	if err != nil {
		response.Fail(60002, "你输入的账号和密码有误", c)
		return
	}

	// 这个时候就判断用户输入密码和数据库的密码是否一致
	// inputPassword = utils.Md5(123456) = 2ec9f77f1cde809e48fabac5ec2b8888
	// dbUser.Password = 2ec9f77f1cde809e48fabac5ec2b8888
	if dbUser != nil && dbUser.Password == utils.Md5(inputPassword) {//--=--------------这里是修改的代码
		token := api.generaterToken(c, dbUser)
		// 根据用户id查询用户的角色
		roles := [2]map[string]any{}
		m1 := map[string]any{"id": 1, "name": "超级管理员", "code": "admin"}
		m2 := map[string]any{"id": 2, "name": "财务", "code": "visitor"}
		roles[0] = m1
		roles[1] = m2
		// 根据用户id查询用户的角色的权限
		permissions := [2]map[string]any{}
		pm1 := map[string]any{"code": 10001, "name": "保存用户"}
		pm2 := map[string]any{"code": 20001, "name": "删除用户"}
		permissions[0] = pm1
		permissions[1] = pm2

		response.Ok(map[string]any{"user": dbUser, "token": token, "roles": roles, "permissions": permissions}, c)
	} else {
		response.Fail(60002, "你输入的账号和密码有误", c)
	}
}

```



# 系统变更

1： 把前端移动项目 web目录

2： 新建一个sql目录

3： 读取表的信息

````sql
```sql
-- 1. 根据库名获取所有表的信息
-- 使用以下SQL语句来获取指定数据库中所有表的信息：
SELECT * FROM information_schema.`TABLES` WHERE TABLE_SCHEMA='kva-admin-db';
-- 其中，“your_database”是你要查询的数据库名称。这条语句将返回一个包含所有表信息的数据表。

-- 2. 根据库名获取所有表名称和表说明
-- 使用以下SQL语句来获取指定数据库中所有表的名称和注释：

SELECT TABLE_NAME, TABLE_COMMENT FROM information_schema.`TABLES` WHERE TABLE_SCHEMA='kva-admin-db';

-- 这条语句将返回一个包含表名和注释的数据表。
-- 3. 根据库名获取所有的字段信息
-- 使用以下SQL语句来获取指定数据库中所有表的所有字段信息：

SELECT 
	TABLE_SCHEMA AS 'schema',
	TABLE_NAME AS 'tablename',
	COLUMN_NAME AS 'cname',
	ORDINAL_POSITION AS 'position',
	COLUMN_DEFAULT AS 'cdefault',
	IS_NULLABLE AS 'nullname',
	DATA_TYPE AS 'dataType',
	CHARACTER_MAXIMUM_LENGTH AS 'maxLen',
	NUMERIC_PRECISION AS 'precision',
	NUMERIC_SCALE AS 'scalename',
	COLUMN_TYPE AS 'ctype',
	COLUMN_KEY AS 'ckey',
	EXTRA AS 'extra',
	COLUMN_COMMENT AS 'comment'
FROM information_schema.`COLUMNS`
WHERE TABLE_SCHEMA='kva-admin-db'
ORDER BY TABLE_NAME, ORDINAL_POSITION;

-- 这条语句将返回一个包含所有字段信息的数据表。
-- 4. 根据库名获取所有的库和表字段的基本信息
-- 使用以下SQL语句来获取指定数据库中所有表和字段的基本信息：

SELECT C.TABLE_SCHEMA AS 'schema',
T.TABLE_NAME AS 'tablename',
T.TABLE_COMMENT AS 'tablecomment',
C.COLUMN_NAME AS 'columnname',
C.COLUMN_COMMENT AS 'columncomment',
C.ORDINAL_POSITION AS 'position',
C.COLUMN_DEFAULT AS 'columndefault',
C.IS_NULLABLE AS 'nullname',
C.DATA_TYPE AS 'dataType',
C.CHARACTER_MAXIMUM_LENGTH AS 'maxlen',
C.NUMERIC_PRECISION AS 'precision',
C.NUMERIC_SCALE AS 'scale',
C.COLUMN_TYPE AS 'ctype',
C.COLUMN_KEY AS 'ckey',
C.EXTRA AS 'extra'
FROM information_schema.`TABLES` T
LEFT JOIN information_schema.`COLUMNS` C ON T.TABLE_NAME=C.TABLE_NAME AND T.TABLE_SCHEMA=C.TABLE_SCHEMA
WHERE T.TABLE_SCHEMA='kva-admin-db'
ORDER BY C.TABLE_NAME, C.ORDINAL_POSITION;
```
````

4: 在utils新建目录adr

 [AES.go](C:\Users\zxc\go\xkginweb\utils\adr\AES.go)  [BASE64.go](C:\Users\zxc\go\xkginweb\utils\adr\BASE64.go)  [DES.go](C:\Users\zxc\go\xkginweb\utils\adr\DES.go)  [MD5.go](C:\Users\zxc\go\xkginweb\utils\adr\MD5.go)  [RSA.go](C:\Users\zxc\go\xkginweb\utils\adr\RSA.go) 

5： 数据载体的划分

# 9、数据载体

pojo

- model、entity、pojo
- po、context
- vo —-返回



# 10、后端校验validate

官网：更多请查看  https://github.com/gookit/validate

1: 安装验证框架

```go
go get github.com/gookit/validate
```

2:  验证三部曲

结构体可以实现 3 个接口方法，方便做一些自定义：

- `ConfigValidation(v *Validation)` 将在创建验证器实例后调用
- `Messages() map[string]string` 可以自定义==验证器==错误消息
- `Translates() map[string]string` 可以自定义字段映射/翻译

3: 验证

验证结构体增加验证器

```go
type UserForm struct {
  Name     string    `validate:"required|minLen:7"`
  Email    string    `validate:"email" message:"email is invalid"`
  Age      int       `validate:"required|int|min:1|max:99" message:"int:age must int| min: age min value is 1"`
  CreateAt int       `validate:"min:1"`
  Safe     int       `validate:"-"`
  UpdateAt time.Time `validate:"required"`
  Code     string    `validate:"customValidator"` // 使用自定义验证器
}
```

开始验证结构体

```go
// 创建 Validation 实例
v := validate.Struct(u)
 if v.Validate() { // 验证成功
    // do something ...
} else {
    fmt.Println(v.Errors) // 所有的错误消息
    fmt.Println(v.Errors.One()) // 返回随机一条错误消息
    fmt.Println(v.Errors.Field("Name")) // 返回该字段的错误消息
}
```



| 验证器/别名                               | 描述信息                                                     |
| :---------------------------------------- | :----------------------------------------------------------- |
| `required`                                | 字段为必填项，值不能为空                                     |
| `required_if/requiredIf`                  | `required_if:anotherfield,value,...` 如果其它字段 *anotherField* 为任一值 *value* ，则此验证字段必须存在且不为空。 |
| `required_unless/requiredUnless`          | `required_unless:anotherfield,value,...` 如果其它字段 *anotherField* 不等于任一值 *value* ，则此验证字段必须存在且不为空。 |
| `required_with/requiredWith`              | `required_with:foo,bar,...` 在其他任一指定字段出现时，验证的字段才必须存在且不为空 |
| `required_with_all/requiredWithAll`       | `required_with_all:foo,bar,...` 只有在其他指定字段全部出现时，验证的字段才必须存在且不为空 |
| `required_without/requiredWithout`        | `required_without:foo,bar,...` 在其他指定任一字段不出现时，验证的字段才必须存在且不为空 |
| `required_without_all/requiredWithoutAll` | `required_without_all:foo,bar,...` 只有在其他指定字段全部不出现时，验证的字段才必须存在且不为空 |
| `-/safe`                                  | 标记当前字段是安全的，无需验证                               |
| `int/integer/isInt`                       | 检查值是 `intX` `uintX` 类型，同时支持大小检查 `"int"` `"int:2"` `"int:2,12"` |
| `uint/isUint`                             | 检查值是 `uintX` 类型(`value >= 0`)                          |
| `bool/isBool`                             | 检查值是布尔字符串(`true`: "1", "on", "yes", "true", `false`: "0", "off", "no", "false"). |
| `string/isString`                         | 检查值是字符串类型，同时支持长度检查 `"string"` `"string:2"` `"string:2,12"` |
| `float/isFloat`                           | 检查值是 float(`floatX`) 类型                                |
| `slice/isSlice`                           | 检查值是 slice 类型(`[]intX` `[]uintX` `[]byte` `[]string` 等). |
| `in/enum`                                 | 检查值()是否在给定的枚举列表(`[]string`, `[]intX`, `[]uintX`)中 |
| `not_in/notIn`                            | 检查值不是在给定的枚举列表中                                 |
| `contains`                                | 检查输入值(`string` `array/slice` `map`)是否包含给定的值     |
| `not_contains/notContains`                | 检查输入值(`string` `array/slice` `map`)是否不包含给定值     |
| `string_contains/stringContains`          | 检查输入的 `string` 值是否不包含给定sub-string值             |
| `starts_with/startsWith`                  | 检查输入的 `string` 值是否以给定sub-string开始               |
| `ends_with/endsWith`                      | 检查输入的 `string` 值是否以给定sub-string结束               |
| `range/between`                           | 检查值是否为数字且在给定范围内                               |
| `max/lte`                                 | 检查输入值小于或等于给定值                                   |
| `min/gte`                                 | 检查输入值大于或等于给定值(for `intX` `uintX` `floatX`)      |
| `eq/equal/isEqual`                        | 检查输入值是否等于给定值                                     |
| `ne/notEq/notEqual`                       | 检查输入值是否不等于给定值                                   |
| `lt/lessThan`                             | 检查值小于给定大小(use for `intX` `uintX` `floatX`)          |
| `gt/greaterThan`                          | 检查值大于给定大小(use for `intX` `uintX` `floatX`)          |
| `int_eq/intEq/intEqual`                   | 检查值为int且等于给定值                                      |
| `len/length`                              | 检查值长度等于给定大小(use for `string` `array` `slice` `map`). |
| `min_len/minLen/minLength`                | 检查值的最小长度是给定大小                                   |
| `max_len/maxLen/maxLength`                | 检查值的最大长度是给定大小                                   |
| `email/isEmail`                           | 检查值是Email地址字符串                                      |
| `regex/regexp`                            | 检查该值是否可以通过正则验证                                 |
| `arr/array/isArray`                       | 检查值是数组`array`类型                                      |
| `map/isMap`                               | 检查值是 `map` 类型                                          |
| `strings/isStrings`                       | 检查值是字符串切片类型(`[]string`)                           |
| `ints/isInts`                             | 检查值是`int` slice类型(only allow `[]int`)                  |
| `eq_field/eqField`                        | 检查字段值是否等于另一个字段的值                             |
| `ne_field/neField`                        | 检查字段值是否不等于另一个字段的值                           |
| `gte_field/gtField`                       | 检查字段值是否大于另一个字段的值                             |
| `gt_field/gteField`                       | 检查字段值是否大于或等于另一个字段的值                       |
| `lt_field/ltField`                        | 检查字段值是否小于另一个字段的值                             |
| `lte_field/lteField`                      | 检查字段值是否小于或等于另一个字段的值                       |
| `file/isFile`                             | 验证是否是上传的文件                                         |
| `image/isImage`                           | 验证是否是上传的图片文件，支持后缀检查                       |
| `mime/mimeType/inMimeTypes`               | 验证是否是上传的文件，并且在指定的MIME类型中                 |
| `date/isDate`                             | 检查字段值是否为日期字符串。（只支持几种常用的格式） eg `2018-10-25` |
| `gt_date/gtDate/afterDate`                | 检查输入值是否大于给定的日期字符串                           |
| `lt_date/ltDate/beforeDate`               | 检查输入值是否小于给定的日期字符串                           |
| `gte_date/gteDate/afterOrEqualDate`       | 检查输入值是否大于或等于给定的日期字符串                     |
| `lte_date/lteDate/beforeOrEqualDate`      | 检查输入值是否小于或等于给定的日期字符串                     |
| `hasWhitespace`                           | 检查字符串值是否有空格                                       |
| `ascii/ASCII/isASCII`                     | 检查值是ASCII字符串                                          |
| `alpha/isAlpha`                           | 验证值是否仅包含字母字符                                     |
| `alpha_num/alphaNum/isAlphaNum`           | 验证是否仅包含字母、数字                                     |
| `alpha_dash/alphaDash/isAlphaDash`        | 验证是否仅包含字母、数字、破折号（ - ）以及下划线（ _ ）     |
| `multi_byte/multiByte/isMultiByte`        | 检查值是多字节字符串                                         |
| `base64/isBase64`                         | 检查值是Base64字符串                                         |
| `dns_name/dnsName/DNSName/isDNSName`      | 检查值是DNS名称字符串                                        |
| `data_uri/dataURI/isDataURI`              | Check value is DataURI string.                               |
| `empty/isEmpty`                           | 检查值是否为空                                               |
| `hex_color/hexColor/isHexColor`           | 检查值是16进制的颜色字符串                                   |
| `hexadecimal/isHexadecimal`               | 检查值是十六进制字符串                                       |
| `json/JSON/isJSON`                        | 检查值是JSON字符串。                                         |
| `lat/latitude/isLatitude`                 | 检查值是纬度坐标                                             |
| `lon/longitude/isLongitude`               | 检查值是经度坐标                                             |
| `mac/isMAC`                               | 检查值是MAC字符串                                            |
| `num/number/isNumber`                     | 检查值是数字字符串. `>= 0`                                   |
| `cn_mobile/cnMobile/isCnMobile`           | 检查值是中国11位手机号码字符串                               |
| `printableASCII/isPrintableASCII`         | Check value is PrintableASCII string.                        |
| `rgbColor/RGBColor/isRGBColor`            | 检查值是RGB颜色字符串                                        |
| `full_url/fullUrl/isFullURL`              | 检查值是完整的URL字符串(*必须以http,https开始的URL*).        |
| `url/URL/isURL`                           | 检查值是URL字符串                                            |
| `ip/IP/isIP`                              | 检查值是IP（v4或v6）字符串                                   |
| `ipv4/isIPv4`                             | 检查值是IPv4字符串                                           |
| `ipv6/isIPv6`                             | 检查值是IPv6字符串                                           |
| `cidr/CIDR/isCIDR`                        | 检查值是 CIDR 字符串                                         |
| `CIDRv4/isCIDRv4`                         | 检查值是 CIDR v4 字符串                                      |
| `CIDRv6/isCIDRv6`                         | 检查值是 CIDR v6 字符串                                      |
| `uuid/isUUID`                             | 检查值是UUID字符串                                           |
| `uuid3/isUUID3`                           | 检查值是UUID3字符串                                          |
| `uuid4/isUUID4`                           | 检查值是UUID4字符串                                          |
| `uuid5/isUUID5`                           | 检查值是UUID5字符串                                          |
| `filePath/isFilePath`                     | 检查值是一个存在的文件路径                                   |
| `unixPath/isUnixPath`                     | 检查值是Unix Path字符串                                      |
| `winPath/isWinPath`                       | 检查值是Windows路径字符串                                    |
| `isbn10/ISBN10/isISBN10`                  | 检查值是ISBN10字符串                                         |
| `isbn13/ISBN13/isISBN13`                  | 检查值是ISBN13字符串                                         |



# 11 、 关于逻辑删除得问题和处理

## 01、gorm默认机制

### gorm.Model

GORM 定义一个 `gorm.Model` 结构体，其包括字段 `ID`、`CreatedAt`、`UpdatedAt`、`DeletedAt ` 。

- 其中这里得deletedAt就是用于逻辑删除控制得字段，如果null 就代表没有删除，如果有时间就说你执行过delete from 才会把删除时写入到数据库表中
- 如果有字段，未来做任何得查询都自动跟上条件deletedAt is null 



### 修改逻辑删除的默认规则

如果你先修改默认的规则，从时间变成0/1这种方式，你必须如下执行

1: 先安装组件

```go
gorm.io/plugin/soft_delete
```

2: 把deletedAt删掉，修改如下：

IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt;default:0" json:"isDeleted"`

```go
package global

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type GVA_MODEL struct {
	ID        uint      `gorm:"primarykey;comment:主键ID" json:"id"` // 主键ID
	CreatedAt time.Time `gorm:"type:datetime(0);autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime(0);autoUpdateTime;comment:更新时间" json:"updatedAt"`
	//DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"` // 删除时间
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt;default:0" json:"isDeleted"`
}

```

前面的那些步骤和流程一个都不能省去，比如：注册表

```go
package orm

import (
	"xkginweb/global"
	bbs2 "xkginweb/model/entity/bbs"
	"xkginweb/model/entity/jwt"
	sys2 "xkginweb/model/entity/sys"
	user2 "xkginweb/model/entity/user"
	video2 "xkginweb/model/entity/video"
)

func RegisterTable() {
	db := global.KSD_DB
	// 注册和声明model
	db.AutoMigrate(user2.XkUser{})
	db.AutoMigrate(user2.XkUserAuthor{})
	// 系统用户，角色，权限表
	db.AutoMigrate(sys2.SysApis{})
	db.AutoMigrate(sys2.SysMenus{})
	db.AutoMigrate(sys2.SysRoleApis{})
	db.AutoMigrate(sys2.SysRoleMenus{})
	db.AutoMigrate(sys2.SysRoles{})
	db.AutoMigrate(sys2.SysUserRoles{})
	db.AutoMigrate(sys2.SysUser{})
	// 视频表
	db.AutoMigrate(video2.XkVideo{})
	db.AutoMigrate(video2.XkVideoCategory{})
	db.AutoMigrate(video2.XkVideoChapterLesson{})
	// 社区
	db.AutoMigrate(bbs2.XkBbs{})
	db.AutoMigrate(bbs2.BbsCategory{})

	// 声明一下jwt模型
	db.AutoMigrate(jwt.JwtBlacklist{})
}

```



3：然后重启查看效果即可。

- 其中这里得isDeleted就是用于逻辑删除控制得字段，如果0就代表没有删除，如果是1就是删除
- 未来你执行任何的删除操作就变成update table set is_deleted = 1,update_time = now() where id = 1
- 未来做任何得查询都自动跟上条件is_deleted = 0



4: 我要把删除和未删除全部查询出来？

往往在做后台管理系统的时候，你就必须要把删除和未删除全部查询出来。那么你就必须加上：.Unscoped() 来进行处理这样会把默认机制打破。不在跟已删除过滤。如下：

```go
// 查询分页
func (service *SysUserService) LoadSysUserPage(info request.PageInfo) (list interface{}, total int64, err error) {
	// 获取分页的参数信息
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 准备查询那个数据库表
	db := global.KSD_DB.Model(&sys.SysUser{})

	// 准备切片帖子数组
	var sysUserList []sys.SysUser

	// 加条件
	if info.Keyword != "" {
		db = db.Where("(username like ? or account like ? )", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}

	// 排序默时间降序降序
	db = db.Order("created_at desc")

	// 查询中枢
	err = db.Unscoped().Count(&total).Error
	if err != nil {
		return sysUserList, total, err
	} else {
		// 执行查询
		err = db.Unscoped().Limit(limit).Offset(offset).Find(&sysUserList).Error
	}

	// 结果返回
	return sysUserList, total, err
}
```





# 12：查询用户的角色

1：你要把所有的角色查询出来。

- 放入表格
- 放入下拉框

2：  查询用户的角色

3： 把查询用户角色和所有角色进行碰撞，如果一致就选中

​	

## 为什么用户角色喜欢弄一个中间表来维护

- 一对一的关系

| userid       | username | roleid | roleName   |
| ------------ | -------- | ------ | ---------- |
| 1            | 飞哥     | 1      | 超级管理员 |
| 2            | 小玉     | 2      | 财务       |
| 一对多的关系 |          |        |            |

| userid | username | roleid | roleName            |
| ------ | -------- | ------ | ------------------- |
| 1      | 飞哥     | 1,3    | 超级管理员,开发人员 |
| 2      | 小玉     | 2,3    | 财务,开发人员       |

那为什么不用上面的方式来维护。因为如果我们变更角色名字的。那么就不方便角色用户数据的维护



### 真正的设计

| userId | userName |      |
| ------ | -------- | ---- |
| 1      | 飞哥     | 1    |
| 2      | 小玉     |      |
| 3      | 小伟     |      |
|        |          |      |
|        |          |      |

| roleId | roleName   |      |
| ------ | ---------- | ---- |
| 1      | 超级管理员 |      |
| 2      | 财务       |      |
| 3      | 开发人员   |      |

map概念

sys_user_roles

| userId | roleid |      |
| ------ | ------ | ---- |
| 1      | 3      |      |
| 1      | 1      |      |
| 2      | 2      |      |
| 2      | 3      |      |
| 3      | 3      |      |
|        |        |      |
|        |        |      |
|        |        |      |







## 授予用户角色

- 把roleids全部拿到
- 并且拿到用户的id
- 然后调用授予角色的接口
  - 根据用户id删除对应角色
  - 然后把新角色全部重新保存进sys_user_roles

```
var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
db.Create(&users)

var sysUserRoles =  []SysUserRoles{{UserId: 1,RoleId:1},{UserId: 1,RoleId:2},{UserId: 1,RoleId:3}}
db.Create(&users)
```





# 13、事务的简单认识

对于要把事务在实际中使用好，需要了解事务的特性。

事务的四大特性主要是：原子性（Atomicity）、一致性（Consistency）、隔离性（Isolation）、持久性（Durability）。

**一、事务的四大特性**

**1.1** **原子性（Atomicity）**

原子性是指事务是一个不可分割的工作单位，事务中的操作要么全部成功，要么全部失败。比如在同一个事务中的SQL语句，要么全部执行成功，要么全部执行失败。



```mysql
begin transaction;
    update account set money = money-100 where name = '张三';
    update account set money = money+100 where name = '李四';
commit transaction;
```

**1.2** **一致性（Consistency）**

官网上事务一致性的概念是：事务必须使数据库从一个一致性状态变换到另外一个一致性状态。

换一种方式理解就是：事务按照预期生效，数据的状态是预期的状态。

举例说明：张三向李四转100元，转账前和转账后的数据是正确的状态，这就叫一致性，如果出现张三转出100元，李四账号没有增加100元这就出现了数据错误，就没有达到一致性。

**1.3** **隔离性（Isolation）**

事务的隔离性是多个用户并发访问数据库时，数据库为每一个用户开启的事务，不能被其他事务的操作数据所干扰，多个并发事务之间要相互隔离。

**1.4** **持久性（Durability）**

持久性是指一个事务一旦被提交，它对数据库中数据的改变就是永久性的，接下来即使数据库发生故障也不应该对其有任何影响。

例如我们在使用JDBC操作数据库时，在提交事务方法后，提示用户事务操作完成，当我们程序执行完成直到看到提示后，就可以认定事务以及正确提交，即使这时候数据库出现了问题，也必须要将我们的事务完全执行完成，否则就会造成我们看到提示事务处理完毕，但是数据库因为故障而没有执行事务的重大错误。



## 什么情况下会用到事务呢？

在开发中如果，牵涉到业务执行方法，处理的写入（insert.update,delete）的时候，如果存在多种写入，你要保证他们执行顺序和整体性你就必须靠事务。因为事务可以把这些写入指令全部放入到一个事务队列中，来进行操作。然后在操作的过程中，如果写入这些sql语言，会执行sql语句指令，但是不并不会马上写入到磁盘数据中。暂时只会全部在内存进行记录和处理。只有触发到commit指令的时候，全会进行比对然后开始写入到数据库表中，如果遇到的rollback就会之前执行的sql全部回滚掉不去持久化数据库表中。





# 14、 gorm框架更新0值失败的问题

参考文献：https://gorm.io/zh_CN/docs/update.html

默认情况下：gorm框架更新结构体的时候，只能更新那些非0列。如果你更新为0的列那么久必须使用map

解决方案：

1：直接把0该换其它的状态。（一般不会使用）

2:   把结构体转化成map的方式来进行处理即可

1: 安装组件

```go
go get github.com/fatih/structs
```

2: 定义结构体

```go
package model

type SysUser struct {
	ID        uint   `json:"ID" structs:"omitnested"`
	UUID      string `json:"uuid" structs:"omitnested" ` // 用户UUID
	Slat      string `json:"slat" structs:"omitnested" ` // 用户登录密码
	Enable    int    `json:"enable" structs:"enable" `
	Account   string `json:"account" structs:"account"`    // 用户登录名
	Password  string `json:"password" structs:"password" ` // 密码加盐
	Username  string `json:"username" structs:"username" ` // 用户昵称
	Avatar    string `json:"avatar" structs:"avatar" `     // 用户头像
	Phone     string `json:"phone" structs:"phone" `       // 用户手机号
	Email     string `json:"email" structs:"email" `
	IsDeleted int    `json:"email" structs:"is_deleted"`
}

```

3: 写个测试

```go
package main

import (
	"fmt"
	"github.com/fatih/structs"
	"strutstomap/model"
)

func main() {

	sysUser := model.SysUser{}
	sysUser.ID = 1
	sysUser.UUID = "1111"
	sysUser.Slat = "1111"
	sysUser.Avatar = "XXXXX"
	sysUser.Email = "xxxx@qq.com"
	sysUser.Username = "飞飞"
	sysUser.Account = "feige"
	sysUser.IsDeleted = 0

	fmt.Println("user to map：", structs.Map(sysUser))

}

```

可以看到id,uuid,slat被忽略掉了。而增加structs的列都会按照你指定的列名作为map的key







# 15、泛型的应用

问题：在项目中，我们定义service其实你会发现基本单表的操作80%~90%的操作几乎一摸一样，只是更改结构体。其他并没变化。但是你在开发时候，如果你不考虑到封装，其实就必须每个表就对应结构体，对应service然后调用

- Create
- Updates
- Delete
- First
- Find

那么有有一种方式可以将其这些基本CURD单表操作全部进行封装，给与模块开发提供遍历。

在java有很多持久层框架，用最多的呢是mybati，这里的mybatis等价于gorm框架.

在java中有封装注著名的框架：mybatis-plus。它其实就在单表操作层面全部进行封装，单表操作基本全部进行封装和简化。

## 封装知识点

- 继承
- 泛型
- 反射

## 继承你的理解是什么？

- 继承其实就拥有父类方法和属性（必须大开头方法和属性名字，这种方法和属性就公开，也就可以被子类继承）
- 继承本质特点:   不劳而获，职责分担

## 封装思维

- 方法封装 —————————–package
- 结构体封装（父类）————–自动





## 什么时候用继承

如果在开发中存在多个结构体，有通用共用的方法的时候，可以考虑把这类方法用父类来尽心封装。找每个结构方法中很多的代码片段会经常重复的编写，你可以把这些重复的代码片段用父类来完成。





# 16、实现菜单管理和处理

![image-20230819201422971](images/image-20230819201422971.png)

![image-20230819201649380](images/image-20230819201649380.png)

1： 表的设计

```sql
CREATE TABLE `sys_menus` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `parent_id` bigint(20) DEFAULT NULL COMMENT '父菜单ID',
  `path` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '路由path',
  `name` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '路由name 用于国际化处理',
  `hidden` tinyint(1) DEFAULT NULL COMMENT '是否在列表隐藏',
  `component` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '对应前端文件路径',
  `sort` bigint(20) DEFAULT NULL COMMENT '排序标记',
  `icon` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '对应前端文件路径',
  `is_deleted` bigint(20) unsigned DEFAULT '0',
  `title` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单名称',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='菜单表';
```

2: 写查询

- 先把根节点查询出来，然后在遍历查询对应子节点
- 把所有的数据都查询出来，然后在来进行数据遍历和分配（推荐）

3：定义结构体

```go
package sys

import (
	"xkginweb/global"
)

type SysMenus struct {
	global.GVA_MODEL
	ParentId  uint   `json:"parentId" gorm:"comment:父菜单ID"`      // 父菜单ID
	Path      string `json:"path" gorm:"comment:路由path"`         // 路由path
	Title     string `json:"title" gorm:"comment:菜单名称"`          // 菜单名称
	Name      string `json:"name" gorm:"comment:路由name 用于国际化处理"` // 路由name 用于国际化处理
	Hidden    bool   `json:"hidden" gorm:"comment:是否在列表隐藏"`      // 是否在列表隐藏
	Component string `json:"component" gorm:"comment:对应前端文件路径"`  // 对应前端文件路径
	Sort      int    `json:"sort" gorm:"comment:排序标记"`           // 排序标记
	Icon      string `json:"component" gorm:"comment:对应前端文件路径"`  // 对应前端文件路径
	// 忽略该字段，- 表示无读写，-:migration 表示无迁移权限，-:all 表示无读写迁移权限
	Children []*SysMenus `gorm:"-" json:"children"`
	TopObj   *SysMenus   `gorm:"-" json:"-"`
}

func (s *SysMenus) TableName() string {
	return "sys_menus"
}

```

4: 定义service

```go
func (service *SysMenusService) FinMenus(keyword string) (sysMenus []*sys.SysMenus, err error) {
	db := global.KSD_DB.Unscoped().Order("sort asc")
	if len(keyword) > 0 {
		db.Where("title like ?", "%"+keyword+"%")
	}
	err = db.Find(&sysMenus).Error
	return sysMenus, err
}

/**
*   开始把数据进行编排--递归
*   Tree(all,0)
 */
func (service *SysMenusService) Tree(allSysMenus []*sys.SysMenus, parentId uint) []*sys.SysMenus {
	var nodes []*sys.SysMenus
	for _, dbMenu := range allSysMenus {
		if dbMenu.ParentId == parentId {
			childrensMenu := service.Tree(allSysMenus, dbMenu.ID)
			if len(childrensMenu) > 0 {
				dbMenu.Children = append(dbMenu.Children, childrensMenu...)
			}
			nodes = append(nodes, dbMenu)
		}
	}
	return nodes
}

```

5：定义路由 api/ [sys_menus.go](C:\Users\zxc\go\xkginweb\api\v1\sys\sys_menus.go) 

```go
// 查询菜单
func (api *SysMenuApi) FindMenus(c *gin.Context) {
	keyword := c.Query("keyword")
	sysMenus, err := sysMenuService.FinMenus(keyword)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.Ok(sysMenuService.Tree(sysMenus, 0), c)
}

```

6: 分配路由地址 router/sys/ [sys_memus.go](C:\Users\zxc\go\xkginweb\router\sys\sys_memus.go) 

```go
package sys

import (
	"github.com/gin-gonic/gin"
	v1 "xkginweb/api/v1"
)

// 登录路由
type SysMenusRouter struct{}

func (r *SysMenusRouter) InitSysMenusRouter(Router *gin.RouterGroup) {
	sysMenuApi := v1.WebApiGroupApp.Sys.SysMenuApi
	// 用组定义--（推荐）
	router := Router.Group("/sys")
	{
		// 获取菜单列表
		router.POST("/menus/tree", sysMenuApi.FindMenus)//------------新增
		// 查询父级菜单
		router.POST("/menus/root", sysMenuApi.FindMenusRoot)
		// 保存
		router.POST("/menus/save", sysMenuApi.SaveData)
		// 修改
		router.POST("/menus/update", sysMenuApi.UpdateById)
		// 启用和未启用 （控制启用，发布，删除）
		router.POST("/menus/update/status", sysMenuApi.UpdateStatus)
		// 删除单个 :id 获取参数的时候id := c.Param("id")，传递的时候/sys/user/del/100
		router.POST("/menus/del/:id", sysMenuApi.DeleteById)
		// 查询明细 /user/get/1/xxx
		router.POST("/menus/get/:id", sysMenuApi.GetById)
	}
}

```

7: 注册路由 [init_router.go](C:\Users\zxc\go\xkginweb\initilization\init_router.go) 

```go
package initilization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"xkginweb/commons/filter"
	"xkginweb/commons/middle"
	"xkginweb/global"
	"xkginweb/router"
)

func InitGinRouter() *gin.Engine {
	// 打印gin的时候日志是否用颜色标出
	//gin.ForceConsoleColor()
	//gin.DisableConsoleColor()
	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 创建gin服务
	ginServer := gin.Default()
	// 提供服务组
	courseRouter := router.RouterWebGroupApp.Course.CourseRouter

	videoRouter := router.RouterWebGroupApp.Video.XkVideoRouter

	userStateRouter := router.RouterWebGroupApp.State.UserStateRouter

	bbsRouter := router.RouterWebGroupApp.BBs.XkBbsRouter
	bbsCategoryRouter := router.RouterWebGroupApp.BBs.BBSCategoryRouter

	loginRouter := router.RouterWebGroupApp.Login.LoginRouter
	logoutRouter := router.RouterWebGroupApp.Login.LogoutRouter
	codeRouter := router.RouterWebGroupApp.Code.CodeRouter

	sysMenusRouter := router.RouterWebGroupApp.Sys.SysMenusRouter
	sysApisRouter := router.RouterWebGroupApp.Sys.SysApisRouter
	sysUserRouter := router.RouterWebGroupApp.Sys.SysUsersRouter
	sysRolesRouter := router.RouterWebGroupApp.Sys.SysRolesRouter
	sysUserRolesRouter := router.RouterWebGroupApp.Sys.SysUserRolesRouter
	sysRoleMenusRouter := router.RouterWebGroupApp.Sys.SysRoleMenusRouter
	sysRoleApisRouter := router.RouterWebGroupApp.Sys.SysRoleApisRouter

	// 解决接口的跨域问题
	ginServer.Use(filter.Cors())
	// 接口隔离，比如登录，健康检查都不需要拦截和做任何的处理
	// 业务模块接口，
	privateGroup := ginServer.Group("/api")
	// 无需jwt拦截
	{
		loginRouter.InitLoginRouter(privateGroup)
		codeRouter.InitCodeRouter(privateGroup)
	}
	// 会被jwt拦截
	privateGroup.Use(middle.JWTAuth()).Use(middle.Casbin())
	{
		logoutRouter.InitLogoutRouter(privateGroup)
		videoRouter.InitXkVideoRouter(privateGroup)
		courseRouter.InitCourseRouter(privateGroup)
		userStateRouter.InitUserStateRouter(privateGroup)
		bbsRouter.InitXkBbsRouter(privateGroup)
		bbsCategoryRouter.InitBBSCategoryRouter(privateGroup)
		sysMenusRouter.InitSysMenusRouter(privateGroup)
		sysUserRouter.InitSysUsersRouter(privateGroup)
		sysRolesRouter.InitSysRoleRouter(privateGroup)
		sysApisRouter.InitSysApisRouter(privateGroup)
		sysUserRolesRouter.InitSysUserRolesRouter(privateGroup)
		sysRoleMenusRouter.InitSysRoleMenusRouter(privateGroup)
		sysRoleApisRouter.InitSysRoleApisRouter(privateGroup)
	}

	fmt.Println("router register success")
	return ginServer
}

func RunServer() {

	// 初始化路由
	Router := InitGinRouter()
	// 为用户头像和文件提供静态地址
	Router.StaticFS("/static", http.Dir("/static"))
	address := fmt.Sprintf(":%d", global.Yaml["server.port"])
	// 启动HTTP服务,courseController
	s := initServer(address, Router)
	global.Log.Debug("服务启动成功：端口是：", zap.String("port", address))
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)

	s2 := s.ListenAndServe().Error()
	global.Log.Info("服务启动完毕", zap.Any("s2", s2))
}

```

8:  请求测试 api/sysmenus.js

```js
import request from '@/request/index.js'


/**
 * 查询菜单列表并分页
 */
export const LoadTreeData = (data)=>{
   return request.post(`/sys/menus/tree?keyword=${data.keyword}`,data)
}

```

9：执行测试

```js
import LoadTreeData from '@/api/sysmenus.js'
const handleLodData = async ()=>{
   const resp = await LoadTreeData()
   console.log(resp)
}
```



# 17、查询角色对应的菜单

1:  根据用户查询角色 api/login/login.go

```go
// 登录的接口处理
func (api *LoginApi) ToLogined(c *gin.Context) {

	//  省略代码.........................
	// 这个时候就判断用户输入密码和数据库的密码是否一致
	// inputPassword = utils.Md5(123456) = 2ec9f77f1cde809e48fabac5ec2b8888
	// dbUser.Password = 2ec9f77f1cde809e48fabac5ec2b8888
	if dbUser != nil && dbUser.Password == adr.Md5Slat(inputPassword, dbUser.Slat) {
		token := api.generaterToken(c, dbUser)
		// 根据用户id查询用户的角色
		userRoles, _ := sysUserRolesService.SelectUserRoles(dbUser.ID)//-----新增
		// 根据用户查询菜单信息
		roleMenus, _ := sysRoleMenusService.SelectRoleMenus(userRoles[0].ID) //----新增
		// 根据用户id查询用户的角色的权限
		permissions, _ := sysRoleApisService.SelectRoleApis(userRoles[0].ID) //----新增
		
		// 查询返回
		response.Ok(map[string]any{"user": dbUser, "token": token, "roles": userRoles, "roleMenus": sysMenuService.Tree(roleMenus, 0), "permissions": permissions}, c)
	} else {
		response.Fail(60002, "你输入的账号和密码有误", c)
	}
}
```

2: 根据用户id查询对应的角色

```go
userRoles, _ := sysUserRolesService.SelectUserRoles(dbUser.ID)//-----新增
```

具体如下：

```go
// 查询用户授权的角色信息
func (service *SysUserRolesService) SelectUserRoles(userId uint) (sysRoles []*sys2.SysRoles, err error) {
	err = global.KSD_DB.Select("t2.*").Table("sys_user_roles t1,sys_roles t2").
		Where("t1.user_id = ? and t1.role_id = t2.id", userId).Scan(&sysRoles).Error
	return sysRoles, err
}

```

3:  切换角色必须要角色对应菜单也进行刷新

默认情况下，我们把用户对应的角色列表查询出来，把第一个作为默认的角色，那么就必须把默认的角色对应菜单和权限全部都全查询返回，然后页面根据服务端返回的菜单和权限进行渲染。

```go

```















# 17、数据加密



- aes go/js
- des go/js
- rsa go/js
- sha1
- md5—-密码加密 、 文件唯一标识

# 18、参数处理



# 19、挤下线

# 自定义的方式实现后端的鉴权处理



## 01、上节课为什么debug不生效

造成的原因是：版本升级带来隐患。

go 1.19.2—-泛型（泛型约束），对你的指定模板类型进行约束.

**什么是泛型约束**

```go
type BaseService[D any, T any] struct{}
```



## 02、解决一个bug问题，关于骨架屏幕不退的问题

把原来的状态管理骨架屏幕的代码移植到App.vue 中

```js
import {useSkeletonStore} from '@/stores/skeleton.js'
const skeletonStore = useSkeletonStore()


onMounted(() => {
  setTimeout(() => {
      skeletonStore.skLoading = false;
  }, 600)
})
```



# 权限分配

核心：所谓权限控制：其实就于未来不同角色可以访问到不同权限（具体就是你在router定义的每个接口的调用访问权限）。

- role 1 — /api/sys/user/load —-A1001—user:load—有记录——又权限可以访问
- role 1 — /api/sys/user/load —-A1001—user:load —无记录——权限不足

## 关于权限API的分配和管理

这部分逻辑和菜单是一模一样的

### 1：创建表

```sql
CREATE TABLE `sys_apis` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `path` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'api路径',
  `description` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'api中文描述',
  `parent_id` bigint(20) unsigned DEFAULT NULL COMMENT '隶属于菜单的api',
  `method` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'POST' COMMENT '方法',
  `is_deleted` bigint(20) unsigned DEFAULT '0',
  `title` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'api路径名称',
  `code` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '权限代号',
  `sort` bigint(20) DEFAULT NULL COMMENT '排序标记',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `path` (`path`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='权限表';
```

### 2： 结构体

```go
package sys

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type SysApis struct {
	ID          uint                  `gorm:"primarykey;comment:主键ID" json:"id" structs:"-"` // 主键ID
	CreatedAt   time.Time             `gorm:"type:datetime(0);autoCreateTime;comment:创建时间" json:"createdAt" structs:"-"`
	UpdatedAt   time.Time             `gorm:"type:datetime(0);autoUpdateTime;comment:更新时间" json:"updatedAt" structs:"-"`
	IsDeleted   soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt;default:0" json:"isDeleted" structs:"is_deleted"`
	Title       string                `json:"title" gorm:"comment:api路径名称"`          // api路径
	Path        string                `json:"path" gorm:"comment:api路径"`             // api路径
	Description string                `json:"description" gorm:"comment:api中文描述"`    // api中文描述
	ParentId    uint                  `json:"parentId" gorm:"comment:隶属于菜单的api"`     // api组
	Method      string                `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	Code        string                `json:"code" gorm:"comment:权限代号"`              // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	// 忽略该字段，- 表示无读写，-:migration 表示无迁移权限，-:all 表示无读写迁移权限
	Children []*SysApis `gorm:"-" json:"children"`
}

func (s *SysApis) TableName() string {
	return "sys_apis"
}

```

### 3: service

```go
package sys

import (
	"xkginweb/global"
	"xkginweb/model/entity/sys"
	"xkginweb/service/commons"
)

// 对用户表的数据层处理
type SysApisService struct {
	commons.BaseService[uint, sys.SysApis]
}

// 添加
func (service *SysApisService) SaveSysApis(sysApis *sys.SysApis) (err error) {
	err = global.KSD_DB.Create(sysApis).Error
	return err
}

// 修改
func (service *SysApisService) UpdateSysApis(sysApis *sys.SysApis) (err error) {
	err = global.KSD_DB.Unscoped().Model(sysApis).Updates(sysApis).Error
	return err
}

// 按照map的方式更新
func (service *SysApisService) UpdateSysApisMap(sysApis *sys.SysApis, mapFileds *map[string]any) (err error) {
	err = global.KSD_DB.Unscoped().Model(sysApis).Updates(mapFileds).Error
	return err
}

// 删除
func (service *SysApisService) DelSysApisById(id uint) (err error) {
	var sysApis sys.SysApis
	err = global.KSD_DB.Where("id = ?", id).Delete(&sysApis).Error
	return err
}

// 批量删除
func (service *SysApisService) DeleteSysApissByIds(sysApiss []sys.SysApis) (err error) {
	err = global.KSD_DB.Delete(&sysApiss).Error
	return err
}

// 根据id查询信息
func (service *SysApisService) GetSysApisByID(id uint) (sysApiss *sys.SysApis, err error) {
	err = global.KSD_DB.Unscoped().Omit("created_at", "updated_at").Where("id = ?", id).First(&sysApiss).Error
	return
}

func (service *SysApisService) FinApiss(keyword string) (sysApis []*sys.SysApis, err error) {
	db := global.KSD_DB.Unscoped().Order("sort asc")
	if len(keyword) > 0 {
		db.Where("title like ?", "%"+keyword+"%")
	}
	err = db.Find(&sysApis).Error
	return sysApis, err
}

/**
*   开始把数据进行编排--递归
*   Tree(all,0)
 */
func (service *SysApisService) Tree(allSysApis []*sys.SysApis, parentId uint) []*sys.SysApis {
	var nodes []*sys.SysApis
	for _, dbApis := range allSysApis {
		if dbApis.ParentId == parentId {
			childrensApis := service.Tree(allSysApis, dbApis.ID)
			if len(childrensApis) > 0 {
				dbApis.Children = append(dbApis.Children, childrensApis...)
			}
			nodes = append(nodes, dbApis)
		}
	}
	return nodes
}

/*
*
数据复制
*/
func (service *SysApisService) CopyData(id uint) (dbData *sys.SysApis, err error) {
	// 2: 查询id数据
	sysApisData, err := service.GetByID(id)
	if err != nil {
		return nil, err
	}
	// 3: 开始复制
	sysApisData.ID = 0
	sysApisData.Path = ""
	sysApisData.Code = ""
	// 4: 保存入库
	data, err := service.Save(sysApisData)

	return data, err
}

```

### 4: 定义接口

```go
package sys

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"xkginweb/commons/response"
	"xkginweb/global"
	"xkginweb/model/entity/sys"
)

type SysApisApi struct {
	global.BaseApi
}

// 拷贝
func (api *SysApisApi) CopyData(c *gin.Context) {
	// 1: 获取id数据 注意定义李媛媛的/:id
	id := c.Param("id")
	data, _ := sysApisService.CopyData(api.StringToUnit(id))
	response.Ok(data, c)
}

// 保存
func (api *SysApisApi) SaveData(c *gin.Context) {
	// 1: 第一件事情就准备数据的载体
	var sysApis sys.SysApis
	err := c.ShouldBindJSON(&sysApis)
	if err != nil {
		// 如果参数注入失败或者出错就返回接口调用这。出错了.
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 创建实例，保存帖子
	err = sysApisService.SaveSysApis(&sysApis)
	// 如果保存失败。就返回创建失败的提升
	if err != nil {
		response.FailWithMessage("创建失败", c)
		return
	}
	// 如果保存成功，就返回创建创建成功
	response.Ok("创建成功", c)
}

// 状态修改
func (api *SysApisApi) UpdateStatus(c *gin.Context) {
	type Params struct {
		Id    uint   `json:"id"`
		Filed string `json:"field"`
		Value any    `json:"value"`
	}
	var params Params
	err := c.ShouldBindJSON(&params)
	if err != nil {
		// 如果参数注入失败或者出错就返回接口调用这。出错了.
		response.FailWithMessage(err.Error(), c)
		return
	}

	flag, _ := sysApisService.UnUpdateStatus(params.Id, params.Filed, params.Value)
	// 如果保存失败。就返回创建失败的提升
	if !flag {
		response.FailWithMessage("更新失败", c)
		return
	}
	// 如果保存成功，就返回创建创建成功
	response.Ok("更新成功", c)
}

// 编辑修改
func (api *SysApisApi) UpdateById(c *gin.Context) {
	// 1: 第一件事情就准备数据的载体
	var sysApis sys.SysApis
	err := c.ShouldBindJSON(&sysApis)
	if err != nil {
		// 如果参数注入失败或者出错就返回接口调用这。出错了.
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 结构体转化成map呢？
	m := structs.Map(sysApis)
	m["is_deleted"] = sysApis.IsDeleted
	err = sysApisService.UpdateSysApisMap(&sysApis, &m)
	// 如果保存失败。就返回创建失败的提升
	if err != nil {
		fmt.Println(err)
		response.FailWithMessage("更新失败", c)
		return
	}
	// 如果保存成功，就返回创建创建成功
	response.Ok("更新成功", c)
}

// 根据id删除
func (api *SysApisApi) DeleteById(c *gin.Context) {
	// 绑定参数用来获取/:id这个方式
	id := c.Param("id")
	// 开始执行
	parseUint, _ := strconv.ParseUint(id, 10, 64)
	err := sysApisService.DelSysApisById(uint(parseUint))
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.Ok("ok", c)
}

// 根据id查询信息
func (api *SysApisApi) GetById(c *gin.Context) {
	// 根据id查询方法
	id := c.Param("id")
	// 根据id查询方法
	parseUint, _ := strconv.ParseUint(id, 10, 64)
	sysUser, err := sysApisService.GetSysApisByID(uint(parseUint))
	if err != nil {
		global.SugarLog.Errorf("查询用户: %s 失败", id)
		response.FailWithMessage("查询用户失败", c)
		return
	}

	response.Ok(sysUser, c)
}

// 批量删除
func (api *SysApisApi) DeleteByIds(c *gin.Context) {
	// 绑定参数用来获取/:id这个方式
	ids := c.Query("ids")
	idstrings := strings.Split(ids, ",")
	var sysApis []sys.SysApis
	for _, id := range idstrings {
		parseUint, _ := strconv.ParseUint(id, 10, 64)
		sysApi := sys.SysApis{}
		sysApi.ID = uint(parseUint)
		sysApis = append(sysApis, sysApi)
	}

	err := sysApisService.DeleteSysApissByIds(sysApis)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.Ok("ok", c)
}

// 查询权限信息
func (api *SysApisApi) FindApisTree(c *gin.Context) {
	keyword := c.Query("keyword")
	sysApis, err := sysApisService.FinApiss(keyword)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.Ok(sysApisService.Tree(sysApis, 0), c)
}

```

### 5: 定义具体的接口路由

```go
package sys

import (
	"github.com/gin-gonic/gin"
	v1 "xkginweb/api/v1"
)

// 登录路由
type SysApisRouter struct{}

func (r *SysApisRouter) InitSysApisRouter(Router *gin.RouterGroup) {
	sysApisApi := v1.WebApiGroupApp.Sys.SysApisApi
	// 用组定义--（推荐）
	router := Router.Group("/sys")
	{
		// 获取菜单列表
		router.POST("/apis/tree", sysApisApi.FindApisTree)
		// 保存
		router.POST("/apis/save", sysApisApi.SaveData)
		// 复制数据
		router.POST("/apis/copy/:id", sysApisApi.CopyData)
		// 修改
		router.POST("/apis/update", sysApisApi.UpdateById)
		// 启用和未启用 （控制启用，发布，删除）
		router.POST("/apis/update/status", sysApisApi.UpdateStatus)
		// 删除单个 :id 获取参数的时候id := c.Param("id")，传递的时候/sys/user/del/100
		router.POST("/apis/del/:id", sysApisApi.DeleteById)
		// 删除多个  获取参数的时候ids := c.Query("ids")，传递的时候/sys/user/dels?ids=1,2,3,4
		router.POST("/apis/dels", sysApisApi.DeleteByIds)
		// 查询明细 /user/get/1/xxx
		router.POST("/apis/get/:id", sysApisApi.GetById)
	}
}

```

### 6： 然后在 [initilization](C:\Users\zxc\go\xkginweb\initilization)的init-router.go进行注册

```go
package initilization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"xkginweb/commons/filter"
	"xkginweb/commons/middle"
	"xkginweb/global"
	"xkginweb/router"
)

func InitGinRouter() *gin.Engine {
	// 打印gin的时候日志是否用颜色标出
	//gin.ForceConsoleColor()
	//gin.DisableConsoleColor()
	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 创建gin服务
	ginServer := gin.Default()
	// 提供服务组
	courseRouter := router.RouterWebGroupApp.Course.CourseRouter

	videoRouter := router.RouterWebGroupApp.Video.XkVideoRouter

	userStateRouter := router.RouterWebGroupApp.State.UserStateRouter

	bbsRouter := router.RouterWebGroupApp.BBs.XkBbsRouter
	bbsCategoryRouter := router.RouterWebGroupApp.BBs.BBSCategoryRouter

	loginRouter := router.RouterWebGroupApp.Login.LoginRouter
	logoutRouter := router.RouterWebGroupApp.Login.LogoutRouter
	codeRouter := router.RouterWebGroupApp.Code.CodeRouter

	sysMenusRouter := router.RouterWebGroupApp.Sys.SysMenusRouter
	sysApisRouter := router.RouterWebGroupApp.Sys.SysApisRouter // -----------------新增代码
	sysUserRouter := router.RouterWebGroupApp.Sys.SysUsersRouter
	sysRolesRouter := router.RouterWebGroupApp.Sys.SysRolesRouter
	sysUserRolesRouter := router.RouterWebGroupApp.Sys.SysUserRolesRouter
	sysRoleMenusRouter := router.RouterWebGroupApp.Sys.SysRoleMenusRouter
	sysRoleApisRouter := router.RouterWebGroupApp.Sys.SysRoleApisRouter

	// 解决接口的跨域问题
	ginServer.Use(filter.Cors())
	// 接口隔离，比如登录，健康检查都不需要拦截和做任何的处理
	// 业务模块接口，
	privateGroup := ginServer.Group("/api")
	// 无需jwt拦截
	{
		loginRouter.InitLoginRouter(privateGroup)
		codeRouter.InitCodeRouter(privateGroup)
	}
	// 会被jwt拦截
	privateGroup.Use(middle.JWTAuth()).Use(middle.RBAC())
	{
		logoutRouter.InitLogoutRouter(privateGroup)
		videoRouter.InitXkVideoRouter(privateGroup)
		courseRouter.InitCourseRouter(privateGroup)
		userStateRouter.InitUserStateRouter(privateGroup)
		bbsRouter.InitXkBbsRouter(privateGroup)
		bbsCategoryRouter.InitBBSCategoryRouter(privateGroup)
		sysMenusRouter.InitSysMenusRouter(privateGroup)
		sysUserRouter.InitSysUsersRouter(privateGroup)
		sysRolesRouter.InitSysRoleRouter(privateGroup)
		sysApisRouter.InitSysApisRouter(privateGroup)// ---------------------------------新增代码
		sysUserRolesRouter.InitSysUserRolesRouter(privateGroup)
		sysRoleMenusRouter.InitSysRoleMenusRouter(privateGroup)
		sysRoleApisRouter.InitSysRoleApisRouter(privateGroup)
	}

	fmt.Println("router register success")
	return ginServer
}

func RunServer() {

	// 初始化路由
	Router := InitGinRouter()
	// 为用户头像和文件提供静态地址
	Router.StaticFS("/static", http.Dir("/static"))
	address := fmt.Sprintf(":%d", global.Yaml["server.port"])
	// 启动HTTP服务,courseController
	s := initServer(address, Router)
	global.Log.Debug("服务启动成功：端口是：", zap.String("port", address))
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)

	s2 := s.ListenAndServe().Error()
	global.Log.Info("服务启动完毕", zap.Any("s2", s2))
}

```

### 7： 前端定义api接口

```go
import request from '@/request/index.js'
import { C2B  } from '../utils/wordtransfer'

/**
 * 查询权限列表并分页
 */
export const LoadTreeData = (data)=>{
   return request.post(`/sys/apis/tree`,data)
}

/**
 * 根据id查询权限信息
 */
export const GetById = ( id )=>{
   return request.post(`/sys/apis/get/${id}`)
}

/**
 * 保存权限
 */
export const SaveData = ( data )=>{
   return request.post(`/sys/apis/save`,data)
}

/**
 * 更新权限信息
 */
export const UpdateData = ( data )=>{
   return request.post(`/sys/apis/update`,data)
}


/**
 * 根据id删除权限信息
 */
export const DelById = ( id )=>{
   return request.post(`/sys/apis/del/${id}`)
}

/**
 * 根据ids批量删除权限信息
 */
export const DelByIds = ( ids )=>{
   return request.post(`/sys/apis/dels?ids=${ids}`)
}

/**
 * 权限启用和未启用
 */
export const UpdateStatus = ( data )=>{
   data.field = C2B(data.field)
   return request.post(`/sys/apis/update/status`,data)
}


/**
 * 复制数据
 */
export const CopyData = ( id )=>{
   return request.post(`/sys/apis/copy/${id}`,{})
}

```

### 8: 开始对接页面 views/sys/Permission.vue

```vue
<template>
  <div class="kva-container">
    <div class="kva-contentbox">
      <home-page-header>
        <div class="kva-form-search">
          <el-form :inline="true" ref="searchForm" :model="queryParams">
            <el-form-item>
              <el-button type="primary"  icon="Plus" @click="handleAdd">添加权限</el-button>
            </el-form-item>
            <el-form-item label="关键词：">
              <el-input v-model="queryParams.keyword" placeholder="请输入菜单名称..." maxlength="10" clearable />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="Search" @click.prevent="handleSearch">搜索</el-button>
              <el-button type="danger" icon="Refresh" @click.prevent="handleReset">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
        <!-- default-expand-all -->
        <el-table
          :data="tableData"
          style="width: 100%; margin-bottom: 20px"
          row-key="id"
          border
          stripe
          :height="settings.tableHeight()"
        >
          <el-table-column fixed prop="id" label="ID" align="center" width="80"  />
          <el-table-column fixed prop="parentId" label="父ID" align="center" width="80" />
          <el-table-column prop="title" label="展示名字" align="center" >
            <template #default="{row}">
                <el-input v-model="row.title" style="text-align:center" @change="handleChange(row,'title')"></el-input>
            </template>
          </el-table-column>
          <el-table-column prop="code" label="编号" align="center" >
            <template #default="{row}">
                <el-input v-model="row.code" style="text-align:center" @change="handleChange(row,'code')"></el-input>
            </template>
          </el-table-column>
          <el-table-column prop="code" label="访问路径" align="center" >
            <template #default="{row}">
                <el-input v-model="row.path" style="text-align:center" @change="handleChange(row,'path')"></el-input>
            </template>
          </el-table-column>
          <el-table-column prop="sort" label="排序"  align="center" width="180">
            <template #default="{row}">
                <el-input-number v-model="row.sort" @change="handleChange(row,'sort')"></el-input-number>
            </template>
          </el-table-column>
          <el-table-column label="是否删除" align="center" width="180">
            <template #default="{row}">
              <el-switch 
              v-model="row.isDeleted" 
              @change="handleChange(row,'isDeleted')" 
              active-color="#ff0000"
               active-text="已删除" 
               inactive-text="未删除" 
               :active-value="1" 
               :inactive-value="0"/>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" align="center" width="160">
            <template #default="scope">
              {{ formatTimeToStr(scope.row.createdAt,"yyyy/MM/dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="更新时间" align="center" width="160">
            <template #default="scope">
              {{ formatTimeToStr(scope.row.updatedAt,"yyyy/MM/dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column fixed="right" align="center" label="操作" width="350">
            <template #default="{row,$index}">
              <el-button text icon="edit" @click="handleEdit(row)"  type="primary">编辑</el-button>
              <el-button text icon="Tickets" @click="handleCopy(row)"  type="success">复制</el-button>
              <el-button text icon="remove" @click="handleRemove(row)"  type="danger">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </home-page-header>
    </div>
    <!--添加和修改菜单-->
    <add-sys-apis ref="addRef"  @load="handleLoadData"></add-sys-apis>
  </div>
</template>

<script  setup>
import { C2B,B2C } from '@/utils/wordtransfer'
import KVA from '@/utils/kva.js'
import settings from '@/settings';
import { formatTimeToStr } from '@/utils/date'
import AddSysApis from '@/views/sys/components/AddSysApis.vue'
import { LoadTreeData,UpdateStatus,CopyData,DelById } from '@/api/sysapis.js';
import { reactive } from 'vue';
import { useUserStore } from '@/stores/user.js'
const userStore = useUserStore()
const addRef = ref(null);

// 搜索属性定义
let queryParams = reactive({
  keyword:""
})

// 数据容器
const tableData = ref([]) 
const searchForm = ref(null)
// 搜索
const handleSearch = ()=> {
  handleLoadData()
}

// 查询列表
const handleLoadData = async ()=>{
  const resp = await LoadTreeData(queryParams)
  tableData.value = resp.data
}

// 添加
const handleAdd = ()=>{
  addRef.value.handleOpen('','save',tableData.value?.length)
}

// 编辑
const handleEdit =  async (row) => {
  // 在打开,再查询，
  addRef.value.handleOpen(row.id,'edit',tableData.value?.length)
}

// 添加子菜单
const handleAddChild = (row) => {
  addRef.value.handleOpen(row,'child',row.children?.length)
}

// 改变序号 sorted,标题 title、启用 status,isDeleted
const handleChange = async (row,field) =>{
  var value = row[field];//row.isDeleted=0 
  var params = {id:row.id,field:field,value:value};
  await UpdateStatus(params); 
  KVA.notifySuccess("更新成功")
  if(field=="sort"){
    tableData.value.sort((a,b)=>a.sort-b.sort);
  }
}



// 物理删除
const handleRemove =  async (row) => {
  try{
    await KVA.confirm("警告","你确定要抛弃我么？",{icon:"error"})
    await DelById(row.id)
    KVA.notifySuccess("操作成功")
    userStore.handlePianaRole(0,"")
    handleLoadData()
  }catch(e){
    KVA.notifyError("操作失败")
  }
}

// 重置搜索表单
const handleReset = () => {
  queryParams.keyword = ""
  searchForm.value.resetFields()
  handleLoadData()
}



// 复制
const handleCopy = async (row) => {
  await CopyData(row.id);
  KVA.notifySuccess("复制成功")
  handleLoadData()
}

// 生命周期加载
onMounted(()=>{
  handleLoadData()

  console.log("C2B",C2B("isDeletedNum"))
  console.log("B2C",B2C("is_deleted_num"))
})


</script>

```

具体后续看视频就不展开了



# 关于cashbin的后端鉴权处理

参考文档：

- https://www.jb51.net/article/213556.htm
- https://blog.csdn.net/baidu_32452525/article/details/118199304
- https://github.com/casbin/casbin
- https://github.com/patrickmn/go-cache
- https://www.jianshu.com/p/b5e0c5fcaa2a
- https://blog.csdn.net/lanyanleio/article/details/127516463
- https://blog.csdn.net/qq_42120178/article/details/117156766（推荐）
- https://casbin.org/editor/ （推荐）



## 01、安装

```
go get github.com/casbin/casbin/v2
go get github.com/casbin/gorm-adapter/v3
```

Gorm Adapter
----

> In v3.0.3, method `NewAdapterByDB` creates table named `casbin_rules`,  
> we fix it to `casbin_rule` after that.  
> If you used v3.0.3 and less, and you want to update it,  
> you might need to *migrate* data manually.
> Find out more at: https://github.com/casbin/gorm-adapter/issues/78

Gorm Adapter is the [Gorm](https://gorm.io/gorm) adapter for [Casbin](https://github.com/casbin/casbin). With this library, Casbin can load policy from Gorm supported database or save policy to it.

Based on [Officially Supported Databases](https://v1.gorm.io/docs/connecting_to_the_database.html#Supported-Databases), The current supported databases are:

- MySQL
- PostgreSQL
- SQL Server
- Sqlite3

> gorm-adapter use ``github.com/glebarez/sqlite`` instead of gorm official sqlite driver ``gorm.io/driver/sqlite`` because the latter needs ``cgo`` support. But there is almost no difference between the two driver. If there is a difference in use, please submit an issue.

- other 3rd-party supported DBs in Gorm website or other places.

## Installation

    go get github.com/casbin/gorm-adapter/v3

## Simple Example

```go
package main

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize a Gorm adapter and use it in a Casbin enforcer:
	// The adapter will use the MySQL database named "casbin".
	// If it doesn't exist, the adapter will create it automatically.
	// You can also use an already existing gorm instance with gormadapter.NewAdapterByDB(gormInstance)
	a, _ := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/") // Your driver and data source.
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", a)
	
	// Or you can use an existing DB "abc" like this:
	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	// a := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/abc", true)

	// Load the policy from DB.
	e.LoadPolicy()
	
	// Check the permission.
	e.Enforce("alice", "data1", "read")
	
	// Modify the policy.
	// e.AddPolicy(...)
	// e.RemovePolicy(...)
	
	// Save the policy back to DB.
	e.SavePolicy()
}
```

## Turn off AutoMigrate

New an adapter will use ``AutoMigrate`` by default for create table, if you want to turn it off, please use API ``TurnOffAutoMigrate(db *gorm.DB) *gorm.DB``. See example: 

```go
db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/casbin"), &gorm.Config{})
TurnOffAutoMigrate(db)
// a,_ := NewAdapterByDB(...)
// a,_ := NewAdapterByDBUseTableName(...)
a,_ := NewAdapterByDBWithCustomTable(...)
```

Find out more details at [gorm-adapter#162](https://github.com/casbin/gorm-adapter/issues/162)

## Customize table columns example

You can change the gorm struct tags, but the table structure must stay the same.

```go
package main

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func main() {
	// Increase the column size to 512.
	type CasbinRule struct {
		ID    uint   `gorm:"primaryKey;autoIncrement"`
		Ptype string `gorm:"size:512;uniqueIndex:unique_index"`
		V0    string `gorm:"size:512;uniqueIndex:unique_index"`
		V1    string `gorm:"size:512;uniqueIndex:unique_index"`
		V2    string `gorm:"size:512;uniqueIndex:unique_index"`
		V3    string `gorm:"size:512;uniqueIndex:unique_index"`
		V4    string `gorm:"size:512;uniqueIndex:unique_index"`
		V5    string `gorm:"size:512;uniqueIndex:unique_index"`
	}

	db, _ := gorm.Open(...)

	// Initialize a Gorm adapter and use it in a Casbin enforcer:
	// The adapter will use an existing gorm.DB instnace.
	a, _ := gormadapter.NewAdapterByDBWithCustomTable(db, &CasbinRule{}) 
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", a)
	
	// Load the policy from DB.
	e.LoadPolicy()
	
	// Check the permission.
	e.Enforce("alice", "data1", "read")
	
	// Modify the policy.
	// e.AddPolicy(...)
	// e.RemovePolicy(...)
	
	// Save the policy back to DB.
	e.SavePolicy()
}
```

## Transaction

You can modify policies within a transaction.See example:

```go
package main

func main() {
	a, err := NewAdapterByDB(db)
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", a)
	err = e.GetAdapter().(*Adapter).Transaction(e, func(e casbin.IEnforcer) error {
		_, err := e.AddPolicy("jack", "data1", "write")
		if err != nil {
			return err
		}
		_, err = e.AddPolicy("jack", "data2", "write")
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// handle if transaction failed
		return
	}
}
```

## ConditionsToGormQuery

`ConditionsToGormQuery()` is a function that converts multiple query conditions into a GORM query statement
You can use the `GetAllowedObjectConditions()` API of Casbin to get conditions,
and choose the way of combining conditions through `combineType`.

`ConditionsToGormQuery()` allows Casbin to be combined with SQL, and you can use it to implement many functions.

### Example: GetAllowedRecordsForUser

* model example: [object_conditions_model.conf](examples/object_conditions_model.conf)
* policy example: [object_conditions_policy.csv](examples/object_conditions_policy.csv)

DataBase example:

| id   | title | author  | publisher  | publish_data        | price | category_id |
| ---- | ----- | ------- | ---------- | ------------------- | ----- | ----------- |
| 1    | book1 | author1 | publisher1 | 2023-04-09 16:23:42 | 10    | 1           |
| 2    | book2 | author1 | publisher1 | 2023-04-09 16:23:44 | 20    | 2           |
| 3    | book3 | author2 | publisher1 | 2023-04-09 16:23:44 | 30    | 1           |
| 4    | book4 | author2 | publisher2 | 2023-04-09 16:23:45 | 10    | 3           |
| 5    | book5 | author3 | publisher2 | 2023-04-09 16:23:45 | 50    | 1           |
| 6    | book6 | author3 | publisher2 | 2023-04-09 16:23:46 | 60    | 2           |


```go
type Book struct {
    ID          int
    Title       string
    Author      string
    Publisher   string
    PublishDate time.Time
    Price       float64
    CategoryID  int
}

func TestGetAllowedRecordsForUser(t *testing.T) {
	e, _ := casbin.NewEnforcer("examples/object_conditions_model.conf", "examples/object_conditions_policy.csv")

	conditions, err := e.GetAllowedObjectConditions("alice", "read", "r.obj.")
	if err != nil {
		panic(err)
	}
	fmt.Println(conditions)

	dsn := "root:root@tcp(127.0.0.1:3307)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("CombineTypeOr")
	rows, err := ConditionsToGormQuery(db, conditions, CombineTypeOr).Model(&Book{}).Rows()
	defer rows.Close()
	var b Book
	for rows.Next() {
		err := db.ScanRows(rows, &b)
		if err != nil {
			panic(err)
		}
		log.Println(b)
	}

	fmt.Println("CombineTypeAnd")
	rows, err = ConditionsToGormQuery(db, conditions, CombineTypeAnd).Model(&Book{}).Rows()
	defer rows.Close()
	for rows.Next() {
		err := db.ScanRows(rows, &b)
		if err != nil {
			panic(err)
		}
		log.Println(b)
	}
}
```


## Getting Help

- [Casbin](https://github.com/casbin/casbin)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.



# 账号挤下线（验证码）



## 01、需求

同一个账号，是不同地方登录。只能让一个有效（最后一次登录有效），之前全部会自动挤下去。

## 02、实现

1： 用户登录，会根据用户id生成一个唯一登录标识，然后放入到服务器，返回给客户端。

2：然后用户未来每个接口请求的时候，都会携带这个唯一标识和用户id

3：然后把用户携带的id和标识，和服务端存储的唯一标识进行比对。

## 03、本地缓存

go-cache

1：基于内存的 K/V 存储/缓存 : (类似于Memcached)，适用于单机应用程序。（在go的进程内存中，挖来一个空间出来用来存储数据，而这个数据可以让别线程进行数据共享。）

2：缓存数据，如果一直放的话，可能溢出。可能影响主进程的执行。所以大部分的缓存设计都会考虑到：淘汰策略

- Least Recently  Used (LRU)：最近最少使用策略，删除最近最少被使用的缓存项。
- First In First Out (FIFO)：先进先出策略，删除最早被加入到缓存中的缓存项。
- Least Frequently Used (LFU)：最不经常使用策略，删除使用频率最低的缓存项。
- Random Replacement (RR)：随机替换策略，根据一个随机算法选择要删除的缓存项。

3: 底层原理就是：全局Map (安全性，学习锁)



### 03-01、安装

```go
go get github.com/patrickmn/go-cache
```

### 03-02、使用

在global包下的 global.go增加缓存对象

```go
package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"xkginweb/commons/parse"
)

var (
	Cache      *cache.Cache //-------------------新增代码
	Log        *zap.Logger
	SugarLog   *zap.SugaredLogger
	Lock       sync.RWMutex
	Yaml       map[string]interface{}
	Config     *parse.Config
	KSD_DB     *gorm.DB
	BlackCache local_cache.Cache
	REDIS      *redis.Client
)

```

然后初始化global.cache。在 [initilization](..\go\xkginweb\initilization) 下新建一个  [init_cache.go](..\go\xkginweb\initilization\init_cache.go) 文件如下：

```go
package initilization

import (
	"github.com/patrickmn/go-cache"
	"time"
	"xkginweb/global"
)

func InitCache() {
	c := cache.New(5*time.Minute, 24*60*time.Minute)
	global.Cache = c
}

```

然后在main.go初始化global.cache对象如下

```go
package main

import (
	"xkginweb/initilization"
)

func main() {
	// 解析配置文件
	initilization.InitViper()
	// 初始化日志 开发的时候建议设置成：debug ，发布的时候建议设置成：info/error
	// info --- console + file
	// error -- file
	initilization.InitLogger("debug")
	// 初始化中间 redis/mysql/mongodb
	initilization.InitMySQL()
	// 初始化缓存
	initilization.InitRedis()
	// 初始化本地缓存
	initilization.InitCache()//-------------------新增代码
	// 定时器
	// 并发问题解决方案
	// 异步编程
	// 初始化路由
	initilization.RunServer()
}

```

然后找到api/v1/login/login.go 在登录的时候，增加uuid写入缓存和返回的处理

```go
// 登录的接口处理
func (api *LoginApi) ToLogined(c *gin.Context) {
	type LoginParam struct {
		Account  string
		Code     string
		CodeId   string
		Password string
	}

	// 1：获取用户在页面上输入的账号和密码开始和数据库里数据进行校验
	param := LoginParam{}
	err2 := c.ShouldBindJSON(&param)
	if err2 != nil {
		response.Fail(60002, "参数绑定有误", c)
		return
	}

	if len(param.Code) == 0 {
		response.Fail(60002, "请输入验证码", c)
		return
	}

	if len(param.CodeId) == 0 {
		response.Fail(60002, "验证码获取失败", c)
		return
	}

	// 开始校验验证码是否正确
	verify := store.Verify(param.CodeId, param.Code, true)
	if !verify {
		response.Fail(60002, "你输入的验证码有误!!", c)
		return
	}

	inputAccount := param.Account
	inputPassword := param.Password

	if len(inputAccount) == 0 {
		response.Fail(60002, "请输入账号", c)
		return
	}

	if len(inputPassword) == 0 {
		response.Fail(60002, "请输入密码", c)
		return
	}

	dbUser, err := sysUserService.GetUserByAccount(inputAccount)
	if err != nil {
		response.Fail(60002, "你输入的账号和密码有误", c)
		return
	}

	// 这个时候就判断用户输入密码和数据库的密码是否一致
	// inputPassword = utils.Md5(123456) = 2ec9f77f1cde809e48fabac5ec2b8888
	// dbUser.Password = 2ec9f77f1cde809e48fabac5ec2b8888
	if dbUser != nil && dbUser.Password == adr.Md5Slat(inputPassword, dbUser.Slat) {
		// 根据用户id查询用户的角色
		userRoles, _ := sysUserRolesService.SelectUserRoles(dbUser.ID)
		if len(userRoles) > 0 {
			// 用户信息生成token -----把
			token := api.generaterToken(c, userRoles[0].RoleCode, userRoles[0].ID, dbUser)
			// 根据用户查询菜单信息
			roleMenus, _ := sysRoleMenusService.SelectRoleMenus(userRoles[0].ID)
			// 根据用户id查询用户的角色的权限
			permissions, _ := sysRoleApisService.SelectRoleApis(userRoles[0].ID)

			// 这个uuid是用于挤下线使用 ,//--------------------------增加代码
			uuid := utils.GetUUID()
			userIdStr := strconv.FormatUint(uint64(dbUser.ID), 10)
			global.Cache.Set("LocalCache:Login:"+userIdStr, uuid, cache.NoExpiration)

			// 查询返回
			response.Ok(map[string]any{
				"user":        dbUser,
				"token":       token,
				"roles":       userRoles,//-------------------增加代码
				"uuid":        uuid,
				"roleMenus":   sysMenuService.Tree(roleMenus, 0),
				"permissions": permissions}, c)
		} else {
			// 查询返回--
			response.Fail(80001, "你暂无授权信息", c)
		}
	} else {
		response.Fail(60002, "你输入的账号和密码有误", c)
	}
}
```

然后在前端的状态管理中增加uuid管理

```js
import { defineStore } from 'pinia'
import request from '@/request'
import router from '@/router'
import { handleLogout } from '../api/logout.js'
import { ChangeRoleIdMenus } from '../api/sysroles.js'
import {useSkeletonStore} from '@/stores/skeleton.js'
import {useMenuTabStore} from '@/stores/menuTab.js'



export const useUserStore = defineStore('user', {
  // 定义状态
  state: () => ({
    routerLoaded:true,
    // 登录用户
    user: {},
    username: '',
    userId: '',
    // 挤下线使用
    uuid:"",
    // 登录token
    token: '',
    // 当前角色
    currRoleName:"",
    currRoleCode:"",
    currRoleId:0,
    // 获取用户对应的角色列表
    roles:[],
    // 获取角色对应的权限
    permissions:[],
    // 获取角色对应的菜单
    menuTree:[]
  }),

  // 就是一种计算属性的机制，定义的是函数，使用的是属性就相当于computed
  getters:{
    isLogin(state){
      return state.token ? true : false
    },

    roleName(state){
      return state.roles && state.roles.map(r=>r.name).join(",")
    },

    permissionCode(state){
      return state.permissions &&  state.permissions.map(r=>r.code)
    },
    
    permissionPath(state){
      return state.permissions &&  state.permissions.map(r=>r.path)
    }
  },

  // 定义动作
  actions: {

   /* 设置token */ 
   setToken(newtoken){
      this.token = newtoken
   },

   /* 获取token*/
   getToken(){
    return this.token
   },

   // 改变用户角色的时候把对应菜单和权限查询出来，进行覆盖---更改
   async handlePianaRole(roleId,roleName,roleCode){
      if(roleId > 0 && roleId != this.currRoleId){
        this.currRoleId = roleId
        this.currRoleName = roleName;
        this.currRoleCode = roleCode
      }

      // 获取到导航菜单，切换以后直接全部清空掉
      const menuTabStore = useMenuTabStore();
      menuTabStore.clear()
      
      // 请求服务端--根据切换的角色找到角色对应的权限和菜单
      const resp = await ChangeRoleIdMenus({roleId:this.currRoleId})
      // 对应的权限和菜单进行覆盖
      this.permissions = resp.data.permissions
      this.menuTree = resp.data.roleMenus.sort((a,b)=>a.sort-b.sort)
      if(roleId > 0){
        // 激活菜单中的第一个路由
        router.push(this.menuTree[0].path)
      }
   },
   
   /* 登出*/
   async logout (){
      // 执行服务端退出
      await handleLogout()
      // 清除状态信息
      this.token = ''
      this.user = {}
      this.username = ''
      this.userId = ''
      this.uuid = ''
      this.roles = []
      this.permissions = []
      this.menuTree = []
      // 清除自身的本地存储
      localStorage.removeItem("ksd-kva-language")
      localStorage.removeItem("kva-pinia-userstore")
      sessionStorage.removeItem("kva-pinia-skeleton")
      // 把骨架屏的状态恢复到true的状态
      useSkeletonStore().setLoading(true)
      localStorage.removeItem("isWhitelist")
      location.reload()
      // 然后跳转到登录
      router.push({ name: 'Login', replace: true })
  },

  /* 登录*/
  async toLogin(loginUser){
      // 查询用户信息，角色，权限，角色对应菜单
      const resp = await request.post("login/toLogin", loginUser,{noToken:true})
      // 这个会回退，回退登录页
      var { user ,token,roles,permissions,roleMenus,uuid } = resp.data
      // 登录成功以后获取到菜单信息, 这里要调用一
      this.menuTree = roleMenus
      // 把数据放入到状态管理中
      this.user = user
      this.userId = user.id
      this.username = user.username
      this.token = token
      this.uuid = uuid
      this.roles = roles
      this.permissions = permissions
      // 把roles列表中的角色的第一个作为，当前角色
      this.currRoleId = roles && roles.length > 0 ? roles[0].id : 0
      this.currRoleName = roles && roles.length > 0 ? roles[0].roleName : ""
      this.currRoleCode = roles && roles.length > 0 ? roles[0].roleCode : ""

      return Promise.resolve(resp)
    }
  },
  persist: {
    key: 'kva-pinia-userstore',
    storage: localStorage,//sessionStorage
  }
})
```

然后在每次请求接口的时候把uuid携带上即可。修改“request/index.js即可，如下：

```js
// 1: 导入axios异步请求组件
import axios from "axios";
// 2: 把axios请求的配置剥离成一个config/index.js的文件
import axiosConfig from "./config";
// 3: 获取路由对象--原因是：在加载的过程中先加载的pinia所以不能useRouter机制。
import router from '@/router'
// 4: elementplus消息框
import KVA from "@/utils/kva.js";
// 5: 获取登录的token信息
import { useUserStore } from '@/stores/user.js'
// 6: 然后创建一个axios的实例
const request = axios.create({ ...axiosConfig })

// request request请求拦截器
request.interceptors.request.use(
    function(config){
        // 这个判断就是为了那些不需要token接口做开关处理，比如：登录，检测等
        if(!config.noToken){
             // 如果 token为空，说明没有登录。你就要去登录了
            const userStore = useUserStore()
            const isLogin = userStore.isLogin
            if(!isLogin){
                router.push("/login")
                return
            }else{
                // 90b7d374acc5476eb9beabe9373b2640
                // 这里给请求头增加参数.request--header，在服务端可以通过request的header可以获取到对应参数
                // 比如go: c.GetHeader("Authorization")
                // 比如java: request.getHeader("Authorization")
                config.headers.Authorization = userStore.getToken()
                config.headers.KsdUUID = userStore.uuid
            }
        }
        return config;
    },function(error){
        // 判断请求超时
        if ( error.code === "ECONNABORTED" && error.message.indexOf("timeout") !== -1) {
            KVA.notifyError('请求超时');
            // 这里为啥不写return
        }
        return Promise.reject(error);
    }
);

// request response 响应拦截器
request.interceptors.response.use(async (response) => {
    // 在这里应该可以获取服务端传递过来头部信息
    // 开始续期
    if(response.headers["new-authorization"]){
        const userStore = useUserStore()   
        userStore.setToken(response.headers["new-authorization"])  
    }

    // cashbin的权限拦截处理
    if(response.data?.code === 80001){
        KVA.notifyError(response.data.message);
        // 如果你想调整页面，就把下面注释打开
        //router.push("/nopermission")
        return     
    }

    if(response.data?.code === 20000){
        return response.data;
    }else{
        // 所有得服务端得错误提示，全部在这里进行处理
        if (response.data?.message) {
            KVA.notifyError(response.data.message);
        }

        // 包括: 没登录，黑名单，挤下线
        if(response.data.code === 4001 ){
            const userStore = useUserStore()   
            userStore.logout()
            return Promise.reject(response.data); 
        }   

        // 返回接口返回的错误信息
        return Promise.reject(response.data); 
    }
},(err) => {
    if (err && err.response) {
        switch (err.response.status) {
            case 400:
                err.message = "请求错误";
                break;
            case 401:
                err.message = "未授权，请登录";
                break;
            case 403:
                err.message = "拒绝访问";
                break;
            case 404:
                err.message = `请求地址出错: ${err.response.config.url}`;
                break;
            case 408:
                err.message = "请求超时";
                break;
            case 500:
                err.message = "服务器内部错误";
                break;
            case 501:
                err.message = "服务未实现";
                break;
            case 502:
                err.message = "网关错误";
                break;
            case 503:
                err.message = "服务不可用";
                break;
            case 504:
                err.message = "网关超时";
                break;
            case 505:
                err.message = "HTTP版本不受支持";
                break;
            default:
        }
    }
    if (err.message) {
        KVA.notifyError(err.message);
    }
     // 判断请求超时
    if ( err.code === "ECONNABORTED" && err.message.indexOf("timeout") !== -1) {
        KVA.notifyError('服务已经离开地球表面，刷新或者重试...');
    }
    // 返回接口返回的错误信息
    return Promise.reject(err); 
})
  
export default request
```

然后在jwt.go中进行比较即可

```go
// 1: 导入axios异步请求组件
import axios from "axios";
// 2: 把axios请求的配置剥离成一个config/index.js的文件
import axiosConfig from "./config";
// 3: 获取路由对象--原因是：在加载的过程中先加载的pinia所以不能useRouter机制。
import router from '@/router'
// 4: elementplus消息框
import KVA from "@/utils/kva.js";
// 5: 获取登录的token信息
import { useUserStore } from '@/stores/user.js'
// 6: 然后创建一个axios的实例
const request = axios.create({ ...axiosConfig })

// request request请求拦截器
request.interceptors.request.use(
    function(config){
        // 这个判断就是为了那些不需要token接口做开关处理，比如：登录，检测等
        if(!config.noToken){
             // 如果 token为空，说明没有登录。你就要去登录了
            const userStore = useUserStore()
            const isLogin = userStore.isLogin
            if(!isLogin){
                router.push("/login")
                return
            }else{
                // 90b7d374acc5476eb9beabe9373b2640
                // 这里给请求头增加参数.request--header，在服务端可以通过request的header可以获取到对应参数
                // 比如go: c.GetHeader("Authorization")
                // 比如java: request.getHeader("Authorization")
                config.headers.Authorization = userStore.getToken()
                config.headers.KsdUUID = userStore.uuid
            }
        }
        return config;
    },function(error){
        // 判断请求超时
        if ( error.code === "ECONNABORTED" && error.message.indexOf("timeout") !== -1) {
            KVA.notifyError('请求超时');
            // 这里为啥不写return
        }
        return Promise.reject(error);
    }
);

// request response 响应拦截器
request.interceptors.response.use(async (response) => {
    // 在这里应该可以获取服务端传递过来头部信息
    // 开始续期
    if(response.headers["new-authorization"]){
        const userStore = useUserStore()   
        userStore.setToken(response.headers["new-authorization"])  
    }

     // 包括: 没登录，黑名单，挤下线
     if(response.data.code === 4001 ){
        KVA.notifyError(response.data.message);
        const userStore = useUserStore()   
        userStore.logout() 
        return
    }   


    // cashbin的权限拦截处理
    if(response.data?.code === 80001){
        KVA.notifyError(response.data.message);
        // 如果你想调整页面，就把下面注释打开
        //router.push("/nopermission")
        return response.data;
    }

    if(response.data?.code === 20000){
        return response.data;
    }else{
        // 所有得服务端得错误提示，全部在这里进行处理
        if (response.data?.message) {
            KVA.notifyError(response.data.message);
        }
        // 返回接口返回的错误信息
        return Promise.reject(response.data); 
    }
},(err) => {
    if (err && err.response) {
        switch (err.response.status) {
            case 400:
                err.message = "请求错误";
                break;
            case 401:
                err.message = "未授权，请登录";
                break;
            case 403:
                err.message = "拒绝访问";
                break;
            case 404:
                err.message = `请求地址出错: ${err.response.config.url}`;
                break;
            case 408:
                err.message = "请求超时";
                break;
            case 500:
                err.message = "服务器内部错误";
                break;
            case 501:
                err.message = "服务未实现";
                break;
            case 502:
                err.message = "网关错误";
                break;
            case 503:
                err.message = "服务不可用";
                break;
            case 504:
                err.message = "网关超时";
                break;
            case 505:
                err.message = "HTTP版本不受支持";
                break;
            default:
        }
    }
    if (err.message) {
        KVA.notifyError(err.message);
    }
     // 判断请求超时
    if ( err.code === "ECONNABORTED" && err.message.indexOf("timeout") !== -1) {
        KVA.notifyError('服务已经离开地球表面，刷新或者重试...');
    }
    // 返回接口返回的错误信息
    return Promise.reject(err); 
})
  
export default request
```

## 退出的时候记得清楚缓存

api/v1/login/logout.go如下：

```go
package login

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xkginweb/commons/jwtgo"
	"xkginweb/commons/response"
	"xkginweb/global"
	"xkginweb/model/entity/jwt"
)

// 登录业务
type LogOutApi struct{}

// 退出接口
func (api *LogOutApi) ToLogout(c *gin.Context) {
	// 获取头部的token信息
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Fail(401, "请求未携带token，无权限访问", c)
		return
	}

	// 同时删除缓存中的uuid的信息
	customClaims, _ := jwtgo.GetClaims(c)
	userIdStr := strconv.FormatUint(uint64(customClaims.UserId), 10)
	global.Cache.Delete("LocalCache:Login:" + userIdStr)

	// 退出的token,加入到黑名单中
	err := jwtService.JsonInBlacklist(jwt.JwtBlacklist{Jwt: token})
	// 保存失败会进到到错误
	if err != nil {
		response.Fail(401, "token作废失败", c)
		return
	}

	// 如果保存到黑名单中说明,已经可以告知前端可以进行执行清理动作了
	response.Ok("token作废成功!", c)
}

```





# 指定用户下线

直接把用户的状态清空即可。



# casbin





# Go的项目发布和部署



## 01、准备工作

- 准备一台阿里云服务器
- 开放8080的安全组

## 02、新建一个web工程

### 安装

要安装 Gin 软件包，需要先安装 Go 并设置 Go 工作区。

1.下载并安装 gin：

```sh
go get -u github.com/gin-gonic/gin
```

2.将 gin 引入到代码中：

```go
import "github.com/gin-gonic/gin"
```

### 代码

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	engine := gin.Default()
	engine.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{"code": 200, "msg": "server is success"})
	})
	engine.Run(":8080")
}

```

执行

```
go run main.go
```

### 开始编译

mac  电脑执行

```go
#编译 Linux 64位可执行程序：
GOOS=linux GOARCH=amd64 go build main.go
GOOS=linux GOARCH=arm64 go build main.go
#编译Windows  64位可执行程序：
GOOS=windows GOARCH=amd64 go build main.go
GOOS=windows GOARCH=arm64 go build main.go
#编译 MacOS 64位可执行程序
GOOS=darwin GOARCH=amd64 go build main.go
GOOS=darwin GOARCH=arm64 go build main.go
```

windows执行如下：如果报错执行下面

```go
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows
```

```go
go build main.go
```

Linux

```sh
chmod +x main
./main
#或者
nohup ./main &
```

windows

直接双击打开main.exe即可。



# Go的项目发布和部署

## 01、准备工作

- 云服务器
  - 阿里云，腾讯云，华为云
  - 自己购买服务器和（电信，移动，联通）进行服务器托管
  - 内网穿透 （花生壳）

- 准备一个项目
  - go
  - java

- 安全组
  - 默认情况下，服务器都有防火墙。默认大部分服务在不设置的情况都是只能允许本机访问。但是你现在要暴露服务给外部都调用。




# 02、购买云服务—阿里云

## 01、注册和登录

- https://www.aliyun.com/?spm=5176.12901015-2.0.0.37404b84FqcBlg
- 登录注册地址：https://account.aliyun.com/login/login.htm

## 02、购买服务器

![image-20230903201713828](images/image-20230903201713828.png)

![image-20230903201746238](images/image-20230903201746238.png)



![image-20230903202310333](images/image-20230903202310333.png)

![image-20230903202328980](images/image-20230903202328980.png)

下单完成以后。就进入到控制面板。你可以看到服务实例如下：

![image-20230903202459312](images/image-20230903202459312.png)





## 03、获取到公网IP

![image-20230903202540553](images/image-20230903202540553.png)

目的就是为了获取这个公网IP：116.62.169.76





# 03、如何连接云服务器呢？

## 03-01、进入云服务器

使用客户端工具，xshell/finalshell/golang的提供ssh的链接工具

## 03-02、finalshell链接云服务器

Windows版下载地址:
http://www.hostbuf.com/downloads/finalshell_install.exe

macOS版下载地址:
http://www.hostbuf.com/downloads/finalshell_install.pkg

### windows的finallshell如何连接云服务器呢？

![image-20230903203240581](images/image-20230903203240581.png)



![image-20230903203316020](images/image-20230903203316020.png)

![image-20230903203509648](images/image-20230903203509648.png)

然后就进入到服务器如下：

![image-20230903203527952](images/image-20230903203527952.png)

然后执行系统组件的更新。

```sh
yum update
```



## 03-03、golang开发工具如何连接远程服务器

![image-20230903204149771](images/image-20230903204149771.png)

然后找到菜单栏中：【tools】—【start ssh session】的选项，点击即可：

![image-20230903204258094](images/image-20230903204258094.png)

然后会出现下面的控制台：

![image-20230903204526976](images/image-20230903204526976.png)

这样就可以和xshell和finalshell一样可以操作远程服务器了。





# 04、如何把本地项目和或者上传到云服务器呢？

## 1：lrzsz

```
yum install lrzsz
```

rz ：是负责把本机系统重文件传递到远程服务器

sz ：是负责本远程服务器文件下载到本机系统

## 2： FileZilla

https://2t6y.mydown.com/tianji/child/f690.html?sfrom=166&DTS=1&keyID=123952

## 03、finallshell

直接把文件拖拽到服务器即可。可以完成后续大部分本机系统文件上传到远程服务器的操作和需求。





# 05、开放安全组接口

大家与没发现、我们可以通过xshell/finallshell/golang提供ssh工具都能够很快速的链接到服务远程服务器。然后进行相关的环境安装和文件的传输工作。

在前面我们提过我们购买服务器有一个防火墙，默认是打开的。那你为什么可以直接连接上呢？

- 原因很简单：因为你购买的云服务器。默认情况，云服务器在它安全组下面已经把ssh的协议默认端口22已经打开了。你才可以进行远程连接和传递文件。



## 如何配置和开放我们服务器安全组呢？

![image-20230903205848651](images/image-20230903205848651.png)

![image-20230903210244791](images/image-20230903210244791.png)

注意：未来你项目中所以的服务(gin/mysql/redis/nginx/kafka)都会有一个端口（8080、3306，6937，80，9000） ，这些服务器未来都把它安装和运行到你远程服务器。但是如果你服务器中这些服务（gin/mysql/redis/nginx/kafka)）要对外可以访问到。就必须在远程服务器的安全组中去把这些服务对应的端口进行配置。代表这些服务端口可以被外部进行访问和交互。





# 06、简单工程的发布和部署



## 06-01、新建一个web工程

![image-20230903210703354](images/image-20230903210703354-1693746423995-1.png)

![image-20230903210723542](images/image-20230903210723542.png)

然后点击创建即可。

然后新建一个一个gomodule文件

![image-20230903210750312](images/image-20230903210750312.png)

然后在file进行项目工程的环境设置

![image-20230903210817974](images/image-20230903210817974.png)

环境的值是：`GOPROXY=https://goproxy.io,direct`  ,然后点击应用和ok即可。



## 06-02、安装

要安装 Gin 软件包，需要先安装 Go 并设置 Go 工作区。

1.下载并安装 gin：

```sh
go get -u github.com/gin-gonic/gin
```

2.将 gin 引入到代码中：

```go
import "github.com/gin-gonic/gin"
```

### 代码

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	// 创建一个gin的web服务
	engine := gin.Default()
	// 开始定义个路由
	engine.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{"code": 200, "msg": "server is success"})
	})
	//启动服务端口是：8080 .注意：8080前面有一个冒号:
	engine.Run(":8080")
}

```

执行

```
go run main.go
```

## 06-03、项目开始编译

### mac/linux  系统执行

```go
#编译 Linux 64位可执行程序：
GOOS=linux GOARCH=amd64 go build main.go
GOOS=linux GOARCH=arm64 go build main.go
#编译Windows  64位可执行程序：
GOOS=windows GOARCH=amd64 go build main.go
GOOS=windows GOARCH=arm64 go build main.go
#编译 MacOS 64位可执行程序
GOOS=darwin GOARCH=amd64 go build main.go
GOOS=darwin GOARCH=arm64 go build main.go
```

### windows打包成windows的服务，如下

```go
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows
go build main.go
```

打包成功以后。会在你工程目录下成功一个main.exe文件，就说明打包成功了

### Windows打包Linux服务，如下

```go
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build main.go
```

#### 为什么要指定set GOOS=linux，

- 因为你开发的时候是在windows系统，打包也是在windows中。所以如果你直接去执行 go build main.go的时候，只会生成windows可- 执行文件。也就只会生成main.exe。exe 是windows的可执行文件。 linux根本就不认识。所以比如果要把项目部署到linux服务器上的。其实本应该你的操作是在必须在你远程服务器工作go环境。然后把项目上传到远程服务器，然后在执行 go build main.go就生成linux的可执行文件main (无后缀)。如果这样操作那非常的繁琐
- 因为go项目发布到远程服务器如果是linux系统的话，其实是不需要安装go的环境页可以运行项目。所以我们能够在windows直接能够完成linux包的构建。那么我们是不是就可以省去在远程服务器安装go环境，然后在打包的过程。



### 打包以后的可执行文件认识

-  [main.exe](C:\Users\zxc\go\ginweb\main.exe)  : 这个是windows下的可执行文件

-  main： 这个是linux下的可执行文件 

可执行文件：就文件里提供一个程序入口。并且已经把当前项目工程下的所有.go结尾的文件全部编译成我们操作系统能够运行和识别的文件。然后压缩到这个文件中。你只要启动这个文件，那么你开放系统就可以运行起来了。

##  06-04、指定名字打包

```sh
go build -o server
```

## 06-05、linux 发布和部署服务

### 1： 新建一个/www

```sh
mkdir -p /www/ginweb
cd /www/ginweb
```

- mkdir -p /www/ginweb  递归创建文件夹
- cd /www/ginweb  进入指定的文件夹

### 2:   然后把项目的可执行文件上传到/www/ginweb

使用finalshell直接把main的可执行文件拖拽到/www目录即可。

### 3:  可执行文件授予执行权限

```sh
chmod +x main
```

### 4：启动 （占用式启动）

```sh
./main
```

如下

![image-20230903213658482](images/image-20230903213658482.png)

如果看到这个结果的说明你服务已经启动成功，如果要退出多按1或者多次：ctrl+c

这种占用式的，启动这个服务就不能去做其它的工作和事情呢，获取如果finalshell的连接会话关闭了。你服务器页会自动关闭。其实我们发布和部署希望不管会话关闭还是操作都应该让服务一直运行着。

### 5：守护方式启动 (后台方式启动)

```sh
nohup ./main &
```

### 6：然后执行访问（重点）

http://116.62.169.76:8080 

之所以能够把访问。是因为云服务中关于8080这个服务的房间已经对外开放了意思就是：8080这个端口在云服务器的安全组中已经公开了。

## 06-06、windows系统的发布和部署

1： 直接把main.exe上传到windows服务器

2： 直接双击打开main.exe即可。



## 07、小结

- 云服务购买，记得对应服务的端口要在安全组就进行开放
- 打包的记得在windows指定goos=linux
- linux系统运行go项目不需要安装go环境





# Go的项目发布和部署-深入探索

## 01、打包命令

```sh
go build [-o 输出名] [-i] [编译标记] [包名]
```

- 如果参数为`XX.go`文件或文件列表，则编译为一个个单独的包。
- 当编译单个`main`包（文件），则生成可执行文件。
- 当编译单个或多个包非主包时，只构建编译包，但丢弃生成的对象（`.a`），仅用作检查包可以构建。
- 当编译包时，会自动忽略`_test.go`的测试文件。

## 02、打包的方式

```
go build ---- 以工程名作为打包名 生成可执行文件的名字是工程名
go build . ---- 以工程名作为打包名, 生成可执行文件的名字是工程名
go build main.go ------ 以工程名入口文件作为打包。生成可执行文件的名字是：main
go build hello.go ------这个就是一个编译，在实际就普通打包，没有太大作用
go build -o server ---- 指定名字打包，生成的linux的可执行文件,可行文件的名字是：server
```

```
1. 普通包 【非main包】
go build add.go 【编译add.go,不生成exe执行文件】
go build -o add.exe add.go 【指定生成exe执行文件，但是不能运行此文件，不是main包】

2. main包【package main】
go build main.go 【生成exe执行文件】
go build -o main.exe main.go 【指定生成main.exe执行文件】

3. 项目文件夹下有多个文件
进入文件的目录
go build 【默认编译当前目录下的所有go文件】
go build add.go subtraction.go 【编译add.go 和 subtraction.go】

注意：
1. 如果是普通包，当你执行go build之后，它不会产生任何文件。【非main包】

2. 如果是main包，当你执行go
build之后，它就会在当前目录下生成一个可执行文件exe。如果你需要在$GOPATH/bin下生成相应的文件，需要执行go
install，或者使用go build -o 路径/xxx.exe xxx.go

3. 如果某个项目文件夹下有多个文件，而你只想编译某个文件，就可在go build之后加上文件名，例如go build
xxx.go；go build命令默认会编译当前目录下的所有go文件。

4. 你也可以指定编译输出的文件名。我们可以指定go build -o
xxxx.exe，默认情况是你的package名（main包），或者是第一个源文件的文件名（main包）。

5.go build会忽略目录下以“_”或“.”开头的go文件。
```

这个例子就是告诉我们一个逻辑，在开放中打包尽量使用go build -o 

```sh
// 从main.go开始打包
go build -o ginserver main.go ---- 从main.go开始打包，打包生成ginserver
go build -o ginserver.exe main.go ---- 从main.go开始打包，打包生成ginserver
// 从工程下开始打包
go build -o ginserver ---- 从main.go开始打包，打包生成ginserver
go build -o ginserver.exe---- 从main.go开始打包，打包生成ginserver
```

## 03、通用参数

```go
-a
    完全编译，不理会-i产生的.a文件(文件会比不带-a的编译出来要大？)
-n
    仅打印输出build需要的命令，不执行build动作（少用）。
-p n
    开多少核cpu来并行编译，默认为本机CPU核数（少用）。
-race
    同时检测数据竞争状态，只支持 linux/amd64, freebsd/amd64, darwin/amd64 和 windows/amd64.
-msan
    启用与内存消毒器的互操作。仅支持linux / amd64，并且只用Clang / LLVM作为主机C编译器（少用）。
-v
    打印出被编译的包名（少用）.
-work
    打印临时工作目录的名称，并在退出时不删除它（少用）。
-x
    同时打印输出执行的命令名（-n）（少用）.
-asmflags 'flag list'
    传递每个go工具asm调用的参数（少用）
-buildmode mode
    编译模式（少用）
    'go help buildmode'
-compiler name
    使用的编译器 == runtime.Compiler
    (gccgo or gc)（少用）.
-gccgoflags 'arg list'
    gccgo 编译/链接器参数（少用）
-gcflags 'arg list'
    垃圾回收参数（少用）.
-installsuffix suffix
    ？？？？？？不明白
    a suffix to use in the name of the package installation directory,
    in order to keep output separate from default builds.
    If using the -race flag, the install suffix is automatically set to race
    or, if set explicitly, has _race appended to it.  Likewise for the -msan
    flag.  Using a -buildmode option that requires non-default compile flags
    has a similar effect.
-ldflags 'flag list'
    '-s -w': 压缩编译后的体积
    -s: 去掉符号表
    -w: 去掉调试信息，不能gdb调试了
-linkshared
    链接到以前使用创建的共享库
    -buildmode=shared.
-pkgdir dir
    从指定位置，而不是通常的位置安装和加载所有软件包。例如，当使用非标准配置构建时，使用-pkgdir将生成的包保留在单独的位置。
-tags 'tag list'
    构建出带tag的版本.
-toolexec 'cmd args'
    ？？？？？？不明白
    a program to use to invoke toolchain programs like vet and asm.
    For example, instead of running asm, the go command will run
    'cmd args /path/to/asm <arguments for asm>'.
```



## 总结

打包时候尽量使用如下的方式

```sh
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux/windows
go build -o 目录/可执行文件的名字(.exe)  main.go
go build -o 目录/可执行文件的名字(.exe) 
```

比如：

```sh
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux/windows
go build -o build/ginweb(.exe)  main.go
go build -o build/ginweb(.exe) 
```

# Go的项目发布和部署-后端



## 01、先安装MSYQL

本机系统中： [xkginweb](C:\Users\zxc\go\xkginweb) 是一个go语言开放的项目工程，其中你必须在发布的时候要分析这个工程依赖那些中间件。依赖：mysql/redis

进入官网: https://mysql.com/downloads/

![image-20230903222053991](images/image-20230903222053991.png)

下载MYSQL5.7

![image-20230903222103674](images/image-20230903222103674.png)

### **Linux下载**



### 1:  创建一个目录/www/mysql

```sh
cd /usr/local
wget https://downloads.mysql.com/archives/get/p/23/file/mysql-5.7.40-linux-glibc2.12-x86_64.tar.gz
```

如果执行发现没有wget命令，你先安装一下 `yum install wget`

### 2: 压缩包解压

```
 tar -zxvf  mysql-5.7.40-linux-glibc2.12-x86_64.tar.gz
```

![image-20230903222332391](images/image-20230903222332391.png)

### 3: 修改目录名字和创建数据存储的目录

 修改名称将目录名称直接修改为mysql

```sh
mv mysql-5.7.40-linux-glibc2.12-x86_64  mysql
```

  (进入MySQL 目录) 创建数据目录 ：

```sh
cd /usr/local/mysql 
mkdir -p /usr/local/mysql/data
```

 给数据目录赋权限：

```sh
chmod -R 777 /usr/local/mysql/data
```

如果出现：chmod: invalid mode: ‘–R’  是减号 有问题，复制出来，在编辑下 减号

![image-20230903222454555](images/image-20230903222454555.png)





### 4: 创建用户 、组、并将用户加入组,修改配置文件

```
groupadd mysql
useradd -g mysql mysql
```

修改MySQL 配置文件：` vi /etc/my.cnf`   （通过上下左右将光标移动到需要编辑的位置按i进行添加，完成按Esc结束编辑 输入":wq!" 保存退出）

```cnf
[mysqld]
bind-address=0.0.0.0
port=3306
user=mysql
basedir=/usr/local/mysql
datadir=/usr/local/mysql/data
socket=/tmp/mysql.sock
log-error=/usr/local/mysql/data/mysql.err
pid-file=/usr/local/mysql/data/mysql.pid
#character config
character_set_server=utf8mb4
symbolic-links=0
explicit_defaults_for_timestamp=true
```

### 5: 配置mysql服务，初始化数据目录和工作目录和查看密码

**进入mysql bin 目录下面**

```sh
cd /usr/local/mysql/bin
```

 **执行命令**

```sh
./mysqld --initialize --user=mysql --datadir=/usr/local/mysql/data/ --basedir=/usr/local/mysql/
```

如果执行出现如下错误：

>./mysqld: error while loading shared libraries: libaio.so.1: cannot open shared object file: No such file or directory
>
>请安装： yum install -y libaio 

然后在执行上面的命令。



查看mysql 密码（红线画的地方是密码）

```sh
cat /usr/local/mysql/data/mysql.err
```

![image-20230903222616543](images/image-20230903222616543.png)



### 6: 添加软连接，并启动mysql服务

```sh
ln -s /usr/local/mysql/support-files/mysql.server /etc/init.d/mysql
ln -s /usr/local/mysql/bin/mysql /usr/bin/mysql
service mysql start
```

退出mysql是:  `exit`

### 7： 登录mysql 修改密码，访问权限

进入mysql  bin目录下面

```sh
-- 复制或者手动输入启动
./mysql -hlocalhost -uroot -p
-- 一步到位把密码抓出来启动
./mysql -uroot -p$(awk '/temporary password/{print $NF}'  /usr/local/mysql/data/mysql.err)
```

![image-20230903222736305](images/image-20230903222736305.png)



修改密码（设置密码尽量设置复杂一点，拒绝弱口令）

```sql
mysql > set password=password('mkxiaoer');
mysql > flush privileges;
```

 授予远程访问，方式1：

```sh
mysql > use mysql;
mysql > update user set Host='%' where User='root';
mysql > flush privileges;
```

 授予远程访问，方式2：

```sql
mysql > grant all on *.* to root@'%' identified by 'mkxiaoer';
mysql > flush privileges;
```

**==还有如果你是云服务器还要在安全组中请把mysql的端口：3306进行配置开放。==**



## 02、在把xkginweb打包

```sh
set CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
go build -o build/xkginweb main.go
```

![image-20230903230521996](images/image-20230903230521996.png)

注意：conf是我手动复制上去的



## 03、把打包项目发布云服务器

使用finalshell上传conf和xkginweb可执行到服务器/www/xkginweb目录下即可。

![image-20230903230500531](images/image-20230903230500531.png)



## 04、启动服务

但是注意go项目是不会把资源目录自动打包到执行文件中，所以你必须手动的把资源目录和可行执行一起上传到服务器。

```sh
chmod +x xkginweb
```

启动

```sh
nohup ./xkginweb &
> 16611 ----启动成功的进程id ,方便让你关闭服务的
> oot@iZbp13fwzug417rudbljt3Z xkginweb]# nohup: 忽略输入并把输出追加到"nohup.out
```

查看

```sh
ps -ef | grep xkginweb
root     16611  1312  0 23:02 pts/0    00:00:00 ./xkginweb
root     17271  1312  0 23:02 pts/0    00:00:00 grep --color=auto xkginweb

```

关闭服务

```sh
kill -9 16611
```



# Go的项目发布和部署-前端

发布前端项目，为什么还要nginx。你思考一个问题。

我们在开放前端项目的时候，我们使用脚手架vite。它可以让前端具有服务性。其实就有端口可以让访问。但是前端项目一旦打包。这个服务性就自动丢失。因为前端打包最终会通过webpack把项目中所以的src下面的源码进行编译生成对应js文件。public 进行拷贝，然后把这些js文件和public生成的文件放入到一个dist目录。最后完成打包的过程。

但是生成打包的内容，是不能直接只有服务性，你必须要依托一个外部的服务，才能让我们打包的前端项目继续运行。



## **1:  先安装nginx**

查看 ： 04、Nginx发布项目和部署.md

## 2: web项目打包

优化配置 vite.config.js

```js
import { fileURLToPath, URL } from 'node:url'
import { defineConfig,loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { viteLogo } from './src/core/config'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  // 根据当前工作目录中的 `mode` 加载 .env 文件
  // 设置第三个参数为 '' 来加载所有环境变量，而不管是否有 `VITE_` 前缀。
  const env = loadEnv(mode, process.cwd(), '')
  viteLogo(process.env)
  return {
    base: env.VITE_MODE === 'production' ? './' : '/',
    // vite 配置
    plugins: [
      vue(),
      AutoImport({
        imports: ['vue','vue-router','pinia','vue-i18n'],
        resolvers: [ElementPlusResolver()],
      }),
      Components({
        resolvers: [ElementPlusResolver()],
      })
    ],
    server:{
      // 如果使用docker-compose开发模式，设置为false
      open: true,
      port: 8777,
      proxy: {
        // 把key的路径代理到target位置
        [env.VITE_BASE_API]: { // 需要代理的路径   例如 '/api'
          target: `${env.VITE_BASE_PATH}/`, // 代理到 目标路径
          changeOrigin: true,
          //rewrite: path => path.replace(new RegExp('^' + env.VITE_BASE_API), ''),
        }
      },
    },
     // 构建配置
     build: {
      // 输出目录，默认是 dist
      outDir: 'dist',
      // 是否开启 sourcemap
      sourcemap: false,
      // 是否开启压缩
      minify: 'terser', // 可选值：'terser' | 'esbuild'
      // 是否开启 brotli 压缩
      brotli: true,
      // 是否将模块提取到单独的 chunk 中，默认是 true
      chunkSizeWarningLimit: 500,
      // 是否提取 CSS 到单独的文件中
      cssCodeSplit: true,
      // 是否开启 CSS 压缩
      cssMinify: true,
      // 是否开启 CSS 去重
      cssInlineLimit: 4096,
      // 启用/禁用 esbuild 的 minification，如果设置为 false 则使用 Terser 进行 minification
      target: 'es2018', // 可选值：'esnext' | 'es2020' | 'es2019' | 'es2018' | 'es2017' | 'es2016' | 'es2015' | 'es5'
      // 是否开启 Rollup 的代码拆分功能
      rollupOptions: {
          output: {
              manualChunks: {},
          },
      },
      terserOptions: { 
        compress: { // 打包时清除 console 和 debug 相关代码
          drop_console: true,
          drop_debugger: true,
        },
      },
      // 是否开启增量式构建
      // https://vitejs.dev/guide/build.html#incremental-build
      incremental: false
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    // 优化配置
    optimizeDeps: {
        // 是否将 Vue、React、@vueuse/core 和 @vueuse/head 作为外部依赖提取出来
        include: ['vue', 'react', '@vueuse/core', '@vueuse/head','axios'],
        // 是否开启预构建，将预构建后的代码提前注入到浏览器缓存中，以减少首次加载的时间
        prebuild: false,
    }
  }
})



```

然后执行打包

```sh
npm run build
```

然后生成一个dist目录如下：

![image-20230903231551910](images/image-20230903231551910.png)

但是这个目录，不具有服务性，因为打包以后就是一个纯纯静态页面和资源而已。所以你必须要找一个web服务器来运行它。这个web就是nginx



### 4：上传dist文件的内容到nginx

把dist打包中内容上传到 nginx的html目录覆盖即可。（/usr/local/nginx/html）目录下。

然后执行即可。只要你不修改nginx的配置文件。是不需要重启nginx。你直接访问就可以看到我们web项目运行成功了。

### 5： 访问前端项目

http://116.62.169.76/#/login?path=/dashboard





#  Go的项目发布和部署-api接口的配置



但是你会发现你现在的web网站访问不了。api接口调用失败。原因是什么？

- 本地环境
  - .env.devepolmenet 
    - vite -server proxy
    - vite.config.js 
      - /api  —— http://127.0.0.1:8990
- 生成环境
  - 生成环境就已经和vite没有半毛钱关系了。那个服务代理就对于打包以后的生产环境就没有任何意义了。
  - 我的服务器接口路径就是：/api 它是一个相对路径。
  - 相对路径是什么。相对访问路径。现在访问路径是： http://116.62.169.76
    - 最终通过axios请求路径全部都是： http://116.62.169.76/api/role/save
    - 最终通过axios请求路径全部都是： http://116.62.169.76/api/role/list

大家发现没有。除非你的服务器启动一个80服务同时以/api开头。程序就可以正常的访问。否则都会失败。而现在我们系统提供接口服务是什么样子：http://116.62.169.76:8990/api/xxxxxxxxx



如果能够让访问能正常调用接口服务，可以使用两种方式：

## 第一种方式：修改.evn.production的的配置改成绝对路径

```js
ENV = 'production'

# 前端vite项目的端口
VITE_CLI_PORT = 8777
# 服务前缀
VITE_BASE_API = http://116.62.169.76:8990/api
```

- 直接把.evn.production中/api的相对路径，改成绝对路径。
- 然后重新打包
- 然后把ｎｇｉｎｘ的ｈｔｍｌ目录下的内容删掉，然后把ｄｉｓｔ生成新的上传上来。不需要重新
- 直接访问即可
- 可能会存在跨域的问题，需要在代码中去解决

## **第二种方式：使用nginx服务代理**

### 1： 进入nginx的配置目录

```sh
cd /usr/local/nginx/conf
```

### 2:  然后修改nginx.conf文件如下

```json
server {
  location /api {
     proxy_pass   http://127.0.0.1:8990/api;
  }
}
```

### 3: 注意改了配置文件一定要重启

```json
nginx -t --- 检查你修改的配置文件是否正确
nginx -s reload  ---重新
```

### 4： 然后修改web工程中的.evn.production的配置如下

```json
ENV = 'production'
# 前端vite项目的端口
VITE_CLI_PORT = 8777
# 服务前缀
VITE_BASE_API = http://116.62.169.76
```

### 5: 然后重新打包

```json
npm run build
```

### 6: 然后打包的内容上传到nginx的html目录下

记得上传之前把旧删掉在上传。

# 安装Nginx服务

Nginx安装

nginx下载：http://nginx.org/en/download.html

### 01、创建nginx服务器目录

```
mkdir -p /www/kuangstudy/nignx
cd /www/kuangstudy/nignx
```

### 02、下载安装

```
wget http://nginx.org/download/nginx-1.20.1.tar.gz
```

### 03、安装编译工具及库文件

```
yum -y install make zlib zlib-devel gcc-c++ libtool  openssl openssl-devel
```

### 04、解压nginx

```
tar -zxvf nginx-1.20.1.tar.gz
```

### 05、创建nginx的临时目录

```
mkdir -p /var/temp/nginx
```

### 06、进入安装包目录

```
cd nginx-1.20.1
```

### 07、编译安装

```
./configure \
--prefix=/usr/local/nginx \
--pid-path=/var/run/nginx.pid \
--lock-path=/var/lock/nginx.lock \
--error-log-path=/var/log/nginx/error.log \
--http-log-path=/var/log/nginx/access.log \
--with-http_gzip_static_module \
--http-client-body-temp-path=/var/temp/nginx/client \
--http-proxy-temp-path=/var/temp/nginx/proxy \
--http-fastcgi-temp-path=/var/temp/nginx/fastgi \
--http-uwsgi-temp-path=/var/temp/nginx/uwsgi \
--http-scgi-temp-path=/var/temp/nginx/scgi \
--with-http_stub_status_module \
--with-http_ssl_module \
--with-http_stub_status_module
```

安装以后的目录信息

```
nginx path prefix: "/usr/local/nginx"  nginx binary file: "/usr/local/nginx/sbin/nginx"  nginx modules path: "/usr/local/nginx/modules"  nginx configuration prefix: "/usr/local/nginx/conf"  nginx configuration file: "/usr/local/nginx/conf/nginx.conf"  nginx pid file: "/var/run/nginx.pid"  nginx error log file: "/var/log/nginx/error.log"  nginx http access log file: "/var/log/nginx/access.log"  nginx http client request body temporary files: "/var/temp/nginx/client"  nginx http proxy temporary files: "/var/temp/nginx/proxy"  nginx http fastcgi temporary files: "/var/temp/nginx/fastgi"  nginx http uwsgi temporary files: "/var/temp/nginx/uwsgi"  nginx http scgi temporary files: "/var/temp/nginx/scgi"
```

![img](../../../../../../L_Learning/%25E6%25B5%258B%25E5%25BC%2580%25E8%25AF%25BE%25E7%25A8%258B/%25E7%258B%2582%25E7%25A5%259E/3-%25E9%25A1%25B9%25E7%259B%25AE%25E5%25AE%259E%25E6%2588%2598%2520-%2520GVA%25E5%2590%258E%25E5%258F%25B0%25E9%25A1%25B9%25E7%259B%25AE%25E7%25AE%25A1%25E7%2590%2586%25E5%25BC%2580%25E5%258F%2591/20230903%25EF%25BC%259A%25E7%25AC%25AC%25E5%259B%259B%25E5%258D%2581%25E4%25B8%2583%25E8%25AF%25BE%25EF%25BC%259A%25E8%2587%25AA%25E5%25BB%25BA%25E9%25A1%25B9%25E7%259B%25AE-%2520%25E9%25A1%25B9%25E7%259B%25AE%25E7%259A%2584%25E5%258F%2591%25E5%25B8%2583%25E5%2592%258C%25E9%2583%25A8%25E7%25BD%25B2/%25E9%25A1%25B9%25E7%259B%25AE%25E7%25AC%2594%25E8%25AE%25B0/assets/kuangstudy09f12de8-430a-417d-a5f1-b8bb772a5c0b.png)

### 08、 make编译

```
make
```

### 09、 安装

```
make install
```

### 10、 进入sbin目录启动nginx

```
cd /usr/local/nginx/sbin
```

执行nginx启动

```
./nginx
#停止：
./nginx -s stop
#重新加载：
./nginx -s reload
```

### 11、打开浏览器，访问虚拟机所处内网ip即可打开nginx默认页面，显示如下便表示安装成功：

```
http://ip
```

### 12、注意事项

1. 如果在云服务器安装，需要开启默认的nginx端口：80
2. 如果在虚拟机安装，需要关闭防火墙
3. 本地win或mac需要关闭防火墙
4. nginx的安装目录是：/usr/local/nginx/sbin

### 13、配置nginx的环境变量

```
vim /etc/profile
```

在文件的尾部追加如下：

```
export NGINX_HOME=/usr/local/nginx
export PATH=$NGINX_HOME/sbin:$PATH
```

重启配置文件

```
source /etc/profile
```

# 01、安装Nginx服务集群发布和部署

## 分析

在项目发布和部署，我们其实会把开发好的go的web项目打包成可执行文件，然后运行在linux系统。一旦运行就会在系统进程里生成一个对应服务。然后可以访问。但是你思考过一个这样的问题。如果你的服务突然之间挂了。那么你网站是其实就立即出现无法访问的状况。

所以在大部分的发布和部署的时候，我们都不会单节点运行，而是更多考虑使用集群的方式进行部署。

- 其实早期的发布和部署，其实都停服去进行发布和部署，然后给与网站停服的标识。
- 灰度交替的方式进行发布和部署（集群模式）
  - 比如三个节点：1 2 3  ，先把服务1停止掉，然后在更新到最新，如果服务1运行没有任何问题，然后在停止第二服务器，以此类推。但是你不能全部停掉在更新，这个不建议。
- jenkins进行发布和部署（开发阶段）（==即见即所得==）
- 面向容器化的发布和部署（docker + k8s）

## 01、准备工作

1：准备服务器资源（32G 、64G） （但是有一个前提：必须在同机房才又意义）

2： 准备web项目（单项目，微服务项目）

3：服务环境（mysql/redis/nginx）



## 02、单项目

没有服务拆分，比如：用户服务（登录，注册）、订单服务、广告服务、内容服务、分类服务、搜索服务等等。全部都在一个项目里。不分开，

**缺点：**

1：可执行文件很大，占用服务器资源很多

2：存在服务不隔离，如果一份服务出现了一次和故障，那么你整个体统就全部崩溃。

**好处**

1：发布方便，快捷



## 03、微服务项目

把服务全部进行拆分，独立成一个系统和服务来进行开发和部署。

缺点：

1：运维的难度增大，需求服务器资源也增多

2：成本上升

3：学习的成本也增大

4：数据一致性也要考虑

优点 ：

1：服务直接是独立，出现异常和故障互补影响，



## 04、web集群怎么规划和部署

web集群：其实就把若个项目运行对外提供服务。要完成集群必须要找到：具有负载均衡的web中间件。这种中间件有如下：

- httpd
- etc
- nginx (推荐)



## 05、集群的规划和误区

集群其实就把服务部署多份，然后使用 web中间件来进行负责均衡管理。你就可以看到集群切换效果了。通杀也也可以解决单点故障的问题。

但是如何进行多服务的部署呢？

- 单机部署 （单机部署多应用）
- 多机部署 （交叉集群部署）

但是不论单机部署多应用还是多机部署交叉的方式，你记住集群节点的部署和规划一定是在：同一个局域网（同一个ip段上）才又意义，才会快。最快肯定是单机部署多应用（全部安装在一个机器上）。

==也就是要解决和一个误区。未来我开发了项目。我想要我的项目支持集群，也让我们的项目能支持并发和高可用。我就买一堆的云服务器让后进行发布和部署。你就以为就很快了==

快就必须遵循几个规则：

1：要么是自建服务器，公司会统一分配网关和ip段

2：如果云服务器，建议大家买在一个区域，然后代码中全部使用内网ip而不是公网ip。





## 06、集群带来的问题和思考

1、集群部署的时候，一定是使用奇数节点，不要使用偶数。（3个集群基本开端）

 

# 02、单机部署多应用的方式

## 单机部署多应用

1：什么情况下会使用单机部署多应用。

- 项目的早期
- 比较单一的系统
- ==其实就一个标准：你项目收入和单台服务器成本是合理的。你都可以考虑使用单机部署多应用。==
- 因为不需要考虑多服务器之间的沟通问题和维护问题。因为维护一个服务器比维护N个服务器要轻松很多。



![image-20230910223849424](images/image-20230910223849424.png)

## 1：准备工作

把现有的项目修改三个端口（8081，8082,8083），然后打包三次。生成对应的三个这种服务的可执行文件

```sh
go build -o build/ginweb-8081 main.go
go build -o build/ginweb-8082 main.go
go build -o build/ginweb-8083 main.go
```

但是注意打包之前，一定要把代码中或者配置中的端口改好以后在打包。

![image-20230910213442187](images/image-20230910213442187.png)

## 2：然后分别把这三个可执行文件，上传到服务器上。

然后赋予三个文件的可执行权限。

```sh
chmod +x /www/ginweb/ginweb-8081
chmod +x /www/ginweb/ginweb-8082
chmod +x /www/ginweb/ginweb-8083
```

## 3: 启动它们

窗口1：

```sh
cd /www/ginweb
./ginweb-8081
```

窗口2：

```sh
cd /www/ginweb
./ginweb-8082
```

窗口3：

```sh
cd /www/ginweb
./ginweb-8083
```

### **守护的方式(确定完毕以后，建议使用如下：)**

```sh
cd /www/ginweb
nohup ./ginweb-8081 &
nohup ./ginweb-8082 &
nohup ./ginweb-8083 &
```

## 4: 然后测试它们三个服务器是否正常

- http://120.55.71.124:8081

- http://120.55.71.124:8082

- http://120.55.71.124:8083

## 5:  开始安装nginx 

请看课件04、Nginx发布项目和部署.md 

```sh
# 启动
nginx 
# 重启
nginx -s reload
# 停止
nginx -s stop
# 检查配置是否正确
nginx -t
# 看安装目录
nginx -V
# 看帮助文档
nginx -h

```

## 6: 开始配置nginx的集群映射

默认配置：nginx 端口是：80

nginx.conf

```conf
#user  nobody;
worker_processes  1;

#error_log  logs/error.log;
#error_log  logs/error.log  notice;
#error_log  logs/error.log  info;

#pid        logs/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       mime.types;
    default_type  application/octet-stream;

    #log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
    #                  '$status $body_bytes_sent "$http_referer" '
    #                  '"$http_user_agent" "$http_x_forwarded_for"';

    #access_log  logs/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    #keepalive_timeout  0;
    keepalive_timeout  65;

    #gzip  on;
    # 建议使用外部包含，防止破坏nginx.conf文件而造成报错
    include web/*.conf;

}

```

web/ginweb.conf如下：

```json
upstream goservers{
   server 127.0.0.1:8081 weight=1;
   server 127.0.0.1:8082 weight=2;
   server 127.0.0.1:8083 weight=1;
}

server {
   listen       80;
   server_name  localhost;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
       index  index.html index.htm;
   }
      
   #	 运行接口项目
   location /api {
        proxy_pass   http://127.0.0.1:8990/api;
   }


   location /web {
        proxy_pass   http://goservers/;
   }

   

#   error_page  404    /404.html;

   # redirect server error pages to the static page /50x.html
   #
#   error_page   500 502 503 504  /50x.html;
#   location = /50x.html {
#       root   html;
#   }

}
```



# 03、二级域名的配置

如果端口被占用的情况nginx可以服用端口吗？可以使用二级域名

1： 因为一个端口一个服务默认情况，根访问只能是一个。如下：

```json
upstream goservers{
   server 127.0.0.1:8081;
   server 127.0.0.1:8082;
   server 127.0.0.1:8083;
}

server {
   listen       80;
   server_name  localhost;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
       index  index.html index.htm;
   }
      
   #	 运行接口项目
   location /api {
        proxy_pass   http://127.0.0.1:8990/api;
   }


   location /web {
        proxy_pass   http://goservers/;
   }

   

#   error_page  404    /404.html;

   # redirect server error pages to the static page /50x.html
   #
#   error_page   500 502 503 504  /50x.html;
#   location = /50x.html {
#       root   html;
#   }

}
```

这个时候如果去访问 http://120.55.71.124:80(80可以不加，80是唯一一个不需要指定的端口。为了443是https也不需要)， 所以后续如果你继续想用80端口来访问你的服务。你必须使用location(路径来隔离)，所以就会出现

- http://120.55.71.124 访问后台首页
- http://120.55.71.124/api 访问后台的api接口服务
- http://120.55.71.124/web 访问前台的web服务

但是往往很多初学者会一个这样想法，能不能做到继续使用80端口。但是不使用location来隔离。注意在同一个listen和server_name的情况下是不可能做到的。除非在买个服务器。也可以使用server_name来进行隔离，比如采用二级域名

```json
upstream goservers{
   server 127.0.0.1:8081;
   server 127.0.0.1:8082;
   server 127.0.0.1:8083;
}

server {
   listen       80;
   server_name  localhost;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
       index  index.html index.htm;
   }
      
   #	 运行接口项目
   location /api {
        proxy_pass   http://127.0.0.1:8990/api;
   }


   location /web {
        proxy_pass   http://goservers/;
   }


}


server {
   listen       80;
   server_name  web.haitang.com;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
         proxy_pass   http://goservers/;
   }
      

#   error_page  404    /404.html;

   # redirect server error pages to the static page /50x.html
   #
#   error_page   500 502 503 504  /50x.html;
#   location = /50x.html {
#       root   html;
#   }

}

server {
   listen       80;
   server_name  api.haitang.com;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
         proxy_pass   http://goservers/;
   } 
}

server {
   listen       80;
   server_name  upload.haitang.com;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
         proxy_pass   http://goservers/upload;
   } 
}
```

也使用端口隔离

```json
upstream goservers{
   server 127.0.0.1:8081;
   server 127.0.0.1:8082;
   server 127.0.0.1:8083;
}

server {
   listen       80;
   server_name  localhost;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
       index  index.html index.htm;
   }
      
   #	 运行接口项目
   location /api {
        proxy_pass   http://127.0.0.1:8990/api;
   }


   location /web {
        proxy_pass   http://goservers/;
   }


}


server {
   listen       81;
   server_name  localhost;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
         proxy_pass   http://goservers/;
   }
      

#   error_page  404    /404.html;

   # redirect server error pages to the static page /50x.html
   #
#   error_page   500 502 503 504  /50x.html;
#   location = /50x.html {
#       root   html;
#   }

}

server {
   listen       82;
   server_name  localhost;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
         proxy_pass   http://goservers/;
   } 
}

server {
   listen       83;
   server_name  localhost;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
         proxy_pass   http://goservers/upload;
   } 
}
```





# 04、使用shell的脚本的方式解决新老更替的问题

```sh
echo 'go 项目开始启动了'
ps aux | grep -ai "ginweb" | grep -v grep | awk '{print $2}' | xargs kill -9
echo '上传的服务全部关闭成功'
chmod +x /www/ginweb/ginweb-8081
chmod +x /www/ginweb/ginweb-8082
chmod +x /www/ginweb/ginweb-8083
echo '-----授权成功----开始启动'
nohup /www/ginweb/ginweb-8081 &
nohup /www/ginweb/ginweb-8082 &
nohup /www/ginweb/ginweb-8083 &
echo '启动完毕'
```

然后赋予start.sh的可执行权限

```sh
chmod +x start.sh
```

然后执行

```sh
cd /www/ginweb
./start.sh
```

==但是建议大家不在生产环境的时候，全部关停，可以写两个脚本，一个脚本关闭一部分，如果没问题在执行第二脚本，把其他的全部关闭在发布最新的。==

# 05、自建机房：多机多部署多应用

1: web服务节点

- 192.168.110.1: 8081
- 192.168.110.2: 8081
- 192.168.110.3: 8081
- 192.168.110.4: 8081

2: nginx服务节点

- 192.168.110.2: 80

3:   nginx.conf配置如下：

```
upstream goservers{
   server 192.168.110.1:8081;
   server 192.168.110.2:8081;
   server 192.168.110.3:8081;
   server 192.168.110.4:8081;
}

server {
   listen       80;
   server_name  localhost;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
       index  index.html index.htm;
   }
      
   #	 运行接口项目
   location /api {
        proxy_pass   http://127.0.0.1:8990/api;
   }


   location /web {
        proxy_pass   http://goservers/;
   }

   

#   error_page  404    /404.html;

   # redirect server error pages to the static page /50x.html
   #
#   error_page   500 502 503 504  /50x.html;
#   location = /50x.html {
#       root   html;
#   }

}
```



# 06、云服务器：多机多部署多应用

1: web服务节点

- 176.168.110.135: 8081
- 176.18.11.135: 8081
- 176.68.10.135: 8081
- 176.18.10.135: 8081

2: nginx服务节点

- 125.58.58.41:80

3:   nginx.conf配置如下：

```json
upstream goservers{
   server 176.168.110.135:8081;
   server 176.18.11.135:8081;
   server 176.68.10.135:8081;
   server 176.18.10.135:8081;
}

server {
   listen       80;
   server_name  localhost;

   # 主要是用来运行我们的web的后台项目
   location / {
       root   html;
       index  index.html index.htm;
   }
      
   #	 运行接口项目
   location /api {
        proxy_pass   http://127.0.0.1:8990/api;
   }


   location /web {
        proxy_pass   http://goservers/;
   }

  

#   error_page  404    /404.html;

   # redirect server error pages to the static page /50x.html
   #
#   error_page   500 502 503 504  /50x.html;
#   location = /50x.html {
#       root   html;
#   }

}
```

# 01、环境隔离



## 01、概述

正规的流程：开发环境（开人员）—————–测试环境（测试人员（功能测试，压测指标，服务参数））——-预发布环境（运维人员）—————生成环境

在前面发布的go项目。在发布过程中我们环境一般来说会分为开发环境、测试环境、和生产环境，不同环境都自己的配置信息，如果在项目我们都共用一个配置文件的话，可能会造成开发环境、测试环境、生产环境。三个环境的配置不不断的切换和变更。

比如：我们项目里其实就提供一个`application.yml`配置文件

开发环境：

```yaml
# 服务端口的配置
server:
  port: 8990
  context: /
# 数据库配置
# "root:mkxiaoer@tcp(127.0.0.1:3306)/ksd-social-db?charset=utf8&parseTime=True&loc=Local", // DSN data source name
database:
  mysql:
    host: 127.0.0.1
    port: 3306
    dbname: kva-admin-db
    username: root
    password: mkxiaoer
    config: charset=utf8&parseTime=True&loc=Local
# nosql数据的配置
nosql:
  redis:
    host: 127.0.0.1:6379
    password: 
    db: 0
  es:
    host: 127.0.0.1
    port: 9300
    password: 456
```

使用viper来进行解析的和生效的。

发布的流程：

- 修改application.yaml的配置信息改成测试环境或者生产环境
- 然后使用go build打包
- 然后把可执行文件放入到服务器
- 然后启动或者集群启动即可。

如果这个时候发布完了。我要继续开发。你必须要生产环境的配置改成开发环境配置，你才能够正常开发。所以这种使用单配置文件来属性的标记注释切换，是非常麻烦和不推荐的。假设我们能够实现可以通过环境的方式来加载他们对应的配置文件那么不完美了么？



## 02、根据环境来隔离配置信息

- 开发环境 development 简化的名字：application-dev.yaml

```yaml
# 服务端口的配置
server:
  port: 8990
  context: /
# 数据库配置
# "root:mkxiaoer@tcp(127.0.0.1:3306)/ksd-social-db?charset=utf8&parseTime=True&loc=Local", // DSN data source name
database:
  mysql:
    host: 127.0.0.1
    port: 3306
    dbname: kva-admin-db
    username: root
    password: mkxiaoer
    config: charset=utf8&parseTime=True&loc=Local
# nosql数据的配置
nosql:
  redis:
    host: 127.0.0.1:6379
    password:
    db: 0
  es:
    host: 127.0.0.1
    port: 9300
    password: 456
```

- 测试环境 test ：简化的名字：application-test.yaml

```yaml
# 服务端口的配置
server:
  port: 8991
  context: /
# 数据库配置
# "root:mkxiaoer@tcp(127.0.0.1:3306)/ksd-social-db?charset=utf8&parseTime=True&loc=Local", // DSN data source name
database:
  mysql:
    host: 127.0.0.1
    port: 3306
    dbname: kva-admin-db
    username: root
    password: mkxiaoer
    config: charset=utf8&parseTime=True&loc=Local
# nosql数据的配置
nosql:
  redis:
    host: 127.0.0.1:6379
    password: mkxiaoer1986.
    db: 0
  es:
    host: 127.0.0.1
    port: 9300
    password: 456
```

- 生成环境 production生产环境：简化的名字：application-prod.yaml

```yaml
# 服务端口的配置
server:
  port: 8992
  context: /
# 数据库配置
# "root:mkxiaoer@tcp(127.0.0.1:3306)/ksd-social-db?charset=utf8&parseTime=True&loc=Local", // DSN data source name
database:
  mysql:
    host: 127.0.0.1
    port: 3306
    dbname: kva-admin-db
    username: root
    password: mkxiaoer
    config: charset=utf8&parseTime=True&loc=Local
# nosql数据的配置
nosql:
  redis:
    host: 127.0.0.1:6379
    password: mkxiaoer1986.
    db: 0
  es:
    host: 127.0.0.1
    port: 9300
    password: 456
```



## 03、使用viper进行环境配置文件的加载



```go
package initilization

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"xkginweb/global"
)

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func InitViper() {
	//获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config := viper.New()
	//config.AddConfigPath(path + "/conf") //设置读取的文件路径
	//config.SetConfigName("application")  //设置读取的文件名
	//config.SetConfigType("yaml")         //设置文件的类型
	//logs.Info("你激活的环境是：" + GetEnvInfo("env"))
	config.SetConfigFile(path + "/conf/application-dev.yaml")
	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = config.Unmarshal(&global.Config); err != nil {
			fmt.Println(err)
		}
	})

	// 这里才是把yaml配置文件解析放入到Config对象的过程---map---config
	if err = config.Unmarshal(&global.Config); err != nil {
		fmt.Println(err)
	}

	// 打印文件读取出来的内容:
	keys := config.AllKeys()
	dataMap := make(map[string]interface{})
	for _, key := range keys {
		fmt.Println("yaml存在的key是: " + key)
		dataMap[key] = config.Get(key)
	}

	global.Yaml = dataMap

}

```

通过修改`InitViper` 然后不停的修改配置文件的环境，就得到不同环境的切换，这样解决在一个文件中来回注释切换的问题了。







# go程序去读取环境变量

在大部分的程序中，都可以读取到系统中设定环境变量的参数信息。通过参数系统你可以设定一些初始化值，可以让未来可执行程序不用把代码写死，而是直接通过环境变量的方式来进行配置。但是配置程序就必须有方案或者库区读取。

假设：环境变量配置的加载

1：如果我在系统环境变量设置一个env=prod

2:  go同os组件读取环境变量

3：然后在viper把读取到环境变量进行拼接

4：这样未来我们只需要去更改系统环境变量了。

就不需要在代码切换，更何况如果我们go build打包以后生成的可执行文件根本就不可以更改。如果更改就必须改源码在在打包。





## 01、如何在系统环境中设置参数呢?



#### windows

- 【我的电脑】–【属性】

![image-20230912204546060](images/image-20230912204546060.png)

![image-20230912204603732](images/image-20230912204603732.png)

![image-20230912204615257](images/image-20230912204615257.png)

![image-20230912204649214](images/image-20230912204649214.png)

然后使用代码读取即可

```go
package main

import (
	"fmt"
	"os"
)

func main() {

	//1:go程序去读取环境变量
	fmt.Println("ok")
	path := os.Getenv("GOPATH")
	root := os.Getenv("GOROOT")
	pathext := os.Getenv("PATHEXT")
	goenv := os.Getenv("goenv")
	fmt.Println(path)
	fmt.Println(root)
	fmt.Println(pathext)
	fmt.Println(goenv)

	//2:go程序去读取命令行参数
}

```

你可以在环境变量中自己定义一个属性比如：goenv = dev  ，如下：

![image-20230912205044180](images/image-20230912205044180.png)

==这个你去运行main.go发现读取不出来，是不是没生效，其实不是，你可以把golang工具全部关闭，然后重启，你会发现就读取到了。那就说明开发工具golang对环境变量进行缓存。所以你第一次读取不得原因也就在这里。当然如果你发布生产是不会有这种问题。==



## idea/golang开发工具如何配置环境变量

首先为什么开发工具要自己去同步缓存一份环境变量呢，==其实就告诉不要轻易在电脑中系统环境中去增加页面的属性，而是在开发工具去增加，这样也可以业务的属性增加和系统里进行分离。防止污染！==如何完成在开发工具中去增加呢？

1：新建一个main.go，执行一下

2：然后找到main.go的执行编辑配置

![image-20230912205611723](images/image-20230912205611723.png)

3: 然后增加业务级别的环境配置参数

![image-20230912205919614](images/image-20230912205919614.png)



## 03、方式一：使用golang的环境配置读取环境变量参数

1： 定义读取的方法

```go
func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

```

2: 在代码进行环境变量的配置

![image-20230912210727257](images/image-20230912210727257.png)

3: 根据环境来加载配配置信息

```go
package initilization

import (
	"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"xkginweb/global"
)

func GetEnvInfo(env string) string { //---------------------新增代码
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func InitViper() {
	//获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config := viper.New()
	//config.AddConfigPath(path + "/conf") //设置读取的文件路径
	//config.SetConfigName("application")  //设置读取的文件名
	//config.SetConfigType("yaml")         //设置文件的类型
	logs.Info("你激活的环境是：" + GetEnvInfo("env"))//---------------------新增代码
	config.SetConfigFile(path + "/conf/application-" + GetEnvInfo("env") + ".yaml")//---------------------新增代码
	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = config.Unmarshal(&global.Config); err != nil {
			fmt.Println(err)
		}
	})

	// 这里才是把yaml配置文件解析放入到Config对象的过程---map---config
	if err = config.Unmarshal(&global.Config); err != nil {
		fmt.Println(err)
	}

	// 打印文件读取出来的内容:
	keys := config.AllKeys()
	dataMap := make(map[string]interface{})
	for _, key := range keys {
		fmt.Println("yaml存在的key是: " + key)
		dataMap[key] = config.Get(key)
	}

	global.Yaml = dataMap

}

```

然后运行

你在配置里不切换不同环境。你就可以看到不同环境的文件加载。



### 04、方式2：使用代码的方式完成环境变量参数设置

```go
package main

import (
	"github.com/spf13/viper"
	"xkginweb/initilization"
)

func main() {
	//// 设置环境变量
	viper.SetDefault("env", "prod")//修改这里即可
	// 解析配置文件
	initilization.InitViper()
	// 初始化日志 开发的时候建议设置成：debug ，发布的时候建议设置成：info/error
	// info --- console + file
	// error -- file
	initilization.InitLogger("debug")
	// 初始化中间 redis/mysql/mongodb
	initilization.InitMySQL()
	// 初始化缓存
	//initilization.InitRedis()
	// 初始化本地缓存
	initilization.InitCache()
	// 定时器
	// 并发问题解决方案
	// 异步编程
	// 初始化路由
	initilization.RunServer()
}

```

这样就替换golang的启动配置配置环境变量。更方便一些。





# 命令行参数的价值和意义（推进）

## 1：分析

通过上面你会发现可==设置参数到环境变量中==，然后通过os/viper提供的方法可以读取环境变量的参数，然后让程序代码可以跟随你设置环境变量的参数来加载不同的配置文件信息。确实可以解决加载不同环境的目录

但是同时也增加一个问题：你可能会把系统环境参数搞得乱七八糟。甚至搞坏。因为在真正得服务器上肯定没有golang项目环境配置隔离机制。在服务器上设定得环境变量参数都是共享得。

## 2：假设

如果我打了一个包，并且可以随时在启动可执行文件得时候来覆盖代码中变量参数，那么不完美么？

1：在代码我们先设定一个参数值，并且给默认值

```sh
viper.SetDefault("env", "prod")
```

2：如果在启动时候，增加一个参数，就默认值覆盖掉

```sh
# 使用默认值进行启动
./xinginweb
# 使用指定环境得方式进行启动
./xinginweb --env=test --xxxx=bbbb
```

```go
env := getEnv("env","dev") 第一个参数：env 是key名字，第二参数是：默认值
viper.SetDefault("env", env)
```

3：那么完美了



## 03、实现步骤

1:  使用flag组件

```go
package main

import (
	"flag"
	"fmt"
)

func main() {

	//2:go程序去读取命令行参数
	// --goenv=test
	var env string
	flag.Parse()
	flag.StringVar(&env, "goenv", "dev", "环境标识")
	fmt.Println(env)
}

```

上面代码得含义是指：如果未来我执行exe文件或者linux得可执行文件得时候，如果携带命令行参数，就会生效，如何指定呢？

```sh
./xkginweb.exe
./xkginweb
如果不指定：就使用默认值dev赋值给env这个变量

./xkginweb.exe --goenv=test
./xkginweb --goenv=prod
如果指定：就使用goenv指定值覆盖默认值，然后赋值给env这个变量
```

## 2：疑问

1：我怎么把上面得代码测试出来呢。难道为了学习就打个包。其实不需要，在golang开发工具下又一个配置也完成命令行参数得设定。

![image-20230912213723985](images/image-20230912213723985.png)

支持得写法如下：

通过以上两种方法定义命令行flag参数后，需要通过调用flag.Parse()来对命令行参数进行解析。
支持的命令行参数格式有一下几种：

- -flag xxx (使用空格，一个 - 符号）
- --flag xxx (使用空格，两个 - 符号)
- -flag=xxx （使用等号， 一个 - 符号）
- --flag = xxx (使用等号， 两个- 符号)

其中，布尔类型的参数必须用等号的方式指定。
flag在解析第一个非flag参数之前停止，或者在终止符"-"之后停止。

参考网站：https://studygolang.com/articles/20370

完整代码

```go
package main
import (
    "fmt"
    "flag"
    "time"
)
func main(){
    var name string
    var age int
    var married bool
    var delay time.Duration
    flag.StringVar(&name,"name","张三","姓名")
    flag.IntVar(&age,"age",18,"年龄")
    flag.BoolVar(&married,"married",false,"婚否")
    flag.DurationVar(&delay, "d", 0, "延迟的时间间隔")
    flag.Parse()
    fmt.Println(name,age,married,delay)
    fmt.Println(flag.Args())
    fmt.Println(flag.NArg())
    fmt.Println(flag.NFlag())

}
```

参考网站：https://studygolang.com/articles/20370



完整得代码



## 小结

os这个组件库，即可用读取到环境变量参数，也读取命令行参数



# Dokcer发布go项目



## 01、Docker的入门参考

狂神docker课程：https://www.bilibili.com/video/BV1og4y1q7M4/

## 02、Docker的安装

```sh
（1）yum 包更新到最新
> yum update

（2）安装需要的软件包， yum-util 提供yum-config-manager功能，另外两个是devicemapper驱动依赖的
> yum install -y yum-utils device-mapper-persistent-data lvm2

（3）设置yum源为阿里云
> yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

（4）安装docker
> yum install docker-ce -y

（5）安装后查看docker版本
> docker -v
```

使用阿里云的镜像加速器，这样后续docker在下载镜像的会快很多。

![image-20230915213150745](images/image-20230915213150745.png)

![image-20230915213225838](images/image-20230915213225838.png)

```sh
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": ["https://0wrdwnn6.mirror.aliyuncs.com"]
}
EOF
sudo systemctl daemon-reload
sudo systemctl restart docker
```

如何证明安装成功

```sh
# 启动状态
systemctl status docker ---绿色
# 查看版本
docker -v 
# 查看安装信息
docker info
# 查看镜像
docker images
# 查看容器
docker ps -a
docker ps 
```



## 03、Docker的核心概念

- 宿主机
  - 当前docker安装的系统就是docker宿主
  - docker是宿主机微型机，
- 镜像
  - 可行文件，安装文件—-redis.tar.gz——-redis镜像
- 容器
  - 你安装以后服务 —- redis server—-redis容器服务



### docker产生的背景

- 电脑大部分的内存和cpu都是很闲置，特别是配置高的服务器，那么一般的隔离都采用ip隔离（一个服务器 32c、64g）,然后在这个电脑中安装虚拟机然后进行N个系统分配，分配网关 ，分配Ip 
- 然后要部署一个项目，就必须把对应环境都要安装一次(java,mysql,redis,nginx) — 20-web节点



## go项目的发布和部署

1： 准备一个docker环境

2：准备一个go项目

3：定义一个`Dockerfile`

```dockerfile
# 依赖环境
FROM golang:1.20-alpine
# 复制执行文件到容器的根目录下
COPY build/testdocker ./
# 置顶容器服务的端口
EXPOSE 8080
# 赋予权限
RUN chmod +x ./testdocker
# 当你在执行docker run的时候会去启动你./testdocker
ENTRYPOINT [ "./testdocker" ]

```

基础环境构建

```dockerfile
# 基础镜像
FROM alpine:3.12
# 维护者
MAINTAINER frank
# 复制执行文件到容器的根目录下
COPY build/testdocker ./
# 置顶容器服务的端口
EXPOSE 8080
# 赋予权限
RUN chmod +x ./testdocker
# 当你在执行docker run的时候会去启动你./testdocker
ENTRYPOINT [ "./testdocker" ]
```

4: 上传项目可执行文件和Dockerfile到服务器

目录关系

![image-20230915223915439](images/image-20230915223915439.png)



5: 然后执行docker build 

```sh
 docker build -t testdocker:1.0 .
```

6： 然后查看打包后的镜像

```sh
docker images
```

如下：

![image-20230915223957646](images/image-20230915223957646.png)

7:  启动一个go容器服务

```sh
# 这个占用的启动，如果ctrl+c服务会停止，你需要手动在启动一次，docker start abfeb2edbef6
# 这种启动的方式一般在调试的时候会使用到。
docker run -it --name gowebpro -p 8080:8080 testdocker:1.0
# 守护方式启动。那么久不方便查看到日志信息，如果你查看日志信息久必须使用 docker logs -f abfeb2edbef6
docker run -di --name gowebpro -p 8080:8080 testdocker:1.0
docker run -di --name gowebpro81 -p 8081:8080 testdocker:1.0
docker run -di --name gowebpro82 -p 8082:8080 testdocker:1.0
```

使用

```sh
docker ps -a
```

查看容器的启动





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



## go卸载安装三方库

### 1、查看包

```go
go list -m all
```



### 2、安装

```go
go get [包名]
```



### 3、卸载

```go
go clean -i [包名]
```



## OS库



## golang防缓存击穿神器【singleflight】

[golang防缓存击穿神器【singleflight】 - 知乎 (zhihu.com)](https://zhuanlan.zhihu.com/p/382965636)

```go
go get golang.org/x/sync/singleflight
```



## Zap日志库

[Zap Logger · Go语言中文文档 (topgoer.com)](https://www.topgoer.com/项目/log/ZapLogger.html)





## Gorm——Unscoped

在 Golang 的 ORM 库 GORM 中，`Unscoped` 是一个非常有用的方法，它可以用来包括通常被 GORM 软删除（Soft Delete）的记录在内的查询操作中。在 GORM 中，软删除是一个标记记录为删除的特性，而不是从数据库中实际删除它。当一个记录被软删除时，它仍然存在于数据库中，但 GORM 的正常查询操作会默认忽略这些记录。

当你想要执行不忽略软删除记录的操作时，比如查询或真正的删除，你可以使用 `Unscoped` 方法。

### `Unscoped` 的用法：

1. 查询包括软删除记录的操作：

```go
// 假设有一个名为 User 的模型，它包含软删除功能
var users []User
// 查询所有用户，包括被软删除的用户
db.Unscoped().Find(&users)
```

2. 永久删除记录：

```go
// 删除一个具体的记录，而不是软删除
db.Unscoped().Delete(&user)
```

在上面的例子中，如果不使用 `Unscoped`，`Delete` 方法将会执行软删除（如果模型被配置为支持软删除），即它只是设置了记录的 `deleted_at` 字段（或者是你在模型中指定的软删除字段）。使用 `Unscoped` 后，`Delete` 方法将会从数据库中永久删除这个记录。

### `Unscoped` 的原理：

在 GORM 中，每个查询链都可以包含多个方法，这些方法可以修改 GORM 内部的查询构建器的状态。当你调用 `Unscoped` 方法时，它会告诉 GORM 在接下来的操作中忽略软删除的过滤条件。这样，即使记录有 `deleted_at` 字段被设置了值，查询也会返回这些记录。

### 注意事项：

- 使用 `Unscoped` 时要谨慎，因为它会返回或删除所有相关记录，包括那些被标记为删除的记录。
- 如果你的模型没有使用 GORM 的软删除功能，`Unscoped` 方法不会有任何影响。
- 在使用 `Unscoped` 进行删除操作时，请确保你的操作是正确的，因为这将会永久从数据库中移除记录，这是一个不可逆的操作。



## Gorm —— 默认零值不处理

参考文献：https://gorm.io/zh_CN/docs/update.html

在 GORM 中，当你使用 `Create` 或 `Save` 方法时，如果结构体中的字段为 Go 语言的零值（比如 `int` 类型的 `0`，`string` 类型的空字符串 `""`，布尔类型的 `false` 等），GORM 默认不会将这些零值写入数据库。这是因为 GORM 不能确定这些零值是你有意设置的，还是仅仅因为它们是 Go 语言中类型的默认零值。

如果你需要将零值写入数据库，你可以使用几种方法来告诉 GORM 你的意图：

### 1. 使用指针类型

将字段的类型声明为指针类型，这样字段就可以是 `nil`（不设置值）或者非 `nil`（设置了值，即使是零值）。

```go
type User struct {
  gorm.Model
  Name     string
  Age      *int
  IsActive *bool
}

// 当你想要设置零值时
age := 0
isActive := false
user := User{Name: "John", Age: &age, IsActive: &isActive}
db.Create(&user)
```

使用指针类型时，如果字段值为 `nil`，GORM 将会忽略该字段，如果字段值为非 `nil`（包括零值），GORM 将会保存该值。

### 2. 使用 `sql.Null*` 类型

对于基本类型，你可以使用 SQL 包中的 `NullString`、`NullInt64`、`NullBool` 等类型，这些类型可以明确表示值是否存在。

```go
import "database/sql"

type User struct {
  gorm.Model
  Name     string
  Age      sql.NullInt64
  IsActive sql.NullBool
}

user := User{Name: "John", Age: sql.NullInt64{Int64: 0, Valid: true}, IsActive: sql.NullBool{Bool: false, Valid: true}}
db.Create(&user)
```

### 3. 使用 `gorm.Expr`

使用 `gorm.Expr` 来构建一个表达式，这样 GORM 就会在 SQL 语句中使用这个表达式的值。

```go
db.Model(&user).Update("age", gorm.Expr("?", 0))
```

### 4. 使用 `Select` 方法

在 `Create` 或 `Save` 方法中使用 `Select`，明确指定要更新的字段。

```go
db.Select("Age", "IsActive").Create(&user)
```

### 5. 配置 GORM 使用默认值

你可以为 GORM 模型的字段设置标签 `default`，以便在插入记录时使用指定的默认值。

```go
type User struct {
  gorm.Model
  Name     string
  Age      int `gorm:"default:0"`
  IsActive bool `gorm:"default:false"`
}

user := User{Name: "John"}
db.Create(&user)
```

在这个例子中，即使 `Age` 和 `IsActive` 是零值，由于设置了默认值标签，GORM 也会在插入时使用这些默认值。

选择哪种方法取决于你的具体需求以及你希望如何处理零值。要注意的是，每种方法都有其适用场景和潜在的陷阱，因此在决定使用哪种方法之前，请确保你了解每种方法的行为以及如何在你的应用程序中正确使用它们。

## Jebrains编辑器 —— 关闭ESlint检查

关闭Vue项目中ESlint检查

file ==> settings ==> 搜索ESlint，去掉Enable的勾





## 结构体转化为map

[golang如何优雅的将struct转换为map - 掘金 (juejin.cn)](https://juejin.cn/post/7187042947618046009)

[golang常用库之mapstructure包 | 多json格式情况解析、GO json 如何转化为 map 和 struct、Go语言结构体标签（Struct Tag）_go mapstructure-CSDN博客](https://blog.csdn.net/inthat/article/details/127121728)

https://blog.csdn.net/Alen_xiaoxin/article/details/124078796



`mapstructure` 是一个 Go 语言库，用于将通用的 map 值解码到相应的 Go 结构体中。这在处理动态格式的数据时非常有用，如 JSON 解码到 `map[string]interface{}` 后，再将其转换到一个更严格类型化的结构体。这种方式在处理配置文件或者网络通信时尤其常见。

### 使用 `mapstructure` 的基本步骤：

1. **安装**：首先，你需要安装 `mapstructure` 库。

```sh
go get github.com/mitchellh/mapstructure
```

2. **定义结构体**：定义一个 Go 结构体，它的字段应该与你想从 map 中解码的数据相匹配。

```go
type Person struct {
    Name    string
    Age     int
    Emails  []string
    Extra   map[string]string
}
```

3. **解码数据**：使用 `mapstructure.Decode` 函数将 map 数据解码到结构体中。

```go
package main

import (
    "fmt"
    "github.com/mitchellh/mapstructure"
)

func main() {
    // 假定你有一个 map 数据结构
    inputData := map[string]interface{}{
        "Name":   "Alice",
        "Age":    25,
        "Emails": []string{"alice@example.com", "alice@another.com"},
        "Extra":  map[string]string{"Twitter": "@alice"},
    }

    var person Person

    // 解码 inputData 到 person 结构体中
    err := mapstructure.Decode(inputData, &person)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Printf("%+v\n", person)
}
```

### `mapstructure` 的高级用法：

- **自定义标签**：你可以使用结构体字段标签来指定 map 中的键如何映射到结构体的字段。

```go
type Person struct {
    Name    string `mapstructure:"name"`
    Age     int    `mapstructure:"age"`
    Emails  []string `mapstructure:"emails"`
    Extra   map[string]string `mapstructure:"extra"`
}
```

- **解码钩子**：`mapstructure` 提供了解码钩子功能，允许你在解码过程中自定义转换数据的逻辑。

```go
config := &mapstructure.DecoderConfig{
    DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
    Result:     &person,
}

decoder, err := mapstructure.NewDecoder(config)
if err != nil {
    // 处理错误
}
err = decoder.Decode(inputData)
```

- **弱类型输入**：`mapstructure` 可以处理弱类型的输入，例如 JSON 解码可能会将所有数字解码为 `float64`，而不管它们在 Go 结构体中是整型还是浮点型。

```go
config := &mapstructure.DecoderConfig{
    WeaklyTypedInput: true,
    Result:           &person,
}

decoder, err := mapstructure.NewDecoder(config)
if err != nil {
    // 处理错误
}
err = decoder.Decode(inputData)
```

### `mapstructure` 的作用：

- **灵活性**：它允许你从灵活的 map 结构中解码数据到严格类型化的结构体，提高了类型安全性和易用性。
- **配置管理**：在处理配置文件时，`mapstructure` 可以简化从 `map[string]interface{}` 到具体配置结构体的转换过程。
- **通用解码**：它可以用于解码来自不同数据源（如 JSON、YAML 或环境变量）的数据。

总之，`mapstructure` 是一个处理动态数据和将其映射到固定结构的有用工具，特别是在配置管理和数据解析的上下文中。它的灵活性和强大的定制选项使得它在 Go 社区中非常受欢迎。



# 流程梳理

1、model数据库表建设

2、初始化数据库 -> 注册表

3、
