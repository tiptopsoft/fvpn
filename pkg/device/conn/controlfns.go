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

package conn

import (
	"net"
	"syscall"
)

// UDP socket read/write buffer size (7MB). The value of 7MB is chosen as it is
// the max supported by a default configuration of macOS. Some platforms will
// silently clamp the value to other maximums, such as linux clamping to
// net.core.{r,w}mem_max (see _linux.go for additional implementation that works
// around this limitation)
const socketBufferSize = 7 << 20

type fn func(network, address string, c syscall.RawConn) error

var fns []fn

func ListenConfig() *net.ListenConfig {
	return &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			for _, f := range fns {
				err := f(network, address, c)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
