package node

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"strings"
)

func (r *router) AddRouter(ip string) error {
	return action(ip, "add")
}

func (r *router) RemoveRouter(ip string) error {
	return action(ip, "delete")
}

func action(ip, action string) error {
	originIP := ip
	if !strings.Contains(ip, "/") {
		ip = fmt.Sprintf("%s/24", ip)
	}

	_, ipNet, err := net.ParseCIDR(ip)

	//route add -net 5.244.24.0 netmask 255.255.255.0 fvpn0
	rule := fmt.Sprintf("route %s -net %s netmask %s %s", action, originIP, ipNet.Mask, r.name)

	return util.ExecCommand("/bin/sh", "-c", rule)
}
