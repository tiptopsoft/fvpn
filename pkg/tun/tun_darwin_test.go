package tun

import (
	"fmt"
	"net"
	"testing"
)

func Test(t *testing.T) {
	//tun, err := New()
	//if err != err {
	//	panic(err)
	//}
	//
	//fmt.Println("tun is: ", tun.name)
	//buff := make([]byte, 1024)
	//for {
	//	n, err := tun.Read(buff)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	fmt.Println(fmt.Sprintf("Read from %s %d byte", tun.name, n))
	//}
	//

	iface, err := net.InterfaceByName("utun3")
	if err != nil {
		panic(err)
	}

	addr1, err := iface.Addrs()
	if err != nil {
		panic(err)
	}
	fmt.Println(addr1[0].Network(), addr1[0].(*net.IPNet).IP)
}
