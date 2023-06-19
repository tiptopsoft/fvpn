package tuntap

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/log"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
	"unsafe"
)

var (
	logger = log.Log()
)

func New(mode Mode, networkId string) (*Tuntap, error) {
	name := fmt.Sprintf("%s%s", NamePrefix, networkId[:10])
	var f = "/dev/net/tun"

	fd, err := unix.Open(f, os.O_RDWR, 0)
	if err != nil {
		panic(err)
		return nil, err
	}

	logger.Infof("tun name: %s, networkId: %s", name, networkId)
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

	////设置IP
	//if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", name, ip, mask, 1420)); err != nil {
	//	return nil, err
	//}

	mac, _, _ := addr.GetMacAddrAndIPByDev(name)
	return &Tuntap{
		Name:      name, // size is 16
		MacAddr:   mac,
		file:      os.NewFile(uintptr(fd), name),
		Fd:        fd,
		NetworkId: networkId,
		//IP:        net.ParseIP(ip),
	}, nil
}

func Delete(networkId string) error {

	return nil
}

// GetTuntap can be used only one time when created.
//func GetTuntap(networkId string) (*Tuntap, error) {
//	name := fmt.Sprintf("%s%s", NamePrefix, networkId[:10])
//	var f = "/dev/net/tun"
//
//	fd, err := unix.Open(f, os.O_RDWR, 0)
//	if err != nil {
//		return nil, err
//	}
//
//	var ifr Ifreq
//	copy(ifr.Name[:], name)
//
//	var errno syscall.Errno
//
//	ifr.Flags = IFF_TAP | IFF_NO_PI
//	_, _, errno = unix.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))
//	fmt.Println("errno: ", errno)
//
//	if errno < 0 {
//		return nil, errors.New("get tun failed")
//	}
//
//	mac, ip, _ := addr.GetMacAddrAndIPByDev(name)
//	return &Tuntap{
//		Name:      name, // size is 16
//		MacAddr:   mac,
//		file:      os.NewFile(uintptr(fd), name),
//		Fd:        fd,
//		NetworkId: networkId,
//		IP:        ip,
//	}, nil
//}
