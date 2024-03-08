/*
* @Author: 梦无矶小仔
* @Date: 2024/3/7 18:47
 */
package sys

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gookit/validate"
	"strconv"
	"strings"
	"time"
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/commons/response"
	"xz-go-frame/global"
	"xz-go-frame/model/entity/sys"
	"xz-go-frame/utils"
)

type SysRoleApisApi struct {
	global.BaseApi
}

type SysRoleApisContext struct {
	RoleId  uint   `json:"roleId" validate:"required"`
	ApisIds string `json:"apisIds" validate:"required"`
}

// Messages 您可以自定义验证器错误消息
func (f SysRoleApisContext) Messages() map[string]string {
	return validate.MS{
		"required": "{field}不能为空",
	}
}

// Translates 你可以自定义字段翻译
func (f SysRoleApisContext) Translates() map[string]string {
	return validate.MS{
		"RoleId": "角色ID",
		"Apis":   "权限ID",
	}
}

// 角色授权api
func (api *SysRoleApisApi) SaveData(c *gin.Context) {
	// 1: 第一件事情就准备数据的载体
	var sysRoleApisContext SysRoleApisContext
	err := c.ShouldBindJSON(&sysRoleApisContext)
	if err != nil {
		// 如果参数注入失败或者出错就返回接口调用这。出错了.
		response.FailWithBindParams(c)
		return
	}

	//2: 开始验证
	validation := validate.Struct(sysRoleApisContext)
	if !validation.Validate() {
		response.FailWithValidatorData(validation, c)
		return
	}

	menuIdsSplit := strings.Split(sysRoleApisContext.ApisIds, ",")
	var sysRoleApis []*sys.SysRoleApis
	for _, menuId := range menuIdsSplit {
		parseUint, _ := strconv.ParseUint(menuId, 10, 64)
		sysRoleApi := sys.SysRoleApis{}
		sysRoleApi.RoleId = sysRoleApisContext.RoleId
		sysRoleApi.ApiId = uint(parseUint)
		sysRoleApis = append(sysRoleApis, &sysRoleApi)
	}

	// 创建实例，保存帖子
	err2 := sysRoleApisService.SaveSysRoleApis(sysRoleApisContext.RoleId, sysRoleApis)
	// 如果保存失败。就返回创建失败的提升
	if err2 != nil {
		response.FailWithMessage("创建失败", c)
		return
	}
	// 如果保存成功，就返回创建创建成功
	response.Ok("创建成功", c)
}

/**查询用户角色对应的菜单*/
func (api *SysRoleApisApi) SelectRoleApis(c *gin.Context) {
	roleIdStr := c.Query("roleId")
	roleId := api.StringToUnit(roleIdStr)
	if roleId == 0 {
		response.FailWithMessage("请选择一个角色操作", c)
		return
	}
	roles, _ := sysRoleApisService.SelectRoleApis(roleId)
	response.Ok(roles, c)
}

/**角色切换更换菜单*/
func (api *SysRoleApisApi) ChangeRoleIdMenus(c *gin.Context) {
	roleIdStr := c.Query("roleId")
	roleId := api.StringToUnit(roleIdStr)
	if roleId == 0 {
		response.FailWithMessage("请选择一个角色操作", c)
		return
	}

	var sysRoles *sys.SysRoles
	err1 := global.XZ_DB.Where("id = ?", roleId).First(&sysRoles).Error
	if err1 != nil {
		response.FailWithMessage("该角色没有权限", c)
		return
	}
	// 根据用户查询菜单信息
	roleMenus, _ := sysRoleMenusService.SelectRoleMenus(roleId)
	// 根据用户id查询用户的角色的权限
	permissions, _ := sysRoleApisService.SelectRoleApis(roleId)
	// 改变角色，切换token
	api.ChangeRoleToken(c, sysRoles)
	// 查询返回
	response.Ok(map[string]any{"roleMenus": sysMenuService.Tree(roleMenus, 0), "permissions": permissions}, c)
}

func (api *SysRoleApisApi) ChangeRoleToken(c *gin.Context, role *sys.SysRoles) {
	// 生成jwt的对象
	myJwt := jwtgo.NewJWT()
	token := c.Request.Header.Get("Authorization")
	// 把token重新生成在返回
	claims := jwtgo.GetUserInfo(c)
	claims.RoleCode = role.RoleCode
	claims.RoleId = role.ID
	eptime, _ := utils.ParseDuration("7d")
	// 用当前时间+eptime 就是 新的token过期时间
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(eptime))
	// 生成新的token
	newToken, _ := myJwt.CreateTokenByOldToken(token, *claims)
	// 输出给浏览器----request--header----给服务端
	// 输出给浏览器----response--header---给浏览器
	c.Header("new-authorization", newToken)
	c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt.Unix(), 10))
}
