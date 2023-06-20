package tunnel

import (
	"crypto/rand"
	"fmt"
	"github.com/ccding/go-stun/stun"
	"math"
	"math/big"

	reuse "github.com/libp2p/go-reuseport"
	"net"
)

var (
	Pool = NewPool()
)

// PortPair used for p2p
type PortPair struct {
	SrcPort uint16
	SrcIP   net.IP
	NatIP   net.IP
	NatPort uint16
}

type PortPairPool struct {
	ch chan *PortPair
}

// init 10
func NewPool() PortPairPool {
	ch := make(chan *PortPair, 10)
	pp := PortPairPool{
		ch: ch,
	}

	go func() {
		for {
			p, err := initPortPair()

			if err != nil {
				logger.Errorf("init port pair failed. %v:", err)
			}
			pp.ch <- p
		}
	}()

	return pp
}

func initPortPair() (*PortPair, error) {
	localPort := RandomPort(10000, 50000)
	conn, err := reuse.ListenPacket("udp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		return nil, err
	}
	client := stun.NewClientWithConnection(conn.(net.PacketConn))
	addr := conn.LocalAddr().(*net.UDPAddr)
	client.SetServerAddr("stun.miwifi.com:3478")
	_, host, err := client.Discover()

	if err != nil {
		return nil, err
	}

	p := new(PortPair)
	p.SrcPort = uint16(addr.Port)
	p.SrcIP = addr.IP
	p.NatIP = net.ParseIP(host.IP())
	p.NatPort = host.Port()

	logger.Debugf("init a portpair.............. src port:%v, nat ip: %v, nat port: %v", p.SrcPort, p.NatIP, p.NatPort)
	return p, nil
}

func RandomPort(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	}
	result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return min + result.Int64()
}
