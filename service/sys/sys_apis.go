/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 16:08
 */
package sys

import (
	"xz-go-frame/global"
	"xz-go-frame/model/entity/sys"
	"xz-go-frame/service/commons"
)

// 对用户表的数据层处理
type SysApisService struct {
	commons.BaseService[uint, sys.SysApis]
}

// 添加
func (service *SysApisService) SaveSysApis(sysApis *sys.SysApis) (err error) {
	err = global.XZ_DB.Create(sysApis).Error
	return err
}

// 修改
func (service *SysApisService) UpdateSysApis(sysApis *sys.SysApis) (err error) {
	err = global.XZ_DB.Unscoped().Model(sysApis).Updates(sysApis).Error
	return err
}

// 按照map的方式更新
func (service *SysApisService) UpdateSysApisMap(sysApis *sys.SysApis, mapFileds *map[string]any) (err error) {
	err = global.XZ_DB.Unscoped().Model(sysApis).Updates(mapFileds).Error
	return err
}

// 根据ID删除
func (service *SysApisService) DelSysApisById(id uint) (err error) {
	var sysApis sys.SysApis
	err = global.XZ_DB.Where("id = ?", id).Delete(&sysApis).Error
	return err
}

// 批量删除
func (service *SysApisService) DeleteSysApissByIds(sysApiss []sys.SysApis) (err error) {
	err = global.XZ_DB.Delete(&sysApiss).Error
	return err
}

// 根据id查询信息
func (service *SysApisService) GetSysApisByID(id uint) (sysApiss *sys.SysApis, err error) {
	// Omit排除某些字段
	err = global.XZ_DB.Unscoped().Omit("created_at", "updated_at").Where("id = ?", id).First(&sysApiss).Error
	return
}

// 根据title关键字查询
func (service *SysApisService) FinApiss(keyword string) (sysApis []*sys.SysApis, err error) {
	db := global.XZ_DB.Unscoped().Order("id asc")
	if len(keyword) > 0 {
		db.Where("title like ?", "%"+keyword+"%")
	}
	err = db.Find(&sysApis).Error
	return sysApis, err
}

/**
*   开始把数据进行编排--递归
*   Tree(all,0)
 */
func (service *SysApisService) Tree(allSysApis []*sys.SysApis, parentId uint) []*sys.SysApis {
	var nodes []*sys.SysApis
	for _, dbApis := range allSysApis {
		if dbApis.ParentId == parentId {
			childrensApis := service.Tree(allSysApis, dbApis.ID)
			if len(childrensApis) > 0 {
				dbApis.Children = append(dbApis.Children, childrensApis...)
			}
			nodes = append(nodes, dbApis)
		}
	}
	return nodes
}

/*
查询父级权限
*/
func (service *SysApisService) FinApisRoot() (sysApis []*sys.SysApis, err error) {
	err = global.XZ_DB.Where("parent_id = ? ", 0).Order("id asc").Find(&sysApis).Error
	return sysApis, err
}
