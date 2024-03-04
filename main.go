package main

import (
	"fmt"
	"xz-go-frame/initlization"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
	//  开始初始化配置文件
	initlization.InitViper()
	fmt.Println("初始化配置文件成功！")
	// 初始化日志 开发的时候建议设置成：debug ，发布的时候建议设置成：info/error
	// info --- console + file
	// error -- file
	initlization.InitLogger("debug")
	// 初始化数据库
	initlization.InitMySQL()
	// 初始化本地缓存
	initlization.InitCache()
	//开始初始化gin路由服务
	initlization.RunServer()
	fmt.Println("启动xz-go-frame后端成功")

}
