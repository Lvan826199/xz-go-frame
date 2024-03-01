/*
@Author: 梦无矶小仔
@Date:   2024/3/1 16:49
*/
package request

import "xz-go-frame/model/entity/commons/request"

type BbsCategoryPageInfo struct {
	request.PageInfo
	Status int8 `form:"status" json:"status"`
}
