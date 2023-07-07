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
	logger    = log.Log()
	relayPeer *Peer
)

// Node is a dev in any os.
type Node struct {
	mode       int
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

func (n *Node) PutPktToOutbound(pkt *packet.Frame) {
	n.queue.outBound.c <- pkt
}

func (n *Node) PutPktToInbound(pkt *packet.Frame) {
	n.queue.inBound.c <- pkt
}

func NewDevice(iface tun.Device, bind nets.Bind) (*Node, error) {
	n := &Node{
		device: iface,
		net:    struct{ bind nets.Bind }{bind: bind},
		cache:  NewCache(),
		mode:   1,
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

	n.tunHandler = WithMiddlewares(n.tunInHandler(), Encode(), AuthCheck())
	n.udpHandler = WithMiddlewares(n.udpInHandler(), AuthCheck(), Decode())
	n.wg.Add(1)

	return n, nil
}

func (n *Node) initRelay() {
	n.relay = n.NewPeer(security.NoisePublicKey{})
	n.relay.node = n
	n.relay.endpoint = nets.NewEndpoint(n.cfg.ClientCfg.Registry)
	n.relay.start()
	relayPeer = n.relay
	err := n.cache.SetPeer(UCTL.UserId, n.relay.endpoint.DstIP().IP.String(), n.relay)
	if err != nil {
		return
	}
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
	logger.Debugf("device name: %s, ip: %s", d.device.Name(), d.device.IPToString())
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
	go n.WriteToUDP()

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
		frame := packet.NewFrame()
		frame.Lock()
		ctx = context.WithValue(ctx, "cache", n.cache)
		frame.UserId = n.userId
		frame.FrameType = util.MsgTypePacket
		size, err := n.device.Read(frame.Buff[:])
		if err != nil {
			continue
		}
		frame.Size = size
		logger.Debugf("node %s receive %d byte", n.device.Name(), size)

		ipHeader, err := util.GetIPFrameHeader(frame.Buff[:])
		if err != nil {
			continue
		}
		frame.SrcIP = ipHeader.SrcIP
		frame.DstIP = ipHeader.DstIP
		err = n.tunHandler.Handle(ctx, frame)

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
		logger.Debugf("udp receive %d byte from %s, data: %v", size, remoteAddr.IP, f.Buff[:size])
		if err != nil {
			continue
		}
		f.Size = size
		f.RemoteAddr = remoteAddr

		hpkt, err := util.GetPacketHeader(f.Buff[:])
		if err != nil {
			logger.Error(err)
			continue
		}

		f.SrcIP = hpkt.SrcIP
		f.DstIP = hpkt.DstIP
		f.UserId = hpkt.UserId
		f.FrameType = hpkt.Flags

		err = n.udpHandler.Handle(ctx, f)
		if err != nil {
			continue
		}

	}
}

func (n *Node) WriteToUDP() {
	for {
		select {
		case pkt := <-n.queue.outBound.c:
			//only packet will find peer, other type will send to regServer
			if pkt.FrameType == util.MsgTypePacket {
				ip := pkt.DstIP
				peer, err := n.cache.GetPeer(pkt.UidString(), ip.String())
				if err != nil || peer == nil {
					peer = n.relay
				}
				peer.queue.outBound.c <- pkt
			}
		default:

		}
	}
}

func (n *Node) WriteToDevice() {
	for {
		select {
		case pkt := <-n.queue.inBound.c:
			if pkt.FrameType == util.MsgTypePacket {
				ip := pkt.DstIP
				peer, err := n.cache.GetPeer(pkt.UidString(), ip.String())
				if err != nil || peer == nil {
					peer = n.relay
				}

				peer.queue.inBound.c <- pkt
			}
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
