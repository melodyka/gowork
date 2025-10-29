package main

import (
	"log"
	"sync"
	"time"

	"work"
)

// 可更换，定义数据
var names = []string{
	"www.baidu.com",
	"www.google.cn",
}

// Worker实现类型
// 可更换，定义结构体
type namePrinter struct {
	name string
}

// 可更换，定义方法
//func (n *namePrinter) Task(goid int) {
//	log.Printf("goroutineID:%d，打印名字为：%s\n", goid, n.name)
//	time.Sleep(time.Second)
//}

func (n *namePrinter) Task(goid int) {

	log.Printf("goroutineID:%d", goid)

	NetWorkStatus(n.name)
	time.Sleep(time.Second)
}

func NetWorkStatus(web string) bool {
	log.Println(web)
	/*
	out, err := exec.Command("ping", web).Output()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(out))
	}*/
	return true
	
}

/*
	func NetWorkStatus(web string) bool {
		cmd := exec.Command("ping", web)
		fmt.Println("NetWorkStatus Start:", time.Now().Unix())
		out, err2 := cmd.CombinedOutput()
		if err2 != nil {
			fmt.Printf("combined out:\n%s\n", string(out))
			log.Fatalf("cmd.Run() failed with %s\n", err2)
		}
		fmt.Printf("combined out:\n%s\n", string(out))
		fmt.Println("NetWorkStatus End  :", time.Now().Unix())
		return true
	}
*/
func main() {
	p := work.New(4)
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		for _, name := range names {
			//任务实例
			np := namePrinter{
				name: name,
			}

			go func() {
				p.Run(&np)
				wg.Done()
			}()
		}
	}
	wg.Wait()
	p.Shutdown()
}
