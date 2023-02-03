package executor

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
	"github.com/interstellar-cloud/star/pkg/util/packet/forward"
	peerack "github.com/interstellar-cloud/star/pkg/util/packet/peer/ack"
	"github.com/interstellar-cloud/star/pkg/util/packet/register/ack"
	socket2 "github.com/interstellar-cloud/star/pkg/util/socket"
	"github.com/interstellar-cloud/star/pkg/util/tuntap"
	"io"
	"net"
)

type EdgeExecutor struct {
	Tap      *tuntap.Tuntap
	Protocol option.Protocol
}

func (ee EdgeExecutor) Execute(socket socket2.Socket) error {

	if ee.Protocol == option.UDP {

		//for {
		udpBytes := make([]byte, 2048)
		_, err := socket.Read(udpBytes)
		if err != nil {
			if err == io.EOF {
				//no data exists, continue read next frame continue
				log.Logger.Errorf("not data exists")
			} else {
				log.Logger.Errorf("read from remote error: %v", err)
			}
		}

		cp, err := common.Decode(udpBytes)

		if err != nil {
			log.Logger.Errorf("decode err: %v", err)
		}

		switch cp.Flags {
		case option.MsgTypeRegisterAck:
			regAck, err := ack.Decode(udpBytes)
			if err != nil {
				return err
			}
			log.Logger.Infof("got registry registry ack: %v", regAck)
			//create tap tuntap

			//设置IP
			if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", ee.Tap.Name, regAck.AutoIP.String(), regAck.Mask.String(), 1420)); err != nil {
				return err
			}
			break
		case option.MsgTypeQueryPeer:
			//get peerInfo
			peerPacketAck, err := peerack.Decode(udpBytes)
			if err != nil {
				return err
			}

			infos := peerPacketAck.PeerInfos
			log.Logger.Infof("got registry peers: (%v)", infos)
			for _, v := range infos {
				addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", v.Host.String(), v.Port))
				if err != nil {
					log.Logger.Errorf("resolve addr failed. err: %v", err)
				}
				option.AddrMap.Store(v.Mac.String(), addr)
			}

			break

		case option.MsgTypePacket:
			forwardPacket, err := forward.Decode(udpBytes)
			if err != nil {
				return err
			}

			log.Logger.Infof("got through packet: %v", forwardPacket)
			//写入到tap
			if _, err := ee.Tap.Socket.Write(udpBytes); err != nil {
				log.Logger.Errorf("write to tap failed. (%v)", err.Error())
			}
			break
		}

	}

	//}
	return nil
}
