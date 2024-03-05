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

type SysMenusService struct {
	commons.BaseService[uint, sys.SysMenus]
}

/*
查询父级菜单
*/
func (service *SysMenusService) FinMenusRoot() (sysMenus []*sys.SysMenus, err error) {
	err = global.XZ_DB.Where("parent_id = ? ", 0).Order("sort asc").Find(&sysMenus).Error
	return sysMenus, err
}

/*
*

  - 查询菜单形成tree数据

  - 格式： {
    id: 2,
    date: '2016-05-04',
    name: 'wangxiaohu',
    address: 'No. 189, Grove St, Los Angeles',
    },
    {
    id: 3,
    date: '2016-05-01',
    name: 'wangxiaohu',
    address: 'No. 189, Grove St, Los Angeles',
    children: [
    {
    id: 31,
    date: '2016-05-01',
    name: 'wangxiaohu',
    address: 'No. 189, Grove St, Los Angeles',
    },
    {
    id: 32,
    date: '2016-05-01',
    name: 'wangxiaohu',
    address: 'No. 189, Grove St, Los Angeles',
    },
    ],
    },
*/
func (service *SysMenusService) FinMenus(keyword string) (sysMenus []*sys.SysMenus, err error) {
	db := global.XZ_DB.Unscoped().Order("sort asc")
	if len(keyword) > 0 {
		db.Where("title like ?", "%"+keyword+"%")
	}
	err = db.Find(&sysMenus).Error
	return sysMenus, err
}

/**
*   开始把数据进行编排--递归
*   Tree(all,0)
 */
func (service *SysMenusService) Tree(allSysMenus []*sys.SysMenus, parentId uint) []*sys.SysMenus {
	var nodes []*sys.SysMenus
	for _, dbMenu := range allSysMenus {
		if dbMenu.ParentId == parentId {
			childrensMenu := service.Tree(allSysMenus, dbMenu.ID)
			if len(childrensMenu) > 0 {
				dbMenu.Children = append(dbMenu.Children, childrensMenu...)
			}
			nodes = append(nodes, dbMenu)
		}
	}
	return nodes
}

// 添加
func (service *SysMenusService) SaveSysMenus(sysMenus *sys.SysMenus) (err error) {
	err = global.XZ_DB.Create(sysMenus).Error
	return err
}

// 按照map的方式过呢更新
func (service *SysMenusService) UpdateSysMenusMap(sysMenus *sys.SysMenus, sysMenusMap *map[string]any) (err error) {
	err = global.XZ_DB.Unscoped().Model(sysMenus).Updates(sysMenusMap).Error
	return err
}

// 删除
func (service *SysMenusService) DelSysMenusById(id uint) (err error) {
	var sysMenus sys.SysMenus
	err = global.XZ_DB.Where("id = ?", id).Delete(&sysMenus).Error
	return err
}

// 批量删除
func (service *SysMenusService) DeleteSysMenussByIds(sysMenuss []sys.SysMenus) (err error) {
	err = global.XZ_DB.Delete(&sysMenuss).Error
	return err
}

// 根据id查询信息
func (service *SysMenusService) GetSysMenusByID(id uint) (sysMenuss *sys.SysMenus, err error) {
	err = global.XZ_DB.Unscoped().Omit("created_at", "updated_at").Where("id = ?", id).First(&sysMenuss).Error
	return
}
