// Copyright 2023 TiptopSoft, Inc.
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
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/device/conn"
	"github.com/tiptopsoft/fvpn/pkg/log"
	"github.com/tiptopsoft/fvpn/pkg/pprof"
	"github.com/tiptopsoft/fvpn/pkg/security"
	"github.com/tiptopsoft/fvpn/pkg/tun"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	logger = log.Log()
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

	//peers is all peers related to this device
	peers struct {
		lock  sync.Mutex
		peers map[security.NoisePublicKey]*Peer //dst
	}

	queue struct {
		outBound     *OutBoundQueue //after encrypt
		inBound      *InBoundQueue  //after decrypt
		encryptBound *EncryptQueue
		decryptBound *DecryptQueue
	}

	pools struct {
		buffPool  *MemoryPool
		framePool *MemoryPool
	}

	netCtl     NetworkManager
	tunHandler Handler
	udpHandler Handler
	relay      *Peer
	userId     [8]byte
	cache      Interface
}

func (n *Node) PutPktToOutbound(pkt *Frame) {
	n.queue.outBound.c <- pkt
}

func (n *Node) PutPktToEncryptBound(pkt *Frame) {
	n.queue.encryptBound.c <- pkt
}

func (n *Node) PutPktToInbound(pkt *Frame) {
	n.queue.inBound.c <- pkt
}

func (n *Node) PutPktToDecryptBound(pkt *Frame) {
	n.queue.decryptBound.c <- pkt
}

func NewNode(device tun.Device, conn conn.Interface, cfg *util.NodeCfg) (*Node, error) {
	n := &Node{
		device: device,
		cache:  NewCache(cfg.Driver),
		mode:   1,
		cfg:    cfg,
	}
	n.net.conn = conn
	n.pools.buffPool, n.pools.framePool = InitPools()
	n.netCtl = NewNetworkManager(util.Info().GetUserId())
	privateKey, err := security.NewPrivateKey()

	n.peers.peers = make(map[security.NoisePublicKey]*Peer, 1)
	if err != nil {
		return nil, err
	}
	n.privateKey = privateKey
	n.pubKey = n.privateKey.NewPubicKey()
	n.queue.outBound = NewOutBoundQueue()
	n.queue.inBound = NewInBoundQueue()

	n.tunHandler = WithMiddlewares(n.tunInHandler(), AuthCheck(), n.AllowNetwork(), Encode())
	n.udpHandler = WithMiddlewares(n.udpInHandler(), AuthCheck(), Decode())

	return n, nil
}

func (n *Node) initRelay() error {
	ip, endpoint, err := getRegistryUrl(n.cfg.RegistryUrl())
	if err != nil {
		return err
	}

	n.relay = n.NewPeer(util.Info().GetUserId(), ip, n.privateKey.NewPubicKey(), n.cache)
	n.relay.isRelay = true
	n.relay.node = n
	n.relay.SetEndpoint(conn.NewEndpoint(endpoint))
	n.relay.SetMode(1)
	n.relay.Start()
	return nil
}

func getRegistryUrl(registryUrl string) (string, string, error) {
	var ip string
	var endpoint string
	if !strings.Contains(registryUrl, ":") {
		addr, err := net.ResolveIPAddr("ip4", registryUrl)
		if err != nil {
			return "", "", err
		}
		ip = addr.IP.String()
		endpoint = fmt.Sprintf("%s:%d", ip, 4000)
	} else {
		addr, err := net.ResolveUDPAddr("udp", registryUrl)
		if err != nil {
			return "", "", err
		}
		ip = addr.IP.String()
		endpoint = fmt.Sprintf("%s:%d", ip, addr.Port)
	}
	return ip, endpoint, nil
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
	port, err := n.net.conn.Open(uint16(n.cfg.Listen))
	logger.Infof("fvpn started at: %d", port)
	if err != nil {
		return err
	}
	//init first
	if err := n.initRelay(); err != nil {
		return err
	}

	go n.ReadFromTun()
	go n.ReadFromUdp()
	go n.WriteToUDP()
	go n.WriteToDevice()

	if n.cfg.PProf.Enable {
		go func() {
			pprof.Pprof()
		}()
	}

	return n.HttpServer()
}

func (n *Node) Close() error {
	close(n.queue.outBound.c)
	return nil
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
	//p.isTry.Store(true)
	p.st = time.Now()
	p.node = n

	p.ip = ip
	p.checkCh = make(chan int, 1)
	p.sendCh = make(chan int, 1)
	p.keepaliveCh = make(chan int, 1)
	p.pubKey = pk
	p.cache = cache
	logger.Debugf("created Peer for : %v, Peer: [%v]", ip, p.GetEndpoint())
	return p
}

func (n *Node) lookupPeer() {

}
