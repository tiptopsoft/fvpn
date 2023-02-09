package util

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net/netip"
)

func GetAddress(address string, port int) (unix.SockaddrInet4, error) {

	ad, err := netip.ParseAddr(address)
	return unix.SockaddrInet4{
		Port: port,
		Addr: ad.As4(),
	}, err
}

func GetMacAddr(buf []byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}
