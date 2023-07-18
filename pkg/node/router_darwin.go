package node

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/util"
)

func (r *router) AddRouter(cidr string) error {
	return r.action(cidr, "add")
}

func (r *router) RemoveRouter(cidr string) error {
	return r.action(cidr, "delete")
}

func (r *router) action(cidr, action string) error {
	//example: sudo route -nv add -net 192.168.10.1 -netmask 255.255.255.0 -interface en0
	rule := fmt.Sprintf("route -nv %s -net %s -interface %s", action, cidr, r.name)
	return util.ExecCommand("/bin/sh", "-c", rule)
}
