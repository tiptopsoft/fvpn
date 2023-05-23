package option

import (
	"os"
	"os/exec"
)

const (
	MsgTypeRegister         uint16 = 1
	MsgTypeDeregister       uint16 = 2
	MsgTypePacket           uint16 = 3
	MsgTypeRegisterAck      uint16 = 4
	MsgTypeRegisterSuper    uint16 = 5
	MsgTypeUnregisterSuper  uint16 = 6
	MsgTypeRegisterSuperAck uint16 = 7
	MsgTypeRegisterSuperNak uint16 = 8
	MsgTypeFederation       uint16 = 9
	MsgTypePeerInfo         uint16 = 10
	MsgTypeQueryPeer        uint16 = 11
	MsgTypeMaxType          uint16 = 11
	MsgTypeReRegisterSuper  uint16 = 12
	MsgTypeNotify           uint16 = 13
)

func ExecCommand(name string, commands ...string) error {
	cmd := exec.Command(name, commands...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
