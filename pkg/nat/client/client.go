package main

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/socket"
	"golang.org/x/sys/unix"
	"time"
)

func main() {

	sock := socket.NewSocket(1234)
	sock.Connect(&unix.SockaddrInet4{
		Port: 8001,
		Addr: [4]byte{0, 0, 0, 0},
	})

	go func() {
		for {
			time.Sleep(time.Second * 3)
			sock.Write([]byte("hello, i am 8001"))
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

	sock1 := socket.NewSocket(1234)
	sock1.Connect(&unix.SockaddrInet4{
		Port: 8002,
		Addr: [4]byte{0, 0, 0, 0},
	})

	go func() {
		for {
			time.Sleep(time.Second * 3)
			sock1.Write([]byte("hello, i am 8002"))
		}
	}()

	data1 := make([]byte, 1024)
	go func() {
		for {
			_, err := sock1.Read(data1)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(" receive 8002 data:", string(data1))
		}
	}()

	time.Sleep(time.Hour * 1)
}
