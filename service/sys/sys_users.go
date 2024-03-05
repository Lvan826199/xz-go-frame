/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 16:10
 */
package sys

import (
	"xz-go-frame/global"
	"xz-go-frame/model/entity/sys"
	"xz-go-frame/service/commons"
)

// 对用户表的数据层处理
type SysUserService struct {
	commons.BaseService[uint, sys.SysUser]
}

// 用于登录
func (service *SysUserService) GetUserByAccount(account string) (sysUser *sys.SysUser, err error) {
	// 根据account进行查询
	err = global.XZ_DB.Unscoped().Where("account = ?", account).First(&sysUser).Error
	if err != nil {
		return nil, err
	}
	return sysUser, nil
}

// 添加
func (service *SysUserService) SaveSysUser(sysUser *sys.SysUser) (err error) {
	err = global.XZ_DB.Create(sysUser).Error
	return err
}

// 修改
func (service *SysUserService) UpdateSysUser(sysUser *sys.SysUser) (err error) {
	err = global.XZ_DB.Unscoped().Model(sysUser).Updates(sysUser).Error
	return err
}

// 按照map的方式更新
func (service *SysUserService) UpdateSysUserMap(sysUser *sys.SysUser, mapField *map[string]any) (err error) {
	err = global.XZ_DB.Unscoped().Model(sysUser).Updates(mapField).Error
	return err
}

// 删除
func (service *SysUserService) DelSysUserById(id uint) (err error) {
	var sysUser sys.SysUser
	err = global.XZ_DB.Where("id = ?", id).Delete(&sysUser).Error
	return err
}

// 批量删除
func (service *SysUserService) DeleteSysUsersByIds(sysUsers []sys.SysUser) (err error) {
	err = global.XZ_DB.Delete(&sysUsers).Error
	return err
}

// 根据id查询信息
func (service *SysUserService) GetSysUserByID(id uint) (sysUsers *sys.SysUser, err error) {
	err = global.XZ_DB.Unscoped().Omit("created_at", "updated_at").Where("id = ?", id).First(&sysUsers).Error
	return
}
