/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 16:29
 */
package global

import (
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"xz-go-frame/commons/parse"
)

var (
	Yaml   map[string]any
	Config *parse.Config
	XZ_DB  *gorm.DB
	// 用于登录登出的缓存 go get github.com/patrickmn/go-cache
	Cache    *cache.Cache
	Log      *zap.Logger
	SugarLog *zap.SugaredLogger
)
