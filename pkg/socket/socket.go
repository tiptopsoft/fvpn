package socket

import (
	"golang.org/x/sys/unix"
	"net"
)

// Socket use to wrap fd
type Socket struct {
	Fd int
}

func (socket Socket) ReadFromUdp(bytes []byte) (n int, addr unix.Sockaddr, err error) {
	return unix.Recvfrom(socket.Fd, bytes, 0)
}

func (socket Socket) WriteToUdp(bytes []byte, addr unix.Sockaddr) (err error) {
	return unix.Sendto(socket.Fd, bytes, 0, addr)
}

func (socket Socket) Read(bytes []byte) (n int, err error) {
	return unix.Read(socket.Fd, bytes)
}

func (socket Socket) Write(bytes []byte) (n int, err error) {
	return unix.Write(socket.Fd, bytes)
}

func (socket Socket) Close() error {
	return unix.Close(socket.Fd)
}

func NewSocket(port int) Interface {
	fd, _ := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
	unix.SetsockoptInt(fd, unix.IP_DONTFRAG, unix.IPPROTO_IP, 0) //调协允许分包
	if port != 0 {
		addr := unix.SockaddrInet4{Port: port}
		copy(addr.Addr[:], net.IPv4zero.To4())
		unix.Bind(fd, &addr)
	}
	return Socket{Fd: fd}
}

func (socket Socket) Connect(addr unix.Sockaddr) error {

	return unix.Connect(socket.Fd, addr)
}

func (socket Socket) Listen(addr unix.Sockaddr) error {
	return unix.Bind(socket.Fd, addr)
}
