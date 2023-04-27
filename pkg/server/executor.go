package server

import (
	"context"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
)

func (r *RegServer) ReadFromUdp() {
	logger.Infof("start a udp loop")
	for {
		ctx := context.Background()
		frame := packet.NewFrame()
		n, addr, err := r.socket.ReadFromUdp(frame.Buff[:])
		ctx = context.WithValue(ctx, "size", n)
		ctx = context.WithValue(ctx, "addr", addr)
		if err != nil || n < 0 {
			logger.Warnf("no data exists")
			continue
		}
		r.h.Handle(ctx, frame)
		r.Outbound <- frame

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
			_, destIP, err := util.GetMacAddr(pkt.Packet)
			if err != nil {
				logger.Errorf("dest ip :%s not on line", destIP)
			}

			nodeInfo, err := r.cache.GetNodeInfo("", destIP.String())
			r.socket.WriteToUdp(pkt.Packet[:], nodeInfo.Addr)
		}
	}
}

func (r *RegServer) serverUdpHandler() handler.HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {
		addr := ctx.Value("addr").(unix.Sockaddr)
		size := ctx.Value("size").(int)
		data := frame.Packet[:]
		pInterface, err := packet.NewPacketWithoutType().Decode(data)
		p := pInterface.(packet.Header)

		if err != nil {
			fmt.Println(err)
		}

		switch p.Flags {

		case option.MsgTypeRegisterSuper:
			//packet = register.NewPacket("")
			//processRegister(addr, data[:size], nil)
			p := register.NewPacket("")
			registerPacket, err := p.Decode(data[:size])

			// build an ack
			f, err := registerAck(addr, registerPacket.(register.RegPacket).SrcMac)
			logger.Infof("build a server ack: %v", f)
			if err != nil {
				logger.Errorf("build resp failed. err: %v", err)
			}
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
			logger.Infof("server got forward packet: %v", data)
			break
		}

		return nil
	}
}
