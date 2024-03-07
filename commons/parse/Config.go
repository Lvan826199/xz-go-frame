/*
* @Author: 梦无矶小仔
* @Date:   2023/12/27 17:02
 */
package parse

// 配置文件解析的总入口
type Config struct {
	// 数据库
	Database Database `mapstructure:"database" json:"database" yaml:"database"`
	Ksd      Ksd      `mapstructure:"ksd" json:"ksd" yaml:"ksd"`
	NoSQL    NoSQL    `mapstructure:"nosql" json:"nosql" yaml:"nosql"`
	Local    Local    `mapstructure:"local" json:"local" yaml:"local"`
}

type Local struct {
	Path       string `mapstructure:"path" json:"path" yaml:"path"`                   // 本地文件访问路径
	Fileserver string `mapstructure:"fileserver" json:"fileserver" yaml:"fileserver"` // 本地文件访问路径
	StorePath  string `mapstructure:"store-path" json:"store-path" yaml:"store-path"` // 本地文件存储路径
}
