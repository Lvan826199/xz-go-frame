/*
@Author: 梦无矶小仔
@Date:   2024/1/11 18:08
*/
package router

import (
	"xz-go-frame/router/bbs"
	"xz-go-frame/router/code"
	"xz-go-frame/router/course"
	"xz-go-frame/router/login"
	"xz-go-frame/router/state"
	"xz-go-frame/router/sys"
	"xz-go-frame/router/video"
)

type WebRouterGroup struct {
	Course course.WebRouterGroup
	Video  video.WebRouterGroup
	Sys    sys.WebRouterGroup
	State  state.WebRouterGroup
	BBs    bbs.WebRouterGroup
	Login  login.WebRouterGroup
	Code   code.WebRouterGroup
}

var RouterWebGroupApp = new(WebRouterGroup)
