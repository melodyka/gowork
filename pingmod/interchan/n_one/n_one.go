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

func producer(wo chan<- Task, startNum int64, nums int64) {
	var i int64
	for i = startNum; i <= startNum+nums; i++ {
		t := Task{
			ID: i,
		}
		wo <- t
	}
	// 单个生产者就可以直接关闭通道了，关闭后，消费者任然可以消费
	//close(wo)
}

func consumer(ro <-chan Task) {
	for t := range ro {
		if t.ID != 0 {
			t.run()
		}
	}
}

func main() {
	var nums int64 = 10
	wg := &sync.WaitGroup{}
	pwg := &sync.WaitGroup{}
	wg.Add(1)
	var i int64
	for i = 0; i < taskNum; i += nums {
		if i >= taskNum {
			break
		}
		wg.Add(1)
		pwg.Add(1)
		go func(i int64) {
			defer wg.Done()
			defer pwg.Done()
			producer(taskCh, i, nums)
		}(i)
	}

	go func() {
		defer wg.Done()
		consumer(taskCh)
	}()

	pwg.Wait()
	go close(taskCh)
	wg.Wait()
	fmt.Println("执行成功")

}
