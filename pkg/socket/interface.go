package socket

import "golang.org/x/sys/unix"

type Interface interface {
	Listen(addr unix.Sockaddr) error
	Connect(addr unix.Sockaddr) error
	Read(bytes []byte) (n int, err error)
	Write(bytes []byte) (n int, err error)
	WriteToUdp(bytes []byte, addr unix.Sockaddr) (err error)
	ReadFromUdp(bytes []byte) (n int, addr unix.Sockaddr, err error)
	Close() error
}
