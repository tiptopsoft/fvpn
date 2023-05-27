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
		Port: 4000,
		Addr: [4]byte{0, 0, 0, 0},
	}

	err := sock.Connect(addr)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			time.Sleep(time.Second * 3)
			sock.Write([]byte("hello, i am 8001"))
		}
	}()

	go func() {
		for {
			data := make([]byte, 1024)
			_, err := sock.Read(data)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(" receive 8001 data:", string(data))
		}
	}()

	//sock1 := socket.NewSocket(1234)
	//sock1.Connect(&unix.SockaddrInet4{
	//	Port: 8002,
	//	Addr: [4]byte{211, 159, 225, 186},
	//})
	//
	//go func() {
	//	for {
	//		time.Sleep(time.Second * 3)
	//		sock1.Write([]byte("hello, i am 8002"))
	//	}
	//}()
	//
	//go func() {
	//	for {
	//		data1 := make([]byte, 1024)
	//		_, err := sock1.Read(data1)
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//
	//		fmt.Println(" receive 8002 data:", string(data1))
	//	}
	//}()

	time.Sleep(time.Hour * 1)
}
