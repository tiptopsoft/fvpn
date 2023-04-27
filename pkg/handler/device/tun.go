package device

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/forward"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"net"
)

var (
	logger = log.Log()
)

func Handle() handler.HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {
		networkId := ctx.Value("networkId").(string)
		tun := ctx.Value("tun").(*tuntap.Tuntap)
		destMac := ctx.Value("mac").(string)
		var err error
		//destMac := util.GetMacAddr(frame.Buff)
		//fmt.Println(fmt.Sprintf("Read %d bytes from device %s, will write to dest %s", n, tun.Name, destMac))
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
		idx = packet.EncodeBytes(newPacket, frame.Buff[:frame.Size], idx)

		frame.Packet = newPacket[:]

		return nil

	}
}
