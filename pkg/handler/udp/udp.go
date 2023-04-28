package udp

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/forward"
	peerack "github.com/topcloudz/fvpn/pkg/packet/peer/ack"
	"github.com/topcloudz/fvpn/pkg/packet/register/ack"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"unsafe"
)

var (
	logger = log.Log()
)

func Handle() handler.HandlerFunc {

	return func(ctx context.Context, frame *packet.Frame) error {
		buff := frame.Buff[:]

		cpInterface, err := packet.NewPacketWithoutType().Decode(buff)
		header := cpInterface.(*packet.Header)
		if err != nil {
			logger.Errorf("decode err: %v", err)
		}

		switch header.Flags {
		case option.MsgTypeRegisterAck:
			regAckInterface, err := ack.NewPacket().Decode(buff)
			regAck := regAckInterface.(ack.RegPacketAck)

			if err != nil {
				//return err
			}
			logger.Infof("got server server ack: (%v)", regAck.AutoIP)
			break
		case option.MsgTypeQueryPeer:
			peerPacketAckIface, err := peerack.NewPacket().Decode(buff)
			peerPacketAck := peerPacketAckIface.(peerack.EdgePacketAck)
			if err != nil {
				//return err
			}
			infos := peerPacketAck.NodeInfos
			logger.Infof("got server peers: (%v)", infos)
			for _, info := range infos {
				address, err := util.GetAddress(info.IP.String(), int(info.Port))
				if err != nil {
					logger.Errorf("resolve addr failed, err: %v", err)
				}
				sock := socket.NewSocket()
				err = sock.Connect(&address)
				if err != nil {
					//return err
				}
				nodeInfo := &cache.NodeInfo{
					Socket:  sock,
					MacAddr: info.Mac,
					IP:      info.IP,
					Port:    info.Port,
				}
				c := ctx.Value("cache").(*cache.Cache)
				tun := ctx.Value("tun").(*tuntap.Tuntap)
				//cache.Nodes[info.Mac.String()] = nodeInfo
				c.SetCache(tun.NetworkId, info.IP.String(), nodeInfo)
			}
			break
		case option.MsgTypePacket:
			forwardPacketInterface, err := forward.NewPacket("").Decode(buff[:])
			forwardPacket := forwardPacketInterface.(forward.ForwardPacket)
			if err != nil {
				//return err
			}
			logger.Infof("got through packet: %v, srcMac: %v", forwardPacket, forwardPacket.SrcMac)

			//写入到tap
			idx := unsafe.Sizeof(forwardPacket)
			//networkId := header.NetworkId
			frame.Packet = buff[idx:]
			frame.NetworkId = string(header.NetworkId[:])

			break
		}

		return nil
	}

}
