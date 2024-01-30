/*
@Author: 梦无矶小仔
@Date:   2024/1/30 11:41
*/
package sys

import "xz-go-frame/global"

type SysRoles struct {
	global.GVA_MODEL
	RoleName string `json:"roleName" gorm:"comment:角色名"`  // 角色名
	RoleCode string `json:"roleCode" gorm:"comment:角色代号"` // 角色代号
}

func (s *SysRoles) TableName() string {
	return "sys_roles"
}
