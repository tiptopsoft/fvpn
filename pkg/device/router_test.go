package device

import (
	"net"
	"testing"
)

func TestEncode(t *testing.T) {
	cidr := "192.168.0.1/24"
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}

}
