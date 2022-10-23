package service

import (
	"errors"
	"fmt"
	"github.com/interstellar-cloud/star/device"
	"io"
	"net"
)

var (
	Unknown = errors.New("unknown")
)

type Server struct {
	Tun   *device.Tuntap
	Addr  *net.UDPAddr
	Udp   bool
	Serve bool
}

func (s *Server) ListenTcp() (net.Conn, error) {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		return nil, err
	}

	return s.handle(listener)
}

func (s *Server) ListenUdp() (*net.UDPConn, error) {
	return net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0),
		Port: DefaultPort})
}

func (s *Server) handle(listener net.Listener) (net.Conn, error) {

	conn, err := listener.Accept()
	//defer conn.Close()
	if err != nil {
		panic(err)
	}
	return conn, nil
}

func (s *Server) Client(tap2net int, netfd io.ReadWriteCloser, tun *device.Tuntap) {

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

		/* write packet */
		n, err = netfd.Write(buf[:n])
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Printf("tap2net:%d,write %d byte to network", tap2net, n))
	}
}

func (s *Server) UdpClient(tap2net int, netfd *net.UDPConn, tun *device.Tuntap) {

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

		/* write packet */
		if s.Udp && s.Serve {
			n, err = netfd.WriteToUDP(buf[:n], s.Addr)
		} else {
			n, err = netfd.Write(buf[:n])
		}

		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Printf("tap2net:%d,write %d byte to network", tap2net, n))
	}
}

func (s *Server) Server(netfd io.ReadWriteCloser, tun *device.Tuntap) {
	for {
		buf := make([]byte, 2000)
		n, err := netfd.Read(buf)
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

func (s *Server) UdpServer(netfd *net.UDPConn, tun *device.Tuntap) {
	for {
		buf := make([]byte, 2000)
		n, addr, err := netfd.ReadFromUDP(buf)
		s.Addr = addr
		s.Udp = true
		if err == io.EOF {
			continue
		}
		fmt.Println(fmt.Printf("Recevied %d byte from net, addr: %s", n, addr))
		//write to tap
		_, err = tun.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fmt.Printf("write %d byte to tap %s", n, tun.Name))
	}
}
