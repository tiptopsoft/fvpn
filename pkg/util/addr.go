package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/errors"
	"go.uber.org/atomic"
	"math/rand"
	"net"
	"runtime"
	"sync"
)

var (
	ipMap              sync.Map
	ipNumber           atomic.Int32
	BROADCAST_MAC      = net.HardwareAddr{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	MULTICAST_MAC      = net.HardwareAddr{0x01, 0x00, 0x5E, 0x00, 0x00, 0x00} // first 3 bytes are meaningful
	IPV6_MULTICAST_MAC = net.HardwareAddr{0x33, 0x33, 0x00, 0x00, 0x00, 0x00}
	NULL_MAC           = net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	DefaultPort        = 4000
)

type Endpoint struct {
	Mac      net.HardwareAddr
	IP       net.IP
	Mask     net.IP
	ipNumber int32
}

// AddrCache 存储到map里
type AddrCache struct {
	Group    [4]byte
	SrcMac   string
	Endpoint Endpoint
}

// New generate a endpoint
func New(offset int32) (*Endpoint, error) {
	num := ipNumber.Load()
	num = string2Long("192.168.0.1")
	num = num + offset
	ip := net.ParseIP(GenerateIP(num))
	ipMask, _, err := net.ParseCIDR("255.255.255.0/24")

	if err != nil {
		return nil, errors.ErrInvalieCIDR
	}
	ep := &Endpoint{
		IP:       ip,
		Mask:     ipMask,
		ipNumber: num,
	}

	return ep, nil
}

// ip到数字
func string2Long(ip string) int32 {
	var long int32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

// GenerateIP 数字到IP
func GenerateIP(ipInt int32) string {
	// need to do two bit shifting and “0xff” masking
	b0 := (ipInt >> 24) & 0xff
	b1 := (ipInt >> 16) & 0xff
	b2 := (ipInt >> 8) & 0xff
	b3 := ipInt & 0xff
	return fmt.Sprintf("%d.%d.%d.%d", b0, b1, b2, b3)
}

// GetMacAddrAndIPByDev get tun mac, ip for register
func GetMacAddrAndIPByDev(name string) (net.HardwareAddr, net.IP, error) {
	var fa *net.Interface
	var err error
	systemType := runtime.GOOS
	switch systemType {
	case "linux":
		fa, err = net.InterfaceByName(name[:14])
		break
	case "darwin":
		fa, err = net.InterfaceByName(name)
		break
	case "windows":
		break
	}

	if err != nil {
		return nil, nil, err
	}

	addrs, err := fa.Addrs()
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return fa.HardwareAddr, ipnet.IP, nil
			}
		}
	}

	return fa.HardwareAddr, nil, nil
}

func GetHostMac() (net.HardwareAddr, error) {
	systemType := runtime.GOOS
	switch systemType {
	case "linux":
		face, err := net.InterfaceByName("eth0")
		return face.HardwareAddr, err
	case "darwin":
		face, err := net.InterfaceByName("en0")
		return face.HardwareAddr, err
	}

	return nil, nil
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

func GetLocalMacAddr() net.HardwareAddr {
	// getMac
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, v := range ifaces {
		if v.HardwareAddr == nil {
			continue
		}
		return v.HardwareAddr
	}

	return nil
}

func IsBroadCast(destMac string) bool {
	if destMac == BROADCAST_MAC.String() {
		return true
	}

	return false
}
