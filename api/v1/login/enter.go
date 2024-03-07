/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 18:30
 */
package login

import (
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/service"
)

type WebApiGroup struct {
	LoginApi
	LogOutApi
}

// 公共实例---服务共享
var (
	sysUserService      = service.ServiceGroupApp.SyserviceGroup.SysUserService
	sysMenuService      = service.ServiceGroupApp.SyserviceGroup.SysMenusService
	sysUserRolesService = service.ServiceGroupApp.SyserviceGroup.SysUserRolesService
	sysRoleApisService  = service.ServiceGroupApp.SyserviceGroup.SysRoleApisService
	sysRoleMenusService = service.ServiceGroupApp.SyserviceGroup.SysRoleMenusService
	jwtService          = jwtgo.JwtService{}
)
