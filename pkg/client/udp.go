package client

import (
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/middleware"
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

func (t Tun) ReadFromUdp(middleware middleware.Middleware) {
	for {
		frame := NewFrame()
		n, err := t.socket.Read(frame.buff[:])
		if n < 0 || err != nil {
			continue
		}

		buff := frame.buff[:n]

		cpInterface, err := packet.NewPacketWithoutType().Decode(buff)
		header := cpInterface.(packet.Header)
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
			infos := peerPacketAck.PeerInfos
			logger.Infof("got server peers: (%v)", infos)
			for _, info := range infos {
				address, err := util.GetAddress(info.Host.String(), int(info.Port))
				if err != nil {
					logger.Errorf("resolve addr failed, err: %v", err)
				}
				sock := socket.NewSocket()
				err = sock.Connect(&address)
				if err != nil {
					//return err
				}
				peerInfo := &cache.Peer{
					Socket:  sock,
					MacAddr: info.Mac,
					IP:      info.Host,
					Port:    info.Port,
				}
				t.cache.Nodes[info.Mac.String()] = peerInfo
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
			frame.packet = buff[idx:]
			frame.networkId = header.NetworkId
			t.inbound <- frame
			break
		}

	}

}

// WriteToDevice write to device from the queue
func (t Tun) WriteToDevice() {
	for {
		pkt := <-t.inbound
		device, err := tuntap.GetTuntap(pkt.networkId)
		if err != nil {
			logger.Errorf("invalid network: %s", pkt.networkId)
		}
		device.Write(pkt.packet[:])
	}
}
