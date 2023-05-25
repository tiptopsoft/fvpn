package main

import "fmt"

// import (
//
//	"fmt"
//
// )
func main() {
	//	//srcAddr := &net.UDPAddr{
	//	//	IP:   net.IPv4zero,
	//	//	Port: 6061,
	//	//}
	//	//
	//	//destAddr1 := &net.UDPAddr{
	//	//	IP:   net.ParseIP("211.125.225.186"),
	//	//	Port: 4000,
	//	//}
	//	//
	//	//destAddr2 := &net.UDPAddr{
	//	//	IP:   net.ParseIP("81.70.36.156"),
	//	//	Port: 4000,
	//	//}
	//
	//	//sock := socket.NewSocket()
	//	//destAddr1 := unix.SockaddrInet4{
	//	//	Port: 4000,
	//	//	SourceIP: [4]byte{211, 125, 225, 186},
	//	//}
	//	//
	//	//destAddr2 := unix.SockaddrInet4{
	//	//	Port: 5000,
	//	//	SourceIP: [4]byte{81, 70, 36., 156},
	//	//}
	//	//err := sock.Connect(&destAddr1)
	//	//err = sock.Connect(&destAddr2)
	//	//if err != nil {
	//	//
	//	//	panic(err)
	//	//}
	//
	//	//conn1, err := net.DialUDP("udp", srcAddr, destAddr1)
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//	//
	//	//fmt.Println(conn1)
	//	//
	//	//conn2, err := net.DialUDP("udp", srcAddr, destAddr2)
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//
	buff := []byte{104, 101, 108, 108, 111, 44, 32, 104, 111, 108, 101, 32, 112, 117, 110, 99, 104, 105, 110, 103, 46, 46, 46, 32}

	fmt.Println(string(buff))
}
