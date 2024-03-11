/*
* @Author: 梦无矶小仔
* @Date: 2024/3/11 14:25
 */
package sys

import (
	"github.com/gin-gonic/gin"
	v1 "xz-go-frame/api/v1"
)

type SysRoleApisRouter struct{}

func (r *SysRoleApisRouter) InitSysRoleApisRouter(Router *gin.RouterGroup) {
	sysRoleApisApi := v1.WebApiGroupApp.Sys.SysRoleApisApi
	// 用组定义--（推荐）
	router := Router.Group("/sys")
	{
		// 保存
		router.POST("/role/api/save", sysRoleApisApi.SaveData)
		// 角色改变
		router.POST("/role/api/change", sysRoleApisApi.ChangeRoleIdMenus)
		// 查询明细 /user/get/1/xxx
		router.POST("/role/api/list", sysRoleApisApi.SelectRoleApis)
	}
}
