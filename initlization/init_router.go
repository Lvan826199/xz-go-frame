package initlization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"xz-go-frame/api/v1/login"
	"xz-go-frame/commons/filter"
	"xz-go-frame/middle"
	"xz-go-frame/router"
	"xz-go-frame/router/code"
)

func WebRouterInit() {
	// 初始化 gin 服务
	ginServer := gin.Default()
	// 提供服务组

	videoRouter := router.RouterWebGroupApp.Video.VideoRouter
	// 解决接口的跨域问题
	ginServer.Use(filter.Cors())

	loginApi := login.LoginApi{}
	// 验证码接口
	codeRouter := code.CodeRouter{}
	codeRouter.InitCodeRouter(ginServer)

	// 登录路由
	ginServer.GET("/login", loginApi.Login)

	// 首页路由
	AdminGroup := ginServer.Group("admin")
	AdminGroup.Use(middle.LoginInterceptor())
	{

		videoRouter.InitVideoRouter(AdminGroup)
	}

	// 启动HTTP服务,可以修改端口
	address := fmt.Sprintf(":%d", 8088)
	ginServer.Run(address)
}
