package node

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"strings"
)

func (r *router) AddRouter(ip string) error {
	return r.action(ip, "add")
}

func (r *router) RemoveRouter(ip string) error {
	return r.action(ip, "delete")
}

func (r *router) action(ip, action string) error {
	var rule string
	if !strings.Contains(ip, "/") {
		rule = fmt.Sprintf("route %s -host %s dev %s", action, ip, r.name)
	} else {
		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			return err
		}
		rule = fmt.Sprintf("route %s -net %s dev %s", action, ipNet.String(), r.name)
	}

	//route add -net 5.244.24.0 netmask 255.255.255.0 fvpn0

	return util.ExecCommand("/bin/sh", "-c", rule)
}
