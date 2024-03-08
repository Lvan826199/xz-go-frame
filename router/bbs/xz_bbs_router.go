package bbs

import (
	"github.com/gin-gonic/gin"
	v1 "xz-go-frame/api/v1"
)

type XkBbsRouter struct{}

func (r *XkBbsRouter) InitXkBbsRouter(Router *gin.RouterGroup) {

	xkBbsApi := v1.WebApiGroupApp.Bbs.XkBbsApi
	// 这个路由多了一个对对post，put请求的中间件处理，而这个中间件做了一些对post和put的参数的处理和一些公共信息的处理
	router := Router.Group("bbs")
	{
		// 保存
		router.POST("save", xkBbsApi.CreateXkBbs)
		// 更新
		router.POST("update", xkBbsApi.UpdateXkBbs)
		// 更新
		router.DELETE("delete/:id", xkBbsApi.DeleteById)
		// 分页查询
		router.POST("page", xkBbsApi.LoadXkBbsPage)
		// 明细查询
		router.GET("get", xkBbsApi.GetXkBbs)
	}
}
