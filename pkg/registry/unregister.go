package registry

import (
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/socket"
	"golang.org/x/sys/unix"
)

func (r *RegStar) processUnregister(addr unix.Sockaddr, socket socket.Socket, data []byte, cp *packet.Header) {
	regPacket, err := r.packet.Decode(data)
	if err := r.unRegister(regPacket); err != nil {
		logger.Errorf("registry failed. err: %v", err)
	}
	// build a ack
	f, err := r.registerAck(addr, regPacket.(register.RegPacket).SrcMac)
	logger.Infof("build a registry ack: %v", f)
	if err != nil {
		logger.Errorf("build resp p failed. err: %v", err)
	}
	err = socket.WriteToUdp(f, addr)
	if err != nil {
		logger.Errorf("registry write failed. err: %v", err)
	}
}

func (r *RegStar) unRegister(packet packet.Interface) error {
	return nil
}
