package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"runtime/pprof"
	"sync"
	"time"
	"work"
)

// tcp 报文前20个是报文头，后面的才是 ICMP 的内容。
// ICMP：组建 ICMP 首部（8 字节） + 我们要传输的内容
// ICMP 首部：type、code、校验和、ID、序号，1 1 2 2 2
// 回显应答：type = 0，code = 0
// 回显请求：type = 8, code = 0
var (
	helpFlag bool
	timeout  int64 = 1000 // 耗时
	size     int   = 32   // 大小
	count    int   = 1    // 请求次数
	typ      uint8 = 8
	code     uint8 = 0
	SendCnt  int                   // 发送次数
	RecCnt   int                   // 接收次数
	MaxTime  int64 = math.MinInt64 // 最大耗时
	MinTime  int64 = math.MaxInt64 // 最短耗时
	SumTime  int64                 // 总计耗时
	wg       sync.WaitGroup
	names3   string
	mod      int
)

// 重写Set定义flag新类型
type Value interface {
	String() string
	Set(string) error
}

type arrayFlags []string

// Value ...
func (i *arrayFlags) String() string {
	return fmt.Sprint(*i)
}

// Set 方法是flag.Value接口, 设置flag Value的方法.
// 通过多个flag指定的值， 所以我们追加到最终的数组上.
func (i *arrayFlags) Set(value string) error {
	//*i = append(*i, value)
	*i = append(*i, value)
	return nil
}

var names2 arrayFlags

// GetCommandArgs 命令行参数

func GetCommandArgs() {
	flag.Int64Var(&timeout, "w", 1000, "请求超时时间")
	flag.IntVar(&size, "l", 32, "发送字节数")
	flag.IntVar(&mod, "mod", 2, "长短模式")
	flag.Var(&names2, "addr", "hosts")
	flag.StringVar(&names3, "file", "names.txt", "hosts")
	flag.IntVar(&count, "n", 5, "请求次数")
	flag.BoolVar(&helpFlag, "h", false, "显示帮助信息")
	flag.Parse()
}

func displayHelp() {
	fmt.Println(`选项：
	-file 读取文件
	- mod ping模式
	-addr hosts   主机名
    -n count       要发送的回显请求数。
    -l size        发送缓冲区大小。
    -w timeout     等待每次回复的超时时间(毫秒)。
    -h         帮助选项`)
}

var names = []string{
	//"www.baidu.com",
	//"cn.bing.com",
	//"www.google.com",

}

// ICMP 序号不能乱
type ICMP struct {
	Type        uint8  // 类型
	Code        uint8  // 代码
	CheckSum    uint16 // 校验和
	ID          uint16 // ID
	SequenceNum uint16 // 序号
}

var s = make(chan int, count)
var r = make(chan int, count)

type namePing struct {
	name string
}

// @brief：耗时统计函数
func timeCost(start time.Time) {
	tc := time.Since(start)
	fmt.Printf("time cost = %v\n", tc)
}

func (n *namePing) Task(goid int) {
	Index_ping(n.name)
}

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

func main() {
	f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	//log.SetFlags(log.Llongfile)

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	GetCommandArgs()
	//打印帮助信息

	if helpFlag {
		displayHelp()
		os.Exit(0)
	}

	defer timeCost(time.Now())

	switch {
	case len(names) != 0:
		break
	case len(names2) != 0:
		names = names2
	case len(names3) != 0:
		names = read_txt(names3)

	}

	p := work.New(len(names))

	//for i := 0; i < 1; i++ {

	for _, name := range names {
		SendCnt, RecCnt = 0, 0
		MinTime, MaxTime, SumTime = math.MaxInt64, math.MinInt64, 0
		np := namePing{
			name: name,
		}

		conn, err := net.DialTimeout("ip:icmp", name, time.Duration(timeout)*time.Millisecond)

		if err != nil {
			log.Println(1)
			log.Println(err.Error())
			continue
		}
		conn.Close()

		go func() {
			p.Run(&np)
			//wg.Done()
		}()

		time.Sleep(time.Second)
		SendCnt = <-s
		RecCnt = <-r

		if mod == 2 {

			fmt.Printf("    数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%.f%% 丢失)，\n", SendCnt, RecCnt, count*2-SendCnt-RecCnt, float64(count*2-SendCnt-RecCnt)/float64(count)*100)
			fmt.Printf("    最短 = %d，最长 = %d，平均 = %d\n", MinTime, MaxTime, SumTime/int64(count))
		}
	}
	//}

	wg.Wait()
	p.Shutdown()
	close(s)
	close(r)
}

func Index_ping(pname string) {

	// 获取目标 IP
	desIP := pname
	//fmt.Println(desIP)
	// 构建连接

	conn, err := net.DialTimeout("ip:icmp", desIP, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer conn.Close()
	// 远程地址
	for i := 0; i < count; i++ {
		wg.Add(1)

		go Start_ping(desIP, conn, i)

	}

	wg.Wait()
	s <- SendCnt
	r <- RecCnt
}

func Start_ping(desIP string, conn net.Conn, i int) {

	// 构建请求
	icmp := &ICMP{
		Type:        typ,
		Code:        code,
		CheckSum:    uint16(0),
		ID:          uint16(i),
		SequenceNum: uint16(i),
	}
	// 将请求转为二进制流
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	// 请求的数据
	data := make([]byte, size)
	// 将请求数据写到 icmp 报文头后
	buffer.Write(data)
	data = buffer.Bytes()
	// ICMP 请求签名（校验和）：相邻两位拼接到一起，拼接成两个字节的数
	checkSum := checkSum(data)
	// 签名赋值到 data 里
	data[2] = byte(checkSum >> 8)
	data[3] = byte(checkSum)
	startTime := time.Now()
	// 设置超时时间

	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	//conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	//conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	// 将 data 写入连接中，
	n, err := conn.Write(data)

	if err != nil {
		log.Println(desIP, err)

	}

	// 发送数 ++
	SendCnt++

	// 接收响应
	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {

		log.Println(desIP, err)

	} else {

		// 接受数 ++
		RecCnt++
		//fmt.Println(n, err) // data：64，ip首部：20，icmp：8个 = 92 个
		// 打印信息
		t := time.Since(startTime).Milliseconds()

		log.Printf("(%d) %s 来自 %d.%d.%d.%d 的回复：字节=%d 时间=%d TTL=%d\n", i, desIP, buf[12], buf[13], buf[14], buf[15], n-28, t, buf[8])
		MaxTime = Max(MaxTime, t)
		MinTime = Min(MinTime, t)
		SumTime += t
	}

	//time.Sleep(time.Second)

	wg.Done()
}

// 求校验和
func checkSum(data []byte) uint16 {
	// 第一步：两两拼接并求和
	length := len(data)
	index := 0
	var sum uint32
	for length > 1 {
		// 拼接且求和
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		length -= 2
		index += 2
	}
	// 奇数情况，还剩下一个，直接求和过去
	if length == 1 {
		sum += uint32(data[index])
	}
	// 第二部：高 16 位，低 16 位 相加，直至高 16 位为 0
	hi := sum >> 16
	for hi != 0 {
		sum = hi + uint32(uint16(sum))
		hi = sum >> 16
	}
	// 返回 sum 值 取反
	return uint16(^sum)
}

func Max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
