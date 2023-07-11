package node

import (
	"net"
)

// NetworkManager Join a network like : 192.168.0.1/24 if you give 192.168.0.1, default is 24
type NetworkManager interface {
	JoinIP(userId, ip string)
	JoinNet(userID, cidr string) error
	Leave(userId, ip string) error
	LeaveNet(userId, cidr string) error
	Access(userId, ip string) bool
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
	networks := &network{
		ips:   make([]net.IP, 20),
		ipNet: make([]*net.IPNet, 20),
	}
	m := &nodeNet{
		networks: make(map[string]*network, 1),
	}
	m.networks[userId] = networks
	return m
}

// JoinIP will add an ip or cidr, if you only give a ip, cidr is: ip.24, default 24 given
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
