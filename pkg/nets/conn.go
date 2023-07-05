package nets

import (
	"net"
)

// Bind a Bind listens a port for IPV6 and IPv4 UDP packets. Also send packets to destination peer.
type Bind interface {
	// Open listen a port using a given port, if not successec, a random port will return, which ia actualPort
	Open(port uint16) (actualPort uint16, conn *net.UDPConn, err error)

	Send(buff []byte, ep Endpoint) (int, error)

	Conn() *net.UDPConn
	//BatchSize is size use to receive or send
	BatchSize() int
}

type Endpoint interface {
	SrcToString() string
	DstToString() string
	SrcIP() *net.UDPAddr
	DstIP() *net.UDPAddr
}

type endpoint struct {
	srcIP *net.UDPAddr
	dstIP *net.UDPAddr
}

func NewEndpoint(dstip string) Endpoint {
	addr, err := net.ResolveUDPAddr("udp", dstip)
	if err != nil {
		return nil
	}
	return &endpoint{
		dstIP: addr,
	}
}

var (
	_ Endpoint = (*endpoint)(nil)
)

func (p *endpoint) SrcToString() string {
	return p.srcIP.String()
}

func (p *endpoint) DstToString() string {
	return p.dstIP.String()
}

func (p *endpoint) SrcIP() *net.UDPAddr {
	return p.srcIP
}

func (p *endpoint) DstIP() *net.UDPAddr {
	return p.dstIP
}
