/*
@Author: 梦无矶小仔
@Date:   2024/1/11 18:08
*/
package router

import "xz-go-frame/router/video"

type WebRouterGroup struct {
	Video video.WebRouterGroup
}

var RouterWebGroupApp = new(WebRouterGroup)
