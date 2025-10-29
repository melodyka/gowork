package tool

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var Viper *viper.Viper

type StrMap map[string]interface{}

func init() {
   //初始化viper库
   Viper = viper.New()
   Viper.AddConfigPath("./")
   Viper.SetConfigName("config") // name of config file (without extension)
   Viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
   //Viper.SetConfigFile("config.yaml")   // path to look for the config file in

   err := Viper.ReadInConfig() // Find and read the config file
   if err != nil { // Handle errors reading the config file
      panic(fmt.Errorf("Fatal error config file: %s \n", err))
   }
}

func Env(name string, defaultValue ...interface{}) interface{} {
   if len(defaultValue) > 0{
      return Get(name,defaultValue[0])
   }
   return Get(name)
}
//添加到显式值中，可以优先获取
func Add(name string, configuration map[string]interface{}) {
   //如果某个键通过set()这个函数设置了的话，那么这个值的优先级最高。
   Viper.Set(name,configuration)
}
//根据传递的key得到配置文件的值
func Get(key string, defaultValue ...interface{}) interface{} {
   value := Viper.Get(key)
   if value == nil {
      return defaultValue[0]
   }
   return value
}
//转换成string类型
func GetString(key string, defaultValue ...interface{}) string {
   return cast.ToString(Get(key, defaultValue...))
}
//转换成int类型
func GetInt(key string, defaultValue ...interface{}) int {
   return cast.ToInt(Get(key, defaultValue...))
}
//转换成uint类型
func GetUint(key string, defaultValue ...interface{}) uint {
   return cast.ToUint(Get(key, defaultValue...))
}
//转换成int64类型
func GetInt64(key string, defaultValue ...interface{}) int64 {
   return cast.ToInt64(Get(key, defaultValue...))
}
//转换成bool类型
func GetBool(key string, defaultValue ...interface{}) bool {
   return cast.ToBool(Get(key, defaultValue...))
}
