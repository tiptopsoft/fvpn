package util

import (
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
