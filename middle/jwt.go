/*
@Author: 梦无矶小仔
@Date:   2024/1/15 11:36
*/
package middle

import (
	"github.com/gin-gonic/gin"
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/commons/response"
)

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
		// 让后续的路由方法可以直接通过c.Get("claims")
		c.Set("claims", customClaims)
		c.Next()
	}

}
