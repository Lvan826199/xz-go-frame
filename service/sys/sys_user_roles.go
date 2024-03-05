/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 16:09
 */
package sys

import (
	"xz-go-frame/global"
	sys2 "xz-go-frame/model/entity/sys"
	"xz-go-frame/service/commons"
)

type SysUserRolesService struct {
	commons.BaseService[uint, sys2.SysUserRoles]
}

// 用户授权
func (service *SysUserRolesService) SaveSysUserRoles(userId uint, sysUserRoles []*sys2.SysUserRoles) (err error) {
	// 事务开启
	tx := global.XZ_DB.Begin()
	// 删除用户对应的角色-------------执行成功了，会立即提交吗？
	if err := tx.Where("user_id = ?", userId).Delete(&sys2.SysUserRoles{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 开始保存用户和角色的关系--------------执行成功了，会立即提交吗？
	if err := tx.Create(sysUserRoles).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 事务提交
	return tx.Commit().Error
}

// 查询用户授权的角色信息
func (service *SysUserRolesService) SelectUserRoles(userId uint) (sysRoles []*sys2.SysRoles, err error) {
	err = global.XZ_DB.Select("t2.*").Table("sys_user_roles t1,sys_roles t2").
		Where("t1.user_id = ? and t1.role_id = t2.id and t2.is_deleted = 0", userId).Scan(&sysRoles).Error
	return sysRoles, err
}
