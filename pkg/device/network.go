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

package device

import (
	"net"
)

// NetworkManager Join a network like : 192.168.0.1/24 if you give 192.168.0.1, default is 24
type NetworkManager interface {
	JoinNet(userId, cidr string) error
	Leave(userId, cidr string) error
	LeaveNet(userId, cidr string) error
	Access(userId, cidr string) bool
}

var (
	_      NetworkManager = (*nodeNet)(nil)
	AllIPs                = "0.0.0.0/0"
)

func NewNetworkManager(userId string) NetworkManager {
	return newNodeNet(userId)
}

type nodeNet struct {
	networks map[string]*network
}

type network struct {
	ips   []net.IP
	ipNet []*net.IPNet
}

func newNodeNet(userId string) *nodeNet {
	networks := &network{}
	m := &nodeNet{
		networks: make(map[string]*network, 1),
	}
	m.networks[userId] = networks
	return m
}

// JoinIP will add an Cidr or Cidr, if you only give a Cidr, Cidr is: Cidr.24, default 24 given
func (n *nodeNet) JoinIP(userId, ip string) {
	IP := net.ParseIP(ip)
	networks := n.networks[userId]
	networks.ips = append(networks.ips, IP)
	n.networks[userId] = networks
}

func (n *nodeNet) JoinNet(userId, cidr string) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	networks := n.networks[userId]
	networks.ipNet = append(networks.ipNet, ipNet)
	n.networks[userId] = networks
	return nil
}

func (n *nodeNet) Leave(userId, ip string) error {
	return nil
}

func (n *nodeNet) LeaveNet(userId, cidr string) error {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	networks := n.networks[userId].ipNet
	for index, net := range networks {
		if net == ipNet {
			networks = append(networks[:index], networks[index+1:]...)
		}
	}

	n.networks[userId].ipNet = networks
	return nil
}

func (n *nodeNet) Access(userId, ip string) bool {
	IP := net.ParseIP(ip)
	networks := n.networks[userId]
	for _, v := range networks.ipNet {
		if v == nil {
			continue
		}
		if v.String() == AllIPs {
			return true
		}
		if v.Contains(IP) {
			return true
		}
	}

	for _, v := range networks.ips {
		if v == nil {
			continue
		}
		if v.Equal(IP) {
			return true
		}
	}

	return false
}
