package super

import (
	"context"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/internal"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/pack"
	"net"
)

var limitChan = make(chan int, 1000)

// udp key : mac_group value:addr
var m map[[4]byte]interface{}

type Node struct {
	Mac   [4]byte
	Proto internal.Protocol
	Conn  net.Conn
	Addr  *net.UDPAddr
}

//RegistryStar use as register
type RegistryStar struct {
	Config   *option.Config
	Handlers []internal.Handler
}

func (r *RegistryStar) Start(port int) error {
	return r.start(port)
}

func (r *RegistryStar) AddHandler(handler internal.Handler) {
	r.Handlers = append(r.Handlers, handler)
}

// Node super node for net, and for user create star
func (r *RegistryStar) start(listen int) error {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: listen,
	})

	if err != nil {
		return err
	}

	defer conn.Close()
	for {
		limitChan <- 1
		go r.handleUdp(conn)
	}

}

func (r *RegistryStar) handleUdp(conn *net.UDPConn) {

	data := make([]byte, 1024)
	_, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println(err)
	}

	p, _ := pack.Decode(data[:24])

	switch p.Flags {
	case option.TAP_REGISTER:
		if err := r.register(p); err != nil {
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
			fmt.Println("super write failed.")
		}
		<-limitChan
		break
	case option.TAP_UNREGISTER:
		unRegister(p)
		break

	case option.TAP_MESSAGE:
		addr, _ := m[p.SourceMac]
		if _, err := conn.WriteToUDP(data, addr.(*net.UDPAddr)); err != nil {
			fmt.Println(err)
		}
		break
	}

}

// register star node register to super
func (r *RegistryStar) register(p *pack.Packet) error {
	ips := p.IPv4
	m[p.SourceMac] = &Node{
		Addr: &net.UDPAddr{
			IP: net.IPv4(ips[0], ips[1], ips[2], ips[3]), Port: int(p.UdpPort),
		},
		Proto: r.Config.Proto,
	}

	return nil
}

func ackBuilder(orginPacket *pack.Packet) ([]byte, error) {
	p := pack.NewPacket()

	p.SourceMac = orginPacket.DestMac
	p.DestMac = orginPacket.SourceMac
	p.Flags = option.TAP_REGISTER_ACK

	return pack.Encode(p)
}

// unRegister star node unregister from super
func unRegister(pack *pack.Packet) {
	delete(m, pack.SourceMac)
}

func (r *RegistryStar) Execute(ctx context.Context, p pack.Packet) error {
	handlers := r.Handlers
	for _, h := range handlers {
		if err := h.Handle(ctx, p); err != nil {
			return err
		}
	}

	return nil
}
