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

package util

import (
	"errors"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"net"
)

type IPHeader struct {
	SrcIP net.IP
	DstIP net.IP
}

// GetIPFrameHeader return srcIP, destIP
func GetIPFrameHeader(buff []byte) (*IPHeader, error) {
	if len(buff) < packet.IPBuffSize {
		return nil, errors.New("invalid ip frame")
	}

	h := new(IPHeader)
	h.SrcIP = net.IPv4(buff[12], buff[13], buff[14], buff[15])
	h.DstIP = net.IPv4(buff[16], buff[17], buff[18], buff[19])
	return h, nil
}

func GetPacketHeader(buff []byte) (packet.Header, error) {
	if len(buff) < packet.HeaderBuffSize {
		return packet.Header{}, errors.New("not invalid packer")
	}
	h, err := packet.Decode(buff[:packet.HeaderBuffSize])
	if err != nil {
		return packet.Header{}, err
	}
	return h, nil
}
