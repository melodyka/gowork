package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func NetWorkStatus() bool {
	cmd := exec.Command("ping", "www.baidu.com")
	fmt.Println("NetWorkStatus Start:", time.Now().Unix())
	/*
		err := cmd.Run()

			if err != nil {
				fmt.Println(err.Error())
				return false
			} else {
				fmt.Println("Net Status , OK")
			}
	*/
	out, err2 := cmd.CombinedOutput()
	if err2 != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err2)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
	fmt.Println("NetWorkStatus End  :", time.Now().Unix())
	return true
}

func main() {
	NetWorkStatus()
	/*
		p := fastping.NewPinger()
		ra, err := net.ResolveIPAddr("ip4:icmp", os.Args[1])
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
	*/

}
