package main

import "github.com/sirupsen/logrus"

func main() {
	//创建一个实例
	log := *logrus.New()	
	//设置为json格式
	log.SetFormatter(&logrus.JSONFormatter{
	TimestampFormat: "2006-01-02 15:04:05",
	})
	//设置日志等级
	log.SetLevel(logrus.InfoLevel)
	//写入日志
	log.WithFields(logrus.Fields{
	"name": "一颗蛋蛋",
	}).Info("这里是logrus快速使用")
}