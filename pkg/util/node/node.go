package node

import (
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"golang.org/x/sys/unix"
	"net"
)

type NodesCache struct {
	Nodes   map[string]*Node
	IPNodes map[string]*Node
}

type Node struct {
	Socket  socket.Socket
	Addr    unix.Sockaddr
	MacAddr net.HardwareAddr
	IP      net.IP
	Port    uint16
	P2P     bool
}

func New() NodesCache {
	return NodesCache{
		Nodes:   make(map[string]*Node, 1),
		IPNodes: make(map[string]*Node, 1),
	}
}

func FindNode(cache NodesCache, destMac string) *Node {
	return cache.Nodes[destMac]
}

func FindNodeByIP(cache NodesCache, ip string) *Node {
	return cache.IPNodes[ip]
}

type Server interface {
	Start(port int) error
	Stop() error
}
