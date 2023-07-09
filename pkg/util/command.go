package util

import (
	"os"
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

func ExecCommand(name string, commands ...string) error {
	cmd := exec.Command(name, commands...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
