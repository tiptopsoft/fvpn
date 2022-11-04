package pack

import (
	"errors"
)

type RegisterPacket struct {
	*CommonPacket
	SrcMac  [4]byte
	DestMac [4]byte
}

func (cp *RegisterPacket) Encode() ([]byte, error) {
	b := make([]byte, 28)

	cmBytes, err := cp.CommonPacket.Encode()
	if err != nil {
		return nil, errors.New("invalid common packets")
	}
	copy(b[0:20], cmBytes)
	copy(b[20:24], cp.SrcMac[:])
	copy(b[24:28], cp.DestMac[:])
	return b, nil
}

func (cp *RegisterPacket) Decode(udpBytes []byte) (interface{}, error) {

	res := &RegisterPacket{}
	cm := &CommonPacket{}
	cm, err := cm.Decode(udpBytes[0:20])
	if err != nil {
		return nil, errors.New("decode common packets failed")
	}

	copy(res.SrcMac[:], udpBytes[20:24])
	copy(res.DestMac[:], udpBytes[25:28])
	res.CommonPacket = cm
	return res, nil
}
