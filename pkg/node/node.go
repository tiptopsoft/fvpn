package node

import (
	"context"
	. "github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/nets"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/security"
	"github.com/topcloudz/fvpn/pkg/tun"
	"github.com/topcloudz/fvpn/pkg/util"
	"sync"
)

var (
	logger = log.Log()
)

// Node is a dev in any os.
type Node struct {
	cfg        *util.Config
	privateKey security.NoisePrivateKey
	pubKey     security.NoisePublicKey
	device     tun.Device
	net        struct {
		bind nets.Bind
	}

	//peers is all peers releated to this device
	peers struct {
		lock  sync.Mutex
		peers map[security.NoisePublicKey]*Peer //dst
	}

	queue struct {
		outBound *OutBoundQueue //after encrypt
		inBound  *InBoundQueue  //after decrypt
	}

	netManager NetManagerFn
	tunHandler Handler
	udpHandler Handler
	relay      *Peer
	wg         sync.WaitGroup
	userId     [8]byte
	cache      CacheFunc
}

func NewDevice(iface tun.Device, bind nets.Bind) (*Node, error) {
	n := &Node{
		device: iface,
		net:    struct{ bind nets.Bind }{bind: bind},
		cache:  NewCache(),
	}
	privateKey, err := security.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	n.privateKey = privateKey
	n.pubKey = n.privateKey.NewPubicKey()
	n.queue.outBound = NewOutBoundQueue()
	n.queue.inBound = NewInBoundQueue()
	//n.queue.handshakeBound = newHandshakeQueue()

	n.tunHandler = WithMiddlewares(tunHandler(), AuthCheck(), Encode())
	n.udpHandler = WithMiddlewares(udpHandler(), AuthCheck(), Decode())
	n.wg.Add(1)

	return n, nil
}

func (n *Node) initRelay() {
	n.relay = n.NewPeer(security.NoisePublicKey{})
	n.relay.node = n
	n.relay.endpoint = nets.NewEndpoint(n.cfg.ClientCfg.Registry)
	n.relay.start()
}

func (n *Node) NewPeer(pk security.NoisePublicKey) *Peer {
	p := new(Peer)
	p.PubKey = pk
	p.queue.outBound = NewOutBoundQueue()
	p.queue.inBound = NewInBoundQueue()
	p.node = n
	return p
}

func (n *Node) nodeRegister() error {
	rPkt := register.NewPacket()
	n.userId = rPkt.UserId
	copy(rPkt.PubKey[:], n.pubKey[:])
	buff, err := register.Encode(rPkt)
	if err != nil {
		return nil
	}

	size := len(buff)
	f := packet.NewFrame()
	copy(f.Packet[:size], buff)
	n.relay.PutPktToOutbound(f)
	return nil
}

func Start(cfg *util.Config) error {
	iface, err := tun.New()
	if err != nil {
		return err
	}

	d, err := NewDevice(iface, nets.NewStdBind())
	d.cfg = cfg
	if err != nil {
		return err
	}

	return d.up()
}

func (n *Node) up() error {
	defer n.wg.Done()
	port, _, err := n.net.bind.Open(0)
	logger.Debugf("fvpn start at: %d", port)
	if err != nil {
		return err
	}

	go n.ReadFromUdp()
	go n.ReadFromTun()

	n.initRelay()
	n.wg.Wait()
	return nil
}

func (n *Node) Close() error {
	close(n.queue.outBound.c)
	return nil
}

func (n *Node) ReadFromTun() {
	for {
		ctx := context.Background()
		f := packet.NewFrame()
		f.UserId = n.userId
		size, err := n.device.Read(f.Buff[:])
		if err != nil {
			continue
		}
		logger.Debugf("node %s receive %n byte", n.device.Name(), size)
		err = n.tunHandler.Handle(ctx, f)

		if err != nil {
			continue
		}
	}
}

func (n *Node) ReadFromUdp() {
	for {
		ctx := context.Background()
		f := packet.NewFrame()
		size, remoteAddr, err := n.net.bind.Conn().ReadFromUDP(f.Buff[:])
		if err != nil {
			continue
		}

		err = n.udpHandler.Handle(ctx, f)
		if err != nil {
			continue
		}
		logger.Debugf("udp receive %d byte from %s, data: %v", size, remoteAddr.IP, f.Buff[:size])
	}
}

func (n *Node) WriteToUDP() {
	for {
		select {
		case pkt := <-n.queue.outBound.c:
			//only packet will find peer, other type will send to regServer
			if pkt.FrameType == util.MsgTypePacket {
				var peer *Peer
				peer.queue.outBound.c <- pkt
			}
		default:

		}
	}
}

// appId is a unique identify for a node
func (n *Node) appId() string {
	return string(n.pubKey[:])
}

//func (d *Node) RoutineEncryption(id int) {
//
//}
//
//func (d *Node) RoutineDescryption(id int) {
//
//}
//
//func (d *Node) RoutineHandshake(id int) {
//
//}
