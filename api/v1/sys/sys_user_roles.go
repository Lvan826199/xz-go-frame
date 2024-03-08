/*
* @Author: 梦无矶小仔
* @Date: 2024/3/7 18:48
 */
package sys

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"xz-go-frame/commons/response"
	"xz-go-frame/model/entity/sys"
)

type SysUserRolesApi struct{}

// 用户授权角色
func (api *SysUserRolesApi) SaveData(c *gin.Context) {
	// 1: 第一件事情就准备数据的载体
	type SysUserRoleContext struct {
		RoleIds string `json:"roleIds"`
		UserId  uint   `json:"userId"`
	}

	// 绑定上下文参数
	var params SysUserRoleContext
	if err := c.ShouldBindJSON(&params); err != nil {
		// 如果参数注入失败或者出错就返回接口调用这。出错了.
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 开始对授权的角色进行分割
	roleIdsSplit := strings.Split(params.RoleIds, ",")
	// 准备一个用户角色中间表的切片对象
	var sysUserRoles []*sys.SysUserRoles
	// 开始变量每个角色信息绑定切片对象
	for _, roleId := range roleIdsSplit {
		parseUint, _ := strconv.ParseUint(roleId, 10, 64)
		sysUserRole := sys.SysUserRoles{}
		sysUserRole.UserId = params.UserId
		sysUserRole.RoleId = uint(parseUint)
		sysUserRoles = append(sysUserRoles, &sysUserRole)
	}

	// 开始进行批量保存
	if err := sysUserRolesService.SaveSysUserRoles(params.UserId, sysUserRoles); err != nil {
		response.FailWithMessage("创建失败", c)
		return
	}
	// 如果保存成功，就返回创建创建成功
	response.Ok("创建成功", c)
}

/*
*
查询用户授权
*/
func (api *SysUserRolesApi) SelectUserRoles(c *gin.Context) {
	roles, _ := sysUserRolesService.SelectUserRoles(c.GetUint("userId"))
	response.Ok(roles, c)
}
