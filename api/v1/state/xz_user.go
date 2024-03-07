/*
* @Author: 梦无矶小仔
* @Date: 2024/3/7 16:38
 */
package state

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/commons/response"
	resp "xz-go-frame/model/entity/commons/response"
	"xz-go-frame/model/entity/state/request"
)

type UserStateApi struct {
}

/*
**
统计某年度用户注册量量信息
*/
func (api *UserStateApi) UserRegState(ctx *gin.Context) {
	// 获取查询的年份
	year := ctx.Query("year")
	// 开始执行
	states, err := userStatService.UserRegStateData(year)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Ok(states, ctx)
}

/*
**
统计某年度用户注册量量信息-明细信息
*/
func (api *UserStateApi) FindUserRegStateDetail(ctx *gin.Context) {
	// 创建一个分页对象
	var pageInfo request.UserStatePageInfo
	// 把前端json的参数传入给PageInfo
	err := ctx.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	// 开始执行
	list, total, err := userStatService.FindUserRegStateDetail(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Ok(resp.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, ctx)
}
