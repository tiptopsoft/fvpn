package registry

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
)

var ipMap sync.Map

type Endpoint struct {
	Mac  net.HardwareAddr
	IP   net.IP
	Mask net.IPMask
}

// New generate a endpoint
func New() (*Endpoint, error) {
	macStr := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", 0x52, 0x54, 0x00, rand.Intn(0xff), rand.Intn(0xff), rand.Intn(0xff))
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return nil, errors.New("new mac failed")
	}

	var ip1 any
	var ok bool
	ip1, ok = ipMap.Load("ip")
	if !ok {
		ip1 = string2Long("192.168.0.1")
	} else {
		ip1 = ip1.(uint64) + 1
		ipMap.Store("ip", ip1)
	}

	ip := net.ParseIP(GenerateIP(ip1.(uint32)))
	_, ipMask, err := net.ParseCIDR("255.255.255.0/24")
	if err != nil {
		return nil, err
	}
	return &Endpoint{
		Mac:  mac,
		IP:   ip,
		Mask: ipMask.Mask,
	}, nil
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
