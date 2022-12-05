package main

import (
	"syscall"
)

//Socket use to wrap FileDescriptor
type Socket struct {
	FileDescriptor int
}

func (socket Socket) Read(bytes []byte) (n int, err error) {
	n, err = syscall.Read(socket.FileDescriptor, bytes)
	if err != nil {
		return 0, err
	}
	return
}

func (socket Socket) Write(bytes []byte) (n int, err error) {
	n, err = syscall.Write(socket.FileDescriptor, bytes)
	if err != nil {
		n = 0
	}
	return n, err
}

func (socket Socket) Close() error {
	return syscall.Close(socket.FileDescriptor)
}
