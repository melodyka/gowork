package main

import (
	"fmt"
	"sync"
	"time"
)

// Task 任务
type Task struct {
	ID int
}

// 消费任务
func (t *Task) run() {
	fmt.Println(t.ID)
}


func createProducerConsumerSystem(numProducers int, numConsumers int, maxTasksPerProducer int, maxTasksPerConsumer int) {
	// 缓冲池
	taskCh := make(chan Task, 10)

	// 停止运行的信号
	done := make(chan struct{})

	var wg sync.WaitGroup

	// 启动多个生产者
	for i := 0; i < numProducers; i++ {
		go func(id int) {
			var i int
			for i = 1; i <= maxTasksPerProducer; i++ {
				t := Task{
					ID: i,
				}
				// 可以防止因为生产者阻塞，而导致关闭信号无法关闭
				select {
				case taskCh <- t:
					fmt.Printf("生产者 %d 生产了: %d\n", id, i)
				case <-done:
					fmt.Printf("生产者 %d 退出\n", id)
					return
				}
			}
			fmt.Printf("生产者 %d 完成生产任务\n", id)
		}(i)
	}

	// 启动多个消费者
	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var consumedCount int = 0
			for {
				select {
				case t, ok := <-taskCh:
					if !ok {
						fmt.Printf("消费者 %d 通道已关闭，实际消费了 %d 个任务，退出\n", id, consumedCount)
						return
					}
					if t.ID != 0 {
						t.run()
						consumedCount++
						// 如果该消费者已经消费了指定数量的任务，则退出
						if consumedCount >= maxTasksPerConsumer {
							fmt.Printf("消费者 %d 已消费 %d 个任务，达到上限，退出\n", id, consumedCount)
							return
						}
					}
				case <-done:
					fmt.Printf("消费者 %d 收到退出信号，实际消费了 %d 个任务\n", id, consumedCount)
					return
				}
			}
		}(i)
	}

	// 等待所有生产者完成生产后关闭通道
	time.Sleep(time.Second * 2) // 给生产者时间完成生产
	close(taskCh)               // 关闭任务通道，通知消费者没有更多任务了

	wg.Wait() // 等待所有消费者完成
	fmt.Println("执行成功")
	fmt.Printf("启动系统: %d个生产者, %d个消费者, 每个生产者生产%d个任务, 每个消费者消费%d个任务\n", 
	numProducers, numConsumers, maxTasksPerProducer, maxTasksPerConsumer)
	
}


func main() {
	// 可以通过参数设置消费者数量和每个消费者消费的任务数
	createProducerConsumerSystem(3, 5, 3, 2)	
}



