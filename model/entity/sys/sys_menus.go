/*
@Author: 梦无矶小仔
@Date:   2024/1/30 11:40
*/
package sys

import "xz-go-frame/global"

type SysMenus struct {
	global.GVA_MODEL
	ParentId  uint   `json:"parentId" gorm:"comment:父菜单ID"`      // 父菜单ID
	Path      string `json:"path" gorm:"comment:路由path"`         // 路由path
	Title     string `json:"title" gorm:"comment:菜单名称"`          // 菜单名称
	Name      string `json:"name" gorm:"comment:路由name 用于国际化处理"` // 路由name 用于国际化处理
	Hidden    bool   `json:"hidden" gorm:"comment:是否在列表隐藏"`      // 是否在列表隐藏
	Component string `json:"component" gorm:"comment:对应前端文件路径"`  // 对应前端文件路径
	Sort      int    `json:"sort" gorm:"comment:排序标记"`           // 排序标记
	Icon      string `json:"icon" gorm:"comment:菜单图标"`           // 菜单图标
	// 忽略该字段，- 表示无读写，-:migration 表示无迁移权限，-:all 表示无读写迁移权限
	Children []*SysMenus `gorm:"-" json:"children"`
	//TopObj   *SysMenus   `gorm:"-" json:"-"`
}

func (s *SysMenus) TableName() string {
	return "sys_menus"
}
