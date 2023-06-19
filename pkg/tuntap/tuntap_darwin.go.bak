package tuntap

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/option"
	"net"
	"os"
	"syscall"
	"unsafe"

	"github.com/topcloudz/fvpn/pkg/addr"
	"golang.org/x/sys/unix"
)

const appleCTLIOCGINFO = (0x40000000 | 0x80000000) | ((100 & 0x1fff) << 16) | uint32(byte('N'))<<8 | 3

// New create a Tuntap
func New(mode Mode, ip, mask, networkId string) (*Tuntap, error) {
	i := 0
	var name string
	var err error
	var fd int
	var socketFD int
	for {
		name = fmt.Sprintf("tap%d", i)
		f := fmt.Sprintf("/dev/%s", name)
		fd, err = unix.Open(f, os.O_RDWR, 0)
		if err != nil {
			panic(err)
			return nil, err
		}

		var ifr = &struct {
			ifName    [16]byte
			ifruFlags int16
			pad       [16]byte
		}{}
		copy(ifr.ifName[:], name)

		var errno syscall.Errno
		if socketFD, err = syscall.Socket(syscall.AF_SYSTEM, syscall.SOCK_DGRAM, 2); err != nil {
			return nil, fmt.Errorf("error in syscall.Socket: %v", err)
		}
		switch mode {

		case TUN:
			ifr.ifruFlags = IFF_TUN | IFF_NO_PI
			_, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(socketFD), uintptr(TUNSETIFF), uintptr(unsafe.Pointer(ifr)))

		case TAP:
			ifr.ifruFlags = IFF_TAP | IFF_NO_PI
			_, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(socketFD), uintptr(unix.SIOCGIFFLAGS), uintptr(unsafe.Pointer(ifr)))
		}

		if errno != 0 {
			return nil, fmt.Errorf("tuntap ioctl failed, errno %v", errno)
		}

		ifr.ifruFlags |= unix.IFF_RUNNING | unix.IFF_UP
		if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(socketFD), uintptr(syscall.SIOCSIFFLAGS), uintptr(unsafe.Pointer(ifr))); errno != 0 {
			err = errno
			return nil, fmt.Errorf("error in syscall.Syscall(syscall.SYS_IOCTL, ...): %v", err)
		}

		//_, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(socketFD), uintptr(TUNSETPERSIST), 1)
		//if errno != 0 {
		//	err = fmt.Errorf("tuntap ioctl TUNSETPERSIST failed, errno %v", errno)
		//	return nil, err
		//}

		//set euid egid
		//if _, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), TUNSETGROUP, uintptr(os.Getegid())); errno < 0 {
		//	err = fmt.Errorf("tuntap set group error, errno %v", errno)
		//}
		//
		//if _, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), TUNSETOWNER, uintptr(os.Geteuid())); errno < 0 {
		//	err = fmt.Errorf("tuntap set group error, errno %v", errno)
		//}
		//设置IP
		if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", name, ip, mask, 1420)); err != nil {
			return nil, err
		}

		if err != nil && i < 255 {
			i++
		} else {
			break
		}
	}

	fmt.Println("Successfully connect to tun/tap interface:", name)

	mac, _, _ := addr.GetMacAddrAndIPByDev(name)
	return &Tuntap{
		Name:    name,
		Mode:    mode,
		MacAddr: mac,
		file:    os.NewFile(uintptr(fd), name),
		Fd:      fd,
		IP:      net.ParseIP(ip),
	}, nil
}
