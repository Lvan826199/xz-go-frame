/*
* @Author: 梦无矶小仔
* @Date: 2024/3/11 14:21
 */
package state

import (
	"github.com/gin-gonic/gin"
	v1 "xz-go-frame/api/v1"
)

type UserStateRouter struct{}

func (r *UserStateRouter) InitUserStateRouter(Router *gin.RouterGroup) {

	userStateApi := v1.WebApiGroupApp.State.UserStateApi
	router := Router.Group("state") //.Use(middleware.OperationRecord())
	{
		// 统计某年的用户注册量
		router.GET("user/reg", userStateApi.UserRegState)
		// 统计某年的用户注册量--明细信息
		router.POST("user/detail", userStateApi.FindUserRegStateDetail)
	}

}
