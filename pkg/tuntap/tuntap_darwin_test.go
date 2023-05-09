package tuntap

import (
	"fmt"
	"net"
	"testing"
)

func TestNew(t *testing.T) {

	fa, err := net.InterfaceByName("en0")
	var ip net.IP
	addrs, err := fa.Addrs()
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP
				fmt.Println(fa.HardwareAddr, ipnet.IP, nil)
			}
		}
	}

	tun, err := New(TAP, ip.String(), "255.255.255.254", "1111")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tun)

}
