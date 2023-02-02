package socket

import (
	"golang.org/x/sys/unix"
	"net"
	"reflect"
)

//Socket use to wrap FileDescriptor
type Socket struct {
	FileDescriptor int
	UdpSocket      *net.UDPConn
}

func (socket Socket) ReadFromUdp(bytes []byte) (n int, addr *net.UDPAddr, err error) {
	return socket.UdpSocket.ReadFromUDP(bytes)
}

func (socket Socket) Read(bytes []byte) (n int, err error) {
	n, err = unix.Read(socket.FileDescriptor, bytes)
	if err != nil {
		return 0, err
	}
	return
}

func (socket Socket) Write(bytes []byte) (n int, err error) {
	n, err = unix.Write(socket.FileDescriptor, bytes)
	if err != nil {
		n = 0
	}
	return n, err
}

func SocketFD(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	//if tls {
	//	tcpConn = reflect.Indirect(tcpConn.Elem())
	//}
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func (socket Socket) Close() error {
	return unix.Close(socket.FileDescriptor)
}

type Executor interface {
	Execute(socket Socket) error
}
