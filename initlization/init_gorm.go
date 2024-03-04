/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 16:32
 */
package initlization

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"xz-go-frame/commons/orm"
	"xz-go-frame/global"
)

func InitMySQL() {
	// 初始化gorm的日志
	newLogger := logger2.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger2.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger2.Info, // Log level
			IgnoreRecordNotFoundError: false,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,        // Don't include params in the SQL log
			Colorful:                  true,         // Disable color
		},
	)

	m := global.Config.Database.Mysql
	//fmt.Println(m.Dsn())
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               m.Dsn(), // DNS data source name
		DefaultStringSize: 191,     //string类型字段的默认长度
		//DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		// GORM 定义了这些日志级别：Silent、Error、Warn、Info
		//Logger: logger.Default.LogMode(logger.Info),
		Logger: newLogger,
	})

	// 如果报错，请检查数据库配置
	if err != nil {
		global.Log.Error("数据连接出错了", zap.String("error", err.Error()))
		panic("数据连接出错了" + err.Error()) // 把程序直接阻断，把数据连接好了在启动
	}

	global.XZ_DB = db // 数据库信息全局变量

	// 初始化数据库
	orm.RegisterTable()

	fmt.Println("数据库初始化完成,开始运行：", db)
	// 日志输出
	global.Log.Debug("数据库连接成功。开始运行", zap.Any("db", db))
}
