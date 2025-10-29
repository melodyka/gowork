package main

/*
import (
	"fmt"

	Viper "github.com/spf13/viper"
)

func main() {
	Viper.AddConfigPath("./")	// 文件地址
	Viper.SetConfigName("config") // 文件名称
	Viper.SetConfigType("yaml") // 文件名称的类型，这里配置文件是config.yaml

	err := Viper.ReadInConfig() // 查找并且读取配置文件信息
	if err != nil { // 如果找不到文件，则会有错误信息展示
	   panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	c := Viper.AllSettings() //得到所有的配置信息
	fmt.Printf("得到的配置信息：%v",c)

}
*/

import (
	"config"
	"fmt"
	"tool"
)

func init() {
	config.Initialize()
}

func main() {
	port := tool.GetString("app.port","5005")
	fmt.Println(port)
}

