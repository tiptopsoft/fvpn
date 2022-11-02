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

	p, _ := packet.Decode(data[:24])
	switch p.Flags {
	case option.TAP_REGISTER:
		if err := register(p); err != nil {
			fmt.Println(err)
		}
		_, err = conn.WriteToUDP(data[25:], addr)
		if err != nil {
			fmt.Println("super write failed.")
		}
		<-limitChan
		break
	}

}

// register star node register to super
func register(pack *packet.Packet) error {

	ips, err := packet.IntToBytes(int(pack.IPv4))
	if err != nil {
		return err
	}

	m.Store(pack.SourceMac, &net.UDPAddr{
		IP: net.IPv4(byte2int(ips[:8]), byte2int(ips[9:16]), byte2int(ips[17:24]), byte2int(ips[25:])), Port: int(pack.UdpPort),
	})

	return nil
}

func byte2int(b []byte) byte {
	a, _ := packet.BytesToInt(b)
	return byte(a)
}

// unRegister star node unregister from super
func unRegister() error {
	return nil
}
