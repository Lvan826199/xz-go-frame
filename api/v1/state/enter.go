/*
* @Author: 梦无矶小仔
* @Date: 2024/3/7 16:38
 */
package state

import "xz-go-frame/service"

type WebApiGroup struct {
	UserStateApi
}

// 公共实例---服务共享
var (
	userStatService = service.ServiceGroupApp.UserStateServiceGroup.UserStateService
)
