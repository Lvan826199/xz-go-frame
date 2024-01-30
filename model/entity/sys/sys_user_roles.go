/*
@Author: 梦无矶小仔
@Date:   2024/1/30 11:41
*/
package sys

type SysUserRoles struct {
	UserId uint `gorm:"column:user_id"`
	RoleId uint `gorm:"column:role_id"`
}

func (s *SysUserRoles) TableName() string {
	return "sys_user_roles"
}
