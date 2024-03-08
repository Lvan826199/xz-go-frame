/*
* @Author: 梦无矶小仔
* @Date: 2024/3/7 18:46
 */
package sys

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"strconv"
	"xz-go-frame/commons/response"
	"xz-go-frame/global"
	"xz-go-frame/model/entity/sys"
)

type SysMenuApi struct {
	global.BaseApi
}

// 拷贝
func (api *SysMenuApi) CopyData(c *gin.Context) {
	// 1: 获取id数据 注意定义李媛媛的/:id
	id := c.Param("id")
	// 先从数据库查询出这条数据
	var sysMenu sys.SysMenus
	err := global.XZ_DB.Where("id = ?", id).First(&sysMenu).Error
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	// 未验证
	sysMenu.ID = 0
	sysMenu.Path = ""
	err2 := global.XZ_DB.Create(&sysMenu).Error
	if err2 != nil {
		response.FailWithMessage("创建失败", c)
		return
	}
	response.Ok(sysMenu, c)
}

// 保存
func (api *SysMenuApi) SaveData(c *gin.Context) {
	// 1: 第一件事情就准备数据的载体
	var sysMenus sys.SysMenus
	err := c.ShouldBindJSON(&sysMenus)
	if err != nil {
		// 如果参数注入失败或者出错就返回接口调用这。出错了.
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 创建实例，保存帖子
	err = sysMenuService.SaveSysMenus(&sysMenus)
	// 如果保存失败。就返回创建失败的提升
	if err != nil {
		response.FailWithMessage("创建失败", c)
		return
	}
	// 如果保存成功，就返回创建创建成功
	response.Ok("创建成功", c)
}

// 状态修改
func (api *SysMenuApi) UpdateStatus(c *gin.Context) {
	type Params struct {
		Id    uint   `json:"id"`
		Filed string `json:"field"`
		Value any    `json:"value"`
	}
	var params Params
	err := c.ShouldBindJSON(&params)
	if err != nil {
		// 如果参数注入失败或者出错就返回接口调用这。出错了.
		response.FailWithMessage(err.Error(), c)
		return
	}
	var sysMenus sys.SysMenus
	affected := global.XZ_DB.Unscoped().Model(&sysMenus).Where("id = ?", params.Id).Update(params.Filed, params.Value).RowsAffected
	flag := affected > 0
	// 如果保存失败。就返回创建失败的提升
	if !flag {
		response.FailWithMessage("更新失败", c)
		return
	}
	// 如果保存成功，就返回创建创建成功
	response.Ok("更新成功", c)
}

// 编辑修改
func (api *SysMenuApi) UpdateById(c *gin.Context) {
	// 1: 第一件事情就准备数据的载体
	var sysMenus sys.SysMenus
	err := c.ShouldBindJSON(&sysMenus)
	if err != nil {
		// 如果参数注入失败或者出错就返回接口调用这。出错了.
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 结构体转化成map呢？
	m := structs.Map(sysMenus)
	m["is_deleted"] = sysMenus.IsDeleted
	err = sysMenuService.UpdateSysMenusMap(&sysMenus, &m)
	// 如果保存失败。就返回创建失败的提升
	if err != nil {
		fmt.Println(err)
		response.FailWithMessage("更新失败", c)
		return
	}
	// 如果保存成功，就返回创建创建成功
	response.Ok("更新成功", c)
}

// 根据id删除
func (api *SysMenuApi) DeleteById(c *gin.Context) {
	// 绑定参数用来获取/:id这个方式
	id := c.Param("id")
	// 开始执行
	parseUint, _ := strconv.ParseUint(id, 10, 64)
	var sysMenus sys.SysMenus
	err := global.XZ_DB.Unscoped().Where("id = ?", parseUint).Delete(&sysMenus)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.Ok("ok", c)
}

// 根据id查询信息
func (api *SysMenuApi) GetById(c *gin.Context) {
	// 根据id查询方法
	id := c.Param("id")
	// 根据id查询方法
	parseUint, _ := strconv.ParseUint(id, 10, 64)
	sysUser, err := sysMenuService.GetSysMenusByID(uint(parseUint))
	if err != nil {
		global.SugarLog.Errorf("查询用户: %s 失败", id)
		response.FailWithMessage("查询用户失败", c)
		return
	}

	response.Ok(sysUser, c)
}

// 查询菜单
func (api *SysMenuApi) FindMenus(c *gin.Context) {
	keyword := c.Query("keyword")
	sysMenus, err := sysMenuService.FinMenus(keyword)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.Ok(sysMenuService.Tree(sysMenus, 0), c)
}

// 查询父菜单
func (api *SysMenuApi) FindMenusRoot(c *gin.Context) {
	sysMenus, err := sysMenuService.FinMenusRoot()
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.Ok(sysMenus, c)
}
