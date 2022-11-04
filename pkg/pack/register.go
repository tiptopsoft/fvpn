package pack

type RegisterPacket struct {
	*CommonPacket
	SrcMac  [4]byte
	DestMac [4]byte
}

func (cp *RegisterPacket) Encode() ([]byte, error) {

	return nil, nil
}

func (cp *RegisterPacket) Decode() (interface{}, error) {

	return cp, nil
}
