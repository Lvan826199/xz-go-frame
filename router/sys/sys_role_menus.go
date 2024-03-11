/*
* @Author: 梦无矶小仔
* @Date: 2024/3/11 14:25
 */
package sys

import (
	"github.com/gin-gonic/gin"
	v1 "xz-go-frame/api/v1"
)

type SysRoleMenusRouter struct{}

func (r *SysRoleMenusRouter) InitSysRoleMenusRouter(Router *gin.RouterGroup) {
	sysRoleMenusApi := v1.WebApiGroupApp.Sys.SysRoleMenusApi
	// 用组定义--（推荐）
	router := Router.Group("/sys")
	{
		// 保存
		router.POST("/role/menu/save", sysRoleMenusApi.SaveData)
		// 查询明细 /user/get/1/xxx
		router.POST("/role/menu/list", sysRoleMenusApi.SelectRoleMenus)
	}
}
