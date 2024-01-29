/*
@Author: 梦无矶小仔
@Date:   2024/1/15 16:40
*/
package jwtgo

import (
	"errors"
	"gorm.io/gorm"
	"xz-go-frame/global"
	"xz-go-frame/model/entity/jwt"
)

type JwtService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList model.JwtBlacklist
//@return: err error

func (jwtService *JwtService) JsonInBlacklist(jwtList jwt.JwtBlacklist) (err error) {
	err = global.XZ_DB.Create(&jwtList).Error
	if err != nil {
		return
	}
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool

func (jwtService *JwtService) IsBlacklist(jwttoken string) bool {
	//_, ok := global.BlackCache.Get(jwt)
	//return ok
	err := global.XZ_DB.Where("jwt = ?", jwttoken).First(&jwt.JwtBlacklist{}).Error
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !isNotFound
}
