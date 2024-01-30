/*
@Author: 梦无矶小仔
@Date:   2024/1/12 18:23
*/
package jwt

import "xz-go-frame/global"

type JwtBlacklist struct {
	global.GVA_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
