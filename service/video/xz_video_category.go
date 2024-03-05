/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 15:42
 */
package video

import (
	"xz-go-frame/global"
	"xz-go-frame/model/entity/bbs"
	"xz-go-frame/model/entity/commons/request"
	"xz-go-frame/model/entity/video"
)

type VideoCategoryService struct {
}

func (service *VideoCategoryService) CreateXkVideoCategory(xkVideoCategory *video.XkVideoCategory) (err error) {
	// 1： 获取数据的连接对象 如果执行成功err是nil，如果失败就把失败告诉
	err = global.XZ_DB.Create(xkVideoCategory).Error
	return err
}

func (service *VideoCategoryService) UpdateXkVideoCategory(xkVideoCategory *video.XkVideoCategory) (err error) {
	err = global.XZ_DB.Model(xkVideoCategory).Updates(xkVideoCategory).Error
	return err
}

// 批量删除
func (cbbs *VideoCategoryService) DeleteVideoCategorysByIds(xkVideoCategorys []video.XkVideoCategory) (err error) {
	err = global.XZ_DB.Delete(&xkVideoCategorys).Error
	return err
}

// 根据ID删除帖子
func (service *VideoCategoryService) DeleteXkVideoCategoryById(id uint) (err error) {
	var xkVideoCategory video.XkVideoCategory
	// Unscoped()可以操作包括那些被软删除的记录在内的数据。
	err = global.XZ_DB.Unscoped().Where("id = ?", id).Delete(&xkVideoCategory).Error
	return err
}

func (service *VideoCategoryService) GetXkVideoCategory(id uint) (xkVideoCategory *video.XkVideoCategory, err error) {
	err = global.XZ_DB.Unscoped().Where("id = ?", id).First(&xkVideoCategory).Error
	return
}

// 查询所有的主分类
func (service *VideoCategoryService) FindCategoryAll() (xkVideoCategory []*video.XkVideoCategory, err error) {
	err = global.XZ_DB.Unscoped().Where("status = 1  and parent_id = 0").Find(&xkVideoCategory).Error
	return
}

// 分页获取
func (service *VideoCategoryService) LoadXkVideoCategoryPage(info request.PageInfo) (list interface{}, total int64, err error) {
	// 获取分页的参数信息
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 准备查询那个数据库表
	db := global.XZ_DB.Unscoped().Model(&video.XkVideoCategory{})

	// 准备切片帖子数组
	var XkVideoCategoryList []video.XkVideoCategory

	// 加条件
	if info.Keyword != "" {
		db = db.Where("title like ?", "%"+info.Keyword+"%")
	}

	// 排序默时间降序降序
	db = db.Order("created_at desc")

	// 查询中枢
	err = db.Count(&total).Error
	if err != nil {
		return XkVideoCategoryList, total, err
	} else {
		// 执行查询
		err = db.Limit(limit).Offset(offset).Find(&XkVideoCategoryList).Error
	}

	// 结果返回
	return XkVideoCategoryList, total, err
}

func (service *VideoCategoryService) FindCategories() (categories []*video.XkVideoCategory, err error) {
	err = global.XZ_DB.Unscoped().Order("sorted asc").Find(&categories).Error
	return categories, err
}

func (service *VideoCategoryService) Tree(allDbCategoires []*video.XkVideoCategory, parentId uint) []*video.XkVideoCategory {
	var nodes []*video.XkVideoCategory //---------准备空教室
	// 开始遍历父类
	for _, dbCategory := range allDbCategoires { //1 parentId = 0 parentId=0 2 3 4 5 6 7 8 9 10
		if dbCategory.ParentId == parentId {
			dbCategory.Children = append(dbCategory.Children, service.Tree(allDbCategoires, dbCategory.ID)...)
			nodes = append(nodes, dbCategory)
		}
	}
	return nodes
}

// 修改状态
func (cbbs *VideoCategoryService) UpdateBbsCategoryStatus(statusReq *request.StatusReq) (err error) {
	err = global.XZ_DB.Model(new(*video.XkVideoCategory)).Where("id=?", statusReq.ID).Update(statusReq.Field, statusReq.Value).Error
	return err
}

// 删除
func (cbbs *VideoCategoryService) DeleteBbsCategory(bbsCategory *bbs.BbsCategory) (err error) {
	err = global.XZ_DB.Delete(&bbsCategory).Error
	return err
}

// 删除
func (cbbs *VideoCategoryService) DeleteVideoCategoryById(id uint) (err error) {
	var videoCategory video.XkVideoCategory
	err = global.XZ_DB.Where("id = ?", id).Delete(&videoCategory).Error
	return err
}

// 批量删除
func (cbbs *VideoCategoryService) DeleteBbsCategoryByIds(videoCategories []video.XkVideoCategory) (err error) {
	err = global.XZ_DB.Delete(&videoCategories).Error
	return err
}

// 两级写死的做法
//func (service *VideoCategoryService) Tree(allDbCategoires []*video.XkVideoCategory) []*video.XkVideoCategory {
//	// 定义一个节点
//	//allDbCategoires =
//	//1	Java	Java	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	0
//	//2	Go	Go	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	0
//	//3	Javascript	Javascript	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	0
//	//4	Spring	Spring	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	1
//	//5	SpringBoot	SpringBoot	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	1
//	//6	Gin	Gin	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	2
//	//7	Beego	Beego	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	2
//	//8	XOrm	XOrm	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	2
//	//9	Gorm	Gorm	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	2
//	//10	GVA	GVA	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	2
//
//	//nodes
//	//1	Java	Java	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	0 ---parentNode-XkVideoCategory
//	//2	Go	Go	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	0
//	//3	Javascript	Javascript	2023-06-04 20:59:19	2023-06-04 20:59:19	1	1	0	0
//
//	var nodes []*video.XkVideoCategory //---------准备空教室
//	for _, dbCategory := range allDbCategoires {
//		if dbCategory.ParentId == 0 {
//			// 这里找到所有的父类
//			nodes = append(nodes, dbCategory)
//		}
//	}
//
//	// 开始遍历父类
//	for _, dbCategory := range allDbCategoires {
//		for _, parentNode := range nodes {
//			if dbCategory.ParentId == parentNode.ID {
//				parentNode.Children = append(parentNode.Children, dbCategory)
//			}
//		}
//	}
//	return nodes
//}
