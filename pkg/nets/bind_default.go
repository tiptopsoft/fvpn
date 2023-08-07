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

package nets

import (
	"fmt"
	"net"
)

type StdNetBind struct {
	v4conn *net.UDPConn
	v6conn *net.UDPConn
}

var (
	_ Bind = (*StdNetBind)(nil)
)

func NewStdBind() Bind {
	return &StdNetBind{}
}

func (s *StdNetBind) Open(port uint16) (uint16, *net.UDPConn, error) {
	listen := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}
	conn, err := net.ListenUDP("udp4", listen)
	if err != nil {
		return 0, nil, nil
	}

	addr := conn.LocalAddr()
	listenAddr, err := net.ResolveUDPAddr(
		addr.Network(),
		addr.String(),
	)
	if err != nil {
		return 0, nil, err
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>listenAddr", listenAddr)

	s.v4conn = conn
	s.v6conn = conn
	return uint16(listenAddr.Port), conn, nil
}

func (s *StdNetBind) Send(buff []byte, ep Endpoint) (int, error) {
	return s.v4conn.WriteToUDP(buff, ep.DstIP())
}

func (*StdNetBind) BatchSize() int {
	return 0
}

func (s *StdNetBind) Conn() *net.UDPConn {
	return s.v4conn
}
