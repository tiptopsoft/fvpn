package peer

import (
	"bytes"
	"encoding/binary"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/security"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
)

type PeerInfo struct {
	IP         net.IP
	RemoteAddr net.UDPAddr
	PubKey     security.NoisePublicKey
}

type PeerPacket struct {
	Header packet.Header
	UserId [8]byte
	Peers  []PeerInfo
}

func NewPeerPacket() PeerPacket {
	h, _ := packet.NewHeader(util.MsgTypeQueryPeer, handler.UCTL.UserId)
	return PeerPacket{
		Header: h,
		Peers:  nil,
	}
}

func Encode(peerPacket PeerPacket) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, peerPacket)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func Decode(buff []byte) (peerPacket PeerPacket, err error) {
	buf := bytes.NewReader(buff)
	err = binary.Read(buf, binary.BigEndian, &peerPacket)
	if err != nil {
		return PeerPacket{}, err
	}

	return peerPacket, nil
}
