/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 18:57
 */
package orm

import (
	"xz-go-frame/global"
	"xz-go-frame/model/entity/bbs"
	"xz-go-frame/model/entity/jwt"
	"xz-go-frame/model/entity/sys"
	"xz-go-frame/model/entity/user"
	"xz-go-frame/model/entity/video"
)

func RegisterTable() {
	db := global.XZ_DB
	// 注册和声明model
	// 用户表注册
	db.AutoMigrate(&user.XzUser{})
	db.AutoMigrate(&user.XzUserAuthor{})

	// 系统用户，角色，权限表
	db.AutoMigrate(sys.SysApis{})
	db.AutoMigrate(sys.SysMenus{})
	db.AutoMigrate(sys.SysRoleApis{})
	db.AutoMigrate(sys.SysRoleMenus{})
	db.AutoMigrate(sys.SysRoles{})
	db.AutoMigrate(sys.SysUserRoles{})
	db.AutoMigrate(sys.SysUser{})

	// 视频表
	db.AutoMigrate(video.XkVideo{})
	db.AutoMigrate(video.XkVideoCategory{})
	db.AutoMigrate(video.XkVideoChapterLesson{})

	// 社区
	db.AutoMigrate(bbs.XkBbs{})
	db.AutoMigrate(bbs.BbsCategory{})

	// JWT黑名单表
	db.AutoMigrate(jwt.JwtBlacklist{})
}
