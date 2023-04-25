package device

import (
	"context"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/forward"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
)

var (
	logger = log.Log()
)

func Handle(ctx context.Context, frame *packet.Frame) error {
	networkId := ctx.Value("networkId").(string)
	tun, err := tuntap.GetTuntap(networkId)
	if err != nil {
		logger.Fatalf("invalid network: %s", networkId)
	}
	n, err := tun.Read(frame.Buff[:])

	if err != nil {
		return err
	}

	destMac := util.GetMacAddr(frame.Buff)
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
	idx = packet.EncodeBytes(newPacket, frame.Buff[:], idx)

	frame.Packet = newPacket[:]

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
	return nil
}
