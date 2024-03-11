/*
@Author: 梦无矶小仔
@Date:   2024/1/12 16:46
*/
package login

import (
	"github.com/gin-gonic/gin"
	v1 "xz-go-frame/api/v1"
)

// 登录路由
type LoginRouter struct{}

func (router *LoginRouter) InitLoginRouter(Router *gin.RouterGroup) {
	//loginApi := login.LoginApi{}
	// 单个定义
	//Router.GET("/login/toLogin", loginApi.ToLogined)
	//Router.GET("/login/toReg", loginApi.ToLogined)
	//Router.GET("/login/forget", loginApi.ToLogined)

	loginApi := v1.WebApiGroupApp.Login.LoginApi
	// 用组定义 ---》 推荐
	loginRouter := Router.Group("/login")
	{
		loginRouter.POST("/toLogin", loginApi.ToLogined)
	}
}
