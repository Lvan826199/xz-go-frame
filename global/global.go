/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 16:29
 */
package global

import (
	"gorm.io/gorm"
	"xz-go-frame/commons/parse"
)

var (
	Yaml   map[string]any
	Config *parse.Config
	XZ_DB  *gorm.DB
)
