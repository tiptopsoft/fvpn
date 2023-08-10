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

package node

import (
	"fmt"
	"net"
	"strings"
)

// Bind a Bind listens a port for IPV6 and IPv4 UDP packets. Also send packets to destination peer.
type Bind interface {
	// Open listen a port using a given port, if not successec, a random port will return, which ia actualPort
	Open(port uint16) (actualPort uint16, conn *net.UDPConn, err error)

	Send(buff []byte, ep Endpoint) (int, error)

	Conn() *net.UDPConn
	//BatchSize is size use to receive or send
	BatchSize() int
}

type Endpoint interface {
	SrcToString() string
	DstToString() string
	SrcIP() net.IP
	DstIP() *net.UDPAddr
	//SetSrcIP(ip net.IP)
}

type endpoint struct {
	srcIP net.IP
	dstIP *net.UDPAddr
}

func NewEndpoint(dstip string) Endpoint {
	if !strings.Contains(dstip, ":") {
		dstip = fmt.Sprintf("%s:%d", dstip, 4000)
	}
	addr, err := net.ResolveUDPAddr("udp4", dstip)
	if err != nil {
		return nil
	}

	destIP := &net.UDPAddr{
		IP:   net.ParseIP(addr.IP.To4().String()),
		Port: addr.Port,
		Zone: "",
	}

	return &endpoint{
		dstIP: destIP,
	}
}

var (
	_ Endpoint = (*endpoint)(nil)
)

func (p *endpoint) SrcToString() string {
	return p.srcIP.String()
}

func (p *endpoint) DstToString() string {
	return p.dstIP.String()
}

func (p *endpoint) SetSrcIP(ip net.IP) {
	p.srcIP = ip
}

func (p *endpoint) SrcIP() net.IP {
	return p.srcIP
}

func (p *endpoint) DstIP() *net.UDPAddr {
	return p.dstIP
}

func ResolveUapAddr(addr string) {

}
