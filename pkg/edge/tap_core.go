package edge

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/forward"
	"github.com/interstellar-cloud/star/pkg/packet/peer/ack"
	"os"
)

// TapHandle  use to handle tap frame, write to udp sock.
// Read a single packet from the TAP interface, process it and write out the corresponding packet to the cooked socket.
func TapHandle(fd uintptr, name string) {

	b := make([]byte, option.STAR_PKT_BUFF_SIZE)
	file := os.NewFile(fd, name)

	n, err := file.Read(b)
	if err != nil {
		log.Logger.Errorf("dev: %s read tap byte failed. ", name)
	}
	log.Logger.Infof("Tap dev: %s receive: %d byte", name, n)

	mac := getMacAddr(b)

	// get dest
	info, ok := m.Load(mac)
	dst := info.(ack.PeerInfo)
	if ok {
		//check it is use supernode or p2p
		if dst.P2p == 1 {
			// p2p
		}

		if dst.P2p == 2 {
			// through supernode
			cp := common.NewPacket()
			cp.Flags = option.MSG_TYPE_PACKET

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

func getMacAddr(buf []byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}

func getDestEdge() {

}
