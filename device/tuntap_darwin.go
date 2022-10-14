package device

import (
	"github.com/vishvananda/netlink"
)

// Tuntap a device for net proxy
type Tuntap struct {
	netlink.Tuntap
}

func New() (*Tuntap, error) {

	return &Tuntap{}, nil
}
