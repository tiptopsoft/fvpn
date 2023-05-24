package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/packet/notify"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
	"net"
)

func (r *RegServer) ReadFromUdp() {
	logger.Infof("start a udp loop")
	for {
		ctx := context.Background()
		frame := packet.NewFrame()
		n, addr, err := r.socket.ReadFromUdp(frame.Buff[:])
		frame.Packet = frame.Buff[:n]
		logger.Debugf("Read from udp %d byte, data: %v", n, frame.Packet)

		packetHeader, err := util.GetPacketHeader(frame.Packet[:12])
		if err != nil {
			logger.Errorf("get header falied. %v", err)
		}
		if packetHeader.Flags == option.MsgTypePacket {
			header, err := util.GetFrameHeader(frame.Packet[12:])
			if err != nil {
				logger.Errorf("get invalid header..:%v", err)
			}
			ctx = context.WithValue(ctx, "header", header)
		}
		networkId := hex.EncodeToString(packetHeader.NetworkId[:])
		ctx = context.WithValue(ctx, "pkt", packetHeader)
		ctx = context.WithValue(ctx, "flag", packetHeader.Flags)
		ctx = context.WithValue(ctx, "networkId", networkId)
		ctx = context.WithValue(ctx, "size", n)
		ctx = context.WithValue(ctx, "srcAddr", addr)
		frame.NetworkId = networkId
		if err != nil || n < 0 {
			logger.Warnf("no data exists")
			continue
		}
		err = r.h.Handle(ctx, frame)
		if err != nil {
			logger.Errorf(err.Error())
		}
		r.Outbound <- frame
		logger.Infof("success handler frame")
	}
}

func (r *RegServer) WriteToUdp() {
	logger.Infof("start a udp write loop")
	for {
		pkt := <-r.Outbound
		packetHeader, err := util.GetPacketHeader(pkt.Packet[:12])
		if err != nil {
			logger.Errorf("get header failed")
		}

		switch packetHeader.Flags {
		case option.MsgTypePacket:
			header, err := util.GetFrameHeader(pkt.Packet[12:]) //why is 12, because we add our header in, header length is 12
			if err != nil {
				logger.Debugf("get header failed, dest ip: %s", header.DestinationIP.String())
			}

			nodeInfo, err := r.cache.GetNodeInfo(pkt.NetworkId, header.DestinationIP.String())
			if nodeInfo == nil || err != nil {
				logger.Debugf("could not found destitation, destIP: %s", header.DestinationIP.String())
			} else {
				r.socket.WriteToUdp(pkt.Packet[:], nodeInfo.Addr)
			}
			break
		case option.MsgTypeRegisterAck:
			r.socket.WriteToUdp(pkt.Packet, pkt.RemoteAddr)
			break
		case option.MsgTypeQueryPeer:
			r.socket.WriteToUdp(pkt.Packet, pkt.RemoteAddr)
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

			r.socket.WriteToUdp(pkt.Packet, nodeInfo.Addr)
		}

	}
}

// serverUdpHandler  core self handler
func (r *RegServer) serverUdpHandler() handler.HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {

		srcAddr := ctx.Value("srcAddr").(unix.Sockaddr)
		networkId := ctx.Value("networkId").(string)
		size := ctx.Value("size").(int)
		data := frame.Packet[:]

		p := ctx.Value("pkt").(header.Header)
		switch p.Flags {

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
			frame.RemoteAddr = srcAddr
			break
		case option.MsgTypeQueryPeer:
			peers, size, err := getPeerInfo(r.cache.GetNodes())
			logger.Infof("server peers: (%v), size: (%v)", peers, size)
			if err != nil {
				logger.Errorf("get peers from server failed. err: %v", err)
			}

			f, err := peerAckBuild(peers, size, networkId)
			if err != nil {
				logger.Errorf("get peer ack from server failed. err: %v", err)
			}

			frame.Packet = f
			frame.RemoteAddr = srcAddr
			break
		case option.MsgTypePacket:
			logger.Infof("server got forward packet size:%d, data: %v", size, data)
			break
		case option.MsgTypeNotify:
			logger.Infof("add packet: %v", frame.Packet[:])
			np, err := notify.Decode(frame.Packet[:])
			if err != nil {
				logger.Errorf("invalid notify packet: %v", err)
			}

			//add nat info to packet
			addr := srcAddr.(*unix.SockaddrInet4)
			natIP := net.ParseIP(fmt.Sprintf("%d.%d.%d.%d", addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3]))
			np.NatAddr = natIP
			np.NatPort = uint16(addr.Port)

			newBuff, err := notify.Encode(np)
			logger.Debugf("new notify buff: %v", newBuff[:])
			if err != nil {
				logger.Errorf("encode failed:%v", err)
			}

			copy(frame.Packet[:], newBuff)
			logger.Debugf("frame packet: %v", frame.Packet[:])
		}

		return nil
	}
}

func getFrameHeader(ctx context.Context) (*util.FrameHeader, error) {
	return ctx.Value("header").(*util.FrameHeader), nil
}
