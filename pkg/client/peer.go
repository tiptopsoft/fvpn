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
	relaySocket *socket.Socket
	devices     map[string]*tuntap.Tuntap //networkId -> *Tuntap
	cache       *cache.Cache
	tunHandler  handler.Handler
	udpHandler  handler.Handler
	Outbound    chan *packet.Frame //read frame from tun

	relayTunnel *tunnel.Tunnel
	manager     *tunnel.Manager
	middlewares []middleware.Middleware
	networks    map[string]string //cidr -> networkId
}

func (p *Peer) Start() error {
	runtime.GOMAXPROCS(2)
	once.Do(func() {
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
		frame.NetworkId = networkId
		frame.Size = n
		logger.Debugf("origin packet size: %d, data: %v", n, frame.Packet[:n])
		h, err := util.GetFrameHeader(frame.Packet)

		dest := h.DestinationIP.String()
		frame.RemoteAddr = dest
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
			peerTunnel := p.getPeerTunnel(pkt.RemoteAddr)
			peerTunnel.Outbound <- pkt
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
