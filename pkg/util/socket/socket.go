package socket

import (
	"github.com/interstellar-cloud/star/pkg/util/option"
	"golang.org/x/sys/unix"
	"net"
	"reflect"
)

//Socket use to wrap FileDescriptor
type Socket struct {
	AppType        option.Protocol
	FileDescriptor int
	UdpSocket      *net.UDPConn
}

func (socket Socket) ReadFromUdp(bytes []byte) (n int, addr *net.UDPAddr, err error) {
	return socket.UdpSocket.ReadFromUDP(bytes)
}

func (socket Socket) WriteToUdp(bytes []byte, addr *net.UDPAddr) (n int, err error) {
	return socket.UdpSocket.WriteToUDP(bytes, addr)
}

func (socket Socket) Read(bytes []byte) (n int, err error) {
	if socket.AppType == option.UDP {
		return socket.UdpSocket.Read(bytes)
	}
	return unix.Read(socket.FileDescriptor, bytes)
}

func (socket Socket) Write(bytes []byte) (n int, err error) {
	if socket.AppType == option.UDP {
		return socket.UdpSocket.Write(bytes)
	} else {
		return unix.Write(socket.FileDescriptor, bytes)
	}
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
