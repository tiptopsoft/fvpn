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

	_, err = conn.WriteToUDP(data, addr)
	if err != nil {
		fmt.Println("super write failed.")
	}

	<-limitChan
}

// register star node register to super
func register(p []byte) error {
	pack, err := packet.Decode(p)
	if err != nil {
		return err
	}
	ips, err := packet.IntToBytes(int(pack.IPv4))
	if err != nil {
		return err
	}

	m.Store(pack.SourceMac, &net.UDPAddr{

		IP: net.IPv4(ips[:8], byte(pack.IPv4[9:16]), byte(pack.IPv4[17:24]), byte(pack.IPv4[25:])), Port: int(pack.UdpPort),
	})

	return nil
}

// unRegister star node unregister from super
func unRegister() error {
	return nil
}
