/*
@Author: 梦无矶小仔
@Date:   2024/3/1 17:14
*/
package response

// 响应参数
type XkVideoResp struct {
	PageNum  int         `form:"pageNum" json:"pageNum"`
	PageSize int         `form:"pageSize" json:"pageSize"`
	Total    int64       `json:"total"`
	List     interface{} `json:"list"`
}
