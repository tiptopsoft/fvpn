package tunnel

import (
	"context"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/middleware"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/notify"
	"github.com/topcloudz/fvpn/pkg/packet/notify/ack"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"sync"
)

var (
	logger = log.Log()
)

// Tunnel is a manager for peer-to-peer， it used for peer to registry, registry to peer, peer-to-peer
type Tunnel struct {
	IsP2P         bool
	socket        socket.Socket // underlay
	p2pSocket     sync.Map      //p2psocket
	devices       map[string]*tuntap.Tuntap
	Inbound       chan *packet.Frame //used from udp
	Outbound      chan *packet.Frame //used for tun
	QueryBound    chan *packet.Frame
	RegisterBound chan *packet.Frame
	P2PBound      chan *P2PNode
	cache         *cache.Cache
	tunHandler    handler.Handler
	udpHandler    handler.Handler
	p2pNode       sync.Map // ip target -> socket
	manager       *Manager
}

type P2PNode struct {
	NodeInfo *cache.Endpoint
	Frame    *packet.Frame
	Socket   socket.Socket
}

func (t *Tunnel) Start() {
	go t.ReadFromUDP()
	go t.WriteToUdp()
	go t.WriteToTun()
}

func (t *Tunnel) Close() {
	//close a tunnel, release all resources
}

func NewTunnel(tunHandler handler.Handler, s socket.Socket, devices map[string]*tuntap.Tuntap, m []middleware.Middleware, manager *Manager) *Tunnel {
	tun := &Tunnel{
		Inbound:       make(chan *packet.Frame, 10000), // data to write to tun
		Outbound:      make(chan *packet.Frame, 10000), // data from tun to write to peer
		QueryBound:    make(chan *packet.Frame, 10000),
		RegisterBound: make(chan *packet.Frame, 10000),
		P2PBound:      make(chan *P2PNode, 10000),
		devices:       devices,
		cache:         cache.New(),
		tunHandler:    tunHandler,
		manager:       manager,
	}
	tun.socket = s
	tun.udpHandler = middleware.WithMiddlewares(tun.Handle(), m...)

	return tun
}

func (t *Tunnel) CacheDevice(networkId string, device *tuntap.Tuntap) {
	if t.devices[networkId] == nil {
		t.devices[networkId] = device
	}
}

// GetSelf get self node from cache
func (t *Tunnel) GetSelf(networkId string) (*cache.Endpoint, error) {
	device := t.devices[networkId]
	if device == nil {
		return nil, fmt.Errorf("you have not to join this network: %s", networkId)
	}

	ip := device.IP
	return t.cache.GetNodeInfo(networkId, ip.String())
}

func (t *Tunnel) findNode(networkId, ip string) (*cache.Endpoint, error) {
	return t.cache.GetNodeInfo(networkId, ip)
}

func (t *Tunnel) WriteToTun() {
	for {
		select {
		case pkt := <-t.Inbound:
			networkId := pkt.NetworkId
			tun := t.GetTun(networkId)
			if tun != nil {
				tun.Write(pkt.Packet[12:])
			}
		default:
		}
	}
}

func (t *Tunnel) GetTun(networkId string) *tuntap.Tuntap {
	return t.devices[networkId]
}

// WriteToUdp data write to remote
func (t *Tunnel) WriteToUdp() {
	for {
		select {
		case pkt := <-t.Outbound:
			if t.IsP2P {
				//here is the default relay tunnel
				t.socket.Write(pkt.Packet)
			} else {
				if pkt.RemoteAddr != "" && !t.manager.GetNotifyStatus(pkt.RemoteAddr) {
					buff, err := t.buildNotifyMessage(pkt.RemoteAddr, pkt.NetworkId)
					if err != nil {
						logger.Errorf("send hand shake failed: %v", err)
						continue
					}

					t.socket.Write(buff)
					t.manager.SetNotifyStatus(pkt.RemoteAddr, true)
				}
				t.socket.Write(pkt.Packet)
			}
		default:
		}
	}
}

func (t *Tunnel) buildNotifyMessage(destIP, networkId string) ([]byte, error) {
	// send self data to remote， to tell remote to connected to.
	pkt := notify.NewPacket(networkId)
	portPair := <-Pool.ch
	//send handshake to remote
	tap := t.GetTun(networkId)
	pkt.SourceIP = tap.IP
	pkt.Port = portPair.SrcPort
	pkt.NatIP = portPair.NatIP
	pkt.NatPort = portPair.NatPort
	pkt.DestAddr = net.ParseIP(destIP)
	return notify.Encode(pkt)
}
func (t *Tunnel) buildNotifyMessageAck(destIP, networkId string) ([]byte, error) {
	// send self data to remote， to tell remote to connected to.
	pkt := ack.NewPacket(networkId)
	portPair := <-Pool.ch
	//send handshake to remote
	tap := t.GetTun(networkId)
	pkt.SourceIP = tap.IP
	pkt.Port = portPair.SrcPort
	pkt.NatIP = portPair.NatIP
	pkt.NatPort = portPair.NatPort
	pkt.DestAddr = net.ParseIP(destIP)
	return ack.Encode(pkt)
}

var m sync.Mutex

func (t *Tunnel) GetSocket(targetIP string) socket.Socket {
	v, b := t.p2pSocket.Load(targetIP)
	if !b {
		return socket.Socket{}
	}

	return v.(socket.Socket)
}

func (t *Tunnel) SaveSocket(target string, s socket.Socket) {
	t.p2pSocket.Store(target, s)
}

// ReadFromUDP read data from remote peer， write back or write to tun
func (t *Tunnel) ReadFromUDP() {
	logger.Debugf("start a udp read from udp socket is: %v", t.socket)

	for {
		ctx := context.Background()
		frame := packet.NewFrame()

		n, remoteAddr, err := t.socket.ReadFromUDP(frame.Buff[:])
		h, _ := util.GetFrameHeader(frame.Buff[:])
		if n < 0 || err != nil {
			logger.Errorf("got data err: %v", err)
			continue
		}
		logger.Debugf("receive data from remote: %v, size: %d, data: %v", remoteAddr, n, frame.Buff[:n])
		ctx = context.WithValue(ctx, "cache", t.cache)
		ctx = context.WithValue(ctx, "destAddr", h.DestinationIP.String())
		err = t.udpHandler.Handle(ctx, frame)
		if err != nil {
			logger.Errorf("Read from udp failed: %v", err)
		}

		//

	}
}
