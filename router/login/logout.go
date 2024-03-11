/*
* @Author: 梦无矶小仔
* @Date: 2024/3/11 14:17
 */
package login

import (
	"github.com/gin-gonic/gin"
	v1 "xz-go-frame/api/v1"
)

// 登出路由
type LogoutRouter struct{}

func (r *LogoutRouter) InitLogoutRouter(Router *gin.RouterGroup) {
	logoutApi := v1.WebApiGroupApp.Login.LogOutApi
	// 用组定义--（推荐）
	router := Router.Group("/login")
	{
		router.POST("/logout", logoutApi.ToLogout)
	}
}
