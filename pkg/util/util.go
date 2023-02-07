package util

import (
	"fmt"
	"math/rand"
	"net"
)

var (
	BROADCAST_MAC      = net.HardwareAddr{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	MULTICAST_MAC      = net.HardwareAddr{0x01, 0x00, 0x5E, 0x00, 0x00, 0x00} // first 3 bytes are meaningful
	IPV6_MULTICAST_MAC = net.HardwareAddr{0x33, 0x33, 0x00, 0x00, 0x00, 0x00}
	NULL_MAC           = net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

func GetMacAddrByDev(name string) (net.HardwareAddr, error) {
	fa, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	return fa.HardwareAddr, nil
}

// RandMac rand gen a mac
func RandMac() (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	buf[0] |= 2
	mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	return mac, nil
}

func GetLocalMacAddr() string {
	// getMac
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, v := range ifaces {
		if v.HardwareAddr == nil {
			continue
		}
		return v.HardwareAddr.String()
	}

	return ""
}

func IsBroadCast(destMac string) bool {
	if destMac == BROADCAST_MAC.String() {
		return true
	}

	return false
}
