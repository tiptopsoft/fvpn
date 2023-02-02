package registry

import (
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"github.com/interstellar-cloud/star/pkg/registry/addr"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"net"
)

func (r *RegStar) processRegister(addr *net.UDPAddr, conn *net.UDPConn, data []byte, cp *common.CommonPacket) {
	var regPacket register.RegPacket
	var err error
	if cp != nil {
		regPacket, err = register.DecodeWithCommonPacket(data, *cp)
	} else {
		regPacket, err = register.Decode(data)
	}

	if err := r.register(addr, regPacket); err != nil {
		log.Logger.Errorf("registry failed. err: %v", err)
	}
	// build a ack
	f, err := ackBuilder(regPacket)
	log.Logger.Infof("build a registry ack: %v", f)
	if err != nil {
		log.Logger.Errorf("build resp p failed. err: %v", err)
	}
	_, err = conn.WriteToUDP(f, addr)
	if err != nil {
		log.Logger.Errorf("registry write failed. err: %v", err)
	}

}

func ackBuilder(rp register.RegPacket) ([]byte, error) {
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

	return ack.Encode(p)
}

// register edge node register to register
func (r *RegStar) register(addr *net.UDPAddr, packet register.RegPacket) error {
	m.Store(packet.SrcMac.String(), addr)
	m.Range(func(key, value any) bool {
		log.Logger.Infof("registry data key: %s, value: %v", key, value)
		return true
	})
	return nil
}
