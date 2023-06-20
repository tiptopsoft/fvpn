package socket

import (
	reuse "github.com/libp2p/go-reuseport"
	"github.com/topcloudz/fvpn/pkg/log"
	"net"
)

var (
	logger = log.Log()
)

// Socket use to wrap fd
type Socket struct {
	conn *net.UDPConn
}

func (socket *Socket) ReadFromUDP(bytes []byte) (n int, addr *net.UDPAddr, err error) {
	return socket.conn.ReadFromUDP(bytes)
}

func (socket *Socket) WriteToUdp(bytes []byte, addr *net.UDPAddr) (n int, err error) {
	return socket.conn.WriteTo(bytes, addr)
}

func (socket *Socket) Read(bytes []byte) (n int, err error) {
	return socket.conn.Read(bytes)
}

func (socket *Socket) Write(bytes []byte) (n int, err error) {
	return socket.conn.Write(bytes)
}

func (socket *Socket) Close() error {
	return socket.conn.Close()
}

func (socket Socket) LocalAddr() *net.UDPAddr {
	addr := socket.conn.LocalAddr()
	return addr.(*net.UDPAddr)
}

func NewSocket(laddr, addr string) (*Socket, error) {
	var conn *net.UDPConn
	var err error
	if laddr == "" {
		dest, _ := net.ResolveUDPAddr("udp", addr)
		conn, err = net.DialUDP("udp", nil, dest)
		if err != nil {
			return nil, err
		}
		return &Socket{conn: conn}, nil
	}
	con, err := reuse.Dial("udp", laddr, addr)
	if err != nil {
		return nil, err
	}
	return &Socket{conn: con.(*net.UDPConn)}, nil
}
