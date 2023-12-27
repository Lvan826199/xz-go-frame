/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 17:02
 */
package parse

// Database 配置展示
// database:--------------------------------struct
//
//	mysql: --------------------------------struct
//		host: 127.0.0.1 -------------------field
//		port: 3306 -------------------field
//		dbname: xz-go-frame-db -------------------field
//		username: root -------------------field
//		password: 123456 -------------------field
//		config: charset=utf8&parseTime=True&loc=Local -------------------field
type Database struct {
	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
}

type Mysql struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	Dbname   string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Config   string `mapstructure:"config" json:"config" yaml:"config"`
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}
