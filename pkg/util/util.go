package util

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/util/errors"
	"math/rand"
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
		return [4]byte{}, errors.ErrGetMac
	}

	copy(b[:], macAddr)
	return b, nil

}

// RandMac rand gen a mac
func RandMac() (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	buf[0] |= 2
	mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	return mac, nil
}
