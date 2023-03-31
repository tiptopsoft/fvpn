package udp

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/interstellar-cloud/star/pkg/cache"
	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/forward"
	peerack "github.com/interstellar-cloud/star/pkg/packet/peer/ack"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"github.com/interstellar-cloud/star/pkg/util"
)

var (
	logger = log.Log()
)

type UdpHandler struct {
	device *tuntap.Tuntap
	cache  cache.PeersCache
}

func New(device *tuntap.Tuntap, cache cache.PeersCache) handler.Handler {
	return UdpHandler{
		device: device,
		cache:  cache,
	}
}

func (uh UdpHandler) Handle(ctx context.Context, buff []byte) error {
	cpInterface, err := packet.NewPacketWithoutType().Decode(buff[:])
	cp := cpInterface.(packet.Header)
	if err != nil {
		logger.Errorf("decode err: %v", err)
	}

	switch cp.Flags {
	case option.MsgTypeRegisterAck:
		regAckInterface, err := ack.NewPacket().Decode(buff[:])
		regAck := regAckInterface.(ack.RegPacketAck)

		if err != nil {
			return err
		}
		logger.Infof("got fvpns fvpns ack: (%v)", buff[:])
		//设置IP
		if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", uh.device.Name, regAck.AutoIP.String(), regAck.Mask.String(), 1420)); err != nil {
			return err
		}
		break
	case option.MsgTypeQueryPeer:
		peerPacketAckIface, err := peerack.NewPacket().Decode(buff[:])
		peerPacketAck := peerPacketAckIface.(peerack.EdgePacketAck)
		if err != nil {
			return err
		}
		infos := peerPacketAck.PeerInfos
		logger.Infof("got fvpns peers: (%v)", infos)
		for _, info := range infos {
			address, err := util.GetAddress(info.Host.String(), int(info.Port))
			if err != nil {
				logger.Errorf("resolve addr failed, err: %v", err)
			}
			sock := socket.NewSocket()
			err = sock.Connect(&address)
			if err != nil {
				return err
			}
			peerInfo := &cache.Peer{
				Socket:  sock,
				MacAddr: info.Mac,
				IP:      info.Host,
				Port:    info.Port,
			}
			uh.cache.Nodes[info.Mac.String()] = peerInfo
		}
		break
	case option.MsgTypePacket:
		forwardPacketInterface, err := forward.NewPacket().Decode(buff[:])
		forwardPacket := forwardPacketInterface.(forward.ForwardPacket)
		if err != nil {
			return err
		}
		logger.Infof("got through packet: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, uh.device.MacAddr)

		if forwardPacket.SrcMac.String() == uh.device.MacAddr.String() {
			//self, drop packet
			logger.Infof("self packet droped: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, uh.device.MacAddr)
		} else {
			//写入到tap
			idx := unsafe.Sizeof(forwardPacket)
			if _, err := uh.device.Write(buff[idx:]); err != nil {
				logger.Errorf("write to tap failed. (%v)", err.Error())
			}
		}
		break
	}

	return nil
}
