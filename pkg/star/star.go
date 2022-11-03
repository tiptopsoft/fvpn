package star

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/device"
	"github.com/interstellar-cloud/star/pkg/option"
	"io"
	"net"
)

var (
	DefaultPort = 3000
)

type StarServer struct {
	Tun   *device.Tuntap
	Addr  *net.UDPAddr
	Type  int
	Serve bool
	Conn  net.Conn
}

func (s *StarServer) Start(port int) error {
	conn, err := s.listen()
	if err != nil {
		return nil
	}

	s.Conn = conn
	return nil
}

func (s *StarServer) listen() (net.Conn, error) {
	var conn net.Conn
	var err error
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		return nil, err
	}

	switch s.Type {
	case option.TCP:
		conn, err = listener.Accept()
	case option.UDP:
		conn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0),
			Port: DefaultPort})
	}

	//defer conn.Close()
	if err != nil {
		panic(err)
	}

	return conn, nil
}

func (s *StarServer) Dial(opts *option.StarConfig) (net.Conn, error) {
	if opts.Port == 0 {
		opts.Port = DefaultPort
	}
	address := fmt.Sprintf("%s:%d", opts.MoonIP, opts.Port)
	fmt.Println("connect to:", address)
	var conn net.Conn
	var err error
	switch s.Type {
	case option.TCP:
		conn, err = net.Dial("tcp", address)
	case option.UDP:
		ip := net.ParseIP(opts.MoonIP)
		conn, err = net.DialUDP("udp", nil, &net.UDPAddr{IP: ip,
			Port: DefaultPort})
	}

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (s *StarServer) Client(tap2net int, netfd io.ReadWriteCloser, tun *device.Tuntap) {

	for {
		var buf [2000]byte
		n, err := tun.Read(buf[:])
		if err == io.EOF {
			continue
		}
		if err != nil {
			panic(err)
		}

		tap2net++
		fmt.Println(fmt.Printf("tap2net:%d, tun received %d byte from %s: ", tap2net, n, tun.Name))

		/* write pack */
		if s.Type == option.UDP && s.Serve {
			fmt.Println("using write 2 udp")
			netfd.(*net.UDPConn).WriteToUDP(buf[:], s.Addr)
		} else {
			fmt.Println("using write udp")
			n, err = netfd.Write(buf[:n])
		}

		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Printf("tap2net:%d,write %d byte to network", tap2net, n))
	}
}

func (s *StarServer) Server(netfd io.ReadWriteCloser, tun *device.Tuntap) {

	for {
		buf := make([]byte, 2000)
		var n int
		var addr *net.UDPAddr
		var err error
		if s.Type == option.UDP {
			n, addr, err = netfd.(*net.UDPConn).ReadFromUDP(buf)
			s.Addr = addr
		} else {
			n, err = netfd.Read(buf)
		}

		if err == io.EOF {
			continue
		}

		fmt.Println(fmt.Printf("Recevied %d byte from net", n))
		//write to tap
		_, err = tun.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fmt.Printf("write %d byte to tap %s", n, tun.Name))
	}

}
