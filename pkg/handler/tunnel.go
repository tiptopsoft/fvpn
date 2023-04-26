package handler

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"sync"
)

type Tun struct {
	socket     socket.Interface // relay or p2p
	p2pSocket  sync.Map         //p2psocket
	device     *tuntap.Tuntap
	Inbound    chan *packet.Frame //used from udp
	Outbound   chan *packet.Frame //used for tun
	cache      *cache.Cache
	tunHandler Handler
	udpHandler Handler
}

func NewTun(tunHandler, udpHandler Handler, socket socket.Interface) *Tun {
	tun := &Tun{
		Inbound:    make(chan *packet.Frame, 15000),
		Outbound:   make(chan *packet.Frame, 15000),
		cache:      cache.New(),
		tunHandler: tunHandler,
		udpHandler: udpHandler,
	}

	tun.socket = socket

	return tun
}

func (t *Tun) ReadFromTun(ctx context.Context, networkId string) {
	logger.Infof("start a tun loop for networkId: %s", networkId)
	for {
		ctx = context.WithValue(ctx, "networkId", networkId)
		networkId := ctx.Value("networkId").(string)
		tun, err := tuntap.GetTuntap(networkId)
		if err != nil {
			logger.Fatalf("invalid network: %s", networkId)
		}
		frame := packet.NewFrame()
		n, err := tun.Read(frame.Buff[:])
		frame.Packet = frame.Buff[:n]
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
		destMac := util.GetMacAddr(pkt.Packet)

		p2pSocket := t.GetSocket(destMac)
		if p2pSocket == nil {
			nodeInfo, err := t.cache.GetNodeInfo(destMac)
			if err != nil {
				logger.Errorf("got nodeInfo failed.")
			}
			t.SaveSocket(destMac, nodeInfo.Socket)
			t.socket.Write(pkt.Packet)
			//启动一个udp goroutine用于处理P2P的轮询
			go func() {
				newTun := NewTun(t.tunHandler, t.udpHandler, nodeInfo.Socket)
				newTun.ReadFromUdp()
				newTun.WriteToDevice()
			}()
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
	logger.Infof("start a udp loop socket is: %s", t.socket)
	for {
		ctx := context.Background()
		frame := packet.NewFrame()

		n, err := t.socket.Read(frame.Buff[:])
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
		device, err := tuntap.GetTuntap(pkt.NetworkId)
		if err != nil {
			logger.Errorf("invalid network: %s", pkt.NetworkId)
		}
		device.Write(pkt.Packet[:])
	}
}
