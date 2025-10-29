package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
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
	timeout  int64 // 耗时
	size     int   // 大小
	count    int   // 请求次数
	typ      uint8 = 8
	code     uint8 = 0
	SendCnt  int                   // 发送次数
	RecCnt   int                   // 接收次数
	MaxTime  int64 = math.MinInt64 // 最大耗时
	MinTime  int64 = math.MaxInt64 // 最短耗时
	SumTime  int64                 // 总计耗时
)

// ICMP 序号不能乱
type ICMP struct {
	Type        uint8  // 类型
	Code        uint8  // 代码
	CheckSum    uint16 // 校验和
	ID          uint16 // ID
	SequenceNum uint16 // 序号
}

var names = []string{
	"CANRSDTGM0B514Z",
	"cn.bing",
}

type namePing struct {
	name string
}

func (n *namePing) Task(goid int) {
	Index_ping(n.name)
}
func main() {
	p := work.New(2)
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 1; i++ {
		for _, name := range names {
			np := namePing{
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
func Index_ping(pname string) {
	timeout = 1000
	size = 32
	count = 1
	//log.SetFlags(log.Llongfile)
	//GetCommandArgs()
	// 打印帮助信息
	/*
	   if helpFlag {
	       displayHelp()
	       os.Exit(0)
	   }
	*/
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
	remoteaddr := conn.RemoteAddr()
	fmt.Printf("正在 Ping %s [%s] 具有 %d 字节的数据:\n", desIP, remoteaddr, size)
	for i := 0; i < count; i++ {
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
		// 将 data 写入连接中，
		n, err := conn.Write(data)
		if err != nil {
			log.Println(err)
			continue
		}
		// 发送数 ++
		SendCnt++
		// 接收响应
		buf := make([]byte, 1024)
		n, err = conn.Read(buf)
		//fmt.Println(data)
		if err != nil {
			log.Println(err)
			continue
		}
		// 接受数 ++
		RecCnt++
		//fmt.Println(n, err) // data：64，ip首部：20，icmp：8个 = 92 个
		// 打印信息
		t := time.Since(startTime).Milliseconds()
		fmt.Printf("来自 %d.%d.%d.%d 的回复：字节=%d 时间=%d TTL=%d\n", buf[12], buf[13], buf[14], buf[15], n-28, t, buf[8])
		MaxTime = Max(MaxTime, t)
		MinTime = Min(MinTime, t)
		SumTime += t
		time.Sleep(time.Second)
	}
	fmt.Printf("\n%s 的 Ping 统计信息:\n", remoteaddr)
	fmt.Printf("    数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%.f%% 丢失)，\n", SendCnt, RecCnt, count*2-SendCnt-RecCnt, float64(count*2-SendCnt-RecCnt)/float64(count*2)*100)
	fmt.Println("往返行程的估计时间(以毫秒为单位):")
	fmt.Printf("    最短 = %d，最长 = %d，平均 = %d\n", MinTime, MaxTime, SumTime/int64(count))

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

// GetCommandArgs 命令行参数
/*
func GetCommandArgs() {
    flag.Int64Var(&timeout, "w", 1000, "请求超时时间")
    flag.IntVar(&size, "l", 32, "发送字节数")
    flag.IntVar(&count, "n", 4, "请求次数")
    flag.BoolVar(&helpFlag, "h", false, "显示帮助信息")
    flag.Parse()
}
*/
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

/*
func displayHelp() {
    fmt.Println(`选项：
    -n count       要发送的回显请求数。
    -l size        发送缓冲区大小。
    -w timeout     等待每次回复的超时时间(毫秒)。
    -h         帮助选项`)
}
*/
