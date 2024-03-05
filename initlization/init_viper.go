/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 14:43
 */
package initlization

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"xz-go-frame/global"
)

func GetEnvInfo(env string) string {
	/*
		例如，如果你有一个配置项叫 server.port，在你的应用程序中通过 viper.GetInt("server.port") 获取端口号。
		如果你调用了 viper.AutomaticEnv() 并且设置了环境变量 SERVER_PORT，Viper 会优先使用 SERVER_PORT 环境变量的值，而不是配置文件中的值。
	*/
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func InitViper(args map[string]any) {
	// 获取项目的执行路径
	path, err := os.Getwd() // D:\Z_Enviroment\GoWorks\src\xz-go-frame
	if err != nil {
		panic(err)
	}
	// 初始化一个viper解析配置对象
	config := viper.New()
	//// 开始设置从哪个目录下去找yaml文件
	//config.AddConfigPath(path + "/configfile") // 设置读取的文件路径
	//// 设置配置文件的名字
	//config.SetConfigName("application") // 设置读取的文件名
	//// 设置配置文件的后缀
	//config.SetConfigType("yaml") // 设置文件类型
	//fmt.Println("你激活的环境是：" + GetEnvInfo("env"))
	global.Log.Info("你激活的环境是：", zap.String("env", GetEnvInfo("env")))
	config.SetConfigFile(path + "/configfile/application-" + GetEnvInfo("env") + ".yaml")

	// 尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	// 监控配置文件的变化
	config.WatchConfig()
	// 当配置文件发生变化时，执行回调函数
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件发生变化：", e.Name)
		// 把改变的值重新放入到config配置中去
		if err = config.Unmarshal(&global.Config); err != nil {
			fmt.Println(err)
		}
	})

	// 这里才是把yaml配置文件解析放入到Config对象的过程---map---config
	if err = config.Unmarshal(&global.Config); err != nil {
		fmt.Println(err)
	}

	// 遍历打印文件读取来的内容
	keys := config.AllKeys()
	dataMap := make(map[string]any)
	for _, key := range keys {
		fmt.Println("yaml存在的key是: "+key+" : ", config.Get(key))
		dataMap[key] = config.Get(key)
	}
	// 用环境变量覆盖
	// 命令行参数覆盖 boot
	port := args["server.port"].(int)
	if port != -1 {
		dataMap["server.port"] = port
	}

	global.Yaml = dataMap

	//// 打印文件读取出来的内容
	//fmt.Println(len(config.Get("mysql.database.host").(string))) // 9
	//fmt.Println(config.Get("mysql.database.host"))               // 127.0.0.1
	//fmt.Println(config.Get("mysql.database.user"))               // root
	//fmt.Println(config.Get("mysql.database.dbname"))             // test
	//fmt.Println(config.Get("mysql.database.pwd"))                // 123456
	//fmt.Println(config.Get("server.port"))                       // 8888
	//fmt.Println(config.Get("server.cookiname"))                  // mysession

}
