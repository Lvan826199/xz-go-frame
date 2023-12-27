/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 16:32
 */
package initlization

import (
	"fmt"
	"xz-go-frame/global"
)

func InitMySQL() {
	m := global.Config.Database.Mysql
	fmt.Println(m.Dsn())
}
