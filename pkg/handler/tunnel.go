package handler

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
)

type Tun struct {
	socket     socket.Interface
	device     *tuntap.Tuntap
	Inbound    chan *packet.Frame //used from udp
	Outbound   chan *packet.Frame //used for tun
	Cache      cache.PeersCache
	tunHandler Handler
	udpHandler Handler
}

func NewTun(tunHandler, udpHandler Handler) *Tun {
	return &Tun{
		Inbound:    make(chan *packet.Frame, 15000),
		Outbound:   make(chan *packet.Frame, 15000),
		Cache:      cache.New(),
		tunHandler: tunHandler,
		udpHandler: udpHandler,
	}
}

func (t Tun) ReadFromTun(ctx context.Context, networkId string) {
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

func (t Tun) WriteToUdp() {
	for {
		pkt := <-t.Outbound
		t.socket.Write(pkt.Packet[:])
	}
}

func (t Tun) ReadFromUdp() {
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
func (t Tun) WriteToDevice() {
	for {
		pkt := <-t.Inbound
		device, err := tuntap.GetTuntap(pkt.NetworkId)
		if err != nil {
			logger.Errorf("invalid network: %s", pkt.NetworkId)
		}
		device.Write(pkt.Packet[:])
	}
}
