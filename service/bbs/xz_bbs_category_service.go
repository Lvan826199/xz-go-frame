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

// 定义bbs的service提供BbsCategory的数据curd的操作
type BbsCategoryService struct{}

// 创建分类
func (service *BbsCategoryService) CreateBbsCategory(bbsCategory *bbs.BbsCategory) (err error) {
	// 1： 获取数据的连接对象 如果执行成功err是nil，如果失败就把失败告诉
	err = global.XZ_DB.Create(bbsCategory).Error
	return err
}

// 更新分类
func (service *BbsCategoryService) UpdateBbsCategory(bbsCategory *bbs.BbsCategory) (err error) {
	err = global.XZ_DB.Model(bbsCategory).Save(bbsCategory).Error
	return err
}

// 修改状态
func (service *BbsCategoryService) UpdateBbsCategoryStatus(statusReq *request.StatusReq) (err error) {
	err = global.XZ_DB.Model(new(*bbs.BbsCategory)).Where("id=?", statusReq.ID).Update(statusReq.Field, statusReq.Value).Error
	return err
}

// 删除
func (service *BbsCategoryService) DeleteBbsCategory(bbsCategory *bbs.BbsCategory) (err error) {
	err = global.XZ_DB.Delete(&bbsCategory).Error
	return err
}

// 根据id删除
func (service *BbsCategoryService) DeleteBbsCategoryById(id uint) (err error) {
	var BbsCategory bbs.BbsCategory
	err = global.XZ_DB.Where("id = ?", id).Delete(&BbsCategory).Error
	return err
}

// 批量删除
func (service *BbsCategoryService) DeleteBbsCategoryByIds(bbsCategorys []bbs.BbsCategory) (err error) {
	err = global.XZ_DB.Delete(&bbsCategorys).Error
	return err
}

// 获取帖子信息
func (service *BbsCategoryService) GetBbsCategory(id uint) (bbsCategory *bbs.BbsCategory, err error) {
	err = global.XZ_DB.Where("id = ?", id).First(&bbsCategory).Error
	return
}

// 获取帖子信息
func (service *BbsCategoryService) FindBbsCategory() (bbsCategorys []bbs.BbsCategory, err error) {
	err = global.XZ_DB.Where("status = 1 and is_delete = 0").Find(&bbsCategorys).Error
	return
}

// 分页获取客户列表
func (service *BbsCategoryService) LoadBbsCategoryPage(info request.BbsCategoryPageInfo) (list interface{}, total int64, err error) {
	// 获取分页的参数信息
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 准备查询那个数据库表
	db := global.XZ_DB.Model(&bbs.BbsCategory{})

	// 准备切片帖子数组
	var BbsCategoryList []bbs.BbsCategory

	// 加条件
	if info.Keyword != "" {
		db = db.Where("title like ?", "%"+info.Keyword+"%")
	}

	if info.Status != -1 {
		db = db.Where("status = ?", info.Status)
	}
	// 排序默时间降序降序
	db = db.Order("sorted asc")

	// 查询中枢
	err = db.Count(&total).Error
	if err != nil {
		return BbsCategoryList, total, err
	} else {
		// 执行查询
		err = db.Limit(limit).Offset(offset).Find(&BbsCategoryList).Error
	}

	// 结果返回
	return BbsCategoryList, total, err
}
