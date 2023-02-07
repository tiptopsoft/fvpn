package tuntap

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util"
	"github.com/interstellar-cloud/star/pkg/util/socket"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"syscall"
	"unsafe"
)

// Tuntap a tuntap for net
type Tuntap struct {
	Fd      uintptr
	Name    string
	Socket  socket.Socket
	Mode    Mode
	MacAddr net.HardwareAddr
}

type Mode int

const (
	TUN Mode = iota
	TAP
)

type Ifreq struct {
	Name  [16]byte
	Flags uint16
}

// New craete a tuntap
func New(mode Mode) (*Tuntap, error) {
	i := 0
	var name string
	var err error
	var file *os.File
	for {
		name = fmt.Sprintf("tap%d", i)
		var f = "/dev/net/tun"

		file, err = os.OpenFile(f, os.O_RDWR, 0)
		if err != nil {
			panic(err)
			return nil, err
		}

		var ifr Ifreq
		copy(ifr.Name[:], name)

		var errno syscall.Errno
		switch mode {

		case TUN:
			ifr.Flags = IFF_TUN | IFF_NO_PI
			_, _, errno = unix.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))

		case TAP:
			ifr.Flags = IFF_TAP | IFF_NO_PI
			_, _, errno = unix.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))
		}

		if errno != 0 {
			return nil, fmt.Errorf("tuntap ioctl failed, errno %v", errno)
		}

		_, _, errno = unix.Syscall(unix.SYS_IOCTL, file.Fd(), uintptr(TUNSETPERSIST), 1)
		if errno != 0 {
			err = fmt.Errorf("tuntap ioctl TUNSETPERSIST failed, errno %v", errno)
		}

		//set euid egid
		if _, _, errno = unix.Syscall(unix.SYS_IOCTL, file.Fd(), TUNSETGROUP, uintptr(os.Getegid())); errno < 0 {
			err = fmt.Errorf("tuntap set group error, errno %v", errno)
		}

		if _, _, errno = unix.Syscall(unix.SYS_IOCTL, file.Fd(), TUNSETOWNER, uintptr(os.Geteuid())); errno < 0 {
			err = fmt.Errorf("tuntap set group error, errno %v", errno)
		}

		if err != nil && i < 255 {
			i++
		} else {
			break
		}

	}

	fmt.Println("Successfully connect to tun/tap interface:", name)

	mac, _ := util.GetMacAddrByDev(name)
	return &Tuntap{
		Fd:      file.Fd(),
		Name:    name,
		Socket:  socket.Socket{FileDescriptor: int(file.Fd())},
		Mode:    mode,
		MacAddr: mac,
	}, nil
}
