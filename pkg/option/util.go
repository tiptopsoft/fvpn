package option

import (
	"net"
)

func GetLocalMac(name string) ([4]byte, error) {
	var b [4]byte
	fa, err := net.InterfaceByName(name)

	if err != nil {
		return [4]byte{}, err
	}
	macAddr := fa.HardwareAddr.String()
	if len(macAddr) == 0 {
		return [4]byte{}, ErrGetMac
	}

	copy(b[:], macAddr)
	return b, nil

}
