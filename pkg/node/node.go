package node

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/nets"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/security"
	"github.com/topcloudz/fvpn/pkg/tun"
	"github.com/topcloudz/fvpn/pkg/util"
	"io"
	"os"
	"sync"
	"time"
)

var (
	logger    = log.Log()
	relayPeer *Peer
)

// Node is a dev in any os.
type Node struct {
	mode       int
	cfg        *util.ClientConfig
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

	netCtl     NetworkManager
	tunHandler Handler
	udpHandler Handler
	relay      *Peer
	wg         sync.WaitGroup
	userId     [8]byte
	cache      CacheFunc
}

func (n *Node) PutPktToOutbound(pkt *Frame) {
	n.queue.outBound.c <- pkt
}

func (n *Node) PutPktToInbound(pkt *Frame) {
	n.queue.inBound.c <- pkt
}

func NewNode(iface tun.Device, bind nets.Bind, cfg *util.ClientConfig) (*Node, error) {
	n := &Node{
		device: iface,
		net:    struct{ bind nets.Bind }{bind: bind},
		cache:  NewCache(),
		mode:   1,
		cfg:    cfg,
	}
	n.netCtl = NewNetworkManager(util.UCTL.UserId)
	privateKey, err := security.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	n.privateKey = privateKey
	n.pubKey = n.privateKey.NewPubicKey()
	n.queue.outBound = NewOutBoundQueue()
	n.queue.inBound = NewInBoundQueue()
	//n.queue.handshakeBound = newHandshakeQueue()

	n.tunHandler = WithMiddlewares(n.tunInHandler(), AuthCheck(), n.AllowNetwork(), Encode())
	n.udpHandler = WithMiddlewares(n.udpInHandler(), AuthCheck(), Decode())
	n.wg.Add(1)

	return n, nil
}

func (n *Node) initRelay() {
	n.relay = n.NewPeer(security.NoisePublicKey{})
	n.relay.node = n
	n.relay.endpoint = nets.NewEndpoint(n.cfg.RegistryUrl())
	n.relay.start()
	n.relay.handshake(n.relay.endpoint.DstIP().IP)
	relayPeer = n.relay
	err := n.cache.SetPeer(util.UCTL.UserId, n.relay.endpoint.DstIP().IP.String(), n.relay)
	if err != nil {
		return
	}
}

func (n *Node) NewPeer(pk security.NoisePublicKey) *Peer {
	p := new(Peer)
	p.st = time.Now()
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
	f := NewFrame()
	copy(f.Packet[:size], buff)
	n.relay.PutPktToOutbound(f)
	return nil
}

func Start(cfg *util.Config) error {
	iface, err := tun.New()
	if err != nil {
		return err
	}

	//send http to get cidr
	client := NewClient(cfg.ClientCfg.ControlUrl())
	appId, err := appId()
	if err != nil {
		return err
	}
	resp, err := client.Init(appId)
	if err != nil {
		return err
	}
	d, err := NewNode(iface, nets.NewStdBind(), cfg.ClientCfg)
	if err != nil {
		return err
	}
	err = d.device.SetIP(resp.Mask, resp.IP)
	if err != nil {
		return err
	}
	logger.Debugf("device name: %s, cidr: %s", d.device.Name(), d.device.IPToString())

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
	go n.WriteToDevice()
	n.initRelay()
	go func() {
		timer := time.NewTimer(time.Second * 5)
		for {
			select {
			case <-timer.C:
				logger.Debugf("sending list packets...")
				n.sendListPackets()
				timer.Reset(time.Second * 5)
			}
		}
	}()

	go func() {
		err := n.HttpServer()
		if err != nil {
			logger.Errorf("start http failed. %v", err)
		}
	}()
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
		frame := NewFrame()
		//frame.Lock()
		ctx = context.WithValue(ctx, "cache", n.cache)
		frame.UserId = n.userId
		frame.FrameType = util.MsgTypePacket
		st1 := time.Now()
		size, err := n.device.Read(frame.Buff[:])
		st := time.Now()
		if err != nil {
			logger.Error(err)
			continue
		}
		logger.Debugf("node %s receive %d byte", n.device.Name(), size)
		ipHeader, err := util.GetIPFrameHeader(frame.Buff[:])
		if err != nil {
			continue
		}
		if ipHeader.DstIP.String() == n.device.Addr().String() {
			continue
		}

		peer, err := n.cache.GetPeer(util.UCTL.UserId, ipHeader.DstIP.String())
		if err != nil || peer == nil {
			if n.cfg.EnableRelay() {
				frame.Peer = n.relay
			}
		} else {
			if peer.p2p {
				frame.Peer = peer
			} else if n.cfg.EnableRelay() {
				frame.Peer = n.relay
			} else {
				//drop packets
				continue
			}
		}

		frame.SrcIP = ipHeader.SrcIP
		frame.DstIP = ipHeader.DstIP

		h, _ := packet.NewHeader(util.MsgTypePacket, util.UCTL.UserId)
		frame.UserId = h.UserId
		h.SrcIP = frame.SrcIP
		h.DstIP = frame.DstIP
		headerBuff, err := packet.Encode(h)
		if err != nil {
			logger.Error(err)
			continue
		}

		copy(frame.Packet[:packet.HeaderBuffSize], headerBuff)
		copy(frame.Packet[packet.HeaderBuffSize:], frame.Buff[:])
		frame.Size = size + packet.HeaderBuffSize

		err = n.tunHandler.Handle(ctx, frame)
		et := time.Since(st)
		et2 := time.Since(st1)
		logger.Debugf("================encode cost: %v, from read: %v", et, et2)

		if err != nil {
			logger.Error(err)
			continue
		}

	}
}

func (n *Node) ReadFromUdp() {
	defer func() {
		fmt.Println("ReadFromUDP has exit.....")
	}()
	for {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "cache", n.cache)
		f := NewFrame()
		size, remoteAddr, err := n.net.bind.Conn().ReadFromUDP(f.Buff[:])
		copy(f.Packet[:size], f.Buff[:size])
		if err != nil {
			logger.Error(err)
			continue
		}
		f.Size = size
		f.RemoteAddr = remoteAddr

		hpkt, err := util.GetPacketHeader(f.Buff[:])
		if err != nil {
			logger.Error(err)
			continue
		}
		logger.Debugf("udp receive %d byte from %s, data type: [%v]", size, remoteAddr.IP, util.GetFrameTypeName(hpkt.Flags))

		f.SrcIP = hpkt.SrcIP //192.168.0.1->192.168.0.2 srcIP =1, dstIP =2
		f.DstIP = hpkt.DstIP
		f.UserId = hpkt.UserId
		f.FrameType = hpkt.Flags

		err = n.udpHandler.Handle(ctx, f)
		if err != nil {
			logger.Error(err)
			continue
		}

	}
}

// sendListPackets send a packet list all nodes in current user
func (n *Node) sendListPackets() {
	//
	h, _ := packet.NewHeader(util.MsgTypeQueryPeer, util.UCTL.UserId)
	hpkt, err := packet.Encode(h)
	if err != nil {
		logger.Errorf("send list packet failed %v", err)
		return
	}
	frame := NewFrame()
	frame.Peer = n.relay
	frame.DstIP = n.relay.endpoint.DstIP().IP
	copy(frame.Packet, hpkt)
	frame.Size = len(hpkt)
	frame.UserId = h.UserId
	n.PutPktToOutbound(frame)
}

func (n *Node) WriteToUDP() {
	for {
		select {
		case pkt := <-n.queue.outBound.c:
			//only packet will find peer, other type will send to regServer
			ip := pkt.DstIP
			logger.Debugf("userId: %v, dst cidr: %v", pkt.UidString(), ip)
			pkt.Peer.queue.outBound.c <- pkt
		default:

		}
	}
}

func (n *Node) WriteToDevice() {
	for {
		select {
		case pkt := <-n.queue.inBound.c:
			if pkt.FrameType == util.MsgTypePacket {
				size, err := n.device.Write(pkt.Packet[packet.HeaderBuffSize:pkt.Size])
				if err != nil {
					return
				}
				logger.Debugf("node write %d byte to %s", size, n.device.Name())
			}

		}
	}
}

// appId is a unique identify for a node
func appId() (string, error) {
	local, err := util.NewLocal(os.O_RDWR)
	//local, err := util.GetLocalConfig()
	if err != nil {
		return "", err
	}
	l, err := local.ReadFile()
	if err != nil {
		return "", err
	}

	if l.AppId == "" {
		local.Close()
		local1, err := util.NewLocal(os.O_RDWR | os.O_CREATE)
		if err != nil {
			return "", err
		}
		defer local1.Close()
		var appId [5]byte
		if _, err := io.ReadFull(rand.Reader, appId[:]); err != nil {
			return "", errors.New("generate appId failed")
		}

		l.AppId = hex.EncodeToString(appId[:])
		err = local1.WriteFile(l)
		if err != nil {
			return "", err
		}
	}

	return l.AppId, nil
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
