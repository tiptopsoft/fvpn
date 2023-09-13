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
	"fmt"
	"net"
	"testing"
)

func Test(t *testing.T) {
	//tun, err := New()
	//if err != err {
	//	panic(err)
	//}
	//
	//fmt.Println("tun is: ", tun.name)
	//buff := make([]byte, 1024)
	//for {
	//	n, err := tun.Read(buff)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	fmt.Println(fmt.Sprintf("Read from %s %d byte", tun.name, n))
	//}
	//

	iface, err := net.InterfaceByName("utun3")
	if err != nil {
		panic(err)
	}

	addr1, err := iface.Addrs()
	if err != nil {
		panic(err)
	}
	fmt.Println(addr1[0].Network(), addr1[0].(*net.IPNet).IP)
}
