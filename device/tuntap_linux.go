package device

import (
	"github.com/vishvananda/netlink"
	"net"
)

// Tuntap a device for net
type Tuntap struct {
	netlink.Tuntap
}

type Attr struct {
	IP   string
	Mask string
}

func New() (*Tuntap, error) {
	la := netlink.NewLinkAttrs()
	la.Name = "tun1"
	tap := netlink.Tuntap{
		LinkAttrs: la,
		Mode:      netlink.TUNTAP_MODE_TUN,
	}

	err := netlink.LinkAdd(&tap)
	if err != nil {
		panic(err)
	}
	addr := &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   []byte("192.168.0.1"),
			Mask: []byte("255.255.255.0"),
		},
	}
	err = netlink.AddrAdd(&tap, addr)
	if err != nil {
		panic(err)
	}
	err = netlink.LinkSetUp(&tap)
	if err != nil {
		panic(err)
	}

	return &Tuntap{tap}, nil
}
