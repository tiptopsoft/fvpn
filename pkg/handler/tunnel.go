package handler

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
)

type Tun struct {
	socket     socket.Interface //relay socket
	device     *tuntap.Tuntap
	Inbound    chan *packet.Frame //used from udp
	Outbound   chan *packet.Frame //used for tun
	cache      *cache.Cache
	tunHandler Handler
	udpHandler Handler
}

func NewTun(tunHandler, udpHandler Handler) *Tun {
	return &Tun{
		Inbound:    make(chan *packet.Frame, 15000),
		Outbound:   make(chan *packet.Frame, 15000),
		cache:      cache.New(),
		tunHandler: tunHandler,
		udpHandler: udpHandler,
	}
}

func (t *Tun) ReadFromTun(ctx context.Context, networkId string) {
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
		destMac := util.GetMacAddr(pkt.Packet)
		sock, err := t.GetSock(destMac)
		if err != nil {
			logger.Errorf("get socket failed:%v", err)
		}
		sock.Write(pkt.Packet[:])
	}
}

func (t *Tun) ReadFromUdp() {
	for {
		ctx := context.Background()
		frame := packet.NewFrame()
		destMac := util.GetMacAddr(frame.Packet)
		sock, err := t.GetSock(destMac)
		if err != nil {
			logger.Errorf("can not get sock")
		}

		n, err := sock.Read(frame.Buff[:])
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

func (t *Tun) GetSock(mac string) (socket.Interface, error) {
	nodeInfo, err := t.cache.GetNodeInfo(mac)
	if err != nil {
		//走releay
		return nil, err
	}

	return nodeInfo.Socket, nil
}
