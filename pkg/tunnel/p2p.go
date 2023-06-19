package tunnel

import (
	"context"
	"github.com/ccding/go-stun/stun"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/socket"

	"net"
	"time"
)

var (
	Pool = NewPool()
)

// NewP2PTunnel use stun server to init serval client used to p2p connection
func (t *Tunnel) NewP2PTunnel(f *packet.Frame) (*Tunnel, error) {

	sock := socket.NewSocket(0)
	var err error
	if err != nil {
		logger.Errorf("close origin socket failed: %v", err)
		return nil, err
	}
	err = sock.Connect(f.Target.Addr)
	if err != nil {
		logger.Errorf("connect p2p address failed. %v", err)
		return nil, err
	}

	//open session, node-> remote addr
	handPacket, _ := header.NewHeader(option.HandShakeMsgType, f.NetworkId)
	buff, _ := header.Encode(handPacket)
	newAddr, _ := sock.LocalAddr()
	destAddr := f.Target
	logger.Debugf(">>>>>>> punching hole, localIP: %v, port: %v, destIP: %v, destPort: %v", newAddr.Addr, newAddr.Port, destAddr.Addr, destAddr.Port)
	if err != nil {
		logger.Errorf("open hole failed: %v", err)
	}
	_, err = sock.Write(buff)
	if err != nil {
		logger.Errorf("send punch hole failed: %v", err)
		return nil, err
	}

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	ch := make(chan int)
	data := make([]byte, 1024)
	go func() {
		n, err := sock.Read(data)
		if err != nil {
			ch <- 0
		}
		logger.Debugf("hole msg size: %d, data: %v", n, data)
		if n > 0 {
			//start a p2p runner
			//go t.p2pRunner(sock, pNode.NodeInfo)
			ch <- 1
		}
	}()

	select {
	case v := <-ch:
		if v == 1 {
			//pNode.NodeInfo.P2P = true
			//pNode.NodeInfo.Socket = newSock
			//t.cache.SetCache(pkt.NetworkId, pNode.NodeInfo.IP.String(), pNode.NodeInfo)
			logger.Debugf("punch hole success")
		} else {
			logger.Debugf("punch hole failed.")
		}
	case <-timeout.Done():
		logger.Debugf("punch hole failed.")
	}

	return nil, nil
}

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

	go func() {
		for i := 0; i < 10; i++ {
			p, err := initPortPair()

			if err != nil {
				ch <- p
			}
		}
	}()

	return PortPairPool{
		ch: ch,
	}
}

func initPortPair() (*PortPair, error) {
	client := stun.NewClient()

	conn, _ := net.ListenUDP("udp", nil)
	stun.NewClientWithConnection(conn)

	addr := conn.LocalAddr().(*net.UDPAddr)

	client.SetServerAddr("stun.miwifi.com:3478")
	//client.SetServerAddr("101.43.97.112:3478")
	_, host, err := client.Discover()

	if err != nil {
		return nil, err
	}

	p := new(PortPair)
	p.SrcPort = uint16(addr.Port)
	p.SrcIP = addr.IP
	p.NatIP = net.ParseIP(host.IP())
	p.NatPort = host.Port()

	return p, nil
}

func GetPortPair() {

}
