/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 16:09
 */
package sys

import (
	"xz-go-frame/global"
	"xz-go-frame/model/entity/sys"
	"xz-go-frame/service/commons"
)

type SysRoleMenusService struct {
	commons.BaseService[uint, sys.SysRoleMenus]
}

// 角色授予菜单
func (service *SysRoleMenusService) SaveSysRoleMenus(roleId uint, sysRolesMenus []*sys.SysRoleMenus) (err error) {
	// 事务加持
	tx := global.XZ_DB.Begin()
	// 删除用户对应的角色
	if err := tx.Where("role_id = ?", roleId).Delete(&sys.SysRoleMenus{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 开始保存用户和角色的关系
	if err := tx.Create(sysRolesMenus).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 查询角色对应的菜单信息
func (service *SysRoleMenusService) SelectRoleMenus(roleId uint) (sysMenus []*sys.SysMenus, err error) {
	err = global.XZ_DB.Select("t2.*").Table("sys_role_menus t1,sys_menus t2").
		Where("t1.menu_id = t2.id  AND t1.role_id = ? AND t2.hidden = 1 and t2.is_deleted = 0", roleId).Scan(&sysMenus).Error
	return sysMenus, err
}
