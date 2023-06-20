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
	notifyack "github.com/topcloudz/fvpn/pkg/packet/notify/ack"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
)

var (
	logger = log.Log()
)

// Tunnel is a manager for peer-to-peer， it used for peer to registry, registry to peer, peer-to-peer
type Tunnel struct {
	socket     *socket.Socket // underlay
	devices    map[string]*tuntap.Tuntap
	Inbound    chan *packet.Frame //used from udp
	Outbound   chan *packet.Frame //used for tun
	cache      *cache.Cache
	tunHandler handler.Handler
	udpHandler handler.Handler
	manager    *Manager
}

func (t *Tunnel) Start() {
	go t.ReadFromUDP()
	go t.WriteToUdp()
	go t.WriteToTun()
}

func (t *Tunnel) Close() {
	//close a tunnel, release all resources TODO
}

func NewTunnel(tunHandler handler.Handler, s *socket.Socket, devices map[string]*tuntap.Tuntap, m []middleware.Middleware, manager *Manager) *Tunnel {
	tun := &Tunnel{
		Inbound:    make(chan *packet.Frame, 10000), // data to write to tun
		Outbound:   make(chan *packet.Frame, 10000), // data from tun to write to peer
		devices:    devices,
		cache:      cache.New(),
		tunHandler: tunHandler,
		manager:    manager,
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
			//will use in relay tunnel
			if pkt.RemoteAddr != "" && t.manager.GetNotifyPortPair(pkt.RemoteAddr) == nil {
				buff, err := t.buildNotifyMessage(pkt.RemoteAddr, pkt.NetworkId)
				if err != nil {
					logger.Errorf("send hand shake failed: %v", err)
					continue
				}

				t.socket.Write(buff)

			}
			t.socket.Write(pkt.Packet)
		default:
		}
	}
}

func (t *Tunnel) buildNotifyMessage(destIP, networkId string) ([]byte, error) {
	// send self data to remote， to tell remote to connected to.
	pkt := notify.NewPacket(networkId)
	portPair := <-Pool.ch
	t.manager.SetNotifyPortPair(destIP, portPair)
	logger.Debugf("cached port pair, source ip: %v, source port: %v, nat ip: %v, nat port: %v", portPair.SrcIP, portPair.SrcPort, portPair.NatIP, portPair.NatPort)
	//send handshake to remote
	tap := t.GetTun(networkId)
	pkt.SourceIP = tap.IP
	pkt.Port = portPair.SrcPort
	pkt.NatIP = portPair.NatIP
	pkt.NatPort = portPair.NatPort
	pkt.DestAddr = net.ParseIP(destIP)
	logger.Debugf("build a notify: source ip: %v, source port: %v, natip: %v, natport: %v", pkt.SourceIP, pkt.Port, pkt.NatIP, pkt.NatPort)
	return notify.Encode(pkt)
}

func (t *Tunnel) buildNotifyMessageAck(destIP, networkId string) ([]byte, error) {
	// send self data to remote， to tell remote to connected to.
	pkt := notifyack.NewPacket(networkId)
	portPair := <-Pool.ch
	t.manager.SetNotifyPortPair(destIP, portPair)
	logger.Debugf("cached port pair, source ip: %v, source port: %v, nat ip: %v, nat port: %v", portPair.SrcIP, portPair.SrcPort, portPair.NatIP, portPair.NatPort)
	//send handshake to remote
	tap := t.GetTun(networkId)
	pkt.SourceIP = tap.IP
	pkt.Port = portPair.SrcPort
	pkt.NatIP = portPair.NatIP
	pkt.NatPort = portPair.NatPort
	pkt.DestAddr = net.ParseIP(destIP)
	logger.Debugf("build a notify ack: source ip: %v, source port: %v, natip: %v, natport: %v", pkt.SourceIP, pkt.Port, pkt.NatIP, pkt.NatPort)
	return notifyack.Encode(pkt)
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
	}
}
