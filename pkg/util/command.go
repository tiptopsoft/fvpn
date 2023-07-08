package util

import (
	"os"
	"os/exec"
)

const (
	MsgTypeRegister        uint16 = 1
	MsgTypePacket          uint16 = 3
	MsgTypeRegisterAck     uint16 = 4
	MsgTypeRegisterSuper   uint16 = 5
	MsgTypeUnregisterSuper uint16 = 6
	MsgTypeQueryPeer       uint16 = 11
	HandShakeMsgType       uint16 = 12
	KeepaliveMsgType       uint16 = 15
	HandShakeMsgTypeAck    uint16 = 13
	MsgTypeNotifyAck       uint16 = 14
)

func ExecCommand(name string, commands ...string) error {
	cmd := exec.Command(name, commands...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
