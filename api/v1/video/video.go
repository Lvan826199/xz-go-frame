package video

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/commons/response"
)

type Video struct {
}

// 查询video
func (videoController *Video) FindVideos(c *gin.Context) {
	claims, _ := c.Get("claims")
	customClaims := claims.(*jwtgo.CustomClaims)
	fmt.Println(customClaims.UserId)
	fmt.Println(customClaims.Username)
	response.Ok("success FindVideos", c)
}

// 获取视频明细
func (videoController *Video) GetByID(c *gin.Context) {
	// 绑定参数用来获取/:id 这个方式
	// id := c.Param("id")
	// 绑定参数 ?ids = 1111
	claims, _ := c.Get("claims")
	customClaims := claims.(*jwtgo.CustomClaims)
	fmt.Println(customClaims.UserId)
	fmt.Println(customClaims.Username)
	response.Ok("success GetByID", c)
}

//// 首页处理
//func (e *Video) VideoIndex(c *gin.Context) {
//	username, _ := c.Get("username")
//	// 可以获取到login放入session的数据
//	c.JSON(http.StatusOK, "我是video的首页: "+username.(string))
//}
//
//// 获取明细
//
//func (e *Video) GetVideoDetailById(c *gin.Context) {
//	username, _ := c.Get("username")
//	// 可以获取到login放入session的数据
//	param := c.Param("id")
//	c.JSON(http.StatusOK, "我是VideoApi的名字,参数:"+param+"  ： "+username.(string))
//}
