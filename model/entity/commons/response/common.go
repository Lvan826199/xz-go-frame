/*
@Author: 梦无矶小仔
@Date:   2024/1/30 16:02
*/

package response

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}
