/*
* @Author: 梦无矶小仔
* @Date:   2024/1/11 14:11
 */
package user

import (
	"database/sql"
	"xz-go-frame/global"
	"xz-go-frame/model/entity/user/request"
	"xz-go-frame/model/entity/user/response"
)

// 对用户表的数据层处理
type UserService struct{}

// 查询作者信息
func (service *UserService) FindAuthorPageMap(req *request.AuthorPageReq) (resp map[string]interface{}, perr error) {
	// 获取globaldb
	var results []map[string]interface{}
	var total int64

	// 执行sql json.
	db := global.XZ_DB.Table("xz_user t1, xz_user_author t2").
		Select("t1.id, t1.avatar, t1.nickname, t1.account, t2.author_name").
		Where("t1.id = t2.user_id AND t2.`status` = 1 AND t1.forbbiden = 0")

	// 开始判断关键词是否存在，
	if len(req.Keyword) > 0 {
		db.Where("(t1.nickname like @name1 or t1.account like ? or t2.author_name like ?)",
			sql.Named("name1", "%"+req.Keyword+"%"), "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 根据用户搜索用户信息，如果是0不参与搜索，否则就是具体用户id
	if req.UserId > 0 {
		db.Where("t1.id = ?", req.UserId)
	}

	perr1 := db.Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Scan(&results).Error
	perr2 := db.Count(&total).Error

	if perr1 != nil || perr2 != nil {
		return resp, perr1
	}

	m := map[string]any{}
	m["pageNum"] = req.PageNum
	m["pageSize"] = req.PageSize
	m["total"] = total
	m["list"] = results
	//结果返回
	return m, nil
}
func (service *UserService) FindAuthorPage(req *request.AuthorPageReq) (resp response.XkAuthorResp, perr error) {
	// 获取globaldb
	var authorResponses []response.AuthorResopnse
	var total int64

	// 执行sql json.
	db := global.XZ_DB.Table("xk_user t1, xk_user_author t2").
		Select("t1.id, t1.avatar, t1.nickname, t1.account, t2.author_name").
		Where("t1.id = t2.user_id AND t2.`status` = 1 AND t1.forbbiden = 0")

	// 开始判断关键词是否存在，
	if len(req.Keyword) > 0 {
		db.Where("(t1.nickname like @name1 or t1.account like ? or t2.author_name like ?)",
			sql.Named("name1", "%"+req.Keyword+"%"), "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 根据用户搜索用户信息，如果是0不参与搜索，否则就是具体用户id
	if req.UserId > 0 {
		db.Where("t1.id = ?", req.UserId)
	}

	perr1 := db.Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Scan(&authorResponses).Error
	perr2 := db.Count(&total).Error

	if perr1 != nil || perr2 != nil {
		return resp, perr1
	}

	resp.PageNum = req.PageNum
	resp.PageSize = req.PageSize
	resp.Total = total
	// 返回列表信息
	resp.List = authorResponses
	//结果返回
	return resp, nil
}

//// nil 是go空值处理，必须是指针类型
//func (service *UserService) GetUserByAccount(account string) (user *user.XzUser, err error) {
//	// 根据account进行查询
//	err = global.XZ_DB.Where("account = ?", account).First(&user).Error
//	if err != nil {
//		return nil, err
//	} else {
//		return user, nil
//	}
//}
