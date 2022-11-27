package registry

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
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

	ip1, ok := ipMap.Load("ip")
	if !ok {
		ip1 = string2Long("192.168.0.1")
	} else {
		ip1 = ip1.(int64) + 1
		ipMap.Store("ip", ip1)
	}

	ip := net.ParseIP(GenerateIP(ip1.(int64)))
	_, ipMask, err := net.ParseCIDR("255.255.255.0")
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
func GenerateIP(ipInt int64) string {
	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}
