package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
	"unsafe"
)

//
//const (
//	IFF_NO_PI = 0x1000
//	IFF_TUN   = 0x1
//	IFF_TAP   = 0x2
//)

type Ifreq struct {
	Name  [16]byte
	Flags uint16
}

// connect to a tuntap
func TunAlloc() Socket {
	f := "/dev/net/tun"
	name := "tun11"
	var fd uintptr
	var errno syscall.Errno
	if file, err := os.OpenFile(f, os.O_RDWR, 0); err != nil {
		panic(err)
	} else {
		fd = file.Fd()
	}

	var ifr Ifreq
	copy(ifr.Name[:], name)
	ifr.Flags = unix.IFF_TUN | unix.IFF_NO_PI
	_, _, errno = unix.Syscall(syscall.SYS_IOCTL, fd, uintptr(unix.TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))
	if errno != 0 {
		panic(errno)
	}

	return Socket{FileDescriptor: int(fd)}
}

func main() {

	s := TunAlloc()

	for {
		b := make([]byte, 2048)

		var FdSet unix.FdSet
		FdSet.Set(s.FileDescriptor)

		maxFd := s.FileDescriptor
		for {
			ret, err := unix.Select(maxFd+1, &FdSet, nil, nil, nil)
			if ret < 0 && err == unix.EINTR {
				continue
			}
			if err != nil {
				panic(err)
			}

			if FdSet.IsSet(s.FileDescriptor) {
				if _, err := s.Read(b); err != nil {
					panic(err)
				}
				fmt.Println(b)
			}

			fmt.Println("aaaaaaaa")
		}
	}
}
