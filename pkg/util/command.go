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

package util

import (
	"fmt"
	"os/exec"
)

const (
	MsgTypeQueryPeer     uint16 = 1
	MsgTypePacket        uint16 = 3
	MsgTypeRegisterAck   uint16 = 4
	MsgTypeRegisterSuper uint16 = 5
	HandShakeMsgType     uint16 = 6
	KeepaliveMsgType     uint16 = 7
	HandShakeMsgTypeAck  uint16 = 8
)

func GetFrameTypeName(key uint16) (name string) {
	switch key {
	case HandShakeMsgType:
		name = "handshake"
	case HandShakeMsgTypeAck:
		name = "handshakeAck"
	case MsgTypeQueryPeer:
		name = "queryPacket"
	case MsgTypePacket:
		name = "MsgPacket"
	case KeepaliveMsgType:
		name = "keeplive"
	case MsgTypeRegisterSuper:
		name = "registryPacket"
	}

	return
}

func ExecCommand(name string, commands ...string) error {
	cmd := exec.Command(name, commands...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(output))
	return nil
}
