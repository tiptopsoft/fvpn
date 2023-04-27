package util

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"net/netip"
)

func GetAddress(address string, port int) (unix.SockaddrInet4, error) {
	ad, err := netip.ParseAddr(address)
	return unix.SockaddrInet4{
		Port: port,
		Addr: ad.As4(),
	}, err
}

// GetMacAddr return dest mac, dest ip, if data provide is null, error returen
func GetMacAddr(buff []byte) (string, net.IP, error) {
	if len(buff) == 0 {
		return "", nil, errors.New("no data exists")
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buff[0], buff[1], buff[2], buff[3], buff[4], buff[5]), net.IPv4(buff[32], buff[33], buff[34], buff[35]), nil
}
