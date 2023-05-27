package main

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/socket"
	"golang.org/x/sys/unix"
	"net"
)

func main() {

	//if len(os.Args) < 2 {
	//	panic("should give a port")
	//}
	//port, _ := strconv.Atoi(os.Args[1])
	//listener, err := net.ListenUDP("udp", &net.UDPAddr{
	//	IP: net.IPv4zero, Port: port})
	//if err != nil {
	//
	//	fmt.Println(err)
	//}
	////LocalAddr返回本地网络地址
	//log.Printf("本地地址：<%s> \n", listener.LocalAddr().String())
	//data := make([]byte, 1024)
	//for {
	//
	//	n, remoteAddr, err := listener.ReadFromUDP(data)
	//	if err != nil {
	//		fmt.Println("err during read: %s", err)
	//	}
	//	log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
	//	//WriteToUDP通过c向地址addr发送一个数据包，b为包的有效负载，返回写入的字节。
	//	//WriteToUDP方***在超过一个固定的时间点之后超时，并返回一个错误。在面向数据包的连接上，写入超时是十分罕见的。
	//	listener.WriteToUDP([]byte(fmt.Sprintf("this is port: %d", port)), remoteAddr)
	//}

	sock := socket.NewSocket(4000)

	for {
		data := make([]byte, 1024)
		_, err := sock.Read(data)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(data))
	}

}

func ResolveAddr(address string) (unix.Sockaddr, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	ip := [4]byte{}
	copy(ip[:], addr.IP.To4())

	result := &unix.SockaddrInet4{
		Port: addr.Port,
		Addr: ip,
	}

	return result, nil
}
