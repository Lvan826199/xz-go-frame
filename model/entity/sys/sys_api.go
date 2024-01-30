/*
@Author: 梦无矶小仔
@Date:   2024/1/30 11:40
*/
package sys

import "xz-go-frame/global"

type SysApis struct {
	global.GVA_MODEL
	Title       string `json:"title" gorm:"comment:api路径名称"`          // api路径
	Path        string `json:"path" gorm:"comment:api路径"`             // api路径
	Description string `json:"description" gorm:"comment:api中文描述"`    // api中文描述
	MenuId      uint   `json:"menuId" gorm:"comment:隶属于菜单的api"`       // api组
	MenuName    string `json:"menuName" gorm:"comment:隶属于菜单的api名字"`   // api组
	Method      string `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
}

func (s *SysApis) TableName() string {
	return "sys_apis"
}
