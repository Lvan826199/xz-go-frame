/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 18:57
 */
package orm

import (
	"xz-go-frame/global"
	"xz-go-frame/model/jwt"
	"xz-go-frame/model/user"
)

func RegisterTable() {
	db := global.XZ_DB
	// 注册和声明model
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(jwt.JwtBlacklist{})
}
