/*
@Author: 梦无矶小仔
@Date:   2024/1/29 18:45
*/
package value_object

type SysRolesVo struct {
	ID       uint   `json:"id"`
	RoleName string `json:"roleName"` // 角色名
	RoleCode string `json:"roleCode"` // 角色代号
}
