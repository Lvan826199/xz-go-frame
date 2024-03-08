/*
* @Author: 梦无矶小仔
* @Date: 2024/3/7 17:12
 */
package sys

import "xz-go-frame/service"

type WebApiGroup struct {
	SysMenuApi
	SysUsersApi
	SysRolesApi
	SysApisApi
	SysUserRolesApi
	SysRoleMenusApi
	SysRoleApisApi
}

var (
	sysMenuService      = service.ServiceGroupApp.SyserviceGroup.SysMenusService
	sysUserService      = service.ServiceGroupApp.SyserviceGroup.SysUserService
	sysRolesService     = service.ServiceGroupApp.SyserviceGroup.SysRolesService
	sysApisService      = service.ServiceGroupApp.SyserviceGroup.SysApisService
	sysUserRolesService = service.ServiceGroupApp.SyserviceGroup.SysUserRolesService
	sysRoleApisService  = service.ServiceGroupApp.SyserviceGroup.SysRoleApisService
	sysRoleMenusService = service.ServiceGroupApp.SyserviceGroup.SysRoleMenusService
)
