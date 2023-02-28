package packet

import (
	"encoding/binary"
	"net"
)

// Packet edge's Packet
/**
  As learn from edge, our packet is form of below:
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
	MAC_SIZE = 6
	IP_SIZE  = 16
)

type Packet struct {
	dstBuff []byte
	srcBuff []byte
}

func New() Packet {
	return Packet{
		dstBuff: make([]byte, 2048),
		srcBuff: make([]byte, 2048),
	}
}

func EncodeBytes(dst, src []byte, idx int) int {
	copy(dst[idx:idx+len(src)], src[:])
	idx += len(src)
	return idx
}

func EncodeUint8(dst []byte, src uint8, idx int) int {
	dst[idx] = src
	idx += 1
	return idx
}

func EncodeUint16(dst []byte, src uint16, idx int) int {
	var b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, src)
	copy(dst[idx:idx+2], b[:])
	idx += 2
	return idx
}

func DecodeUint8(dst *byte, src []byte, idx int) int {
	*dst = src[idx]
	idx += 1
	return idx
}

func DecodeUint16(dst *uint16, src []byte, idx int) int {
	v := binary.BigEndian.Uint16(src[idx : idx+2])
	*dst = v
	idx += 2
	return idx
}

func DecodeBytes(dst *[]byte, src []byte, idx int) int {
	copy(*dst, src[idx:idx+len(*dst)])
	idx += len(*dst)
	return idx
}

func DecodeMacAddr(src []byte, idx int) (net.HardwareAddr, int) {
	mac := make([]byte, MAC_SIZE)
	idx = DecodeBytes(&mac, src, idx)
	return mac, idx
}
