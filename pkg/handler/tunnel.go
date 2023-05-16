package handler

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/peer"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"sync"
	"time"
)

type Tun struct {
	socket     socket.Interface // underlay
	p2pSocket  sync.Map         //p2psocket
	device     map[string]*tuntap.Tuntap
	Inbound    chan *packet.Frame //used from udp
	Outbound   chan *packet.Frame //used for tun
	QueryBound chan *packet.Frame
	cache      *cache.Cache
	tunHandler Handler
	udpHandler Handler
	NetworkId  string
}

func NewTun(tunHandler, udpHandler Handler, socket socket.Interface) *Tun {
	tun := &Tun{
		Inbound:    make(chan *packet.Frame, 10000),
		Outbound:   make(chan *packet.Frame, 10000),
		QueryBound: make(chan *packet.Frame, 10000),
		device:     make(map[string]*tuntap.Tuntap),
		cache:      cache.New(),
		tunHandler: tunHandler,
		udpHandler: udpHandler,
	}
	tun.socket = socket
	return tun
}

func (t *Tun) CacheDevice(networkId string, device *tuntap.Tuntap) {
	if t.device[networkId] == nil {
		t.device[networkId] = device
	}
}

func (t *Tun) ReadFromTun(ctx context.Context, networkId string) {
	time.Sleep(1 * time.Second)
	logger.Infof("start a tun loop for networkId: %s", networkId)
	ctx = context.WithValue(ctx, "networkId", networkId)
	tun := t.device[networkId]
	ctx = context.WithValue(ctx, "tun", tun)
	if tun == nil {
		logger.Fatalf("invalid network: %s", networkId)
	}
	for {
		frame := packet.NewFrame()
		n, err := tun.Read(frame.Buff[:])
		frame.Packet = frame.Buff[:n]
		frame.Size = n
		logger.Debugf("origin packet size: %d, data: %v", n, frame.Packet)
		header, err := util.GetFrameHeader(frame.Packet)
		if err != nil {
			logger.Debugf("no packet...")
			continue
		}
		ctx = context.WithValue(ctx, "header", header)
		err = t.tunHandler.Handle(ctx, frame)
		if err != nil {
			logger.Errorf("tun handle packet failed: %v", err)
		}
		// 放入chan
		t.Outbound <- frame
	}
}

func (t *Tun) WriteToUdp() {
	for {
		pkt := <-t.Outbound
		//这里先尝试P2p, 没有P2P使用relay server
		header, err := util.GetFrameHeader(pkt.Packet[12:]) //why 12? because header length is 12.
		logger.Infof("packet will be write to : mac: %s, ip: %s, content: %v", header.DestinationAddr, header.DestinationIP.String(), pkt.Packet)
		if err != nil {
			continue
		}

		p2pSocket := t.GetSocket(header.DestinationAddr.String())
		if p2pSocket == nil {
			err := t.addQueryRemoteNodes(pkt.NetworkId)
			if err != nil {
				logger.Errorf("add query queue failed. err: %v", err)
			}
			t.socket.Write(pkt.Packet)
			nodeInfo, err := t.cache.GetNodeInfo(t.NetworkId, header.DestinationIP.String())
			if err != nil {
				logger.Debugf("got nodeInfo failed")
			} else {
				t.SaveSocket(header.DestinationAddr.String(), nodeInfo.Socket)
				//启动一个udp goroutine用于处理P2P的轮询
				go func() {
					newTun := NewTun(t.tunHandler, t.udpHandler, nodeInfo.Socket)
					newTun.ReadFromUdp()
					newTun.WriteToDevice()
				}()
			}
		} else {
			p2pSocket.Write(pkt.Packet)
		}

	}
}

func (t *Tun) GetSocket(mac string) socket.Interface {
	v, b := t.p2pSocket.Load(mac)
	if !b {
		return nil
	}

	return v.(socket.Interface)
}

func (t *Tun) SaveSocket(mac string, s socket.Interface) {
	t.p2pSocket.Store(mac, s)
}

func (t *Tun) ReadFromUdp() {
	logger.Infof("start a udp loop socket is: %v", t.socket)
	for {
		ctx := context.Background()
		frame := packet.NewFrame()

		n, err := t.socket.Read(frame.Buff[:])
		logger.Debugf("receive data from remote, size: %d, data: %v", n, frame.Buff[:])
		if n < 0 || err != nil {
			continue
		}
		err = t.udpHandler.Handle(ctx, frame)
		if err != nil {
			logger.Errorf("Read from udp failed: %v", err)
			continue
		}
		t.Inbound <- frame
	}

}

// WriteToDevice write to device from the queue
func (t *Tun) WriteToDevice() {
	for {
		pkt := <-t.Inbound
		device := t.device[pkt.NetworkId]
		if device == nil {
			logger.Errorf("invalid network: %s", pkt.NetworkId)
		}
		logger.Debugf("write to device data :%v", pkt.Packet[12:])
		_, err := device.Write(pkt.Packet[12:]) // start 12, because header length 12
		if err != nil {
			logger.Errorf("write to device err: %v", err)
		}
	}
}

// 添加一个networkID，查询该networkId下节点，更新cache
func (t *Tun) addQueryRemoteNodes(networkId string) error {
	pkt := peer.NewPacket(networkId)
	buff, err := peer.Encode(pkt)
	if err != nil {
		return err
	}
	frame := packet.NewFrame()
	frame.NetworkId = networkId
	frame.Packet = buff
	frame.Buff = buff
	t.QueryBound <- frame
	return nil
}

// QueryRemoteNodes when packet from regserver, this method will be called
func (t *Tun) QueryRemoteNodes(networkId string) {
	for {
		pkt := <-t.QueryBound
		t.socket.Write(pkt.Packet)
	}

}
