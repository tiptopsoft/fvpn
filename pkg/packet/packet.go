package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// Packet star's Packet
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
	PACKET_SIZE = 0x18 //24
)

var (
	INVALIED_PACKET = errors.New("invalid packet")
)

type Packet struct {
	Version     uint8  //1
	TTL         uint8  //1
	Flags       uint16 //2
	Group       uint32 //4 every group 4 byte
	SourceMac   uint32 //4
	DestMac     uint32 //4
	SocketFlags uint16 //2
	IPv4        uint32 //4
	UdpPort     uint16 //2

}

// Encode transfer packet to byte stream
func Encode() ([]byte, error) {
	p := Packet{}
	var b []byte
	if bs, err := IntToBytes(int(p.Version)); err != nil {
		return nil, err
	} else {
		b = append(b, bs...)
	}

	if bs, err := IntToBytes(int(p.TTL)); err != nil {
		return nil, err
	} else {
		b = append(b, bs...)
	}

	if bs, err := IntToBytes(int(p.Flags)); err != nil {
		return nil, err
	} else {
		b = append(b, bs...)
	}

	if bs, err := IntToBytes(int(p.Group)); err != nil {
		return nil, err
	} else {
		b = append(b, bs...)
	}

	if bs, err := IntToBytes(int(p.SourceMac)); err != nil {
		return nil, err
	} else {
		b = append(b, bs...)
	}

	if bs, err := IntToBytes(int(p.DestMac)); err != nil {
		return nil, err
	} else {
		b = append(b, bs...)
	}

	if bs, err := IntToBytes(int(p.SocketFlags)); err != nil {
		return nil, err
	} else {
		b = append(b, bs...)
	}

	if bs, err := IntToBytes(int(p.IPv4)); err != nil {
		return nil, err
	} else {
		b = append(b, bs...)
	}

	if bs, err := IntToBytes(int(p.UdpPort)); err != nil {
		return nil, err
	} else {
		b = append(b, bs...)
	}

	if len(b) > 20 {
		return nil, INVALIED_PACKET
	}
	return b, nil
}

func Decode(b []byte) (*Packet, error) {
	p := &Packet{}

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

	if sMac, err := BytesToInt(b[8:12]); err != nil {
		return nil, err
	} else {
		p.SourceMac = uint32(sMac)
	}

	if dMac, err := BytesToInt(b[12:16]); err != nil {
		return nil, err
	} else {
		p.DestMac = uint32(dMac)
	}

	if sFlags, err := BytesToInt(b[16:18]); err != nil {
		return nil, err
	} else {
		p.SocketFlags = uint16(sFlags)
	}

	if ipv4, err := BytesToInt(b[18:22]); err != nil {
		return nil, err
	} else {
		p.IPv4 = uint32(ipv4)
	}

	if udpPort, err := BytesToInt(b[22:24]); err != nil {
		return nil, err
	} else {
		p.DestMac = uint32(udpPort)
	}

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
