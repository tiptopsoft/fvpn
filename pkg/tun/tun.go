package tun

import (
	"net"
	"os"
	"sync"
)

type Mode int

const (
	SYS_IOCTL     = 29
	TUNSETIFF     = 0x400454ca
	TUNSETPERSIST = 0x400454cb
	TUNSETGROUP   = 0x400454ce
	TUNSETOWNER   = 0x400454cc

	IFF_NO_PI      = 0x1000
	IFF_TUN        = 0x1
	IFF_TAP        = 0x2
	TUN       Mode = iota
	TAP
)

var DefaultNamePrefix = "fvpn"

type Device interface {
	Name() string
	Read(buff []byte) (int, error)
	Write(buff []byte) (int, error)
	SetIP(net, ip string) error
	SetMTU(mtu int) error
	IPToString() string
	Addr() net.IP
}

var (
	_ Device = (*NativeTun)(nil)
)

// NativeTun a tuntap for net
type NativeTun struct {
	lock      sync.Mutex
	file      *os.File
	Fd        int
	name      string
	NetworkId string
	IP        net.IP
}

func (tun *NativeTun) Name() string {
	return tun.name
}

func (tun *NativeTun) IPToString() string {
	return tun.IP.String()
}

func (tun *NativeTun) Addr() net.IP {
	return tun.IP
}

// Close close the device
func (tun *NativeTun) Close() error {
	return tun.file.Close()
}
