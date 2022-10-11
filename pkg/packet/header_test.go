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

package packet

import (
	"fmt"
	"net"
	"testing"
)

func TestEncode(t *testing.T) {
	h, _ := NewHeader(3, "123444444444abcdef")
	h.SrcIP = net.ParseIP("5.244.24.141")
	h.DstIP = net.ParseIP("192.168.0.1")
	buff, _ := Encode(h)
	fmt.Println("len: ", len(buff))
	fmt.Println(buff)

	h1, _ := Decode(buff)
	fmt.Println(h1.SrcIP)
	fmt.Println(h1.DstIP)

}
