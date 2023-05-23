package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var tag string

const HAND_SHAKE_MSG = "我是打洞消息"

func main() {

	if len(os.Args) < 2 {
		//Args保管了命令行参数，第一个是程序名。
		fmt.Println("请输入一个客户端标志")
		os.Exit(0) //Exit让当前程序以给出的状态码code退出。一般来说，状态码0表示成功，非0表示出错。程序会立刻终止，defer的函数不会被执行
	}
	//当前进程标记字符串，便于显示
	tag = os.Args[1]
	p := os.Args[2]
	port, _ := strconv.Atoi(p)
	//UDPAddr代表一个UDP终端地址
	//IPv4zero本地地址，只能作为源地址（曾用作广播地址）
	srcAddr := &net.UDPAddr{
		IP: net.IPv4zero, Port: port} //注意端口必须固定
	//ParseIP将s解析为IP地址，并返回该地址。如果s不是合法的IP地址文本表示，ParseIP会返回nil。
	//字符串可以是小数点分隔的IPv4格式（如"74.125.19.99"）或IPv6格式（如"2001:4860:0:2001::68"）格式。
	dstAddr := &net.UDPAddr{
		IP: net.ParseIP("211.159.225.186"), Port: 9527}
	//DialTCP在网络协议net上连接本地地址laddr和远端地址raddr。net必须是"udp"、"udp4"、"udp6"；如果laddr不是nil，将使用它作为本地地址，否则自动选择一个本地地址。
	//(conn)UDPConn代表一个UDP网络连接，实现了Conn和PacketConn接口
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {

		fmt.Println(err)
	}
	if _, err = conn.Write([]byte("hello,I'm new peer:" + tag)); err != nil {

		log.Panic(err)
	}
	data := make([]byte, 1024)
	//ReadFromUDP从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址。
	//ReadFromUDP方***在超过一个固定的时间点之后超时，并返回一个错误。
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {

		fmt.Printf("error during read: %s", err)
	}
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	fmt.Printf("local:%s server:%s another:%s\n", srcAddr, remoteAddr, &anotherPeer)
	//开始打洞
	bidirectionHole(srcAddr, &anotherPeer)
}

func parseAddr(addr string) net.UDPAddr {

	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{

		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {

	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {

		fmt.Println("send handshake:", err)
	}
	go func() {

		for {

			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + tag + "]")); err != nil {

				log.Println("send msg fail", err)
			}
		}
	}()

	for {

		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {

			log.Printf("error during read:%s\n", err)
		} else {

			log.Printf("收到数据：%s\n", data[:n])
		}
	}
}
