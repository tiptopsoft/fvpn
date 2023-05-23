package udp

import (
	"context"
	"encoding/hex"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	peerack "github.com/topcloudz/fvpn/pkg/packet/peer/ack"
	"github.com/topcloudz/fvpn/pkg/packet/register/ack"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/util"
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

		switch headerBuff.Flags {
		case option.MsgTypeRegisterAck:
			regAck, err := ack.Decode(buff)
			if err != nil {
				//return err
			}
			logger.Infof("register success, got server server ack: (%v)", regAck.AutoIP)
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
				logger.Debugf("got remote node: mac: %v, ip: %s,  natIP: %s, natPort: %d", info.Mac, info.IP, info.NatIp, info.Port)
				address, err := util.GetAddress(info.NatIp.String(), int(info.NatPort))
				if err != nil {
					logger.Errorf("resolve addr failed, err: %v", err)
				} //p2pSocket := t.GetSocket(pkt.NetworkId)
				node, err := c.GetNodeInfo(frame.NetworkId, info.IP.String())
				if node == nil || err != nil {
					sock := socket.NewSocket()
					err = sock.Connect(&address)

					if err != nil {
						logger.Errorf("open hole failed. %v", err)
						continue
					}

					//open session, node-> remote addr
					logger.Debugf("send data nat device, natIP: %s, natPort: %d", info.NatIp, info.NatPort)
					//err := sock.WriteToUdp([]byte("hello"), &address)
					sock.WriteToUdp([]byte("hello"), &address)
					if err != nil {
						logger.Errorf("%v", err)
					}
					logger.Debugf("open session to remote udp")
					nodeInfo := &cache.NodeInfo{
						Socket:  sock,
						MacAddr: info.Mac,
						IP:      info.IP,
						Port:    info.Port,
						P2P:     true,
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
		}

		return nil
	}

}
