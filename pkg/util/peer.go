package util

import "net"

type Peers map[string]*Peer

type Peer struct {
	Conn    net.Conn
	Addr    net.UDPAddr
	MacAddr net.HardwareAddr
	IP      net.IP
	Port    uint16
}

func FindPeers(Peers Peers, destMac string) *Peer {
	return Peers[destMac]
}
