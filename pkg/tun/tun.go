package tun

import (
	"github.com/topcloudz/fvpn/pkg/packet"
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

// Read is a hack to work around the first 4 bytes "packet
// information" because there doesn't seem to be an IFF_NO_PI for darwin.
func (tun *NativeTun) Read(buff []byte) (n int, err error) {
	n, err = tun.file.Read(buff)
	return n - 4, err
}

func (tun *NativeTun) Name() string {
	return tun.name
}

func (tun *NativeTun) ReadToFrame(f *packet.Frame) (n int, err error) {
	n, err = tun.file.Read(f.Buff)
	f.Packet = f.Buff[4:]
	return n - 4, err
}

func (tun *NativeTun) Write(buff []byte) (int, error) {
	tun.lock.Lock()
	defer tun.lock.Unlock()
	n, err := tun.file.Write(buff[4:])
	return n - 4, err
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
