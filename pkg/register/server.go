package register

import (
	"context"
	"fmt"
	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/handler/auth"
	"github.com/interstellar-cloud/star/pkg/handler/encrypt"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/option"
	"github.com/interstellar-cloud/star/pkg/packet/common"
	"github.com/interstellar-cloud/star/pkg/packet/register"
	"github.com/interstellar-cloud/star/pkg/packet/register/ack"
	"net"
	"sync"
)

var limitChan = make(chan int, 1)

// mac:Pub
var m sync.Map

type Node struct {
	Mac   [4]byte
	Proto option.Protocol
	Conn  net.Conn
	Addr  *net.UDPAddr
}

//RegStar use as register
type RegStar struct {
	*option.RegConfig
	handler.Executor
	conn net.Conn
}

func (r *RegStar) Start(address string) error {
	return r.start(address)
}

// Node register node for net, and for user create edge
func (r *RegStar) start(address string) error {
	var ctx = context.Background()
	var conn net.Conn
	r.Executor = handler.NewExecutor()
	if r.OpenAuth {
		r.AddHandler(ctx, &auth.AuthHandler{})
	}

	if r.OpenEncrypt {
		r.AddHandler(ctx, &encrypt.StarEncrypt{})
	}

	switch r.Protocol {
	case option.UDP:
		addr, err := ResolveAddr(address)
		if err != nil {
			return err
		}

		conn, err = net.ListenUDP("udp", addr)

		log.Logger.Infof("registry start at: %s", address)

		if err != nil {
			return err
		}
		defer conn.Close()
		for {
			limitChan <- 1
			go r.handleUdp(ctx, conn.(*net.UDPConn))
		}
	default:
		log.Logger.Info("this is a tcp server")
	}

	return nil
}

func (r *RegStar) handleUdp(ctx context.Context, conn *net.UDPConn) {
	for {
		data := make([]byte, 2048)
		_, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
		}

		p, err := common.Decode(data)
		if err != nil {
			fmt.Println(err)
		}

		//exec executor
		if err := r.Execute(ctx, data); err != nil {
			fmt.Println(err)
		}

		switch p.Flags {

		case option.MSG_TYPE_REGISTER_SUPER:
			r.processRegister(addr, conn, data, nil)
			break

		}
	}

}

// register edge node register to register
func (r *RegStar) register(addr *net.UDPAddr, packet register.RegPacket) error {
	m.Store(packet.SrcMac.String(), addr)
	m.Range(func(key, value any) bool {
		log.Logger.Infof("registry data key: %s, value: %v", key, value)
		return true
	})
	return nil
}

func (r *RegStar) unRegister(packet register.RegPacket) error {
	m.Delete(packet.SrcMac)
	return nil
}

func ackBuilder(cp common.CommonPacket) ([]byte, error) {

	RecMac := "01:01:03:02:03:01"
	ip := "192.168.1.1"
	Mask := "255.255.255.0"

	p := ack.NewPacket()
	mac, err := net.ParseMAC(RecMac)
	if err != nil {
		log.Logger.Errorf("invalid mac:%s", RecMac)
	}

	p.RegMac = mac
	p.AutoIP = net.ParseIP(ip)
	p.Mask = net.ParseIP(Mask)
	cp.Flags = option.MSG_TYPE_REGISTER_ACK
	p.CommonPacket = cp

	return ack.Encode(p)
}

func ResolveAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}
