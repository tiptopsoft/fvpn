package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/addr"
	"github.com/interstellar-cloud/star/pkg/node"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/forward"
	"github.com/interstellar-cloud/star/pkg/socket"
	"github.com/interstellar-cloud/star/pkg/tuntap"
	"github.com/interstellar-cloud/star/pkg/util"
	"net"
)

type TapExecutor struct {
	device *tuntap.Tuntap
	cache  node.NodesCache
}

func (t TapExecutor) Execute(skt socket.Interface) error {
	device := t.device
	b := make([]byte, option.STAR_PKT_BUFF_SIZE)
	size, err := device.Read(b)
	destMac := util.GetMacAddr(b)
	fmt.Println(fmt.Sprintf("Read %d bytes from device %s, will write to dest %s", size, device.Name, destMac))
	if err != nil {
		logger.Errorf("tap read failed. (%v)", err)
		return err
	}
	broad := addr.IsBroadCast(destMac)
	//broad frame, go through supernode
	fp := forward.NewPacket()
	fp.SrcMac, err = addr.GetMacAddrByDev(device.Name)
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
	idx = packet.EncodeBytes(newPacket, b[:size], idx)
	if broad {
		write2Net(skt, newPacket[:idx])
	} else {
		// go p2p
		logger.Infof("find peer in edge, destMac: %v", destMac)
		p := node.FindNode(t.cache, destMac)
		if p == nil {
			write2Net(skt, newPacket[:idx])
			logger.Warnf("peer not found, go through super node")
		} else {
			write2Net(p.Socket, newPacket[:idx])
		}
	}
	return nil
}
