// Copyright 2023 TiptopSoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tun

import (
	"errors"
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/log"
	"github.com/tiptopsoft/fvpn/pkg/util"
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
	logger = log.Log()
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

	tun := &NativeTun{
		file: os.NewFile(uintptr(fd), name),
		Fd:   0,
		name: name,
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

func (tun *NativeTun) Write(buff []byte) (int, error) {
	size := len(buff) + 4
	buf := make([]byte, size)
	copy(buf[4:], buff[:])
	buf[0] = 0x00
	buf[1] = 0x00
	buf[2] = 0x00
	switch buf[4] >> 4 {
	case 4:
		buf[3] = unix.AF_INET
	case 6:
		buf[3] = unix.AF_INET6
	default:
		return 0, unix.EAFNOSUPPORT
	}

	n, err := tun.file.Write(buf[:size])
	return n, err
}

func socketCloexec(family, sotype, proto int) (fd int, err error) {
	syscall.ForkLock.Lock()
	defer syscall.ForkLock.Unlock()

	fd, err = unix.Socket(family, sotype, proto)
	return
}

func (tun *NativeTun) JoinNetwork(network string) error {
	if tun.IP == nil {
		return errors.New("ip should set first")
	}
	return util.ExecCommand("/bin/sh", "-c", fmt.Sprintf("route add -net %s %s", network, tun.IP))
}

func (tun *NativeTun) SetIP(network, ip string) error {
	//set ip
	tun.IP = net.ParseIP(ip)
	return util.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s %s", tun.Name(), ip, ip))
}

func (tun *NativeTun) SetMTU(mtu int) error {
	return nil
}
