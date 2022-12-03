package socket

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/forward"
	"github.com/interstellar-cloud/star/pkg/packet/peer/ack"
	"io"
)

type TapExecutor struct {
	Name string
}

// Execute TapExecutor  use to handle tap frame, write to udp sock.
// Read a single packet from the TAP interface, process it and write out the corresponding packet to the cooked socket.
func (te TapExecutor) Execute(socket Socket) error {

	for {
		b := make([]byte, option.STAR_PKT_BUFF_SIZE)
		n, err := socket.Read(b)
		if err != nil && err == io.EOF {
			continue
		} else {
			log.Logger.Errorf("dev: %s read tap byte failed. ", te.Name)
		}
		log.Logger.Infof("Tap dev: %s receive: %d byte", te.Name, n)

		mac := getMacAddr(b)

		// get dest
		info, ok := option.AddrMap.Load(mac)
		dst := info.(ack.PeerInfo)
		if ok {
			//check it is use supernode or p2p
			if dst.P2p == 1 {
				// p2p
			}

			if dst.P2p == 2 {
				// through supernode
				cp := common.NewPacket()
				cp.Flags = option.MsgTypePacket

				fp := forward.NewPacket()
				fp.CommonPacket = cp
				bs, err := forward.Encode(fp)
				if err != nil {
					log.Logger.Errorf("encode forward failed. err: %v", err)
				}

				idx := 0
				packet.EncodeBytes(b, bs, idx)

				packet.SendPacket(b, mac)
			}

		}
	}

}

func getMacAddr(buf []byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}
