package middle

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取会话
		session := sessions.Default(c)
		// 获取用户登录信息
		user := session.Get("user")
		// 如果用户没有登录，直接重定向返回登录
		if user == nil {
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort() // 拦截，到这里就不会往下执行了
		}
		// 取出会话信息
		username := user.(string)
		// 把session用户信息，放入到context中，给后续路由进行使用
		// 好处就是： router中方法不需要再次获取session再来拿会话中的信息
		c.Set("username", username)
		c.Next() //放行，默认就会放行
	}
}
