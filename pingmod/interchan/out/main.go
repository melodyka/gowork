package main

import (
	"fmt"
	"sync"
	"time"
)

// Out 输出
type Out struct {
	data chan interface{}
}

// 单例模式
var (
	out     *Out
	once    sync.Once
)

// NewOut 初始化
func NewOut() *Out {
	once.Do(func() {
		out = &Out{
			data: make(chan interface{}, 65535), // 这里必须设置缓冲区
		}
	})
	return out
}

// Println out 的写入方法
func Println(i interface{}) {
	// 确保out已初始化
	NewOut()
	out.data <- i
}

// OutPut 将 out 内所有数据全部输出
func (o *Out) OutPut() {
	for {
		select {
		case i := <-o.data:
			fmt.Println(i)
		}
	}
}

// 启动输出协程
func StartOutput() {
	go NewOut().OutPut()
}

func main() {
	// 启动输出协程
	StartOutput()
	
	// 测试输出
	for i := 1; i <= 10; i++ {
		Println(fmt.Sprintf("消息 %d", i))
		time.Sleep(time.Millisecond * 100)
	}
	
	// 让程序运行一段时间以便看到输出
	time.Sleep(time.Second * 2)
}



