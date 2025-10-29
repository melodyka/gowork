package main

//[Hook](https://www.cnblogs.com/rickiyang/p/11074164.html)
import (
	"hook"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)


func initLog() {
	uuids := uuid.NewV1()
	log.AddHook(hook.NewTraceIdHook(uuids.String() +" "))
}

func main() {
	initLog()
	log.WithFields(log.Fields{
		"age": 12,
		"name":   "xiaoming",
		"sex": 1,
	}).Info("小明来了")

	log.WithFields(log.Fields{
		"age": 13,
		"name":   "xiaohong",
		"sex": 0,
	}).Error("小红来了")

	log.WithFields(log.Fields{
		"age": 14,
		"name":   "xiaofang",
		"sex": 1,
	}).Fatal("小芳来了")
}
