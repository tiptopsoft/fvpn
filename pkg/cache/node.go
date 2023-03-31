package cache

import (
	"github.com/interstellar-cloud/star/pkg/socket"
	"golang.org/x/sys/unix"
	"net"
)

type PeersCache struct {
	Nodes   map[string]*Peer
	IPNodes map[string]*Peer
}

// Peer use as other Node
type Peer struct {
	Socket  socket.Interface
	Addr    unix.Sockaddr
	MacAddr net.HardwareAddr
	IP      net.IP
	Port    uint16
	P2P     bool
}

func New() PeersCache {
	return PeersCache{
		Nodes:   make(map[string]*Peer, 1),
		IPNodes: make(map[string]*Peer, 1),
	}
}

func FindPeer(cache PeersCache, destMac string) *Peer {
	return cache.Nodes[destMac]
}

func FindPeerByIP(cache PeersCache, ip string) *Peer {
	return cache.IPNodes[ip]
}

type Server interface {
	Start(port int) error
	Stop() error
}
