package client

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/middleware"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/forward"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
)

type Tun struct {
	socket   socket.Interface
	device   *tuntap.Tuntap
	inbound  chan *Frame //used from udp
	outbound chan *Frame //used for tun
	cache    cache.PeersCache
}

func NewTun() *Tun {
	return &Tun{
		inbound:  make(chan *Frame, 15000),
		outbound: make(chan *Frame, 15000),
		cache:    cache.New(),
	}
}

func (t Tun) ReadFromTun(networkId string, middleware middleware.Middleware) {

	for {
		tun, err := tuntap.GetTuntap(networkId)
		if err != nil {
			logger.Fatalf("invalid network: %s", networkId)
		}
		frame := NewFrame()
		n, err := tun.Read(frame.buff[:])

		if n < 0 || err != nil {
			continue
		}

		destMac := util.GetMacAddr(frame.buff)
		fmt.Println(fmt.Sprintf("Read %d bytes from device %s, will write to dest %s", n, tun.Name, destMac))
		//broad frame, go through supernode
		fp := forward.NewPacket(networkId)
		fp.SrcMac, err = addr.GetMacAddrByDev(tun.Name)

		if err != nil {
			logger.Errorf("get src mac failed, err: %v", err)
		}
		fp.DstMac, err = net.ParseMAC(destMac)
		if err != nil {
			logger.Errorf("get src mac failed, err: %v", err)
		}

		bs, err := fp.Encode()
		if err != nil {
			logger.Errorf("encode forward failed, err: %v", err)
		}

		idx := 0
		newPacket := make([]byte, 2048)
		idx = packet.EncodeBytes(newPacket, bs, idx)
		idx = packet.EncodeBytes(newPacket, frame.buff[:], idx)

		frame.packet = newPacket[:]
		t.outbound <- frame

		//if broad {
		//	tun.write2Net(newPacket[:idx])
		//} else {
		//	// go p2p
		//	logger.Infof("find peer in client, destMac: %v", destMac)
		//	p := cache.FindPeer(dh.cache, destMac)
		//	if p == nil {
		//		dh.write2Net(newPacket[:idx])
		//		logger.Warnf("peer not found, go through super cache")
		//	} else {
		//		dh.write2Net(newPacket[:idx])
		//	}
		//}
		//return nil
	}

}

func (t Tun) WriteToUdp() {
	for {
		pkt := <-t.outbound
		t.socket.Write(pkt.packet[:])
	}
}
