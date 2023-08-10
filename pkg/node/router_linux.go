// Copyright 2023 Tiptopsoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package node

import (
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"net"
)

func (r *router) AddRouter(cidr string) error {
	//first remove
	if err := r.RemoveRouter(cidr); err != nil {
		return err
	}
	return r.action(cidr, "add")
}

func (r *router) RemoveRouter(cidr string) error {
	return r.action(cidr, "delete")
}

func (r *router) action(cidr, action string) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	rule := fmt.Sprintf("route %s -net %v dev %s", action, ipNet, r.name)
	//example: route add -net 5.244.24.0/24 dev fvpn0
	return util.ExecCommand("/bin/sh", "-c", rule)
}
