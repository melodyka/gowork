package main

import (
	"fmt"
	"sync"
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

// 生产者需要生产的任务数量
const taskNum int64 = 100

func producer(wo chan<- Task) {
	var i int64
	for i = 1; i <= taskNum; i++ {
		t := Task{
			ID: i,
		}
		wo <- t
	}
	// 单个生产者就可以直接关闭通道了，关闭后，消费者任然可以消费
	close(wo)
}

func consumer(ro <-chan Task) {
	for t := range ro {
		if t.ID != 0 {
			t.run()
		}
	}
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		producer(taskCh)
	}(wg)
	for i := 0; i < int(taskNum); i++ {
		if i%100 == 0 {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				consumer(taskCh)
			}(wg)
		}
	}
	wg.Wait()
	fmt.Println("执行成功")

}
