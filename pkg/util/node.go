package util

import (
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"golang.org/x/sys/unix"
	"net"
)

type Nodes map[string]*Node

type Node struct {
	Socket  socket.Socket
	Addr    unix.Sockaddr
	MacAddr net.HardwareAddr
	IP      net.IP
	Port    uint16
}

func FindNode(nodes Nodes, destMac string) *Node {
	return nodes[destMac]
}

type Server interface {
	Start(port int) error
	Stop() error
}
