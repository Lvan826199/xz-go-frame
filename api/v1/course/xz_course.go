/*
* @Author: 梦无矶小仔
* @Date: 2024/3/11 14:13
 */
package course

import "github.com/gin-gonic/gin"

type CourseApi struct{}

// 查询video
func (api *CourseApi) FindVideos(c *gin.Context) {

	// service new
	// model new
}

// 获取明细
func (api *CourseApi) GetByID(c *gin.Context) {
	// 绑定参数用来获取/:id这个方式
	//id := c.Param("id")
	// 绑定参数 ?ids=1111
}
