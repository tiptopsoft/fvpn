package registry

import (
	"github.com/interstellar-cloud/star/pkg/registry/addr"
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
	"github.com/interstellar-cloud/star/pkg/util/packet/register"
	"github.com/interstellar-cloud/star/pkg/util/packet/register/ack"
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"net"
)

func (r *RegStar) processRegister(remoteAddr *net.UDPAddr, socket socket.Socket, data []byte, cp *common.CommonPacket) {
	var regPacket register.RegPacket
	var err error
	if cp != nil {
		regPacket, err = register.DecodeWithCommonPacket(data, *cp)
	} else {
		regPacket, err = register.Decode(data)
	}

	// build an ack
	f, err := r.ackBuilder(*remoteAddr, socket, regPacket)
	log.Logger.Infof("build a registry ack: %v", f)
	if err != nil {
		log.Logger.Errorf("build resp p failed. err: %v", err)
	}
	_, err = socket.WriteToUdp(f, remoteAddr)
	if err != nil {
		log.Logger.Errorf("registry write failed. err: %v", err)
	}
	log.Logger.Infof("write a registry ack to remote: %v, data: %v", remoteAddr, f)

}

func (r *RegStar) ackBuilder(peerAddr net.UDPAddr, socket socket.Socket, rp register.RegPacket) ([]byte, error) {
	endpoint, err := addr.New(rp.SrcMac.String())
	if err != nil {
		return nil, err
	}

	p := ack.NewPacket()
	p.RegMac = endpoint.Mac
	p.AutoIP = endpoint.IP
	p.Mask = endpoint.Mask
	rp.CommonPacket.Flags = option.MsgTypeRegisterAck
	p.CommonPacket = rp.CommonPacket

	peer := &util.Peer{
		Conn:    socket.UdpSocket,
		Addr:    peerAddr,
		MacAddr: endpoint.Mac,
		IP:      endpoint.IP,
		Port:    0,
	}

	r.Peers[endpoint.Mac.String()] = peer
	return ack.Encode(p)
}
