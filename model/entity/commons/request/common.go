/*
@Author: 梦无矶小仔
@Date:   2024/1/30 16:02
*/
package request

// PageInfo 分页输入参数结构
type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页条数
	Keyword  string `json:"keyword" form:"keyword"`   // 搜索关键字
}

type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}
