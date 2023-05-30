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
	MsgTypeNotifyAck        uint16 = 19
	MsgTypePunchHole        uint16 = 14
	PacketFromTap           uint16 = 15
	PacketFromUdp           uint16 = 16
	MsgTypePing             uint16 = 17
	MsgTypePingAck          uint16 = 18

	RestrictNat  uint8 = 1
	SymmetricNAT uint8 = 2
)

func ExecCommand(name string, commands ...string) error {
	cmd := exec.Command(name, commands...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
