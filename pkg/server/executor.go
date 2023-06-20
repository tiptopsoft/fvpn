package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/handshake"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/packet/notify"
	notifyack "github.com/topcloudz/fvpn/pkg/packet/notify/ack"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/security"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
	"net"
)

func (r *RegServer) ReadFromUdp() {
	logger.Infof("start a udp loop")
	for {
		ctx := context.Background()
		frame := packet.NewFrame()
		n, addr, err := r.socket.ReadFromUDP(frame.Buff[:])
		frame.Packet = frame.Buff[:n]
		logger.Debugf("Read from udp %d byte, data: %v", n, frame.Packet)

		packetHeader, err := util.GetPacketHeader(frame.Packet[:12])
		if err != nil {
			logger.Errorf("get header falied. %v", err)
			continue
		}
		networkId := hex.EncodeToString(packetHeader.NetworkId[:])
		//ctx = context.WithValue(ctx, "pkt", packetHeader)
		ctx = context.WithValue(ctx, "flag", packetHeader.Flags)
		ctx = context.WithValue(ctx, "networkId", networkId)
		ctx = context.WithValue(ctx, "size", n)
		ctx = context.WithValue(ctx, "srcAddr", transferSockAddr(addr))
		frame.NetworkId = networkId
		if err != nil || n < 0 {
			logger.Warnf("no data exists")
			continue
		}
		err = r.h.Handle(ctx, frame)
		if err != nil {
			logger.Errorf(err.Error())
			continue
		}
		r.Outbound <- frame
	}
}

func (r *RegServer) WriteToUdp() {
	logger.Infof("start a udp write loop")
	for {
		pkt := <-r.Outbound
		frameType := pkt.FrameType
		switch frameType {
		case option.MsgTypePacket:
			frameHeader, err := util.GetFrameHeader(pkt.Packet[12:]) //why is 12, because we add our header in, header length is 12
			if err != nil {
				logger.Debugf("get header failed, dest ip: %s", frameHeader.DestinationIP.String())
			}
			//
			nodeInfo, err := r.cache.GetNodeInfo(pkt.NetworkId, frameHeader.DestinationIP.String())
			if nodeInfo == nil || err != nil {
				logger.Debugf("could not found destitation, destIP: %s", frameHeader.DestinationIP.String())
			} else {
				logger.Infof("packet will relay to: %v", nodeInfo.Addr)
				r.socket.WriteToUDP(pkt.Packet[:], transferUdpAddr(nodeInfo.Addr))
			}

			break
		case option.MsgTypeRegisterAck:
			r.socket.WriteToUDP(pkt.Packet, transferUdpAddr(pkt.SrcAddr))
			break
		case option.MsgTypeQueryPeer:
			logger.Debugf("query nodes result: %v, write to: %v", pkt.Packet, pkt.SrcAddr)
			_, err := r.socket.WriteToUDP(pkt.Packet, transferUdpAddr(pkt.SrcAddr))
			if err != nil {
				logger.Errorf("write query to dest failed: %v", err)
			}
			break
		case option.MsgTypeNotify:
			//write to dest
			np, err := notify.Decode(pkt.Packet[:])
			logger.Debugf("got notify packet: %v, destAddr: %s, networkId: %s", pkt.Packet[:], np.DestAddr.String(), pkt.NetworkId)
			if err != nil {
				logger.Errorf("invalid notify packet: %v", err)
			}

			nodeInfo, err := r.cache.GetNodeInfo(pkt.NetworkId, np.DestAddr.String())
			if nodeInfo == nil || err != nil {
				logger.Errorf("node not on line, err: %v", err)
				break
			}

			r.socket.WriteToUDP(pkt.Packet, transferUdpAddr(nodeInfo.Addr))

		case option.MsgTypeNotifyAck:
			//write to dest
			np, err := notifyack.Decode(pkt.Packet[:])
			logger.Debugf("got notify ack packet: %v, destAddr: %s, networkId: %s", pkt.Packet[:], np.DestAddr.String(), pkt.NetworkId)
			if err != nil {
				logger.Errorf("invalid notify ack packet: %v", err)
			}

			nodeInfo, err := r.cache.GetNodeInfo(pkt.NetworkId, np.DestAddr.String())
			if nodeInfo == nil || err != nil {
				logger.Errorf("node not on line, err: %v", err)
				break
			}

			r.socket.WriteToUDP(pkt.Packet, transferUdpAddr(nodeInfo.Addr))
		}

	}
}

func transferUdpAddr(address unix.Sockaddr) *net.UDPAddr {
	addr := address.(*unix.SockaddrInet4)
	ip := net.ParseIP(fmt.Sprintf("%d.%d.%d.%d", addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3]))
	return &net.UDPAddr{IP: ip, Port: addr.Port}
}

func transferSockAddr(address *net.UDPAddr) *unix.SockaddrInet4 {
	addr := &unix.SockaddrInet4{
		Port: address.Port,
		Addr: [4]byte{},
	}
	copy(addr.Addr[:], address.IP.To4())
	return addr
}

// serverUdpHandler  core self handler
func (r *RegServer) serverUdpHandler() handler.HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {

		srcAddr := ctx.Value("srcAddr").(*unix.SockaddrInet4)
		networkId := ctx.Value("networkId").(string)
		size := ctx.Value("size").(int)
		data := frame.Packet[:]

		flag := ctx.Value("flag").(uint16)

		switch flag {

		case option.MsgTypeRegisterSuper:
			regPkt, err := register.Decode(frame.Packet)
			if err != nil {
				logger.Errorf("register failed, err:%v", err)
				return err
			}
			err = r.registerAck(srcAddr, regPkt.SrcMac, regPkt.SrcIP, networkId)
			h, err := header.NewHeader(option.MsgTypeRegisterAck, networkId)
			if err != nil {
				logger.Errorf("build resp failed. err: %v", err)
			}
			f, _ := header.Encode(h)
			frame.Packet = f
			frame.SrcAddr = srcAddr
			frame.FrameType = option.MsgTypeRegisterAck
			break
		case option.MsgTypeQueryPeer:
			peers, size, err := getPeerInfo(r.cache.GetNodes())
			logger.Debugf("server peers: (%v), size: (%v)", peers, size)
			if err != nil {
				logger.Errorf("get peers from server failed. err: %v", err)
			}

			f, err := peerAckBuild(peers, size, networkId)
			if err != nil {
				logger.Errorf("get peer ack from server failed. err: %v", err)
			}

			frame.Packet = f
			frame.SrcAddr = srcAddr
			frame.FrameType = option.MsgTypeQueryPeer
			break
		case option.MsgTypePacket:
			logger.Infof("server got forward packet size:%d, data: %v", size, data)
			frame.FrameType = option.MsgTypePacket
			break
		case option.MsgTypeNotify:
			logger.Debugf("notify frame packet: %v", frame.Packet[:])
			frame.FrameType = option.MsgTypeNotify
		case option.MsgTypeNotifyAck:
			logger.Debugf("notify ack frame packet: %v", frame.Packet[:])
			frame.FrameType = option.MsgTypeNotify
		case option.HandShakeMsgType:
			handPkt, err := handshake.Decode(frame.Packet)
			if err != nil {
				logger.Errorf("invalid handshake packet: %v", err)
				return err
			}
			privateKey, err := security.NewPrivateKey()
			if err != nil {
				return err
			}
			pubKey := privateKey.NewPubicKey()
			r.PrivateKey = privateKey
			r.PubKey = pubKey

			r.cipher = security.NewCipher(r.PrivateKey, handPkt.PubKey)
		}

		return nil
	}
}
