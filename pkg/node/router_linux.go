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
	rule := fmt.Sprintf("route %s -net %s dev %s", action, cidr, r.name)
	//example: route add -net 5.244.24.0/24 dev fvpn0
	return util.ExecCommand("/bin/sh", "-c", rule)
}
