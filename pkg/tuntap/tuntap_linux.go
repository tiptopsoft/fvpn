package tuntap

import (
	"errors"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/option"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
	"unsafe"
)

func New(mode Mode, ip, mask, networkId string) error {
	name := fmt.Sprintf("%s%s", NamePrefix, networkId)
	var f = "/dev/net/tun"

	fd, err := unix.Open(f, os.O_RDWR, 0)
	if err != nil {
		panic(err)
		return err
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
		return fmt.Errorf("tuntap ioctl failed, errno %v", errno)
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

	//设置IP
	if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", name, ip, mask, 1420)); err != nil {
		return err
	}

	return err
}

func Delete(networkId string) error {

	return nil
}

func GetTuntap(networkId string) (*Tuntap, error) {
	name := fmt.Sprintf("%s%s", NamePrefix, networkId)
	var f = "/dev/net/tun"

	fd, err := unix.Open(f, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	var ifr Ifreq
	copy(ifr.Name[:], name)

	var errno syscall.Errno

	ifr.Flags = IFF_TAP | IFF_NO_PI
	_, _, errno = unix.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))

	if errno < 0 {
		return nil, errors.New("Get tuntap failed.")
	}

	mac, _ := addr.GetMacAddrByDev(name)
	return &Tuntap{
		Name:    name,
		MacAddr: mac,
		file:    os.NewFile(uintptr(fd), name),
		Fd:      fd,
	}, nil
}
