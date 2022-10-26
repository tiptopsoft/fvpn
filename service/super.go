package service

import (
	"fmt"
	"net"
	"sync"
)

var limitChan = make(chan int, 1000)

// udp key : mac_group value:addr
var m sync.Map

// Node super node for net, and for user create star
func start(listen int) error {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: listen,
	})

	if err != nil {
		return err
	}

	defer conn.Close()
	for {
		limitChan <- 1
		go handleUdp(conn)
	}

}

func handleUdp(conn *net.UDPConn) {

	data := make([]byte, 1024)
	_, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println(err)
	}

	_, err = conn.WriteToUDP(data, addr)
	if err != nil {
		fmt.Println("super write failed.")
	}

	<-limitChan
}
