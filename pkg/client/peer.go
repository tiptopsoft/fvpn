package client

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/handler/device"
	"github.com/topcloudz/fvpn/pkg/middleware"
	"github.com/topcloudz/fvpn/pkg/middleware/infra"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/peer"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tunnel"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
	"runtime"
	"sync"
)

var (
	once        sync.Once
	DefaultPort = 6663
)

type Peer struct {
	*option.Config
	Protocol    option.Protocol
	relaySocket socket.Socket
	relayAddr   *unix.SockaddrInet4
	devices     map[string]*tuntap.Tuntap //networkId -> *Tuntap
	cache       *cache.Cache
	tunHandler  handler.Handler
	udpHandler  handler.Handler
	Outbound    chan *packet.Frame //read frame from tun
	Inbound     chan *packet.Frame // write frame to tun

	relayTunnel *tunnel.Tunnel
	manager     *tunnel.Manager
	//tunnels     map[string]*tunnel.Tunnel // map addr->tunnel p2p tunnels
	middlewares []middleware.Middleware
	networks    map[string]string //cidr -> networkId
}

func (p *Peer) Start() error {
	runtime.GOMAXPROCS(2)
	once.Do(func() {
		p.relaySocket = socket.NewSocket(0)
		p.Protocol = option.UDP
		if err := p.conn(); err != nil {
			logger.Errorf("failed to connect to server: %v", err)
		}
		p.devices = make(map[string]*tuntap.Tuntap, 1)
		p.Outbound = make(chan *packet.Frame, 10000)
	})

	p.manager = tunnel.NewManager()
	p.middlewares = p.initMiddleware()
	p.tunHandler = middleware.WithMiddlewares(device.Handle(), p.middlewares...)
	p.relayTunnel = tunnel.NewTunnel(p.tunHandler, p.relaySocket, p.devices, p.middlewares, p.manager)
	p.relayTunnel.Start()

	go p.WriteToUDP()
	go p.WriteToTun()
	return p.runHttpServer()
}

// initMiddleware TODO add impl
func (p *Peer) initMiddleware() []middleware.Middleware {
	return infra.Middlewares(p.OpenAuth, p.OpenEncrypt)
}

// ReadFromTun every tap will start a loop read from tap,and write to remote
func (p *Peer) ReadFromTun(tun *tuntap.Tuntap, networkId string) {
	logger.Debugf("start peer read from tun loop.....")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "tun", tun)
	ctx = context.WithValue(ctx, "networkId", networkId)
	for {
		frame := packet.NewFrame()
		n, err := tun.Read(frame.Buff[:])
		frame.Packet = frame.Buff[:n]
		frame.Size = n
		logger.Debugf("origin packet size: %d, data: %v", n, frame.Packet[:n])
		h, err := util.GetFrameHeader(frame.Packet)
		if err != nil {
			logger.Debugf("no packet...")
			continue
		}
		ctx = context.WithValue(ctx, "header", h)
		err = p.tunHandler.Handle(ctx, frame)
		if err != nil {
			logger.Errorf("tun handle packet failed: %v", err)
		}

		p.Outbound <- frame
	}
}

func (p *Peer) WriteToTun() {
	for {
		select {
		case pkt := <-p.Inbound:
			networkId := pkt.NetworkId
			tun := p.GetTun(networkId)
			if tun != nil {
				tun.Write(pkt.Packet)
			}
		default:
			return
		}
	}
}

func (p *Peer) GetTun(networkId string) *tuntap.Tuntap {
	return p.devices[networkId]
}

// WriteToUDP  data from tun write to destination
func (p *Peer) WriteToUDP() {
	logger.Debugf("start peer write to udp loop.....")
	for {
		select {
		case pkt := <-p.Outbound:
			packetHeader, err := util.GetPacketHeader(pkt.Packet[:])
			if err != nil {
				logger.Errorf("%v", "buff not encoded by fvpn")
				return
			}

			logger.Debugf("pkt type: %v", packetHeader.Flags)
			frameHeader, err := util.GetFrameHeader(pkt.Packet[12:]) //why 12? because packer.Header length is 12.
			dest := frameHeader.DestinationIP.String()

			//if p2p use a p2p tunnel, if not use relay tunnel
			peerTunnel := p.getPeerTunnel(dest)
			peerTunnel.Outbound <- pkt
			//if pkt.Type == option.PacketFromTap {
			//
			//	frameHeader, err := util.GetFrameHeader(pkt.Packet[12:]) //why 12? because packer.Header length is 12.
			//	logger.Debugf("packet will be write to : mac: %s, ip: %s, content: %v", frameHeader.DestinationAddr, frameHeader.DestinationIP.String(), pkt.Packet)
			//	if err != nil {
			//		logger.Errorf("%v", err)
			//		return
			//	}
			//
			//	//target
			//	ip := frameHeader.DestinationIP.String()
			//	target, err := p.cache.GetNodeInfo(pkt.NetworkId, ip)
			//	if err != nil {
			//		//err := t.sendQueryPeer(pkt.NetworkId)
			//		//if err != nil {
			//		//	logger.Errorf("%v", err)
			//		//}
			//		return
			//	}
			//
			//	if target.NatType == option.SymmetricNAT {
			//		//use relay server
			//		logger.Debugf("use relay server to connect to: %v", target.IP.String())
			//		_, err := p.relaySocket.Write(pkt.Packet[:])
			//		if err != nil {
			//			return
			//		}
			//	} else if target.P2P {
			//		logger.Debugf("use p2p to connect to: %v, remoteAddr: %v, sock: %v", target.IP, target.Addr, target.Socket)
			//		if _, err := target.Socket.Write(pkt.Packet); err != nil {
			//			logger.Errorf("send p2p data failed. %v", err)
			//		}
			//	} else {
			//		//同时通过relay server发送数据
			//		p.relaySocket.Write(pkt.Packet[:])
			//
			//		//同时进行punch hole
			//		//go .sendNotifyMessage(pkt.NetworkId, t.relayAddr, ip, option.MsgTypeNotify)
			//	}
			//} else {
			//	p.relaySocket.Write(pkt.Packet)
			//}

		default:

		}
	}

}

// SendRegister register register a node to center.
func (p *Peer) SendRegister(tun *tuntap.Tuntap) error {
	var err error
	//header, err := packet.NewHeader(option.MsgTypeRegisterSuper, tun.NetworkId)
	srcMac, srcIP, err := addr.GetMacAddrAndIPByDev(tun.Name)
	if err != nil {
		return err
	}

	if srcIP == nil {
		return errors.New("device ip not set")
	}
	regPkt := register.NewPacket(tun.NetworkId, srcMac, srcIP)
	copy(regPkt.SrcMac[:], tun.MacAddr)
	if err != nil {
		return err
	}

	data, err := register.Encode(regPkt)
	if err != nil {
		return err
	}
	frame := packet.NewFrame()
	frame.Packet = data
	frame.FrameType = option.MsgTypeRegister
	frame.FrameType = option.PacketFromUdp

	p.Outbound <- frame
	return nil
}

func (p *Peer) sendQueryPeer(networkId string) error {
	pkt := peer.NewPacket(networkId)
	data, err := peer.Encode(pkt)
	if err != nil {
		logger.Errorf("query data failed: %v", err)
	}

	frame := packet.NewFrame()
	frame.Packet = data
	frame.FrameType = option.MsgTypeQueryPeer
	frame.Type = option.PacketFromUdp
	p.Outbound <- frame

	return nil
}

func (p *Peer) getPeerTunnel(dest string) *tunnel.Tunnel {
	t := p.manager.GetTunnel(dest)
	if t != nil {
		return t
	}

	return p.relayTunnel
}

func (p *Peer) AppId() string {
	buf := addr.GetLocalMacAddr()
	appId := hex.EncodeToString(buf)
	return appId
}
