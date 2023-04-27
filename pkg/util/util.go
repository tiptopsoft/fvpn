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

// GetMacAddr return dest mac, dest ip, if data provide is null, error return
func GetMacAddr(buff []byte) (string, net.IP, error) {
	if len(buff) == 0 {
		return "", nil, errors.New("no data exists")
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buff[6], buff[7], buff[8], buff[9], buff[10], buff[11]), net.IPv4(buff[30], buff[31], buff[32], buff[33]), nil
}
