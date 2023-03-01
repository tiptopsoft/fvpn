package edge

import (
	"errors"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/node"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
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
		log.Infof("star net skt receive size: %d, data: (%v)", size, udpBytes[:size])
		if err != nil {
			if err == io.EOF {
				//no data exists, continue read next frame continue
				log.Errorf("not data exists")
			} else {
				log.Errorf("read from remote error: %v", err)
			}
		}

		cpInterface, err := common.NewPacketWithoutType().Decode(udpBytes[:size])
		cp := cpInterface.(common.CommonPacket)
		if err != nil {
			log.Errorf("decode err: %v", err)
		}

		switch cp.Flags {
		case option.MsgTypeRegisterAck:
			regAckInterface, err := ack.NewPacket().Decode(udpBytes[:size])
			regAck := regAckInterface.(ack.RegPacketAck)

			if err != nil {
				return err
			}
			log.Infof("got registry registry ack: (%v)", udpBytes[:size])
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
			log.Infof("got registry peers: (%v)", infos)
			for _, info := range infos {
				address, err := util.GetAddress(info.Host.String(), int(info.Port))
				if err != nil {
					log.Errorf("resolve addr failed, err: %v", err)
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
			log.Infof("got through packet: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, device.MacAddr)

			if forwardPacket.SrcMac.String() == device.MacAddr.String() {
				//self, drop packet
				log.Infof("self packet droped: %v, srcMac: %v, current tap macAddr: %v", forwardPacket, forwardPacket.SrcMac, device.MacAddr)
			} else {
				//写入到tap
				idx := unsafe.Sizeof(forwardPacket)
				if _, err := device.Write(udpBytes[idx:size]); err != nil {
					log.Errorf("write to tap failed. (%v)", err.Error())
				}
				log.Infof("net write to tap as tap response to client. size: %d", size-int(idx))
			}
			break
		}
	}
	return nil
}
