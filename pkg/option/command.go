package option

import (
	"os"
	"os/exec"
)

const (
	MSG_TYPE_REGISTER           uint16 = 1
	MSG_TYPE_DEREGISTER         uint16 = 2
	MSG_TYPE_PACKET             uint16 = 3
	MSG_TYPE_REGISTER_ACK       uint16 = 4
	MSG_TYPE_REGISTER_SUPER     uint16 = 5
	MSG_TYPE_UNREGISTER_SUPER   uint16 = 6
	MSG_TYPE_REGISTER_SUPER_ACK uint16 = 7
	MSG_TYPE_REGISTER_SUPER_NAK uint16 = 8
	MSG_TYPE_FEDERATION         uint16 = 9
	MSG_TYPE_PEER_INFO          uint16 = 10
	MSG_TYPE_QUERY_PEER         uint16 = 11
	MSG_TYPE_MAX_TYPE           uint16 = 11
	MSG_TYPE_RE_REGISTER_SUPER  uint16 = 12
)

func ExecCommand(name string, commands ...string) error {
	cmd := exec.Command(name, commands...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
