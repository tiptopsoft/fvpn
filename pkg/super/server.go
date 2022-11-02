package super

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/internal"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet"
	"net"
	"sync"
)

var limitChan = make(chan int, 1000)

// udp key : mac_group value:addr
var m sync.Map

//RelayServer use as register
type RelayServer struct {
	Config   *option.Config
	Handlers []internal.StarHandler
}

func (s *RelayServer) Start(port int) error {
	return start(port)
}

func (s *RelayServer) AddHandler(handler internal.StarHandler) {
	s.Handlers = append(s.Handlers, handler)
}

// Node super node for net, and for user create star
func start(listen int) error {

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
		go handleUdp(conn)
	}

}

func handleUdp(conn *net.UDPConn) {

	data := make([]byte, 1024)
	_, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println(err)
	}

	frame, _ := packet.Decode(data[:24])
	switch frame.Flags {
	case option.TAP_REGISTER:
		if err := register(frame); err != nil {
			fmt.Println(err)
		}

		// build a ack
		f, err := ackBuilder(frame)
		if err != nil {
			fmt.Println("build resp frame failed.")
		}
		copy(data[0:24], f)
		_, err = conn.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println("super write failed.")
		}
		<-limitChan
		break
	case option.TAP_UNREGISTER:
		unRegister(frame)
		break

	case option.TAP_MESSAGE:
		break
	}

}

// register star node register to super
func register(pack *packet.Frame) error {

	ips := pack.IPv4

	m.Store(pack.SourceMac, &net.UDPAddr{
		IP: net.IPv4(ips[0], ips[1], ips[2], ips[3]), Port: int(pack.UdpPort),
	})

	return nil
}

func ackBuilder(orginPacket *packet.Frame) ([]byte, error) {
	p := packet.NewPacket()

	p.SourceMac = orginPacket.DestMac
	p.DestMac = orginPacket.SourceMac
	p.Flags = option.TAP_REGISTER_ACK

	return packet.Encode(p)
}

// unRegister star node unregister from super
func unRegister(pack *packet.Frame) {
	m.Delete(pack.SourceMac)
}
