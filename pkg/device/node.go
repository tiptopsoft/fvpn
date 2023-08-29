// Copyright 2023 Tiptopsoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package device

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/tiptopsoft/fvpn/pkg/device/conn"
	"github.com/tiptopsoft/fvpn/pkg/log"
	"github.com/tiptopsoft/fvpn/pkg/packet"
	"github.com/tiptopsoft/fvpn/pkg/packet/register"
	"github.com/tiptopsoft/fvpn/pkg/security"
	"github.com/tiptopsoft/fvpn/pkg/tun"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"io"
	"sync"
	"time"
)

var (
	logger    = log.Log()
	limitChan = make(chan bool, 1000)
)

// Node is a dev in any os.
type Node struct {
	lock       sync.Mutex
	mode       int
	cfg        *util.NodeCfg
	privateKey security.NoisePrivateKey
	pubKey     security.NoisePublicKey
	device     tun.Device
	net        struct {
		conn conn.Interface
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
	cache      Interface
}

func (n *Node) PutPktToOutbound(pkt *Frame) {
	n.queue.outBound.c <- pkt
}

func (n *Node) PutPktToInbound(pkt *Frame) {
	n.queue.inBound.c <- pkt
}

func NewNode(iface tun.Device, conn conn.Interface, cfg *util.NodeCfg) (*Node, error) {
	n := &Node{
		device: iface,
		cache:  NewCache(cfg.Driver),
		mode:   1,
		cfg:    cfg,
	}
	n.net.conn = conn
	n.netCtl = NewNetworkManager(util.UCTL.UserId)
	privateKey, err := security.NewPrivateKey()

	n.peers.peers = make(map[security.NoisePublicKey]*Peer, 1)
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
	n.relay = n.NewPeer(util.UCTL.UserId, n.cfg.RegistryUrl(), n.privateKey.NewPubicKey(), n.cache)
	n.relay.isRelay = true
	n.relay.node = n
	n.relay.SetEndpoint(conn.NewEndpoint(n.cfg.RegistryUrl()))
	n.relay.SetMode(1)
	n.relay.Start()
	err := n.cache.Set(util.UCTL.UserId, n.relay.GetEndpoint().DstIP().IP.String(), n.relay)
	if err != nil {
		return
	}
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
	f.Peer = n.relay
	copy(f.Packet[:size], buff)
	n.PutPktToOutbound(f)
	return nil
}

func Start(cfg *util.Config) error {
	iface, err := tun.New()
	if err != nil {
		return err
	}

	//send http to get Cidr
	client := NewClient(cfg.NodeCfg.ControlUrl())
	appId, err := appId()
	if err != nil {
		return err
	}
	resp, err := client.Init(appId)
	if err != nil {
		logger.Error(err)
		return err
	}
	d, err := NewNode(iface, conn.New(cfg.NodeCfg.IPV6.Enable), cfg.NodeCfg)
	if err != nil {
		return err
	}
	err = d.device.SetIP(resp.Mask, resp.IP)
	if err != nil {
		return err
	}
	logger.Debugf("device name: %s, Cidr: %s", d.device.Name(), d.device.IPToString())

	return d.up()
}

func (n *Node) up() error {
	defer n.wg.Done()
	port, err := n.net.conn.Open(uint16(n.cfg.Listen))
	logger.Infof("fvpn started at: %d", port)
	if err != nil {
		return err
	}
	//init first
	n.initRelay()

	go n.ReadFromTun()
	go n.ReadFromUdp()
	go n.WriteToUDP()
	go n.WriteToDevice()
	go func() {
		timer := time.NewTimer(time.Second * 5)
		for {
			select {
			case <-timer.C:
				//logger.Debugf("sending list packets...")
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
		ctx = context.WithValue(ctx, "cache", n.cache)
		frame.UserId = n.userId
		frame.FrameType = util.MsgTypePacket
		size, err := n.device.Read(frame.Buff[:])
		frame.ST = time.Now()
		if err != nil {
			logger.Error(err)
			continue
		}
		ipHeader, err := util.GetIPFrameHeader(frame.Buff[:])
		if err != nil {
			logger.Error(err)
			continue
		}
		if ipHeader.DstIP.String() == n.device.Addr().String() {
			continue
		}
		logger.Debugf("node %s receive %d byte, srcIP: %v, dstIP: %v", n.device.Name(), size, ipHeader.SrcIP, ipHeader.DstIP)

		peer, err := n.cache.Get(util.UCTL.UserId, ipHeader.DstIP.String())
		if err != nil || peer == nil {
			if n.cfg.EnableRelay() {
				frame.Peer = n.relay
			} else {
				//drop
				continue
			}
		} else {
			if !peer.GetP2P() {
				frame.Peer = n.relay
			} else {
				frame.Peer = peer
			}
		}

		logger.Debugf("frame's Peer is :%v", frame.Peer.GetEndpoint().DstToString())
		frame.SrcIP = n.device.Addr()
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
		if !n.cfg.Encrypt.Enable {
			frame.Encrypt = false
		}

		err = n.tunHandler.Handle(ctx, frame)
		if err != nil {
			logger.Error(err)
			continue
		}

	}
}

func (n *Node) ReadFromUdp() {
	logger.Debugf("start thread to handle udp packet")
	defer func() {
		logger.Debugf("udp thread exited")
	}()
	for {
		limitChan <- true
		go n.udpProcess()
	}
}

func (n *Node) udpProcess() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "cache", n.cache)
	frame := NewFrame()
	size, remoteAddr, err := n.net.conn.Conn().ReadFromUDP(frame.Buff[:])
	frame.ST = time.Now()
	copy(frame.Packet[:size], frame.Buff[:size])
	if err != nil {
		logger.Error(err)
	}
	frame.Size = size
	frame.RemoteAddr = remoteAddr

	hpkt, err := util.GetPacketHeader(frame.Buff[:])
	if err != nil {
		logger.Error(err)
	}
	dataType := util.GetFrameTypeName(hpkt.Flags)
	if dataType == "" {
		//drop
		logger.Debugf("got invalid data. size: %d", size)
	}
	logger.Debugf("udp receive %d byte from %s, data type: [%v]", size, remoteAddr, dataType)

	frame.SrcIP = hpkt.SrcIP //192.168.0.1->192.168.0.2 srcIP =1, dstIP =2
	frame.DstIP = hpkt.DstIP
	frame.UserId = hpkt.UserId
	frame.FrameType = hpkt.Flags

	frame.Peer, err = n.cache.Get(frame.UidString(), frame.SrcIP.String())
	if err != nil || !frame.Peer.GetP2P() {
		frame.Peer = n.relay
	}

	if !n.cfg.Encrypt.Enable {
		frame.Encrypt = false
	}

	err = n.udpHandler.Handle(ctx, frame)
	if err != nil {
		logger.Error(err)
	}
	dt := time.Since(frame.ST)
	logger.Debugf("udp receive process finished, dataType: [%v], cost: [%v]", dataType, dt)
	<-limitChan
}

// sendListPackets send a packet list all nodes in current user
func (n *Node) sendListPackets() {
	h, _ := packet.NewHeader(util.MsgTypeQueryPeer, util.UCTL.UserId)
	hpkt, err := packet.Encode(h)
	if err != nil {
		logger.Errorf("send list packet failed %v", err)
		return
	}
	frame := NewFrame()
	frame.Peer = n.relay
	frame.DstIP = n.relay.GetEndpoint().DstIP().IP
	copy(frame.Packet, hpkt)
	frame.Size = len(hpkt)
	frame.UserId = h.UserId
	frame.FrameType = util.MsgTypeQueryPeer
	n.PutPktToOutbound(frame)
}

func (n *Node) WriteToUDP() {
	for {
		select {
		case pkt := <-n.queue.outBound.c:
			dt := time.Since(pkt.ST)
			dataType := util.GetFrameTypeName(pkt.FrameType)
			logger.Debugf("before give to Peer, data type: [%v], cost: [%v]", dataType, dt)
			//pkt.Peer.PutPktToOutbound(pkt)
			peer := pkt.Peer
			send, err := pkt.Peer.node.net.conn.Send(pkt.Packet[:pkt.Size], pkt.Peer.GetEndpoint())
			if err != nil {
				logger.Error(err)
				continue
			}
			t := time.Since(pkt.ST)
			logger.Debugf("node has send [%v] packets to %s from p2p: [%v], data type: [%v], cost: [%v]", send, peer.GetEndpoint().DstToString(), peer.p2p, util.GetFrameTypeName(pkt.FrameType), t)
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

				t := time.Since(pkt.ST)
				logger.Debugf("node write %d byte to %s, cost: [%v]", size, n.device.Name(), t)
			}

		}
	}
}

// appId is a unique identify for a node
func appId() (string, error) {
	l, err := util.GetLocalConfig()
	if err != nil {
		return "", err
	}

	if l.AppId == "" {
		var appId [5]byte
		if _, err := io.ReadFull(rand.Reader, appId[:]); err != nil {
			return "", errors.New("generate appId failed")
		}

		l.AppId = hex.EncodeToString(appId[:])
		err := util.UpdateLocalConfig(l)
		if err != nil {
			return "", err
		}
	}

	return l.AppId, nil
}

func (n *Node) NewPeer(uid, ip string, pk security.NoisePublicKey, cache Interface) *Peer {
	n.peers.lock.Lock()
	defer n.peers.lock.Unlock()
	peer, _ := cache.Get(uid, ip)
	if peer != nil {
		return peer
	}

	logger.Debugf("will create Peer for userId: %v, ip: %v", uid, ip)

	p := new(Peer)
	p.isTry.Store(true)
	p.st = time.Now()
	p.ip = ip
	p.checkCh = make(chan int, 1)
	p.sendCh = make(chan int, 1)
	p.keepaliveCh = make(chan int, 1)
	p.pubKey = pk
	p.cache = cache
	p.queue.outBound = NewOutBoundQueue()
	p.queue.inBound = NewInBoundQueue()

	cache.Set(uid, ip, p)
	logger.Debugf("created Peer for : %v, Peer: [%v]", ip, p.GetEndpoint())
	return p
}
