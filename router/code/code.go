/*
@Author: 梦无矶小仔
@Date:   2024/1/11 16:23
*/
package code

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/api/v1"
)

type CodeRouter struct{}

func (r *CodeRouter) InitCodeRouter(Router *gin.RouterGroup) {

	codeApi := v1.WebApiGroupApp.Code.CodeApi
	// 这个路由多了一个对对post，put请求的中间件处理，而这个中间件做了一些对post和put的参数的处理和一些公共信息的处理
	router := Router.Group("code")
	{
		router.GET("get", codeApi.CreateCaptcha)
		router.GET("verify", codeApi.VerifyCaptcha)
	}
}
