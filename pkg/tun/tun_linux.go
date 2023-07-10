package tun

import (
	"fmt"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
	"unsafe"
)

var (
	logger = log.Log()
)

type Ifreq struct {
	Name  [16]byte
	Flags uint16
}

func New(offset int32) (*NativeTun, error) {
	name := fmt.Sprintf("%s%d", DefaultNamePrefix, 0)
	var f = "/dev/net/tun"

	fd, err := unix.Open(f, os.O_RDWR, 0)
	if err != nil {
		panic(err)
		return nil, err
	}

	logger.Infof("tun name: %s", name)
	var ifr Ifreq
	copy(ifr.Name[:], name)

	var errno syscall.Errno
	ifr.Flags = IFF_TUN | IFF_NO_PI
	_, _, errno = unix.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))

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

	endpoint, _ := util.New(offset)
	logger.Debugf("ip: %v, mask: %v", endpoint.IP, endpoint.Mask)
	if err = util.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s", name, fmt.Sprintf("%s/%d", endpoint.IP.String(), 24))); err != nil {
		return nil, err
	}

	return &NativeTun{
		name: name, // size is 16
		file: os.NewFile(uintptr(fd), name),
		Fd:   fd,
		IP:   endpoint.IP,
	}, nil
}

// Read is a hack to work around the first 4 bytes "packet
// information" because there doesn't seem to be an IFF_NO_PI for darwin.
func (tun *NativeTun) Read(buff []byte) (n int, err error) {
	n, err = tun.file.Read(buff)
	return n, err
}

func (tun *NativeTun) Write(buff []byte) (int, error) {
	tun.lock.Lock()
	defer tun.lock.Unlock()
	n, err := tun.file.Write(buff[:])
	return n, err
}
