package addr

import (
	"encoding/hex"
	"fmt"
	"net"
	"testing"
)

func TestNew(t *testing.T) {
	RecMac := "01:01:03:02:03:01"
	src, err := net.ParseMAC(RecMac)
	if err != nil {
		t.Errorf("%v", err)
	}
	endpoint, _ := New(src)
	fmt.Println(endpoint)
}

func TestTransfer(t *testing.T) {

	s := "8056c2e21c123456"
	for _, v := range s {
		fmt.Println(string(v))
		//转byte
	}

	buff, _ := hex.DecodeString("80")
	fmt.Println(buff)

}

func TestGetMacAddrAndIPByDev(t *testing.T) {
	fa, err := net.InterfaceByName("en0")

	addrs, err := fa.Addrs()
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(fa.HardwareAddr, ipnet.IP, nil)
			}
		}
	}

	fmt.Println(fa.HardwareAddr.String(), nil, err)
}
