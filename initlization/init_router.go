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
	time.Sleep(10 * time.Microsecond)

	s2 := s.ListenAndServe().Error()
	fmt.Println("服务启动完毕", s2)

}
