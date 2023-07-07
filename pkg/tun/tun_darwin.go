package tun

import (
	"errors"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"syscall"
)

const (
	utunControlName = "com.apple.net.utun_control"
	utunPrefix      = "utun"
)

var (
	logger      = log.Log()
	FakeGateway = "5.244.24.141/24"
	FakeIP      = net.ParseIP("5.244.24.141")
)

func New() (Device, error) {
	ifIndex := 0
	var name string
	var fd int
	var err error
	for {
		if ifIndex > 15 {
			return nil, errors.New("create utun device failed")
		}
		name = fmt.Sprintf("%s%d", utunPrefix, ifIndex)

		fd, err = socketCloexec(unix.AF_SYSTEM, unix.SOCK_DGRAM, 2)
		if err != nil {
			return nil, err
		}

		ctlInfo := &unix.CtlInfo{}
		copy(ctlInfo.Name[:], []byte(utunControlName))
		err = unix.IoctlCtlInfo(fd, ctlInfo)
		if err != nil {
			unix.Close(fd)
			return nil, fmt.Errorf("IoctlGetCtlInfo: %w", err)
		}

		sc := &unix.SockaddrCtl{
			ID:   ctlInfo.Id,
			Unit: uint32(ifIndex) + 1,
		}

		err = unix.Connect(fd, sc)
		if err != nil {
			unix.Close(fd)
			logger.Debugf("connect fd failed: %v, index: %d", err, sc.Unit)
			ifIndex++
			continue
		}

		err = unix.SetNonblock(fd, true)
		if err != nil {
			unix.Close(fd)
			logger.Debugf("set non block failed:%v", err)
			ifIndex++
			continue
		}

		break
	}

	//set ip
	if err = util.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s %s", name, FakeGateway, FakeIP.String())); err != nil {
		return nil, err
	}

	tun := &NativeTun{
		file:      os.NewFile(uintptr(fd), name),
		Fd:        0,
		name:      name,
		NetworkId: "",
		IP:        FakeIP,
	}

	logger.Debugf("create tun %s success", name)
	return tun, nil
}

// Read is a hack to work around the first 4 bytes "packet
// information" because there doesn't seem to be an IFF_NO_PI for darwin.
func (tun *NativeTun) Read(buff []byte) (n int, err error) {
	size := len(buff) + 4
	buf := make([]byte, size)
	n, err = tun.file.Read(buf)
	//
	if n <= 0 {
		return 0, err
	}
	copy(buff[:], buf[4:size])
	return n - 4, err
}

func socketCloexec(family, sotype, proto int) (fd int, err error) {
	syscall.ForkLock.Lock()
	defer syscall.ForkLock.Unlock()

	fd, err = unix.Socket(family, sotype, proto)
	return
}

func (tun *NativeTun) JoinNetwork(network string) error {
	return util.ExecCommand("/bin/sh", "-c", fmt.Sprintf("route add -net %s %s", network, FakeIP))
}
