package device

import (
	"errors"
	"fmt"
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
func New(opts *option.StarConfig, mode Mode) (*Tuntap, error) {

	dev, err := net.InterfaceByName(opts.Name)
	if err != nil && err.Error() != NoSuchInterface.Error() {
		return nil, err
	}

	var f = "/dev/net/tun"
	file, err := os.OpenFile(f, os.O_RDWR, 0)
	if err != nil {
		panic(err)
		return nil, err
	}

	var ifr Ifreq
	copy(ifr.Name[:], opts.Name)

	var errno syscall.Errno
	switch mode {

	case TUN:
		ifr.Flags = syscall.IFF_TUN | syscall.IFF_NO_PI
		_, _, errno = unix.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))

	case TAP:
		ifr.Flags = syscall.IFF_TAP | syscall.IFF_NO_PI
		_, _, errno = unix.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))
	}
	if errno != 0 {
		return nil, fmt.Errorf("tuntap ioctl failed, errno %v", errno)
	}

	_, _, errno = unix.Syscall(unix.SYS_IOCTL, file.Fd(), uintptr(unix.TUNSETPERSIST), 1)
	if errno != 0 {
		return nil, fmt.Errorf("tuntap ioctl TUNSETPERSIST failed, errno %v", errno)
	}

	//set euid egid
	if _, _, errno = unix.Syscall(unix.SYS_IOCTL, file.Fd(), syscall.TUNSETGROUP, uintptr(os.Getegid())); errno < 0 {
		return nil, fmt.Errorf("tuntap set group error, errno %v", errno)
	}

	if _, _, errno = unix.Syscall(unix.SYS_IOCTL, file.Fd(), syscall.TUNSETOWNER, uintptr(os.Geteuid())); errno < 0 {
		return nil, fmt.Errorf("tuntap set group error, errno %v", errno)
	}

	// set tun tap up
	if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ip link set %s up", opts.Name)); err != nil {
		panic(err)
		return nil, err
	}

	if dev == nil {
		if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ip addr add %s dev %s", opts.IP, opts.Name)); err != nil {
			panic(err)
			return nil, err
		}
	}

	fmt.Println("Successfully connect to tun/tap interface:", opts.Name)

	return &Tuntap{
		file.Fd(),
		opts.Name,
		os.NewFile(file.Fd(), opts.Name),
		mode, "",
	}, nil
}

//Remove delete a tuntap
func Remove(opts *option.StarConfig) error {
	if err := option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ip link delete %s", opts.Name)); err != nil {
		return err
	}

	return nil
}
