package origin

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	//ListenUDP创建一个接收目的地是本地地址laddr的UDP数据包的网络连接。net必须是"udp"、"udp4"、"udp6"；如果laddr端口为0，函数将选择一个当前可用的端口，可以用Listener的Addr方法获得该端口。返回的*UDPConn的ReadFrom和WriteTo方法可以用来发送和接收UDP数据包（每个包都可获得来源地址或设置目标地址）。
	//IPv4zero:本地地址，只能作为源地址（曾用作广播地址）
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.IPv4zero, Port: 9527})
	if err != nil {

		fmt.Println(err)
	}
	//LocalAddr返回本地网络地址
	log.Printf("本地地址：<%s> \n", listener.LocalAddr().String())
	peers := make([]net.UDPAddr, 0, 2)
	data := make([]byte, 1024)
	for {

		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {

			fmt.Println("err during read: %s", err)
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
		peers = append(peers, *remoteAddr)
		if len(peers) == 2 {

			log.Printf("进行UDP打洞，建立 %s <--> %s 的链接\n", peers[0].String(), peers[1].String())
			//WriteToUDP通过c向地址addr发送一个数据包，b为包的有效负载，返回写入的字节。
			//WriteToUDP方***在超过一个固定的时间点之后超时，并返回一个错误。在面向数据包的连接上，写入超时是十分罕见的。
			listener.WriteToUDP([]byte(peers[1].String()), &peers[0])
			listener.WriteToUDP([]byte(peers[0].String()), &peers[1])
			time.Sleep(time.Second * 8)
			log.Println("中转服务器退出，仍不影响peers间通信")
			return
		}
	}
}
