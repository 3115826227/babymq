package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		ID   uint32 `json:"id"`
		Name string `json:"name"`
		Port int    `json:"port"`
	} `json:"server"`
}

func setConfigYaml() {
	//设置配置文件的名字
	viper.SetConfigName("config")
	//设置配置文件读取路径
	viper.AddConfigPath("./etc") //idea跑的时候直接读取项目etc/目录
	//viper.AddConfigPath("/etc")  //部署到docker容器中挂载到/etc目录下
	//设置配置文件类型
	viper.SetConfigType("yaml")
}

var (
	Cfg Config
)

func InitConfig() {
	setConfigYaml()
	//读取配置文件内容
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}
	Cfg = c
}
