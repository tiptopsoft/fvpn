package util

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util/errors"
	"math/rand"
	"net"
)

var (
	BROADCAST_MAC      = net.HardwareAddr{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	MULTICAST_MAC      = net.HardwareAddr{0x01, 0x00, 0x5E, 0x00, 0x00, 0x00} // first 3 bytes are meaningful
	IPV6_MULTICAST_MAC = net.HardwareAddr{0x33, 0x33, 0x00, 0x00, 0x00, 0x00}
	NULL_MAC           = net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

func GetLocalMac(name string) ([4]byte, error) {
	var b [4]byte
	fa, err := net.InterfaceByName(name)

	if err != nil {
		return [4]byte{}, err
	}
	macAddr := fa.HardwareAddr.String()
	if len(macAddr) == 0 {
		return [4]byte{}, errors.ErrGetMac
	}

	copy(b[:], macAddr)
	return b, nil

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

func IsBroadCast(destMac net.HardwareAddr) bool {
	if destMac.String() == BROADCAST_MAC.String() {
		return true
	}

	return false
}
