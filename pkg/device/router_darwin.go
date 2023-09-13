// Copyright 2023 TiptopSoft, Inc.
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

package device

import (
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/util"
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
