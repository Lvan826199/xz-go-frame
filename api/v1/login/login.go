package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mojocn/base64Captcha"
	"time"
	"xz-go-frame/commons/jwtgo"
	"xz-go-frame/commons/response"
	"xz-go-frame/model/user"
	service "xz-go-frame/service/user"
	"xz-go-frame/utils"
)

// 登录业务
type LoginApi struct {
}

// 1、定义验证码的store -- 默认是存储在go的内存中
var store = base64Captcha.DefaultMemStore

// 登录接口的处理
func (api *LoginApi) ToLogined(c *gin.Context) {
	type LoginParam struct {
		Account  string
		Code     string
		CodeId   string
		Password string
	}
	// 1、获取用户在页面上输入账号和密码，开始在数据库里对数据进行校验
	userService := service.UserService{}
	param := LoginParam{}
	err2 := c.ShouldBindJSON(&param)
	if err2 != nil {
		response.Fail(60002, "参数绑定有误", c)
		return
	}

	//if len(param.Code) == 0 {
	//	response.Fail(60002, "请输入验证码", c)
	//	return
	//}
	//
	//if len(param.CodeId) == 0 {
	//	response.Fail(60002, "验证码获取失败", c)
	//	return
	//}
	//
	//// 开始校验验证码是否正确
	//verify := store.Verify(param.CodeId, param.Code, true)
	//if !verify {
	//	response.Fail(60002, "你输入的验证码有误!!", c)
	//	return
	//}
	inputAccount := param.Account
	inputPassword := param.Password

	if len(inputAccount) == 0 {
		response.Fail(60002, "请输入账号", c)
		return
	}

	if len(inputPassword) == 0 {
		response.Fail(60002, "请输入密码", c)
		return
	}

	dbUser, err := userService.GetUserByAccount(inputAccount)
	if err != nil {
		response.Fail(60002, "您输入的账号或密码错误", c)
		return
	}

	// 校验用户的账号密码输入是否和数据库一致
	if dbUser != nil && dbUser.Password == inputPassword {
		token := api.generaterToken(c, dbUser)
		// 根据用户id查询用户的角色
		roles := [2]map[string]any{}
		m1 := map[string]any{"id": 1, "name": "超级管理员"}
		m2 := map[string]any{"id": 2, "name": "财务"}
		roles[0] = m1
		roles[1] = m2
		// 根据用户id查询用户的角色的权限
		permissions := [2]map[string]any{}
		pm1 := map[string]any{"code": 10001, "name": "保存用户"}
		pm2 := map[string]any{"code": 20001, "name": "删除用户"}
		permissions[0] = pm1
		permissions[1] = pm2

		response.Ok(map[string]any{"user": dbUser, "token": token, "roles": roles, "permissions": permissions}, c)

	} else {
		response.Fail(60002, "你输入的账号和密码有误", c)
	}
}

/*
根据用户信息创建一个token
*/
func (api *LoginApi) generaterToken(c *gin.Context, dbUser *user.User) string {
	// 设置token续期的缓冲时间
	bf, _ := utils.ParseDuration("1d")
	ep, _ := utils.ParseDuration("7d")

	// 1、jwt 生成token
	myJwt := jwtgo.NewJWT()
	// 2、生成token
	token, err2 := myJwt.CreateToken(jwtgo.CustomClaims{
		dbUser.ID,
		dbUser.Name,
		int64(bf / time.Second),
		jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"XZ-USER"},               // 受众
			Issuer:    "MWJ-ADMIN",                               // 签发者
			IssuedAt:  jwt.NewNumericDate(time.Now()),            // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 过期时间 7天
		},
	})
	fmt.Println("当前时间是：", time.Now().Unix())
	fmt.Println("缓冲时间是：", int64(bf/time.Second))
	fmt.Println("签发时间：" + time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("生效时间：" + time.Now().Add(-1000).Format("2006-01-02 15:04:05"))
	fmt.Println("过期时间：" + time.Now().Add(ep).Format("2006-01-02 15:04:05"))
	if err2 != nil {
		response.Fail(60002, "登录失败,token颁发不成功！", c)
	}
	return token
}

// Login 用session时候的代码
//func (e *LoginApi) Login(c *gin.Context) {
//	// session ---- 是一种所有请求之间的数据共享机制，为什么会出现session，是因为http请求是一种无状态。
//	// 什么叫无状态：就是指，用户在浏览器输入方位地址的时候，地址请求到服务区，到响应服务，并不会存储任何数据在客户端或者服务端，
//	// 也是就：一次request---response就意味着内存消亡，也就以为整个过程请求和响应过程结束。
//	// 但是往往在开发中，我们可能要存存储一些信息，让各个请求之间进行共享。所有就出现了session会话机制
//	// session会话机制其实是一种服务端存储技术，底层原理是一个map
//	// 比如：我登录的时候，要把用户信息存储session中，然后给 map[key]any =
//	// key = sdf365454klsdflsd --sessionid
//
//	// 初始化session对象
//	session := sessions.Default(c)
//	// 存放用户信息到session中
//	session.Set("user", "xiaozai") // map[sessionid] == map[user][xiaozai]
//	// 记住一定要调用save方法，否则内存不会写入进去
//	session.Save()
//	username := session.Get("user")
//	c.JSON(http.StatusOK, "我是gin:登录用户名："+username.(string))
//}
