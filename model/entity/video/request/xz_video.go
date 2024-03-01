/*
@Author: 梦无矶小仔
@Date:   2024/3/1 17:14
*/
package request

// 接受参数
type XkVideoReq struct {
	PageNum     int    `form:"pageNum" json:"pageNum"`
	PageSize    int    `form:"pageSize" json:"pageSize"`
	CategoryId  int    `form:"categoryId" json:"categoryId"`
	CategoryCid int    `form:"categoryCid" json:"categoryCid"`
	Keyword     string `form:"keyword" json:"keyword"`
	StartTime   string `form:"startTime" json:"startTime"`
	EndTime     string `form:"endTime" json:"endTime"`
	Status      int8   `form:"status" json:"status"`
}
