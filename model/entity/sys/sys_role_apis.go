/*
@Author: 梦无矶小仔
@Date:   2024/1/30 11:41
*/
package sys

type SysRoleApis struct {
	RoleId uint `gorm:"comment:角色ID"`
	ApiId  uint `gorm:"comment:ApiID"`
}

func (s *SysRoleApis) TableName() string {
	return "sys_role_apis"
}
