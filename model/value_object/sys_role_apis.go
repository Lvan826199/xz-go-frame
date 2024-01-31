/*
@Author: 梦无矶小仔
@Date:   2024/1/29 18:45
*/
package value_object

type SysApisVo struct {
	ID     uint   `json:"id"`
	Path   string `json:"path"`   // 路径
	Title  string `json:"title"`  // 标题
	MenuId string `json:"menuId"` // 菜单id
	Method string `json:"method"` // 请求方式
}
