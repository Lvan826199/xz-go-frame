/*
* @Author: 梦无矶小仔
* @Date: 2024/3/11 14:12
 */
package course

import (
	"github.com/gin-gonic/gin"
	v1 "xz-go-frame/api/v1"
)

type CourseRouter struct{}

func (r *CourseRouter) InitCourseRouter(Router *gin.RouterGroup) {

	courseApi := v1.WebApiGroupApp.Course.CourseApi
	// 这个路由多了一个对对post，put请求的中间件处理，而这个中间件做了一些对post和put的参数的处理和一些公共信息的处理
	router := Router.Group("course") //.Use(middleware.OperationRecord())
	{
		router.GET("find", courseApi.FindVideos)
		router.GET("get", courseApi.GetByID)
	}
}
