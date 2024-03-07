/*
* @Author: 梦无矶小仔
* @Date: 2024/3/7 16:08
 */
package jwtgo

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/global"
)

func GetClaims(c *gin.Context) (*CustomClaims, error) {
	token := c.Request.Header.Get("Authorization")
	// 生成jwt的对象
	myJwt := NewJWT()
	// 解析token
	customClaims, err := myJwt.ParseToken(token)
	if err != nil {
		global.Log.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在Authorization且claims是否为规定结构")
	}
	return customClaims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.UserId
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.UserId
	}
}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserRoleCode(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.RoleCode
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.RoleCode
	}
}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserRoleId(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.RoleId
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse.RoleId
	}
}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*CustomClaims)
		return waitUse
	}
}
