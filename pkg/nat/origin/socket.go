package origin

import (
	"github.com/gogf/greuse"
	"net"
)

// Socket use to wrap fd
type Socket struct {
	conn *net.UDPConn
}

func NewSocket(conn *net.UDPConn) Socket {
	return Socket{conn: conn}
}

func (socket Socket) ReadFromUdp(buff []byte) (n int, addr *net.UDPAddr, err error) {
	return socket.conn.ReadFromUDP(buff)
}

func (socket Socket) WriteToUdp(buff []byte, addr *net.UDPAddr) (int, error) {
	return socket.conn.WriteToUDP(buff, addr)
}

func (socket Socket) Read(buff []byte) (n int, err error) {
	return socket.conn.Read(buff)
}

func (socket Socket) Write(buf []byte) (n int, err error) {
	return socket.conn.Write(buf)
}

func (socket Socket) Close() error {
	return socket.conn.Close()
}

//func NewSocket() Socket {
//	fd, _ := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
//	unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
//	unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
//	//unix.Bind(fd, &unix.SockaddrInet4{
//	//	Port: 4000,
//	//	SourceIP: [4]byte{0, 0, 0, 0},
//	//})
//
//	addr := unix.SockaddrInet4{Port: 4000}
//	copy(addr.SourceIP[:], net.IPv4zero.To4())
//	unix.Bind(fd, &addr)
//	return Socket{Fd: fd}
//}

func Dial(network, address string) (net.Conn, error) {
	return greuse.Dial(network, ":6061", address)
}
