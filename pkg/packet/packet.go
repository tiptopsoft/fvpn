package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/interstellar-cloud/star/pkg/option"
)

// Frame star's Frame
/**
  As learn from star, our packet is form of below:
 Version 1

    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   ! Version=1     ! TTL           ! Flags                         !
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 4 ! Community                                                     :
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 8 ! ... Community ...                                             :
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
12 ! ... Community ...                                             :
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
16 ! ... Community ...                                             :
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
20 ! ... Community ...                                             !
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
24 ! Source MAC Address                                            :
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
28 :                               ! Destination MAC Address       :
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
32 :                                                               !
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
36 ! Socket Flags (v=IPv4)         ! Destination UDP Port          !
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
40 ! Destination IPv4 Address                                      !
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
44 ! Compress'n ID !  Transform ID !
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
48 ! Payload
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

   Socket Flags can be ipv6

Now , we just impl ipv4, and have only one group.
*/

const (
	FRAME_SIZE = 0x18 //24
)

var (
	INVALIED_FRAME = errors.New("invalid frame")
)

type Frame struct {
	Version     uint8   //1
	TTL         uint8   //1
	Flags       uint16  //2
	Group       uint32  //4 every group 4 byte
	SourceMac   [4]byte //4
	DestMac     [4]byte //4
	SocketFlags uint16  //2 user v4 or v6
	IPv4        [4]byte //4
	UdpPort     uint16  //2
}

func NewPacket() *Frame {
	return &Frame{
		Version: uint8(option.Version),
		TTL:     uint8(option.DefaultTTL),
	}
}

// Encode transfer packet to byte stream
func Encode(frame *Frame) ([]byte, error) {
	var b [24]byte
	if bs, err := IntToBytes(int(frame.Version)); err != nil {
		return nil, err
	} else {
		copy(b[:1], bs)
	}

	if bs, err := IntToBytes(int(frame.TTL)); err != nil {
		return nil, err
	} else {
		copy(b[1:2], bs)
	}

	if bs, err := IntToBytes(int(frame.Flags)); err != nil {
		return nil, err
	} else {
		copy(b[2:4], bs)
	}

	if bs, err := IntToBytes(int(frame.Group)); err != nil {
		return nil, err
	} else {
		copy(b[4:8], bs)
	}

	copy(b[8:12], frame.SourceMac[:])

	copy(b[12:16], frame.DestMac[:])

	if bs, err := IntToBytes(int(frame.SocketFlags)); err != nil {
		return nil, err
	} else {
		copy(b[16:18], bs)
	}

	copy(b[18:22], frame.IPv4[:])

	if bs, err := IntToBytes(int(frame.UdpPort)); err != nil {
		return nil, err
	} else {
		copy(b[22:24], bs)
	}

	return b[:], nil
}

func Decode(b []byte) (*Frame, error) {
	if len(b) < FRAME_SIZE {
		return nil, INVALIED_FRAME
	}
	p := &Frame{}

	if v, err := BytesToInt(b[0:1]); err != nil {
		return nil, err
	} else {
		p.Version = uint8(v)
	}

	if ttl, err := BytesToInt(b[1:2]); err != nil {
		return nil, err
	} else {
		p.TTL = uint8(ttl)
	}

	if flags, err := BytesToInt(b[2:4]); err != nil {
		return nil, err
	} else {
		p.Flags = uint16(flags)
	}

	if group, err := BytesToInt(b[4:8]); err != nil {
		return nil, err
	} else {
		p.Group = uint32(group)
	}

	copy(p.SourceMac[:], b[8:12])
	copy(p.DestMac[:], b[12:16])

	if sFlags, err := BytesToInt(b[16:18]); err != nil {
		return nil, err
	} else {
		p.SocketFlags = uint16(sFlags)
	}

	copy(p.IPv4[:], b[18:22])
	copy(p.DestMac[:], b[22:24])

	return p, nil
}

func IntToBytes(n int) ([]byte, error) {
	data := int32(n)
	bytesBuf := bytes.NewBuffer([]byte{})
	if err := binary.Write(bytesBuf, binary.BigEndian, data); err != nil {
		return nil, err
	}
	return bytesBuf.Bytes(), nil
}

func BytesToInt(b []byte) (int, error) {
	bytesBuffer := bytes.NewBuffer(b)
	var data int32
	if err := binary.Read(bytesBuffer, binary.BigEndian, &data); err != nil {
		return 0, err
	}

	return int(data), nil
}
