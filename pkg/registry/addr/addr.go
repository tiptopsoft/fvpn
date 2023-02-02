package addr

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"go.uber.org/atomic"
	"math/rand"
	"net"
	"sync"
)

var (
	ipMap    sync.Map
	ipNumber atomic.Uint32
)

type Endpoint struct {
	Mac  net.HardwareAddr
	IP   net.IP
	Mask net.IP
}

//AddrCache 存储到map里
type AddrCache struct {
	Group    [4]byte
	SrcMac   string
	Endpoint Endpoint
}

// New generate a endpoint
func New(srcMac string) (*Endpoint, error) {
	macStr := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", 0x52, 0x54, 0x00, rand.Intn(0xff), rand.Intn(0xff), rand.Intn(0xff))
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return nil, errors.New("new mac failed")
	}

	var ac any
	var ok bool
	ac, ok = ipMap.Load(srcMac)
	if !ok {
		num := ipNumber.Load()
		if num == 0 {
			num = string2Long("192.168.0.1")
		} else {
			num++
		}
		ip := net.ParseIP(GenerateIP(num))
		ipMask, _, err := net.ParseCIDR("255.255.255.0/24")

		if err != nil {
			log.Logger.Errorf("invalid cidr.")
			return nil, errors.New("invalid cidr")
		}
		ac = AddrCache{
			Group:  [4]byte{},
			SrcMac: srcMac,
			Endpoint: Endpoint{
				Mac:  mac,
				IP:   ip,
				Mask: ipMask,
			},
		}
		ipNumber.Store(num)
		ipMap.Store(srcMac, ac)
	} else {
		cache := ac.(AddrCache)
		ip := net.ParseIP(GenerateIP(ipNumber.Load()))
		cache.Endpoint.IP = ip
		ipMap.Store(srcMac, cache)
	}

	res := ac.(AddrCache)
	return &res.Endpoint, nil
}

//ip到数字
func string2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

//数字到IP
func GenerateIP(ipInt uint32) string {
	// need to do two bit shifting and “0xff” masking
	b0 := (ipInt >> 24) & 0xff
	b1 := (ipInt >> 16) & 0xff
	b2 := (ipInt >> 8) & 0xff
	b3 := ipInt & 0xff
	return fmt.Sprintf("%d.%d.%d.%d", b0, b1, b2, b3)
}
