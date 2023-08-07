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

package peer

import (
	"fmt"
	"net"
	"testing"
)

func TestEncode1(t *testing.T) {

	p := NewPeerPacket()
	p.Header.SrcIP = net.ParseIP("121.1.1.1")
	buff, err := Encode(p)
	if err != nil {
		panic(buff)
	}

	p1, _ := Decode(buff)
	fmt.Println(p1.Header.SrcIP)
}
