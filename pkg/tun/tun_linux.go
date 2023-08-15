// Copyright 2023 Tiptopsoft, Inc.
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
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/log"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"golang.org/x/sys/unix"
	"net"
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

func New() (Device, error) {
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

	//endpoint, _ := util.New(offset)
	//logger.Debugf("ip: %v, mask: %v", endpoint.IP, endpoint.Mask)
	//if err = util.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s", name, fmt.Sprintf("%s/%d", endpoint.IP.String(), 24))); err != nil {
	//	return nil, err
	//}

	return &NativeTun{
		name: name, // size is 16
		file: os.NewFile(uintptr(fd), name),
		Fd:   fd,
		//IP:   endpoint.IP,
	}, nil
}

// Read is a hack to work around the first 4 bytes "packet
// information" because there doesn't seem to be an IFF_NO_PI for darwin.
func (tun *NativeTun) Read(buff []byte) (n int, err error) {
	n, err = tun.file.Read(buff)
	return n, err
}

func (tun *NativeTun) SetIP(network, ip string) error {
	tun.IP = net.ParseIP(ip)
	logger.Debugf("set ip network: %s, ip: %s", network, ip)
	//ex: ifconfig fvpn0 192.168.0.2 netmask 255.255.255.0
	if err := util.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s ", tun.Name(), ip, network)); err != nil {
		return err
	}
	return nil
}

func (tun *NativeTun) SetMTU(mtu int) error {
	return nil
}

func (tun *NativeTun) Write(buff []byte) (int, error) {
	tun.lock.Lock()
	defer tun.lock.Unlock()
	n, err := tun.file.Write(buff[:])
	return n, err
}
