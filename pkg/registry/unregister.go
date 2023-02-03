package registry

import (
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
	"github.com/interstellar-cloud/star/pkg/util/packet/register"
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"net"
)

func (r *RegStar) processUnregister(addr *net.UDPAddr, socket socket.Socket, data []byte, cp *common.CommonPacket) {
	var regPacket register.RegPacket
	var err error
	if cp != nil {
		regPacket, err = register.DecodeWithCommonPacket(data, *cp)
	} else {
		regPacket, err = register.Decode(data)
	}

	if err := r.unRegister(regPacket); err != nil {
		log.Logger.Errorf("registry failed. err: %v", err)
	}
	// build a ack
	f, err := r.ackBuilder(*addr, socket, regPacket)
	log.Logger.Infof("build a registry ack: %v", f)
	if err != nil {
		log.Logger.Errorf("build resp p failed. err: %v", err)
	}
	_, err = socket.WriteToUdp(f, addr)
	if err != nil {
		log.Logger.Errorf("registry write failed. err: %v", err)
	}
}

func (r *RegStar) unRegister(packet register.RegPacket) error {
	return nil
}
