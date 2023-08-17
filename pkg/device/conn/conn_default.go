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

package conn

import (
	"fmt"
	"net"
)

type Default struct {
	v4conn *net.UDPConn
	v6conn *net.UDPConn
	ipv6   bool
}

var (
	_ Interface = (*Default)(nil)
)

func New(enable bool) Interface {
	return &Default{
		ipv6: enable,
	}
}

func (s *Default) Open(port uint16) (uint16, error) {
	var err error
	ipv4Addr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: int(port),
	}
	s.v4conn, err = listen("udp4", ipv4Addr)
	if err != nil {
		return 0, nil
	}

	addr := s.v4conn.LocalAddr()
	listenAddr, err := net.ResolveUDPAddr(
		addr.Network(),
		addr.String(),
	)
	if err != nil {
		logger.Debugf("bind v4 failed, %v", err)
	}

	ipv6Addr := &net.UDPAddr{IP: net.IPv6zero, Port: int(port)}

	s.v6conn, err = listen("udp6", ipv6Addr)

	if err != nil {
		return 0, fmt.Errorf("open bind failed, error: %v", err)
	}

	return uint16(listenAddr.Port), nil
}

func (s *Default) Send(buff []byte, ep Endpoint) (int, error) {
	if s.ipv6 {
		return s.send6(buff, ep)
	}

	return s.send4(buff, ep)
}

func (s *Default) send4(buff []byte, ep Endpoint) (int, error) {
	return s.v4conn.WriteToUDP(buff, ep.DstIP())
}

func (s *Default) send6(buff []byte, ep Endpoint) (int, error) {
	return s.v6conn.WriteToUDP(buff, ep.DstIP())
}

func (s *Default) Conn() *net.UDPConn {
	return s.v4conn
}

func listen(network string, addr *net.UDPAddr) (*net.UDPConn, error) {
	return net.ListenUDP(network, addr)
}
