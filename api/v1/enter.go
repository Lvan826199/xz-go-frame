/*
@Author: 梦无矶小仔
@Date:   2024/1/12 16:38
*/
package v1

import (
	"xz-go-frame/api/v1/bbs"
	"xz-go-frame/api/v1/code"
	"xz-go-frame/api/v1/login"
	"xz-go-frame/api/v1/state"
	"xz-go-frame/api/v1/video"
)

type WebApiGroup struct {
	Video video.WebApiGroup
	Code  code.WebApiGroup
	//Sys    sys.WebApiGroup
	State state.WebApiGroup
	//Upload upload.WebApiGroup
	Bbs   bbs.WebApiGroup
	Login login.WebApiGroup
}

var WebApiGroupApp = new(WebApiGroup)
