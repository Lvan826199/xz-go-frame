/*
@Author: 梦无矶小仔
@Date:   2024/1/12 16:38
*/
package v1

import "xz-go-frame/api/v1/video"

type WebApiGroup struct {
	video.WebApiGroup
}

var WebApiGroupApp = new(WebApiGroup)
