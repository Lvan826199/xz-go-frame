/*
@Author: 梦无矶小仔
@Date:   2024/3/1 16:49
*/
package request

import "xz-go-frame/model/entity/commons/request"

type BbsPageInfo struct {
	request.PageInfo
	CategoryId int8 `json:"categoryId" gorm:"not null;default:0;comment:文章分类ID"`
	Status     int8 `json:"status" gorm:"not null;default:1;comment:0 未发布 1 发布"`
}

type BbsCategorySaveReq struct {
	ID          uint   `json:"id" form:"id"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Parent_id   uint   `json:"parentId" form:"parentId"`
	Sorted      int8   `json:"sorted" form:"sorted"`
	Status      int8   `json:"status" form:"status"`
	IsDelete    int8   `json:"isDelete" form:"isDelete"`
}

type StatusReq struct {
	ID    uint   `json:"id"`
	Value int8   `json:"value"`
	Field string `json:"field"`
}
