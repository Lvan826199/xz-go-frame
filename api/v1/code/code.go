/*
@Author: 梦无矶小仔
@Date:   2024/1/11 14:32
*/

package code

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"xz-go-frame/commons/response"
)

// 验证码生成
type CodeApi struct{}

// 1、定义验证的store -- 默认是存储在go的内存中
var store = base64Captcha.DefaultMemStore

// 2、创建验证码
func (api *CodeApi) CreateCaptcha(c *gin.Context) {
	// 2、生成验证码的类型，这个默认类型是一个数字的验证码
	driver := &base64Captcha.DriverDigit{Height: 70, Width: 240, Length: 6, MaxSkew: 0.8, DotCount: 120}
	// 3、调用NewCaptcha方法生成具体的验证码对象
	captcha := base64Captcha.NewCaptcha(driver, store)
	// 4、调用Generate()生成具体base64验证码的图片地址和id
	// id 是用于后续校验使用，后续根据id和用户输入的验证码去调用store的get方法，就可以得到你输入的验证码是否正确，正确的true，错误的false
	id, baseURL, _, err := captcha.Generate()
	if err != nil {
		response.Fail(40001, "验证生成错误!", c)
		return
	}

	response.Ok(map[string]any{"id": id, "baseURL": baseURL}, c)
}

// 3、开始校验用户输入的验证码是否是正确的
func (api *CodeApi) VerifyCaptcha(c *gin.Context) {

	type BaseCaptcha struct {
		Id   string `form:"id"`
		Code string `form:"code"`
	}
	baseCaptcha := BaseCaptcha{}
	// 开始把用户输入的id和code进行绑定
	err2 := c.ShouldBindQuery(&baseCaptcha)

	if err2 != nil {
		response.Fail(402, "参数绑定失败", c)
		return
	}

	// 开始校验验证码是否正确
	verify := store.Verify(baseCaptcha.Id, baseCaptcha.Code, true)

	if verify {
		response.Ok("success", c)
	} else {
		response.Fail(403, "您输入的验证码有误！", c)
	}

}
