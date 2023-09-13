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

package packet

import (
	"encoding/binary"
)

const (
	FvpnPktBuffSize = 2048
)

const (
	HeaderBuffSize    = 44
	HandshakeBuffSize = 76
	IPBuffSize        = 20
)

type Packet struct {
	dstBuff []byte
	srcBuff []byte
}

func EncodeBytes(dst, src []byte, idx int) int {
	copy(dst[idx:idx+len(src)], src[:])
	idx += len(src)
	return idx
}

func EncodeUint8(dst []byte, src uint8, idx int) int {
	dst[idx] = src
	idx += 1
	return idx
}

func EncodeUint16(dst []byte, src uint16, idx int) int {
	var b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, src)
	copy(dst[idx:idx+2], b[:])
	idx += 2
	return idx
}

func DecodeUint8(dst *byte, src []byte, idx int) int {
	*dst = src[idx]
	idx += 1
	return idx
}

func DecodeUint16(dst *uint16, src []byte, idx int) int {
	v := binary.BigEndian.Uint16(src[idx : idx+2])
	*dst = v
	idx += 2
	return idx
}

func DecodeBytes(dst *[]byte, src []byte, idx int) int {
	copy(*dst, src[idx:idx+len(*dst)])
	idx += len(*dst)
	return idx
}
