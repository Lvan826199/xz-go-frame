/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 15:23
 */
package service

import (
	"xz-go-frame/service/bbs"
	"xz-go-frame/service/state"
)

// 实例创建聚合
type ServicesGroup struct {
	//SyserviceGroup        sys.ServiceGroup
	XkBbsServiceGroup bbs.ServiceGroup
	//XkVideoServiceGroup   video.ServiceGroup
	UserStateServiceGroup state.ServiceGroup
	//UserServiceGroup      user.ServiceGroup
}

// 单例设计模式
var ServiceGroupApp = new(ServicesGroup)
