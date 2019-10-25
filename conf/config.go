package conf

import (
	"github.com/spf13/viper"
)

//配置项
var Con *viper.Viper

func init() {
	Con = viper.New()
	//添加读取的配置文件路径
	Con.AddConfigPath("./conf")
	//设置读取的配置文件
	Con.SetConfigName("conf")
	//设置配置文件类型
	Con.SetConfigType("yaml")
	if err := Con.ReadInConfig(); err != nil {
		panic(err)
	}
}
