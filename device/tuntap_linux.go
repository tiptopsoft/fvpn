package device

import "github.com/vishvananda/netlink"

// Tuntap a device for net
type Tuntap struct {
	netlink.Tuntap
}

type Attr struct {
	IP   string
	Mask string
}

func New() *Tuntap {
	la := netlink.NewLinkAttrs()
	la.Name = "tun1"
	tap := netlink.Tuntap{
		LinkAttrs, la,
	}
	err := netlink.LinkAdd(tap)
	if err != nil {
		panic(err)
	}
	return &Tuntap{tap}
}
