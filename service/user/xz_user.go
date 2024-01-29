/*
* @Author: 梦无矶小仔
* @Date:   2024/1/11 14:11
 */
package user

import (
	"xz-go-frame/global"
	"xz-go-frame/model/entity/user"
)

// 对用户表的数据层处理
type UserService struct{}

// nil 是go空值处理，必须是指针类型
func (service *UserService) GetUserByAccount(account string) (user *user.User, err error) {
	// 根据account进行查询
	err = global.XZ_DB.Where("account = ?", account).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}
