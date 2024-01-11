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
