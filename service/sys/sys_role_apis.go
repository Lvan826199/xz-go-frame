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

type SysRoleApisService struct {
	commons.BaseService[uint, sys.SysRoleApis]
}

// 角色授予api
func (service *SysRoleApisService) SaveSysRoleApis(roleId uint, sysRolesApis []*sys.SysRoleApis) (err error) {
	tx := global.XZ_DB.Begin()
	// 删除用户对应的角色
	if err := tx.Where("role_id = ?", roleId).Delete(&sys.SysRoleApis{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 开始保存用户和角色的关系
	if err := tx.Create(sysRolesApis).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 查询角色授权的信息
func (service *SysRoleApisService) SelectRoleApis(roleId uint) (sysApiss []*sys.SysApis, err error) {
	err = global.XZ_DB.Select("t2.* ").Table("sys_role_apis t1,sys_apis t2").
		Where("t1.role_id = ? AND t1.api_id = t2.id and t2.is_deleted = 0", roleId).Scan(&sysApiss).Error
	return sysApiss, err
}
