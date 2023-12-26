package bbs

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/api/bbs"
)

type BbsRouter struct {
}

func (bBsRouter BbsRouter) InitBBSRouter(group *gin.RouterGroup) {
	// 帖子路由
	bbsApi := bbs.BbsApi{}
	bbsApiGroup := group.Group("bbs")
	{
		// 函数封装
		bbsApiGroup.GET("/indexx", bbsApi.BbsIndex)
		bbsApiGroup.GET("/get/:id", bbsApi.GetBbsDetailById)
	}
}
