package node

import (
	"fmt"
	"net"
	"testing"
)

func TestNodeNet_JoinIP(t *testing.T) {
	ip := "192.168.0.2/24"
	_, ipNet, err := net.ParseCIDR(ip)
	fmt.Println(ipNet, err)
}
