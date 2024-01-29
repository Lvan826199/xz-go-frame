/*
@Author: 梦无矶小仔
@Date:   2024/1/15 11:36
*/
package middle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/commons/response"
	jwtdb "xz-go-frame/model/entity/jwt"
	"xz-go-frame/utils"
)

var jwtService = jwtgo.JwtService{}

// 定义一个JWTAuth的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过http header中的token解析来认证
		// 获取token
		// 我们这里jwt鉴权取头部信息 Authorization 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		//token := c.Request.Header.Get("x-token")
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Fail(701, "请求未携带token，无权限访问", c)
			c.Abort()
			return
		}
		// 生成jwt的对象
		myJwt := jwtgo.NewJWT()
		// parseToken 解析token包含的信息
		customClaims, err := myJwt.ParseToken(token)
		// 如果解析失败就出现异常
		if err != nil {
			response.Fail(60001, "token失效了", c)
			c.Abort()
			return
		}

		// 增加jwt-token的续期功能
		// 判断过期时间 - now  < buffertime 就开始续期 ep 1d -- no
		fmt.Println("customClaims.ExpiresAt", customClaims.ExpiresAt)
		fmt.Println("time.Now().Unix()", time.Now().Unix())
		fmt.Println("customClaims.ExpiresAt - time.Now().Unix()", customClaims.ExpiresAt.Unix()-time.Now().Unix())
		fmt.Println("customClaims.BufferTime", customClaims.BufferTime)

		if customClaims.ExpiresAt.Unix()-time.Now().Unix() < customClaims.BufferTime {
			// 1、生成一个新的token
			// 2、用c把新的token返回页面
			fmt.Println("开始续期.....")
			// 获取7天的过期时间
			eptime, _ := utils.ParseDuration("7d")
			// 用当前时间 + eptime 就是新的token过期时间
			customClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(eptime))
			// 生成新的token
			newToken, _ := myJwt.CreateTokenByOldToken(token, *customClaims)
			// 输出给浏览器 --- request --- header --- 给服务端
			// 输出给浏览器 --- reponse --- header --- 给浏览器
			c.Header("new-authorization", newToken)
			c.Header("new-expires-at", strconv.FormatInt(customClaims.ExpiresAt.Unix(), 10))
			// 如果生成新的token了，旧的token怎么办？jwt没有提供一个机制让旧token失效。
			_ = jwtService.JsonInBlacklist(jwtdb.JwtBlacklist{Jwt: token})
		}

		// 让后续的路由方法可以直接通过c.Get("claims")
		c.Set("claims", customClaims)
		c.Next()
	}

}
