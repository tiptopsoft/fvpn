package service

import (
	"errors"
	"github.com/interstellar-cloud/star/device"
	"net"
)

var (
	Unknown = errors.New("unknown")
)

type Server struct {
	Tun *device.Tuntap
}

func (s *Server) Listen() (net.Conn, error) {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		return nil, err
	}

	return s.handle(listener)

}

func (s *Server) handle(listener net.Listener) (net.Conn, error) {

	conn, err := listener.Accept()
	//defer conn.Close()
	if err != nil {
		panic(err)
	}

	//go func() {
	//	for {
	//		buf := make([]byte, 2000)
	//		n, err := conn.Read(buf)
	//		if err != nil {
	//			panic(err)
	//		}
	//		fmt.Println(fmt.Printf("Recevied %d byte from net", n))
	//		//write to tap
	//		_, err = s.Tun.Write(buf[:n])
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//		fmt.Println(fmt.Printf("write %d byte to tap %s", n, s.Tun.Name))
	//	}
	//}()

	return conn, nil
}
