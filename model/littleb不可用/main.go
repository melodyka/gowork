package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func read_txt(path string) (names []string) {

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		name := scanner.Text()
		names = append(names, name)
		//fmt.Println("Name:", name)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}
	return names
}

func pingIP(ip string,domain string) {
	defer wg.Done()

	conn, err := net.DialTimeout("icmp", ip, time.Second*2)
	if err != nil {
		fmt.Println(ip,domain, "is unreachable")
		return
	}
	defer conn.Close()

	fmt.Println(ip, "is reachable")
}

func main() {
	names := read_txt("names.txt")

	for _, name := range names {
		ipAddr, err := net.ResolveIPAddr("ip", name)
		if err != nil {
			fmt.Println("Error resolving IP address for", name, ":", err)
			continue
		}

		wg.Add(1)
		go pingIP(ipAddr.String(),name)
	}

	wg.Wait()
}
