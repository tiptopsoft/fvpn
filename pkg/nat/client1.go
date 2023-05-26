package main

import (
	"fmt"
	"golang.org/x/sys/unix"
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
	tag = os.Args[1]
	sock := NewSocket(6061)
	err := sock.Connect(&unix.SockaddrInet4{
		Port: 9527,
		Addr: [4]byte{211, 159, 225, 186},
	})

	if err != nil {
		fmt.Println(err)
	}

	if _, err = sock.Write([]byte("hello,I'm new peer:" + tag)); err != nil {

		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := sock.ReadFromUdp(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	anotherPeer := parseAddr(string(data[:n]))
	fmt.Printf("server:%v another:%v\n", remoteAddr, anotherPeer)
	go func() {

		for {
			time.Sleep(time.Second * 5)
			if _, err = sock.Write([]byte("hahahahahaha,I'm new peer:" + tag)); err != nil {

				log.Panic(err)
			}
		}
	}()
	go func() {

		for {

			data := make([]byte, 1024)
			n, _, err := sock.ReadFromUdp(data)
			if err != nil {

				log.Printf("error during read:%s\n", err)
			} else {

				log.Printf("原sock收到数据：%s\n", data[:n])
			}
		}
	}()
	bidirectionHole(&anotherPeer)
}

func parseAddr(addr string) net.UDPAddr {

	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{

		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func bidirectionHole(anotherAddr *net.UDPAddr) {
	sock := NewSocket(0)
	addr := &unix.SockaddrInet4{
		Port: anotherAddr.Port,
		Addr: [4]byte{},
	}
	copy(addr.Addr[:], anotherAddr.IP.To4())
	err := sock.Connect(addr)
	if err != nil {
		fmt.Println("connnect failed:", err)
	}

	if _, err := sock.Write([]byte(HAND_SHAKE_MSG)); err != nil {
		fmt.Println("send handshake:", err)
	}
	go func() {

		for {

			time.Sleep(10 * time.Second)
			if err = sock.WriteToUdp([]byte("from ["+tag+"]"), addr); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()

	for {

		data := make([]byte, 1024)
		n, _, err := sock.ReadFromUdp(data)
		if err != nil {

			log.Printf("error during read:%s\n", err)
		} else {

			log.Printf("收到数据：%s\n", data[:n])
		}
	}
}

// Socket use to wrap fd
type Socket struct {
	Fd int
}

func (socket Socket) ReadFromUdp(bytes []byte) (n int, addr unix.Sockaddr, err error) {
	return unix.Recvfrom(socket.Fd, bytes, 0)
}

func (socket Socket) WriteToUdp(bytes []byte, addr unix.Sockaddr) (err error) {
	return unix.Sendto(socket.Fd, bytes, 0, addr)
}

func (socket Socket) Read(bytes []byte) (n int, err error) {
	return unix.Read(socket.Fd, bytes)
}

func (socket Socket) Write(bytes []byte) (n int, err error) {
	return unix.Write(socket.Fd, bytes)
}

func (socket Socket) Close() error {
	return unix.Close(socket.Fd)
}

func NewSocket(port int) Socket {
	fd, _ := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)

	if port != 0 {
		addr := unix.SockaddrInet4{Port: port}
		copy(addr.Addr[:], net.IPv4zero.To4())
		unix.Bind(fd, &addr)
	}

	return Socket{Fd: fd}
}

func (socket Socket) Connect(addr unix.Sockaddr) error {

	return unix.Connect(socket.Fd, addr)
}

func (socket Socket) Listen(addr unix.Sockaddr) error {
	return unix.Bind(socket.Fd, addr)
}
