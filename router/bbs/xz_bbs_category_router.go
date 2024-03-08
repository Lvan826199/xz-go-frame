/*
* @Author: 梦无矶小仔
* @Date: 2024/3/8 17:36
 */
package bbs

import (
	"github.com/gin-gonic/gin"
	v1 "xz-go-frame/api/v1"
)

type BBSCategoryRouter struct{}

func (r *BBSCategoryRouter) InitBBSCategoryRouter(Router *gin.RouterGroup) {

	bbsCategoryApi := v1.WebApiGroupApp.Bbs.BbsCategoryApi
	// 这个路由多了一个对对post，put请求的中间件处理，而这个中间件做了一些对post和put的参数的处理和一些公共信息的处理
	router := Router.Group("bbscategory")
	{
		// 保存
		router.POST("save", bbsCategoryApi.CreateBbsCategory)
		// 更新
		router.POST("update", bbsCategoryApi.UpdateBbsCategory)
		// 更新状态
		router.POST("update/status", bbsCategoryApi.UpdateBbsCategoryStatus)
		// 删除
		router.DELETE("delete/:id", bbsCategoryApi.DeleteById)
		// 批量删除
		router.DELETE("deletes", bbsCategoryApi.DeleteByIds)
		// 分页查询
		router.POST("page", bbsCategoryApi.LoadBbsCategoryPage)
		// 明细查询
		router.GET("get", bbsCategoryApi.GetBbsCategory)
		// 明细查询
		router.GET("find", bbsCategoryApi.FindBbsCategory)
	}
}
