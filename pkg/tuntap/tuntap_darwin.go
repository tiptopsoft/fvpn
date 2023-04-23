package tuntap

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/topcloudz/fvpn/pkg/addr"
	"golang.org/x/sys/unix"
)

// New create a Tuntap
func New(mode Mode) (*Tuntap, error) {
	i := 0
	var name string
	var err error
	var fd int
	for {
		name = fmt.Sprintf("tap%d", i)
		var f = "/dev/net/tun"

		fd, err = unix.Open(f, os.O_RDWR, 0)
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
			_, _, errno = unix.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))

		case TAP:
			ifr.Flags = IFF_TAP | IFF_NO_PI
			_, _, errno = unix.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))
		}

		if errno != 0 {
			return nil, fmt.Errorf("tuntap ioctl failed, errno %v", errno)
		}

		_, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(TUNSETPERSIST), 1)
		if errno != 0 {
			err = fmt.Errorf("tuntap ioctl TUNSETPERSIST failed, errno %v", errno)
		}

		//set euid egid
		if _, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), TUNSETGROUP, uintptr(os.Getegid())); errno < 0 {
			err = fmt.Errorf("tuntap set group error, errno %v", errno)
		}

		if _, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), TUNSETOWNER, uintptr(os.Geteuid())); errno < 0 {
			err = fmt.Errorf("tuntap set group error, errno %v", errno)
		}

		if err != nil && i < 255 {
			i++
		} else {
			break
		}
	}

	fmt.Println("Successfully connect to tun/tap interface:", name)

	mac, _ := addr.GetMacAddrByDev(name)
	return &Tuntap{
		Name:    name,
		Mode:    mode,
		MacAddr: mac,
		file:    os.NewFile(uintptr(fd), name),
		Fd:      fd,
	}, nil
}
