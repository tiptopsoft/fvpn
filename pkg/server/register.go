package server

import (
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/packet/register/ack"
	"golang.org/x/sys/unix"
	"net"
)

func (r *RegStar) processRegister(remoteAddr unix.Sockaddr, data []byte, cp *packet.Header) {
	packet, err := r.packet.Decode(data)

	// build an ack
	f, err := r.registerAck(remoteAddr, packet.(register.RegPacket).SrcMac)
	logger.Infof("build a server ack: %v", f)
	if err != nil {
		logger.Errorf("build resp failed. err: %v", err)
	}
	err = r.socket.WriteToUdp(f, remoteAddr)
	if err != nil {
		logger.Errorf("server write failed. err: %v", err)
	}
	logger.Infof("write a server ack to remote: %v, data: %v", remoteAddr, f)

}

func (r *RegStar) registerAck(peerAddr unix.Sockaddr, srcMac net.HardwareAddr) ([]byte, error) {
	endpoint, err := addr.New(srcMac)
	if err != nil {
		return nil, err
	}
	p := ack.NewPacket()
	p.RegMac = endpoint.Mac
	p.AutoIP = endpoint.IP
	p.Mask = endpoint.Mask

	ackNode := &cache.Peer{
		Socket:  r.socket,
		Addr:    peerAddr,
		MacAddr: endpoint.Mac,
		IP:      endpoint.IP,
		Port:    0,
	}

	r.cache.Nodes[endpoint.Mac.String()] = ackNode
	r.cache.IPNodes[endpoint.IP.String()] = ackNode
	return p.Encode()
}
