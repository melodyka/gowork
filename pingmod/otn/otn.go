package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

// ---------- ICMP 结构 ----------
type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	ID          uint16
	SequenceNum uint16
}

const (
	timeout = 1000 // ms
	size    = 32   // payload 大小
	count   = 4    // 每个目标发送次数
)

var (
	typ uint8  = 8  // Echo
	code uint8 = 0
	pid  uint16 = uint16(os.Getpid() & 0xffff)
)

// ---------- 目标 ----------
var targets = []string{
	"8.8.8.8",
	"www.baidu.com",
}

// ---------- 任务 ----------
type PingTask struct{ target string }

// ---------- 主逻辑 ----------
func main() {
	log.SetFlags(log.Ldate | log.Ltime)
defer timeCost(time.Now())
	runtime.GOMAXPROCS(runtime.NumCPU()) // 让调度器可以用全部核心
	pool := newWorkerPool(runtime.NumCPU())

	for _, t := range targets {
		pool.Submit(&PingTask{target: t})
	}
	pool.Wait()
}
func timeCost(start time.Time) {
	tc := time.Since(start)
	fmt.Printf("time cost = %v\n", tc)
}
// ---------- WorkerPool ----------
type workerPool struct {
	taskCh chan *PingTask
	wg     sync.WaitGroup
}

func newWorkerPool(n int) *workerPool {
	p := &workerPool{taskCh: make(chan *PingTask)}
	p.wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			// 每个 worker 绑定一条 OS 线程
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()

			for task := range p.taskCh {
				task.Run()
			}
			p.wg.Done()
		}()
	}
	return p
}

func (p *workerPool) Submit(t *PingTask) { p.taskCh <- t }
func (p *workerPool) Wait() {
	close(p.taskCh)
	p.wg.Wait()
}

// ---------- 任务执行 ----------
func (t *PingTask) Run() {
	var wg sync.WaitGroup
	wg.Add(count)

	for seq := 0; seq < count; seq++ {
		go func(seq int) {
			defer wg.Done()
			sendOnePing(t.target, seq)
		}(seq)
	}
	wg.Wait()
}

// ---------- 单次 Ping ----------
func sendOnePing(target string, seq int) {
	conn, err := net.DialTimeout("ip4:icmp", target, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return
	}
	defer conn.Close()

	// 构造 ICMP
	icmp := ICMP{
		Type:        typ,
		Code:        code,
		CheckSum:    0,
		ID:          pid,
		SequenceNum: uint16(seq),
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, icmp)
	buf.Write(make([]byte, size))
	pkt := buf.Bytes()
	binary.BigEndian.PutUint16(pkt[2:4], checkSum(pkt))

	start := time.Now()
	conn.SetWriteDeadline(time.Now().Add(timeout * time.Millisecond))
	conn.SetReadDeadline(time.Now().Add(timeout * time.Millisecond))

	if _, err = conn.Write(pkt); err != nil {
		return
	}

	reply := make([]byte, 1500)
	n, err := conn.Read(reply)
	if err != nil || n < 28 || reply[20] != 0 {
		return
	}

	rtt := time.Since(start).Milliseconds()
	srcIP := net.IPv4(reply[12], reply[13], reply[14], reply[15])
	log.Printf("(%d) %s 来自 %s 的回复：字节=%d 时间=%d TTL=%d",
		seq, target, srcIP, size, rtt, reply[8])
}

// ---------- 校验和 ----------
func checkSum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for sum>>16 != 0 {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	return ^uint16(sum)
}