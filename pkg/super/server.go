package super

import (
	"context"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/internal"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"net"
	"sync"
)

var limitChan = make(chan int, 1000)

// mac:Pub
var m sync.Map

type Node struct {
	Mac   [4]byte
	Proto internal.Protocol
	Conn  net.Conn
	Addr  *net.UDPAddr
}

//RegStar use as register
type RegStar struct {
	*RegConfig
	Handlers []internal.Handler
	conn     net.Conn
}

func (r *RegStar) Start(port int) error {
	return r.start(port)
}

func (r *RegStar) AddHandler(handler internal.Handler) {
	r.Handlers = append(r.Handlers, handler)
}

// Node super node for net, and for user create edge
func (r *RegStar) start(listen int) error {

	var err error
	var conn net.Conn
	switch r.Protocol {
	case internal.UDP:
		conn, err = net.ListenUDP("udp", &net.UDPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: listen,
		})

		if err != nil {
			return err
		}
		defer conn.Close()
		for {
			limitChan <- 1
			go r.handleUdp(conn.(*net.UDPConn))
		}
	default:
		fmt.Println("this is a tcp server")
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
	case option.MSG_TYPE_RE_REGISTER_SUPER:
		rpacket, err := register.NewPacket().Decode(data)
		if err := r.register(addr, rpacket); err != nil {
			fmt.Println(err)
		}

		// build a ack
		f, err := ackBuilder()
		if err != nil {
			fmt.Println("build resp p failed.")
		}
		copy(data[0:24], f)
		_, err = conn.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println("super write failed.")
		}
		<-limitChan
		break

	}

}

// register edge node register to super
func (r *RegStar) register(addr *net.UDPAddr, packet *register.RegPacket) error {
	m.Store(packet.SrcMac, addr)
	return nil
}

func (r *RegStar) unRegister(packet *register.RegPacket) error {
	m.Delete(packet.SrcMac)
	return nil
}

func ackBuilder() ([]byte, error) {
	return nil, nil
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
