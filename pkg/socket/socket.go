package socket

import (
	"golang.org/x/sys/unix"
)

//Socket use to wrap FileDescriptor
type Socket struct {
	FileDescriptor int
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

func (socket Socket) Close() error {
	return unix.Close(socket.FileDescriptor)
}

type Executor interface {
	Execute(socket Socket) error
}
