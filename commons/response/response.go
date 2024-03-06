/*
@Author: 梦无矶小仔
@Date:   2024/1/11 14:53
*/
package response

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"net/http"
	"strconv"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	VALIDATOR_MAP        = map[string]string{"code": "701", "msg": "属性验证有误"}
	BINDING_PAMATERS_MAP = map[string]string{"code": "702", "msg": "参数绑定有误"}
)

func Result(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

var (
	CODE       = 20000
	ERROR_CODE = 40001
	MSG        = "success"
)

/*
Ok
请求响应成功
*/
func Ok(data any, c *gin.Context) {
	Result(CODE, MSG, data, c)
}

/*
OkSuccess 请求响应成功
*/
func OkSuccess(c *gin.Context) {
	Result(CODE, MSG, nil, c)
}

/*
Fail
请求响应失败（无响应数据）
*/
func Fail(code int, msg string, c *gin.Context) {
	Result(code, msg, map[string]any{}, c)
}

/*
FailWithData
请求响应失败（有响应数据）
*/
func FailWithData(code int, msg string, data any, c *gin.Context) {
	Result(code, msg, data, c)
}

func FailWithMessage(msg string, c *gin.Context) {
	Result(ERROR_CODE, msg, map[string]any{}, c)
}

func FailWithPermission(c *gin.Context) {
	Result(80001, "权限不足", map[string]any{}, c)
}

func FailWithError(err error, c *gin.Context) {
	Result(ERROR_CODE, err.Error(), map[string]any{}, c)
}

func FailWithValidatorData(validate *validate.Validation, c *gin.Context) {
	all := validate.Errors.All()
	one := validate.Errors.One()
	code, _ := strconv.ParseInt(VALIDATOR_MAP["code"], 10, 32)
	Result(int(code), one, all, c)
}

func FailWithBindParams(c *gin.Context) {
	code, _ := strconv.ParseInt(BINDING_PAMATERS_MAP["code"], 10, 32)
	Result(int(code), BINDING_PAMATERS_MAP["msg"], nil, c)
}
