package server

import (
	"context"
	"encoding/hex"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
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
		logger.Debugf("Read from udp %d byte", n)

		packetHeader := util.GetPacketHeader(frame.Packet[:12])
		if packetHeader.Flags == option.MsgTypePacket {
			header, err := util.GetFrameHeader(frame.Packet[12:])
			if err != nil {
				logger.Errorf("get invalid header..:%v", err)
			}
			ctx = context.WithValue(ctx, "header", header)
		}
		ctx = context.WithValue(ctx, "pkt", packetHeader)
		ctx = context.WithValue(ctx, "flag", packetHeader.Flags)
		ctx = context.WithValue(ctx, "networkId", hex.EncodeToString(packetHeader.NetworkId[:]))
		ctx = context.WithValue(ctx, "size", n)
		ctx = context.WithValue(ctx, "srcAddr", addr)
		if err != nil || n < 0 {
			logger.Warnf("no data exists")
			continue
		}
		err = r.h.Handle(ctx, frame)
		if err != nil {
			logger.Errorf(err.Error())
		}
		r.Outbound <- frame
		logger.Infof("succes handler frame")
	}
}

func (r *RegServer) WriteToUdp() {
	logger.Infof("start a udp write loop")
	for {
		pkt := <-r.Outbound
		if pkt.RemoteAddr != nil {
			r.socket.WriteToUdp(pkt.Packet[:], pkt.RemoteAddr)
		} else {
			//if util.IsBroadCast(fp.DstMac.String()) {
			packetHeader := util.GetPacketHeader(pkt.Packet[:12])
			if packetHeader.Flags == option.MsgTypePacket { //转发的流量
				header, err := util.GetFrameHeader(pkt.Packet[12:]) //whe is 12, because we add our header in, header length is 12
				if err != nil {
					logger.Debugf("dest ip :%s not on line", header.DestinationIP.String())
				}

				nodeInfo, err := r.cache.GetNodeInfo(pkt.NetworkId, header.DestinationAddr.String())
				if nodeInfo == nil || err != nil {
					logger.Debugf("not found destitation")
				} else {
					r.socket.WriteToUdp(pkt.Packet[:], nodeInfo.Addr)
				}
			} else {
				//ignore
			}

		}
	}
}

func (r *RegServer) serverUdpHandler() handler.HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {

		srcAddr := ctx.Value("srcAddr").(unix.Sockaddr)
		networkId := ctx.Value("networkId").(string)
		size := ctx.Value("size").(int)
		data := frame.Packet[:]

		p := ctx.Value("pkt").(*packet.Header)
		switch p.Flags {

		case option.MsgTypeRegisterSuper:
			p := register.NewPacket(networkId, net.HardwareAddr{}, net.IP{})
			reg, err := p.Decode(frame.Packet)
			regPkt := reg.(register.RegPacket)
			if err != nil {
				logger.Errorf("register failed, err:%v", err)
				return err
			}
			err = r.registerAck(srcAddr, regPkt.SrcMac, regPkt.SrcIP, networkId)
			header, err := packet.NewHeader(option.MsgTypeRegisterAck, networkId)
			if err != nil {
				logger.Errorf("build resp failed. err: %v", err)
			}
			f, _ := header.Encode()
			frame.Packet = f
			break
		case option.MsgTypeQueryPeer:
			peers, size, err := getPeerInfo(r.cache.GetNodes())
			logger.Infof("server peers: (%v), size: (%v)", peers, size)
			if err != nil {
				logger.Errorf("get peers from server failed. err: %v", err)
			}

			f, err := peerAckBuild(peers, size)
			if err != nil {
				logger.Errorf("get peer ack from server failed. err: %v", err)
			}
			frame.Packet = f
			break
		case option.MsgTypePacket:
			logger.Infof("server got forward packet size:%d, data: %v", size, data)
			break
		}

		return nil
	}
}

func getFrameHeader(ctx context.Context) (*util.FrameHeader, error) {
	return ctx.Value("header").(*util.FrameHeader), nil
}
