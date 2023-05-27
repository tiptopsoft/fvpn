package udp

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/packet/notify"
	peerack "github.com/topcloudz/fvpn/pkg/packet/peer/ack"
	"github.com/topcloudz/fvpn/pkg/packet/register/ack"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
)

var (
	logger = log.Log()
)

func Handle() handler.HandlerFunc {

	return func(ctx context.Context, frame *packet.Frame) error {
		buff := frame.Buff[:]

		headerBuff, err := header.Decode(buff)
		if err != nil {
			logger.Errorf("decode err: %v", err)
		}

		frame.NetworkId = hex.EncodeToString(headerBuff.NetworkId[:])
		c := ctx.Value("cache").(*cache.Cache)

		frame.FrameType = headerBuff.Flags
		switch headerBuff.Flags {
		case option.MsgTypeRegisterAck:
			regAck, err := ack.Decode(buff)
			if err != nil {
				//return err
				fmt.Println(err)
			}
			logger.Debugf("register success, got server server ack: (%v)", regAck)
			break
		case option.MsgTypeQueryPeer:
			logger.Debugf("start get query response")
			peerPacketAck, err := peerack.Decode(buff)
			if err != nil {
				//return err
			}
			infos := peerPacketAck.NodeInfos
			logger.Infof("got server peers: (%v)", infos)

			for _, info := range infos {
				//logger.Debugf("got remote node: mac: %v, ip: %s,  natIP: %s, natPort: %d", info.Mac, info.IP, info.NatIp, info.NatPort)
				address, err := util.GetAddress(info.NatIp.String(), int(info.NatPort))
				if err != nil {
					logger.Errorf("resolve addr failed, err: %v", err)
				}
				node, err := c.GetNodeInfo(frame.NetworkId, info.IP.String())
				if node == nil || err != nil {
					sock := socket.NewSocket(6061)
					nodeInfo := &cache.NodeInfo{
						Socket:  sock,
						MacAddr: info.Mac,
						IP:      info.IP,
						Port:    info.Port,
						P2P:     false,
						Addr:    &address,
					}

					//cache.Nodes[info.Mac.String()] = nodeInfo
					c.SetCache(frame.NetworkId, info.IP.String(), nodeInfo)
				}
			}
			break
		case option.MsgTypePacket:
			frame.Packet = buff[:]
			break
		case option.MsgTypeNotify:
			np, err := notify.Decode(buff)
			if err != nil {
				logger.Errorf("got invalid NotifyPacket: %v", err)
			}

			addr := &unix.SockaddrInet4{
				Port: int(np.NatPort),
			}

			copy(addr.Addr[:], np.NatIP.To4())
			info := &cache.NodeInfo{
				Socket:    nil,
				NetworkId: frame.NetworkId,
				Addr:      addr,
				MacAddr:   nil,
				IP:        np.SourceIP,
				Port:      np.Port,
				P2P:       false,
				Status:    false,
				NatType:   np.NatType,
				NatIP:     np.NatIP,
				NatPort:   np.NatPort,
			}
			//}

			frame.Packet = buff[:]
			frame.Target = info
			logger.Debugf("got notify packet: %v", info)
			//c.SetCache(frame.NetworkId, info.IP.String(), info)
		}

		return nil
	}

}
