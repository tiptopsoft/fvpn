package node

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/util"
	"strings"
)

func (r *router) AddRouter(ip string) error {
	return r.action(ip, "add")
}

func (r *router) RemoveRouter(ip string) error {
	return r.action(ip, "delete")
}

func (r *router) action(ip, action string) error {
	if !strings.Contains(ip, "/") {
		ip = fmt.Sprintf("%s/24", ip)
	}
	rule := fmt.Sprintf("route %s %s %s", action, ip, r.ip)

	return util.ExecCommand("/bin/sh", "-c", rule)
}
