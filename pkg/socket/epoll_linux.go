package socket

import (
	"fmt"
	"net"
	"syscall"
)

func NewSocket(conn *net.UDPConn) (*Socket, error) {
	file, err := conn.File()

	if err != nil {
		return nil, err
	}
	return &Socket{FileDescriptor: int(file.Fd())}, nil
}

// Listen use linux epoll
func Listen(ip string, port int) (*Socket, error) {
	socket := &Socket{}

	socketFileDescriptor, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket (%v)", err)
	}

	if err = syscall.SetNonblock(socketFileDescriptor, true); err != nil {
		return nil, fmt.Errorf("set nonblock1: (%v)", err)
	}

	socket.FileDescriptor = socketFileDescriptor
	socketAddress := &syscall.SockaddrInet4{Port: port}
	copy(socketAddress.Addr[:], net.ParseIP(ip))

	if err = syscall.Bind(socket.FileDescriptor, socketAddress); err != nil {
		return nil, fmt.Errorf("failed to bind socket (%v)", err)
	}

	if err = syscall.Listen(socket.FileDescriptor, syscall.SOMAXCONN); err != nil {
		return nil, fmt.Errorf("failed to listen on socket (%v)", err)
	}

	return socket, nil

}
