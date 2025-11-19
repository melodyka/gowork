package main

import (
	"fmt"
	"time"
)

// Task 任务
type Task struct {
	ID int64
}

// 消费任务
func (t *Task) run() {
	fmt.Println(t.ID)
}

// 缓冲池
var taskCh = make(chan Task, 10)

// 停止运行的信号
var done = make(chan struct{})

// 生产者需要生产的任务数量
const taskNum int64 = 10

func producer(wo chan<- Task, done chan struct{}) {
	var i int64
	for {
		if i >= taskNum {
			i = 0
		}
		i++
		t := Task{
			ID: i,
		}

		select {
		case wo <- t:
		case <-done:
			fmt.Println("生产者退出")
			return
		}
	}
}

// 单个生产者就可以直接关闭通道了，关闭后，消费者任然可以消费
//close(wo)

func consumer(ro <-chan Task, done chan struct{}) {
	for {
		select {
		case t := <-ro:
			if t.ID != 0 {
				t.run()
			}
		case <-done:
			for t := range ro {
				if t.ID != 0 {
					t.run()
				}
			}
			fmt.Println("消费者退出")
			return
		}
	}
}

func main() {

	go producer(taskCh, done)
	go producer(taskCh, done)
	go producer(taskCh, done)

	go consumer(taskCh, done)
	go consumer(taskCh, done)
	go consumer(taskCh, done)

	time.Sleep(time.Second * 5)
	close(done)
	close(taskCh)

	fmt.Println("执行成功")

}
