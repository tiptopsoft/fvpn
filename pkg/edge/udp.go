package edge

import (
	"errors"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/node"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/forward"
	peerack "github.com/interstellar-cloud/star/pkg/packet/peer/ack"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"github.com/interstellar-cloud/star/pkg/util"
	"io"
	"unsafe"
)

type SocketExecutor struct {
	device   *tuntap.Tuntap
	Protocol option.Protocol
	cache    node.NodesCache
}

func (s SocketExecutor) Execute(skt socket.Interface) error {
	device := s.device
	if s.Protocol == option.UDP {
		udpBytes := make([]byte, 2048)
		size, err := skt.Read(udpBytes)
		if size < 0 {
			return errors.New("no data exists")
		}
		logger.Infof("star net skt receive size: %d, data: (%v)", size, udpBytes[:size])
		if err != nil {
			if err == io.EOF {
				//no data exists, continue read next frame continue
				logger.Errorf("not data exists")
			} else {
				logger.Errorf("read from remote error: %v", err)
			}
		}

		cpInterface, err := packet.NewPacketWithoutType().Decode(udpBytes[:size])
		cp := cpInterface.(packet.Header)
		if err != nil {
			logger.Errorf("decode err: %v", err)
		}

		switch cp.Flags {
		case option.MsgTypeRegisterAck:
			regAckInterface, err := ack.NewPacket().Decode(udpBytes[:size])
			regAck := regAckInterface.(ack.RegPacketAck)

			if err != nil {
				return err
			}
			logger.Infof("got registry registry ack: (%v)", udpBytes[:size])
			//设置IP
			if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", device.Name, regAck.AutoIP.String(), regAck.Mask.String(), 1420)); err != nil {
				return err
			}
			break
		case option.MsgTypeQueryPeer:
			peerPacketAckIface, err := peerack.NewPacket().Decode(udpBytes[:size])
			peerPacketAck := peerPacketAckIface.(peerack.EdgePacketAck)
			if err != nil {
				return err
			}
			infos := peerPacketAck.PeerInfos
			logger.Infof("got registry peers: (%v)", infos)
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
				peerInfo := &node.Node{
					Socket:  sock,
					MacAddr: info.Mac,
					IP:      info.Host,
					Port:    info.Port,
				}
				s.cache.Nodes[info.Mac.String()] = peerInfo
			}
			break
		case option.MsgTypePacket:
			forwardPacketInterface, err := forward.NewPacket().Decode(udpBytes[:size])
			forwardPacket := forwardPacketInterface.(forward.ForwardPacket)
			if err != nil {
				return err
			}
			logger.Infof("got through packet: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, device.MacAddr)

			if forwardPacket.SrcMac.String() == device.MacAddr.String() {
				//self, drop packet
				logger.Infof("self packet droped: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, device.MacAddr)
			} else {
				//写入到tap
				idx := unsafe.Sizeof(forwardPacket)
				if _, err := device.Write(udpBytes[idx:size]); err != nil {
					logger.Errorf("write to tap failed. (%v)", err.Error())
				}
				logger.Infof("net write to tap as tap response to client. size: %d", size-int(idx))
			}
			break
		}
	}
	return nil
}
