/*
* @Author: 梦无矶小仔
* @Date: 2024/3/11 14:25
 */
package sys

import (
	"github.com/gin-gonic/gin"
	v1 "xz-go-frame/api/v1"
)

type SysUserRolesRouter struct{}

func (router *SysUserRolesRouter) InitSysUserRolesRouter(Router *gin.RouterGroup) {
	sysUserRolesApi := v1.WebApiGroupApp.Sys.SysUserRolesApi
	// 用组定义--（推荐）
	sysMenusRouter := Router.Group("/sys")
	{
		// 保存
		sysMenusRouter.POST("/user/role/save", sysUserRolesApi.SaveData)
		// 查询明细 /user/get/1/xxx
		sysMenusRouter.POST("/user/role/list", sysUserRolesApi.SelectUserRoles)
	}
}
