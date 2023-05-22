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

		header, err := header.Decode(buff)
		if err != nil {
			logger.Errorf("decode err: %v", err)
		}

		frame.NetworkId = hex.EncodeToString(header.NetworkId[:])

		switch header.Flags {
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
				logger.Debugf("got remote node: mac: %v, ip: %s", info.Mac, info.IP)
				address, err := util.GetAddress(info.NatIp.String(), int(info.NatPort))
				if err != nil {
					logger.Errorf("resolve addr failed, err: %v", err)
				}
				sock := socket.NewSocket()
				err = sock.Connect(&address)
				if err != nil {
					//return err
					logger.Errorf("%v", err)
					continue
				}
				nodeInfo := &cache.NodeInfo{
					Socket:  sock,
					MacAddr: info.Mac,
					IP:      info.IP,
					Port:    info.Port,
					P2P:     true,
				}
				c := ctx.Value("cache").(*cache.Cache)
				//cache.Nodes[info.Mac.String()] = nodeInfo
				c.SetCache(frame.NetworkId, info.IP.String(), nodeInfo)
			}
			break
		case option.MsgTypePacket:
			frame.Packet = buff[:]
			break
		}

		return nil
	}

}
