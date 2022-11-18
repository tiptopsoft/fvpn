package register

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"net"
)

func (r *RegStar) processRegister(addr *net.UDPAddr, conn *net.UDPConn, data []byte, cp *common.CommonPacket) {
	var regPacket register.RegPacket
	var err error
	if cp != nil {
		regPacket, err = register.DecodeWithCommonPacket(data, *cp)
	} else {
		regPacket, err = register.Decode(data)
	}

	if err := r.register(addr, regPacket); err != nil {
		fmt.Println(err)
	}
	// build a ack
	f, err := ackBuilder(regPacket.CommonPacket)
	log.Logger.Infof("build a register ack: %v", f)
	if err != nil {
		fmt.Println("build resp p failed.")
	}
	_, err = conn.WriteToUDP(f, addr)
	if err != nil {
		fmt.Println("register write failed.")
	}

	<-limitChan
}