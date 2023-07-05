package nets

import "net"

type StdNetBind struct {
	v4conn *net.UDPConn
	v6conn *net.UDPConn
}

var (
	_ Bind = (*StdNetBind)(nil)
)

func NewStdBind() Bind {
	return &StdNetBind{}
}

func (s *StdNetBind) Open(port uint16) (uint16, *net.UDPConn, error) {
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return 0, nil, nil
	}

	addr := conn.LocalAddr()
	listenAddr, err := net.ResolveUDPAddr(
		addr.Network(),
		addr.String(),
	)
	if err != nil {
		return 0, nil, err
	}

	s.v4conn = conn
	s.v6conn = conn
	return uint16(listenAddr.Port), conn, nil
}

func (s *StdNetBind) Send(buff []byte, ep Endpoint) (int, error) {
	return s.v4conn.WriteToUDP(buff, ep.DstIP())
}

func (*StdNetBind) BatchSize() int {
	return 0
}

func (s *StdNetBind) Conn() *net.UDPConn {
	return s.v4conn
}
