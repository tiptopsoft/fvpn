package pack

import (
	"fmt"
	"testing"
)

func TestCommonPacket_Decode(t *testing.T) {

}

func TestCommonPacket_Encode(t *testing.T) {
	bs, _ := IntToBytes(20000)
	var b [4]byte
	copy(b[:], bs[:])
	cp := &CommonPacket{
		Version: Version,
		TTL:     DefaultTTL,
		Flags:   TAP_REGISTER,
		Group:   b,
	}

	result, _ := cp.Encode()
	fmt.Println(result)

	cp1 := &CommonPacket{}
	cp1, _ = cp1.Decode(result)
	fmt.Println(cp1)

}
