package handler

import (
	"context"
	"encoding/hex"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/packet/notify"
	"github.com/topcloudz/fvpn/pkg/packet/peer"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
	"sync"
	"time"
)

type Tun struct {
	socket     socket.Interface // underlay
	p2pSocket  sync.Map         //p2psocket
	device     map[string]*tuntap.Tuntap
	Inbound    chan *packet.Frame //used from udp
	Outbound   chan *packet.Frame //used for tun
	QueryBound chan *packet.Frame
	P2PBound   chan *cache.NodeInfo
	P2pChannel chan *P2PSocket
	cache      *cache.Cache
	tunHandler Handler
	udpHandler Handler
	NetworkId  string
	p2pNode    sync.Map
}

func NewTun(tunHandler, udpHandler Handler, s socket.Interface) *Tun {
	tun := &Tun{
		Inbound:    make(chan *packet.Frame, 10000),
		Outbound:   make(chan *packet.Frame, 10000),
		QueryBound: make(chan *packet.Frame, 10000),
		P2PBound:   make(chan *cache.NodeInfo, 10000),
		P2pChannel: make(chan *P2PSocket, 100),
		device:     make(map[string]*tuntap.Tuntap),
		cache:      cache.New(),
		tunHandler: tunHandler,
		udpHandler: udpHandler,
	}
	tun.socket = s
	return tun
}

func (t *Tun) CacheDevice(networkId string, device *tuntap.Tuntap) {
	if t.device[networkId] == nil {
		t.device[networkId] = device
	}
}

func (t *Tun) ReadFromTun(ctx context.Context, networkId string) {
	time.Sleep(1 * time.Second)
	logger.Infof("start a tun loop for networkId: %s", networkId)
	ctx = context.WithValue(ctx, "networkId", networkId)
	tun := t.device[networkId]
	ctx = context.WithValue(ctx, "tun", tun)
	if tun == nil {
		logger.Fatalf("invalid network: %s", networkId)
	}
	for {
		frame := packet.NewFrame()
		n, err := tun.Read(frame.Buff[:])
		frame.Packet = frame.Buff[:n]
		frame.Size = n
		logger.Debugf("origin packet size: %d, data: %v", n, frame.Packet[:n])
		header, err := util.GetFrameHeader(frame.Packet)
		if err != nil {
			logger.Debugf("no packet...")
			continue
		}
		ctx = context.WithValue(ctx, "header", header)
		err = t.tunHandler.Handle(ctx, frame)
		if err != nil {
			logger.Errorf("tun handle packet failed: %v", err)
		}

		ip := t.device[networkId].IP
		//find self, send self to remote to tell remote to connect to self
		nodeInfo, err := t.cache.GetNodeInfo(networkId, ip.String())
		if err != nil {
			logger.Errorf("self not register yet: %v", err)
		} else {
			frame.Self = nodeInfo
		}

		t.Outbound <- frame

	}
}

func (t *Tun) WriteToUdp() {
	for {
		pkt := <-t.Outbound
		//这里先尝试P2p, 没有P2P使用relay server
		header, err := util.GetFrameHeader(pkt.Packet[12:]) //why 12? because header length is 12.
		logger.Debugf("packet will be write to : mac: %s, ip: %s, content: %v", header.DestinationAddr, header.DestinationIP.String(), pkt.Packet)
		if err != nil {
			continue
		}

		//target
		target, err := t.cache.GetNodeInfo(pkt.NetworkId, header.DestinationIP.String())
		if err != nil {
			if err := t.AddQueryRemoteNodes(pkt.NetworkId); err != nil {
				logger.Errorf("add query task failed: %v", err)
			}
			continue
		}
		if target.NatType == option.SymmetricNAT {
			//use relay server
			logger.Debugf("use relay server to connect to: %v", target.IP.String())
			t.socket.Write(pkt.Packet[:])
		} else if target.P2P {
			logger.Debugf("use p2p to connect to: %v, remoteAddr: %v, sock: %v", target.IP, target.Addr, target.Socket)
			if err := target.Socket.WriteToUdp(pkt.Packet, target.Addr); err != nil {
				logger.Errorf("send p2p data failed. %v", err)
			}
		} else {
			//write to notify
			np := notify.NewPacket(pkt.NetworkId)
			self := pkt.Self
			np.SourceIP = self.IP
			np.Port = self.Port
			np.NatType = util.NatType
			np.NatIP = self.NatIP
			np.NatPort = self.NatPort
			np.DestAddr = header.DestinationIP
			buff, err := notify.Encode(np)
			if err != nil {
				logger.Errorf("build notify packet failed: %v", err)
			}
			logger.Debugf("send a notify packet to: %v, data: %v", header.DestinationIP.String(), buff)

			t.socket.Write(buff[:])
			//同时进行punch hole
			node, err := t.cache.GetNodeInfo(pkt.NetworkId, header.DestinationIP.String())
			if err != nil {
				logger.Errorf("node has not been query back. %v", err)
			}

			if v, ok := t.p2pNode.Load(node.IP.String()); !ok || v == nil {
				logger.Infof("add %s to p2pBound", node.IP.String())
				t.P2PBound <- node
				t.p2pNode.Store(node.IP.String(), node)
			}

			//同时通过relay server发送数据
			t.socket.Write(pkt.Packet[:])
		}
	}
}

func (t *Tun) GetSocket(mac string) socket.Interface {
	v, b := t.p2pSocket.Load(mac)
	if !b {
		return nil
	}

	return v.(socket.Interface)
}

func (t *Tun) SaveSocket(mac string, s socket.Interface) {
	t.p2pSocket.Store(mac, s)
}

func (t *Tun) ReadFromUdp() {
	logger.Infof("start a udp loop socket is: %v", t.socket)
	for {
		ctx := context.Background()
		frame := packet.NewFrame()

		n, _, err := t.socket.ReadFromUdp(frame.Buff[:])
		logger.Debugf("receive data from remote, size: %d, data: %v", n, frame.Buff[:n])
		if n < 0 || err != nil {
			logger.Errorf("got data err: %v", err)
			continue
		}
		ctx = context.WithValue(ctx, "cache", t.cache)
		err = t.udpHandler.Handle(ctx, frame)
		if err != nil {
			logger.Errorf("Read from udp failed: %v", err)
			continue
		}
		//forward packet to device
		if frame.FrameType == option.MsgTypePacket {
			t.Inbound <- frame
		}

		if frame.FrameType == option.MsgTypeNotify {
			//will connect to target by udp
			ip := frame.Target.IP
			if v, ok := t.p2pNode.Load(ip.String()); !ok || v == nil {
				logger.Infof("add %s to p2pBound", ip.String())
				t.P2PBound <- frame.Target
				t.p2pNode.Store(ip.String(), frame.Target)
			}
		}

	}

}

// WriteToDevice write to device from the queue
func (t *Tun) WriteToDevice() {
	for {
		pkt := <-t.Inbound
		device := t.device[pkt.NetworkId]
		if device == nil {
			logger.Errorf("invalid network: %s", pkt.NetworkId)
			continue
		}
		logger.Debugf("write to device data :%v", pkt.Packet[12:])
		_, err := device.Write(pkt.Packet[12:]) // start 12, because header length 12
		if err != nil {
			logger.Errorf("write to device err: %v", err)
		}
	}
}

// 添加一个networkID，查询该networkId下节点，更新cache
func (t *Tun) AddQueryRemoteNodes(networkId string) error {
	pkt := peer.NewPacket(networkId)
	buff, err := peer.Encode(pkt)
	if err != nil {
		return err
	}
	frame := packet.NewFrame()
	frame.NetworkId = networkId
	frame.Packet = buff
	frame.Buff = buff
	t.QueryBound <- frame
	return nil
}

// QueryRemoteNodes when packet from regserver, this method will be called
func (t *Tun) QueryRemoteNodes() {
	for {
		pkt := <-t.QueryBound
		t.socket.Write(pkt.Packet)
		logger.Debugf("wrote a pkt to query remote ndoes")
	}

}

func (t *Tun) PunchHole() {
	for {
		node := <-t.P2PBound
		if node.NatType == option.SymmetricNAT {
			logger.Debugf("node %v is symmetrict nat, use relay server", node)
			continue
		}

		address := node.Addr
		sock := socket.NewSocket(6061)
		logger.Infof("new socket: %v, origin socket: %v", sock, t.socket)
		err := sock.Connect(address)
		if err != nil {
			logger.Errorf("init p2p failed. address: %v, err: %v", address, err)
			continue
		}

		//open session, node-> remote addr
		hbuf, _ := header.NewHeader(option.MsgTypePunchHole, node.NetworkId)
		buff, _ := header.Encode(hbuf)
		err = sock.WriteToUdp(buff, address)
		if err != nil {
			logger.Errorf("open hole failed: %v", err)
		}

		addr := address.(*unix.SockaddrInet4)
		logger.Debugf(">>>>>>>>>>>>>>>>>>>>>punch message addr: %v natip: %v, natport: %d, ip: %v, port: %v, socket: %v", address, addr.Addr, addr.Port, node.IP, node.Port, sock)
		node.Status = true
		node.P2P = true
		node.Socket = sock
		node.Addr = address

		p2pInfo := &P2PSocket{
			Socket:   sock,
			NodeInfo: node,
		}
		t.P2pChannel <- p2pInfo
	}
}

type P2PSocket struct {
	Socket   socket.Interface
	NodeInfo *cache.NodeInfo
}

func (t *Tun) P2PSocketLoop() {
	for {
		p2pInfo := <-t.P2pChannel
		go t.p2pLoop(p2pInfo)

	}
}

func (t *Tun) p2pLoop(p2pInfo *P2PSocket) {
	for {
		sock := p2pInfo.Socket
		frame := packet.NewFrame()
		n, remoteAddr, err := sock.ReadFromUdp(frame.Buff[:])
		logger.Debugf(">>>>>>>>>>>>>>read from p2p data: %v", frame.Buff[:n])
		if err != nil {
			logger.Errorf("sock read failed: %v, remoteAddr: %v", err, remoteAddr)
		}

		frame.Packet = frame.Buff[:n]
		h, err := util.GetPacketHeader(frame.Packet)

		if err != nil {
			logger.Debugf("not invalid header: %v", string(frame.Packet))
		}

		frame.NetworkId = hex.EncodeToString(h.NetworkId[:])
		logger.Debugf(">>>>>>>>>>>>>>>>p2p header, networkId: %s", frame.NetworkId)
		t.cache.SetCache(frame.NetworkId, p2pInfo.NodeInfo.IP.String(), p2pInfo.NodeInfo)
		logger.Debugf(">>>>>>>>>>>>>>>>p2p node cached: %v, networkId: %s", p2pInfo.NodeInfo, frame.NetworkId)
		//加入inbound
		if h.Flags != option.MsgTypePunchHole {
			t.Inbound <- frame
		}
		logger.Debugf("p2p sock read %d byte, data: %v, remoteAddr: %v", n, frame.Packet[:n], remoteAddr)
	}

}
