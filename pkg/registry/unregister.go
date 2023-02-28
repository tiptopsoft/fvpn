package registry

import (
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/socket"
	"golang.org/x/sys/unix"
)

func (r *RegStar) processUnregister(addr unix.Sockaddr, socket socket.Socket, data []byte, cp *common.CommonPacket) {
	regPacket, err := r.packet.Decode(data)
	if err := r.unRegister(regPacket); err != nil {
		log.Logger.Errorf("registry failed. err: %v", err)
	}
	// build a ack
	f, err := r.registerAck(addr, regPacket.(register.RegPacket).SrcMac)
	log.Logger.Infof("build a registry ack: %v", f)
	if err != nil {
		log.Logger.Errorf("build resp p failed. err: %v", err)
	}
	err = socket.WriteToUdp(f, addr)
	if err != nil {
		log.Logger.Errorf("registry write failed. err: %v", err)
	}
}

func (r *RegStar) unRegister(packet packet.Interface) error {
	return nil
}
