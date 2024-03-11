/*
* @Author: 梦无矶小仔
* @Date: 2024/3/11 14:24
 */
package sys

type WebRouterGroup struct {
	SysMenusRouter
	SysUsersRouter
	SysApisRouter
	SysRolesRouter
	SysUserRolesRouter
	SysRoleMenusRouter
	SysRoleApisRouter
}
