package device

import (
	"errors"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"golang.org/x/sys/unix"
	"io"
	"net"
	"os"
	"syscall"
	"unsafe"
)

// Tuntap a device for net
type Tuntap struct {
	Fd   uintptr
	Name string
	io.ReadWriteCloser
	Mode    Mode
	MacAddr string
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

var (
	NoSuchInterface = errors.New("route ip+net: no such network interface")
)

// New craete a tuntap
func New(mode Mode) (*Tuntap, error) {
	i := 0
	var name string
	var err error
	var file *os.File
	for {
		name = fmt.Sprintf("tap%d", i)
		_, err = net.InterfaceByName(name)
		if err != nil && err.Error() == NoSuchInterface.Error() {
			//build
			var f = "/dev/tun"
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
				return nil, fmt.Errorf("tuntap ioctl TUNSETPERSIST failed, errno %v", errno)
			}

			//set euid egid
			if _, _, errno = unix.Syscall(unix.SYS_IOCTL, file.Fd(), TUNSETGROUP, uintptr(os.Getegid())); errno < 0 {
				return nil, fmt.Errorf("tuntap set group error, errno %v", errno)
			}

			if _, _, errno = unix.Syscall(unix.SYS_IOCTL, file.Fd(), TUNSETOWNER, uintptr(os.Geteuid())); errno < 0 {
				return nil, fmt.Errorf("tuntap set group error, errno %v", errno)
			}

			// set tun tap up
			if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ip link set %s up", name)); err != nil {
				panic(err)
				return nil, err
			}
			break
		}

		i++
		if i < 255 {
			continue
		} else {
			return nil, errors.New("unable to create a tuntap device")
		}
	}

	//if dev == nil {
	//	if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ip addr add %s dev %s", opts.IP, opts.Name)); err != nil {
	//		panic(err)
	//		return nil, err
	//	}
	//}

	log.Logger.Infof("Successfully connect to tun/tap interface: %v", name)

	return &Tuntap{
		file.Fd(),
		name,
		os.NewFile(file.Fd(), name),
		mode, "",
	}, nil
}

//Remove delete a tuntap
//func Remove(opts *option.StarConfig) error {
//	if err := option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ip link delete %s", opts.Name)); err != nil {
//		return err
//	}
//
//	return nil
//}
