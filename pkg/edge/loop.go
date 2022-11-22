package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/device"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	peerack "github.com/interstellar-cloud/star/pkg/packet/peer/ack"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"io"
	"net"
	"sync"
)

var m sync.Map

func (es *EdgeStar) process(conn net.Conn) error {

	if es.Protocol == option.UDP {
		for {
			udpBytes := make([]byte, 2048)
			_, _, err := conn.(*net.UDPConn).ReadFromUDP(udpBytes)
			if err != nil {
				if err == io.EOF {
					//no data exists, continue read next frame.
					continue
				} else {
					log.Logger.Errorf("read from remote error: %v", err)
				}
			}

			cp, err := common.Decode(udpBytes)

			if err != nil {
				log.Logger.Errorf("decode err: %v", err)
			}

			switch cp.Flags {
			case option.MSG_TYPE_REGISTER_ACK:
				regAck, err := ack.Decode(udpBytes)
				if err != nil {
					return err
				}
				log.Logger.Infof("got registry registry ack: %v", regAck)
				//create tap device
				if tap, err := device.New(device.TAP); err != nil {
					return err
				} else {
					es.tap = tap

					//设置IP
					if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ip addr change %s dev %s", regAck.AutoIP.String(), tap.Name)); err != nil {
						return err
					}
				}
				ch <- 2
				break
			case option.MSG_TYPE_PEER_INFO:
				//get peerInfo
				peerPacketAck, err := peerack.Decode(udpBytes)
				if err != nil {
					return err
				}

				infos := peerPacketAck.PeerInfos
				for _, v := range infos {
					addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", v.Host.String(), v.Port))
					if err != nil {
						log.Logger.Errorf("resolve addr failed. err: %v", err)
					}
					m.Store(v.Mac.String(), addr)
				}

				ch <- 3

				break
			}

		}

	}
	return nil
}
