/*
@Author: 梦无矶小仔
@Date:   2024/1/30 16:01
*/
package request

import "xz-go-frame/model/entity/commons/request"

// UserStatePageInfo 根据年月查询用户分页
type UserStatePageInfo struct {
	request.PageInfo
	Ym string `json:"ym" form:"ym"`
}
