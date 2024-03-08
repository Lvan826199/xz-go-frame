/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 16:09
 */
package sys

import (
	"xz-go-frame/global"
	"xz-go-frame/model/entity/commons/request"
	"xz-go-frame/model/entity/sys"
	"xz-go-frame/service/commons"
)

// 对用户表的数据层处理
type SysRolesService struct {
	commons.BaseService[uint, sys.SysRoles]
}

// 添加
func (service *SysRolesService) SaveSysRoles(sysRoles *sys.SysRoles) (err error) {
	err = global.XZ_DB.Create(sysRoles).Error
	return err
}

// 修改
func (service *SysRolesService) UpdateSysRoles(sysRoles *sys.SysRoles) (err error) {
	err = global.XZ_DB.Unscoped().Model(sysRoles).Updates(sysRoles).Error
	return err
}

// 按照map的方式更新
func (service *SysRolesService) UpdateSysRolesMap(sysRoles *sys.SysRoles, sysRolesMap *map[string]any) (err error) {
	err = global.XZ_DB.Unscoped().Model(sysRoles).Updates(sysRolesMap).Error
	return err
}

// 按照map的方式更新
func (service *SysUserService) UpdateSysRolesMap(sysRoles *sys.SysRoles, mapFields *map[string]any) (err error) {
	err = global.XZ_DB.Unscoped().Model(sysRoles).Updates(mapFields).Error
	return err
}

// 删除
func (service *SysRolesService) DelSysRolesById(id uint) (err error) {
	var sysRoles sys.SysRoles
	err = global.XZ_DB.Where("id = ?", id).Delete(&sysRoles).Error
	return err
}

// 批量删除
func (service *SysRolesService) DeleteSysRolessByIds(sysRoless []sys.SysRoles) (err error) {
	err = global.XZ_DB.Delete(&sysRoless).Error
	return err
}

// 根据id查询信息
func (service *SysRolesService) GetSysRolesByID(id uint) (sysRoless *sys.SysRoles, err error) {
	err = global.XZ_DB.Unscoped().Omit("created_at", "updated_at").Where("id = ?", id).First(&sysRoless).Error
	return
}

// 查询角色信息
func (service *SysRolesService) FindRoles() (sysRoless []*sys.SysRoles, err error) {
	err = global.XZ_DB.Order("id asc").Find(&sysRoless).Error
	return
}

// 查询分页
func (service *SysRolesService) LoadSysRolesPage(info request.PageInfo) (list interface{}, total int64, err error) {
	// 获取分页的参数信息
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 准备查询那个数据库表,这里为什么不使用Model呢，因为我要使用别名
	db := global.XZ_DB.Model(sys.SysRoles{})

	// 准备切片帖子数组
	var sysRoless []sys.SysRoles
	// 加条件
	if info.Keyword != "" {
		db = db.Where("(role_name like ?)", "%"+info.Keyword+"%")
	}

	// 排序默时间降序降序
	db = db.Order("created_at desc")

	// 查询中枢
	err = db.Unscoped().Count(&total).Error
	if err != nil {
		return sysRoless, total, err
	} else {
		// 执行查询
		err = db.Unscoped().Limit(limit).Offset(offset).Find(&sysRoless).Error
	}

	// 结果返回
	return sysRoless, total, err
}
