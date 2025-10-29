package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
	"work"

	"github.com/tatsushid/go-fastping"
)

// 可更换，定义数据
var names = []string{
	"www.baidu.com",

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

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", web)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
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
	p := work.New(2)
	var wg sync.WaitGroup
	wg.Add(2)

	for i := 0; i < 2; i++ {
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
