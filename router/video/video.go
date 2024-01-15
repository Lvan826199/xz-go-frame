package video

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/api/v1/video"
)

type VideoRouter struct {
}

func (videoRouter *VideoRouter) InitVideoRouter(group *gin.RouterGroup) {
	// 帖子路由
	videoApi := video.Video{}
	videoGroup := group.Group("video")
	{

		videoGroup.GET("find", videoApi.FindVideos)
		videoGroup.GET("get", videoApi.GetByID)
		//videoGroup.GET("/index", videoApi.VideoIndex)
		//videoGroup.GET("/get/:id", videoApi.GetVideoDetailById)
	}
}
