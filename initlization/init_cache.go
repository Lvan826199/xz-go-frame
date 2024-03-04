/*
@Author: 梦无矶小仔
@Date:   2024/3/4 17:26
*/
package initlization

import (
	"github.com/patrickmn/go-cache"
	"time"
	"xz-go-frame/global"
)

func InitCache() {
	c := cache.New(5*time.Minute, 24*60*time.Minute)
	global.Cache = c
}
