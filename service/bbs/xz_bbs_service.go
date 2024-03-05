/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 14:38
 */
package bbs

import (
	"xz-go-frame/global"
	"xz-go-frame/model/entity/bbs"
	"xz-go-frame/model/entity/bbs/request"
)

// 定义bbs的service提供xkbbs的数据curd的操作

type BbsService struct{}

func (service *BbsService) CreateXkBbs(xkBbs *bbs.XkBbs) (err error) {
	// 1、获取数据的连接对象，如果执行成功err是nil，如果失败就把失败抛出
	err = global.XZ_DB.Create(xkBbs).Error
	return err
}

func (service *BbsService) UpdateXkBbs(xkBbs *bbs.XkBbs) (err error) {
	err = global.XZ_DB.Save(xkBbs).Error
	//err = global.XZ_DB.Model(xkBbs).Updates(xkBbs).Error
	return err
}

func (service *BbsService) DeleteXkBbs(xkBbs *bbs.XkBbs) (err error) {
	err = global.XZ_DB.Delete(&xkBbs).Error
	return err
}

func (service *BbsService) DeleteXkBbsById(id uint) (err error) {
	var xkBbs bbs.XkBbs
	err = global.XZ_DB.Where("id = ?", id).Delete(&xkBbs).Error
	return err
}

func (service *BbsService) GetXkBbs(id uint) (xkBbs *bbs.XkBbs, err error) {
	err = global.XZ_DB.Where("id = ?", id).First(&xkBbs).Error
	return
}

// 分页获取信息

func (service *BbsService) LoadXkBbsPage(info request.BbsPageInfo) (list any, total int64, err error) {
	// 获取分页的参数信息
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 准备查询的数据库表
	db := global.XZ_DB.Model(&bbs.XkBbs{})

	// 准备切片帖子数组
	var XkBbsList []bbs.XkBbs

	if info.CategoryId != -1 {
		db = db.Where("category_id = ?", info.CategoryId)
	}

	if info.Status != -1 {
		db = db.Where("status = ?", info.Status)
	}

	// 加条件
	if info.Keyword != "" {
		db = db.Where("(title like ?  or category_name like ? or user_id like ? or username like ?)", "%"+info.Keyword+"%")
	}

	// 排序默时间降序降序
	db = db.Order("created_at desc")

	// 查询中枢
	err = db.Count(&total).Error

	if err != nil {
		return XkBbsList, total, err
	} else {
		// 执行查询
		err = db.Limit(limit).Offset(offset).Find(&XkBbsList).Error
	}
	// 返回查询结果
	return XkBbsList, total, err
}
