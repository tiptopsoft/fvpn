package executor

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util/log"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"github.com/interstellar-cloud/star/pkg/util/packet"
	"github.com/interstellar-cloud/star/pkg/util/packet/forward"
	"github.com/interstellar-cloud/star/pkg/util/socket"
)

type TapExecutor struct {
	Name   string
	Socket socket.Socket
}

// Execute TapExecutor use to handle tap frame, write to udp sock.
// Read a single packet from the TAP interface, process it and write out the corresponding packet to the cooked socket.
func (te TapExecutor) Execute(socket socket.Socket) error {

	b := make([]byte, option.STAR_PKT_BUFF_SIZE)
	n, err := socket.Read(b)
	log.Logger.Info(fmt.Sprintf("Read from tap %s: length: %d", te.Name, len(b)))
	if err != nil {
		log.Logger.Errorf("tap read failed. (%v)", err)
		return err
	}

	mac := getMacAddr(b)
	log.Logger.Infof("Tap dev: %s receive: %d byte, mac: %v", te.Name, n, mac)
	// get dest
	_, ok := option.AddrMap.Load(mac)
	//dst := info.(ack.PeerInfo)
	if !ok {
		// through supernode
		fp := forward.NewPacket()
		bs, err := forward.Encode(fp)
		if err != nil {
			log.Logger.Errorf("encode forward failed. err: %v", err)
		}

		idx := 0
		packet.EncodeBytes(bs, b, idx)
		write2Net(te.Socket, bs)
	} else {
		// go p2p
	}
	return nil
}

//use host socket write so destination
func write2Net(socket socket.Socket, b []byte) {
	if _, err := socket.Write(b); err != nil {
		log.Logger.Errorf("write to remote failed. (%v)", err)
	}
}

func getMacAddr(buf []byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}
