package pack

const (
	TAP_REGISTER       = 10
	TAP_REGISTER_ACK   = 11
	TAP_MESSAGE        = 12
	TAP_LIST_EDGE_STAR = 13
	TAP_BROADCAST      = 14
	TAP_UNREGISTER     = 15
)

var (
	Version          uint8  = 1
	DefaultTTL       uint8  = 100
	IPV4             uint16 = 0x01
	IPV6             uint16 = 0x02
	COMMON_FRAM_SIZE        = 20
	DefaultPort      uint16 = 3000
)

//CommonPacket  every time sends base frame.
type CommonPacket struct {
	Version uint8   //1
	TTL     uint8   //1
	Flags   uint16  //2
	Group   [4]byte //4
}

func (cp *CommonPacket) Encode() ([]byte, error) {

	var b [8]byte
	b[0] = cp.Version
	copy(b[1:2], []byte{cp.TTL})
	if bs, err := IntToBytes(int(cp.Flags)); err != nil {
		return nil, err
	} else {
		copy(b[2:4], bs[2:4])
	}

	copy(b[4:8], cp.Group[:])
	return b[:], nil
}

func (cp *CommonPacket) Decode(udpByte []byte) (*CommonPacket, error) {
	cp.Version = udpByte[0]
	cp.TTL = udpByte[1]
	cp.Flags = BytesToInt16(udpByte[2:4])
	copy(cp.Group[:], udpByte[4:8])

	return cp, nil
}
