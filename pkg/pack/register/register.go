package register

import (
	"errors"
	"github.com/interstellar-cloud/star/pkg/pack/common"
	"unsafe"
)

// RegPacket register a edge to register
type RegPacket struct {
	*common.CommonPacket
	SrcMac [4]byte
}

func NewPacket() *RegPacket {
	return &RegPacket{}
}

func (cp *RegPacket) Encode() ([]byte, error) {
	b := make([]byte, unsafe.Sizeof(RegPacket{}))
	commonBytes, err := cp.CommonPacket.Encode()
	if err != nil {
		return nil, errors.New("encode common packet failed")
	}
	copy(b[0:20], commonBytes)
	copy(b[20:24], cp.SrcMac[:])
	return b, nil
}

func (reg *RegPacket) Decode(udpBytes []byte) (*RegPacket, error) {

	res := &RegPacket{}
	cp, err := common.NewPacket().Decode(udpBytes[0:20])
	if err != nil {
		return nil, errors.New("decode common packet failed")
	}
	res.CommonPacket = cp
	copy(res.SrcMac[:], udpBytes[20:24])
	return res, nil
}
