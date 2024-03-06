/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 14:29
 */
package bbs

import "xz-go-frame/service"

type WebApiGroup struct {
	XkBbsApi
	BbsCategoryApi
}

// 公共实例---服务共享
var (
	// 创建实例，保存帖子
	//xkBbsService      = new(bbs2.XkBbsService)
	//bbsCatgoryService = new(bbs2.BBSCategoryService)
	xkBbsService      = service.ServiceGroupApp.XkBbsServiceGroup.BbsService
	bbsCatgoryService = service.ServiceGroupApp.XkBbsServiceGroup.BbsCategoryService
)
