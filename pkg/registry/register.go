package registry

import (
	"github.com/interstellar-cloud/star/pkg/addr"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/node"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"golang.org/x/sys/unix"
	"net"
)

func (r *RegStar) processRegister(remoteAddr unix.Sockaddr, data []byte, cp *common.CommonPacket) {
	packet, err := r.packet.Decode(data)

	// build an ack
	f, err := r.registerAck(remoteAddr, packet.(register.RegPacket).SrcMac)
	log.Logger.Infof("build a registry ack: %v", f)
	if err != nil {
		log.Logger.Errorf("build resp failed. err: %v", err)
	}
	err = r.socket.WriteToUdp(f, remoteAddr)
	if err != nil {
		log.Logger.Errorf("registry write failed. err: %v", err)
	}
	log.Logger.Infof("write a registry ack to remote: %v, data: %v", remoteAddr, f)

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
	p.CommonPacket = common.NewPacket(option.MsgTypeRegisterAck)

	ackNode := &node.Node{
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
