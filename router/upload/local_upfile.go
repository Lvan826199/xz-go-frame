/*
* @Author: 梦无矶小仔
* @Date: 2024/3/11 14:22
 */
package upload

import (
	"github.com/gin-gonic/gin"
	v1 "xz-go-frame/api/v1"
)

type LocalUploadRouter struct{}

func (r *LocalUploadRouter) InitLocalUploadRouter(Router *gin.RouterGroup) {

	localUploadApi := v1.WebApiGroupApp.Upload.LocalUploadApi
	// 这个路由多了一个对对post，put请求的中间件处理，而这个中间件做了一些对post和put的参数的处理和一些公共信息的处理
	router := Router.Group("local") //.Use(middleware.OperationRecord())
	{
		// 保存
		router.POST("upload/file", localUploadApi.UploadFile)
		router.POST("upload/wangeditor", localUploadApi.UploadFileWangEditor)
	}

}
