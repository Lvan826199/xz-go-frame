package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"xz-go-frame/initlization"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
	var env string
	var port int
	// flag 包用于解析命令行参数。
	flag.StringVar(&env, "env", "dev", "环境标识")
	flag.IntVar(&port, "server.port", -1, "测试端口")
	flag.Parse()
	args := map[string]any{"env": env, "server.port": port}
	// 设置环境变量
	viper.SetDefault("env", env)
	//  开始初始化配置文件
	initlization.InitViper(args)
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
