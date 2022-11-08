package register

import (
	"context"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"net"
	"sync"
)

var limitChan = make(chan int, 1)

// mac:Pub
var m sync.Map

type Node struct {
	Mac   [4]byte
	Proto option.Protocol
	Conn  net.Conn
	Addr  *net.UDPAddr
}

//RegStar use as register
type RegStar struct {
	*option.RegConfig
	Handlers []handler.Handler
	conn     net.Conn
}

func (r *RegStar) Start(address string) error {
	return r.start(address)
}

func (r *RegStar) AddHandler(handler handler.Handler) {
	r.Handlers = append(r.Handlers, handler)
}

// Node register node for net, and for user create edge
func (r *RegStar) start(address string) error {
	var conn net.Conn
	switch r.Protocol {
	case option.UDP:
		addr, err := ResolveAddr(address)
		if err != nil {
			return err
		}

		conn, err = net.ListenUDP("udp", addr)

		log.Logger.Infof("registry start at: %s", address)

		if err != nil {
			return err
		}
		defer conn.Close()
		for {
			limitChan <- 1
			go r.handleUdp(conn.(*net.UDPConn))
		}
	default:
		log.Logger.Info("this is a tcp server")
	}

	return nil
}

func (r *RegStar) handleUdp(conn *net.UDPConn) {

	data := make([]byte, 2048)
	_, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println(err)
	}

	p, err := common.NewPacket().Decode(data[:24])
	if err != nil {
		fmt.Println(err)
	}
	switch p.Flags {
	case option.MSG_TYPE_REGISTER_SUPER:
		rpacket, err := register.NewPacket().Decode(data)
		if err := r.register(addr, rpacket); err != nil {
			fmt.Println(err)
		}

		// build a ack
		f, err := ackBuilder(p)
		if err != nil {
			fmt.Println("build resp p failed.")
		}
		copy(data[0:24], f)
		_, err = conn.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println("register write failed.")
		}
		<-limitChan
		break

	}

}

// register edge node register to register
func (r *RegStar) register(addr *net.UDPAddr, packet register.RegPacket) error {
	m.Store(packet.SrcMac, addr)
	m.Range(func(key, value any) bool {
		log.Logger.Infof("registry data key: %s, value: %v", key, value)
		return true
	})
	return nil
}

func (r *RegStar) unRegister(packet register.RegPacket) error {
	m.Delete(packet.SrcMac)
	return nil
}

func ackBuilder(cp common.CommonPacket) ([]byte, error) {

	RecMac := "01:01:03:02:03:01"
	ip := "192.168.1.1"
	Mask := "255.255.255.0"

	p := ack.NewPacket()
	copy(p.Mask[:], Mask)
	copy(p.AutoIP[:], ip)
	copy(p.RegMac[:], RecMac)

	cp.Flags = option.MSG_TYPE_REGISTER_ACK
	p.CommonPacket = cp

	return p.Encode(p)
}

func (r *RegStar) Execute(ctx context.Context, p packet.Packet) error {
	handlers := r.Handlers
	for _, h := range handlers {
		if err := h.Handle(ctx, p); err != nil {
			return err
		}
	}

	return nil
}

func ResolveAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}
