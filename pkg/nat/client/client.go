package main

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/socket"
	"golang.org/x/sys/unix"
	"time"
)

func main() {

	sock := socket.NewSocket(1234)
	addr := &unix.SockaddrInet4{
		Port: 8001,
		Addr: [4]byte{0, 0, 0, 0},
	}

	addr1 := &unix.SockaddrInet4{
		Port: 8001,
		Addr: [4]byte{0, 0, 0, 0},
	}

	go func() {
		for {
			time.Sleep(time.Second * 3)
			sock.WriteToUdp([]byte("hello, i am 8001"), addr)
		}
	}()

	data := make([]byte, 1024)
	go func() {
		for {
			_, err := sock.Read(data)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(" receive 8001 data:", string(data))
		}
	}()

	//sock.Connect(&unix.SockaddrInet4{
	//	Port: 8002,
	//	Addr: [4]byte{0, 0, 0, 0},
	//})

	go func() {
		for {
			time.Sleep(time.Second * 3)
			sock.WriteToUdp([]byte("hello, i am 8002"), addr1)
		}
	}()

	data1 := make([]byte, 1024)
	go func() {
		for {
			_, err := sock.Read(data1)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(" receive 8002 data:", string(data1))
		}
	}()

	time.Sleep(time.Hour * 1)
}
