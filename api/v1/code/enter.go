/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 18:16
 */
package code

import (
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/service"
)

type WebApiGroup struct {
	CodeApi
}

// 公共实例---服务共享
var (
	sysUserService = service.ServiceGroupApp.SyserviceGroup.SysUserService
	jwtService     = jwtgo.JwtService{}
)
