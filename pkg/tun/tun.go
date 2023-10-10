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
	"net"
	"os"
	"sync"
)

type Mode int

const (
	SYS_IOCTL     = 29
	TUNSETIFF     = 0x400454ca
	TUNSETPERSIST = 0x400454cb
	TUNSETGROUP   = 0x400454ce
	TUNSETOWNER   = 0x400454cc

	IFF_NO_PI      = 0x1000
	IFF_TUN        = 0x1
	IFF_TAP        = 0x2
	TUN       Mode = iota
	TAP
)

var DefaultNamePrefix = "fvpn"

type Device interface {
	Name() string
	Read(buff []byte, offset int) (int, error)
	Write(buff []byte, offset int) (int, error)
	SetIP(net, ip string) error
	SetMTU(mtu int) error
	IPToString() string
	Addr() net.IP
}

// NativeTun a tuntap for net
type NativeTun struct {
	lock      sync.Mutex
	file      *os.File
	Fd        int
	name      string
	NetworkId string
	IP        net.IP
}

func (tun *NativeTun) Name() string {
	return tun.name
}

func (tun *NativeTun) IPToString() string {
	return tun.IP.String()
}

func (tun *NativeTun) Addr() net.IP {
	return tun.IP
}

// Close this method close the device
func (tun *NativeTun) Close() error {
	return tun.file.Close()
}
