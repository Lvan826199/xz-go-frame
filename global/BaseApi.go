/*
@Author: 梦无矶小仔
@Date:   2024/1/30 16:28
*/
package global

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type BaseApi struct{}

/*
string转uint类型
*/
func (api *BaseApi) StringToUnit(id string) uint {
	parseUint, _ := strconv.ParseUint(id, 10, 64)
	return uint(parseUint)
}

/*
获取当前登录的用户信息
*/

func (api *BaseApi) GetLoginUserId(c *gin.Context) uint {
	userId := c.GetUint("userId")
	return userId
}
