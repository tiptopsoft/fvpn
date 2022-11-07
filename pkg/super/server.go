package super

import (
	"context"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/internal"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/pack"
	"github.com/interstellar-cloud/star/pkg/pack/common"
	"github.com/interstellar-cloud/star/pkg/pack/register"
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

// Node super node for net, and for user create edge
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

	data := make([]byte, 2048)
	_, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println(err)
	}

	p := &common.CommonPacket{}
	p, err = p.Decode(data[:24])
	if err != nil {
		fmt.Println(err)
	}

	rp := &register.RegPacket{}
	rp.Decode(data[25:])

	switch p.Flags {
	case common.TAP_REGISTER:
		if err := r.register(rp); err != nil {
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
	case common.TAP_UNREGISTER:
		//unRegister(p)
		break

	case common.TAP_MESSAGE:
		//addr, _ := m[p.]
		//if _, err := conn.WriteToUDP(data, addr.(*net.UDPAddr)); err != nil {
		//	fmt.Println(err)
		//}
		break
	}

}

// register edge node register to super
func (r *RegistryStar) register(p *register.RegPacket) error {
	//ips := p.IPv4
	//m[p.SrcMac] = &Node{
	//	Addr: &net.UDPAddr{
	//		IP: net.IPv4(ips[0], ips[1], ips[2], ips[3]), Port: int(p.UdpPort),
	//	},
	//	Proto: r.Config.Proto,
	//}

	return nil
}

func ackBuilder() ([]byte, error) {
	return nil, nil
}

// unRegister edge node unregister from super
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
