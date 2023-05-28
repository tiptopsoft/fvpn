package handler

import (
	"context"
	"encoding/hex"
	"fmt"
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
	"log"
	"sync"
	"time"
)

type Tun struct {
	socket     socket.Interface // underlay
	relayAddr  *unix.SockaddrInet4
	p2pSocket  sync.Map //p2psocket
	device     map[string]*tuntap.Tuntap
	Inbound    chan *packet.Frame //used from udp
	Outbound   chan *packet.Frame //used for tun
	QueryBound chan *packet.Frame
	P2PBound   chan *cache.NodeInfo
	P2pChannel chan *P2PSocket
	cache      *cache.Cache
	tunHandler Handler
	udpHandler Handler
	p2pNode    sync.Map // ip target -> socket
}

func NewTun(tunHandler, udpHandler Handler, s socket.Interface, relayAddr *unix.SockaddrInet4) *Tun {
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
		relayAddr:  relayAddr,
	}
	tun.socket = s

	//timer
	t := time.NewTimer(time.Second * 30)
	go func() {
		for {
			<-t.C
			for id := range tun.device {
				err := util.SendQueryPeer(id, tun.socket)
				if err != nil {
					logger.Errorf("send query nodes failed. %v", err)
				}
				err = util.SendRegister(tun.device[id], tun.socket)
				if err != nil {
					logger.Errorf("send register failed. %v", err)
				}
			}

			t.Reset(time.Second * 30)
		}
	}()
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
		h, err := util.GetFrameHeader(frame.Packet)
		if err != nil {
			logger.Debugf("no packet...")
			continue
		}
		ctx = context.WithValue(ctx, "header", h)
		err = t.tunHandler.Handle(ctx, frame)
		if err != nil {
			logger.Errorf("tun handle packet failed: %v", err)
		}

		//find self, send self to remote to tell remote to connect to self
		nodeInfo, err := t.GetSelf(networkId)
		if err != nil {
			logger.Errorf("%v", err)
		} else {
			frame.Self = nodeInfo
		}
		t.Outbound <- frame
	}
}

// GetSelf get self node from cache
func (t *Tun) GetSelf(networkId string) (*cache.NodeInfo, error) {
	device := t.device[networkId]
	if device == nil {
		return nil, fmt.Errorf("you have not to join this network: %s", networkId)
	}

	ip := device.IP
	return t.cache.GetNodeInfo(networkId, ip.String())
}

func (t *Tun) findNode(networkId, ip string) (*cache.NodeInfo, error) {
	return t.cache.GetNodeInfo(networkId, ip)
}

func (t *Tun) WriteToUdp() {
	for {
		pkt := <-t.Outbound
		//这里先尝试P2p, 没有P2P使用relay server
		h, err := util.GetFrameHeader(pkt.Packet[12:]) //why 12? because h length is 12.
		logger.Debugf("packet will be write to : mac: %s, ip: %s, content: %v", h.DestinationAddr, h.DestinationIP.String(), pkt.Packet)
		if err != nil {
			continue
		}

		//target
		target, err := t.cache.GetNodeInfo(pkt.NetworkId, h.DestinationIP.String())
		if err != nil {
			err := util.SendQueryPeer(pkt.NetworkId, t.socket)
			if err != nil {
				logger.Errorf("%v", err)
			}
			continue
		}
		if target.NatType == option.SymmetricNAT {
			//use relay server
			logger.Debugf("use relay server to connect to: %v", target.IP.String())
			t.socket.Write(pkt.Packet[:])
		} else if target.P2P {
			logger.Debugf("use p2p to connect to: %v, remoteAddr: %v, sock: %v", target.IP, target.Addr, target.Socket)
			if _, err := target.Socket.Write(pkt.Packet); err != nil {
				logger.Errorf("send p2p data failed. %v", err)
			}
		} else {

			//同时进行punch hole
			ip := h.DestinationIP.String()
			node, err := t.cache.GetNodeInfo(pkt.NetworkId, ip)
			if err != nil {
				logger.Errorf("node has not been query back. %v", err)
			}

			if v, ok := t.p2pNode.Load(node.IP.String()); !ok || v == nil {
				go t.p2pHole(pkt, node)
				t.p2pNode.Store(node.IP.String(), node)
			}

			//同时通过relay server发送数据
			t.socket.Write(pkt.Packet[:])
		}
	}
}

func (t *Tun) p2pHole(pkt *packet.Frame, node *cache.NodeInfo) (socket.Interface, error) {
	//这里先尝试P2p, 没有P2P使用relay server
	hBuf, err := util.GetFrameHeader(pkt.Packet[12:]) //why 12? because header length is 12.
	if err != nil {
		return nil, err
	}
	//同时进行punch hole
	ip := hBuf.DestinationIP.String()
	//node, err := t.cache.GetNodeInfo(pkt.NetworkId, ip)
	if err != nil {
		logger.Errorf("node has not been query back. %v", err)
	}
	//write to notify
	np := notify.NewPacket(pkt.NetworkId)

	self := pkt.Self
	np.SourceIP = self.IP
	np.Port = self.Port
	np.NatType = util.NatType
	np.NatIP = self.NatIP
	np.NatPort = self.NatPort
	np.DestAddr = hBuf.DestinationIP
	buff, err := notify.Encode(np)
	if err != nil {
		logger.Errorf("build notify packet failed: %v", err)
	}
	logger.Debugf("send a notify packet to: %v, data: %v", ip, buff)
	t.socket.Write(buff[:])

	//punch hole
	sock := socket.NewSocket(6061)
	err = sock.Connect(node.Addr)
	if err != nil {
		return nil, err
	}

	//open session, node-> remote addr
	holePacket, _ := header.NewHeader(option.MsgTypePunchHole, node.NetworkId)
	buff, _ = header.Encode(holePacket)
	_, err = sock.Write(buff)
	if err != nil {
		logger.Errorf("open hole failed: %v", err)
		return nil, err
	}

	addr := node.Addr.(*unix.SockaddrInet4)
	logger.Debugf(">>>>>>>>>>>>>>>>>>>>>punch message addr: %v natip: %v, natport: %d, ip: %v, port: %v, socket: %v", addr, addr.Addr, addr.Port, node.IP, node.Port, node.Socket)
	go func() {
		logger.Infof("start a udp socket for p2p addr: %v", addr)
		for {
			frame := packet.NewFrame()
			n, remoteAddr, err := sock.ReadFromUdp(frame.Buff[:])
			if n < 0 || err != nil {
				break
			}
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

			logger.Debugf(">>>>>>>>>>>>>>>>p2p node cached: %v, networkId: %s", node, frame.NetworkId)
			node.P2P = true
			node.Socket = sock
			t.cache.SetCache(frame.NetworkId, node.IP.String(), node)
			//加入inbound
			if h.Flags != option.MsgTypePunchHole {
				t.Inbound <- frame
			}
			logger.Debugf("p2p sock read %d byte, data: %v, remoteAddr: %v", n, frame.Packet[:n], remoteAddr)
		}

		logger.Errorf("udp socket loop error occurd,  exit: %v", addr)
	}()

	return nil, nil
}

func (t *Tun) GetSocket(targetIP string) socket.Interface {
	v, b := t.p2pSocket.Load(targetIP)
	if !b {
		return nil
	}

	return v.(socket.Interface)
}

func (t *Tun) SaveSocket(tagetIP string, s socket.Interface) {
	t.p2pSocket.Store(tagetIP, s)
}

func (t *Tun) ReadFromUdp() {
	logger.Debugf("start a udp loop socket is: %v", t.socket)
	for {
		ctx := context.Background()
		frame := packet.NewFrame()

		n, remoteAddr, err := t.socket.ReadFromUdp(frame.Buff[:])
		if n < 0 || err != nil {
			logger.Errorf("got data err: %v", err)
			continue
		}
		logger.Debugf("receive data from remote: %v, size: %d, data: %v", remoteAddr, n, frame.Buff[:n])
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

				if frame.Self == nil {
					nodeInfo, err := t.GetSelf(frame.NetworkId)
					if err != nil {
						logger.Errorf("%v", err)
					}

					frame.Self = nodeInfo
				}

				if frame.Self != nil {
					go t.p2pHole(frame, frame.Target)
					t.p2pNode.Store(ip.String(), frame.Target)
				}
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
func (t *Tun) AddQueryRemoteNodes(networkId string) {
	pkt := peer.NewPacket(networkId)
	buff, err := peer.Encode(pkt)
	if err != nil {
		log.Printf("query data failed: %v", err)
	}
	//frame := packet.NewFrame()
	//frame.NetworkId = networkId
	//frame.Packet = buff
	//frame.Buff = buff
	_, err = t.socket.Write(buff)
	if err != nil {
		log.Printf("query data failed: %v", err)
	}
	//t.QueryBound <- frame
}

// QueryRemoteNodes when packet from regserver, this method will be called
//func (t *Tun) QueryRemoteNodes() {
//	for {
//		pkt := <-t.QueryBound
//		_, err := t.socket.Write(pkt.Packet)
//		if err != nil {
//			logger.Errorf("write failed: %v", err)
//		}
//		logger.Debugf("wrote a pkt to query remote nodes data: %v", pkt.Packet)
//	}
//
//}

func (t *Tun) PunchHole() {
	for {
		node := <-t.P2PBound
		if node.NatType == option.SymmetricNAT {
			logger.Debugf("node %v is symmetrict nat, use relay server", node)
			continue
		}

		address := node.Addr
		node.Socket = socket.NewSocket(6061)
		err := node.Socket.Connect(address)
		if err != nil {
			logger.Errorf("init p2p failed. address: %v, err: %v", address, err)
			continue
		}

		//open session, node-> remote addr
		headerBuf, _ := header.NewHeader(option.MsgTypePunchHole, node.NetworkId)
		buff, _ := header.Encode(headerBuf)
		_, err = node.Socket.Write(buff)
		if err != nil {
			logger.Errorf("open hole failed: %v", err)
			continue
		}

		addr := address.(*unix.SockaddrInet4)
		logger.Debugf(">>>>>>>>>>>>>>>>>>>>>punch message addr: %v natip: %v, natport: %d, ip: %v, port: %v, socket: %v", address, addr.Addr, addr.Port, node.IP, node.Port, node.Socket)
		node.Addr = address
		node.Status = true
		p2pInfo := &P2PSocket{
			Socket:   node.Socket,
			NodeInfo: node,
		}

		t.P2pChannel <- p2pInfo
	}
}

type P2PSocket struct {
	Socket   socket.Interface
	NodeInfo *cache.NodeInfo
}

//func (t *Tun) P2PSocketLoop() {
//	for {
//		p2pInfo := <-t.P2pChannel
//		go t.p2pLoop(p2pInfo)
//
//	}
//}

//func (t *Tun) p2pLoop(p2pInfo *P2PSocket) {
//
//	for {
//		sock := p2pInfo.Socket
//		frame := packet.NewFrame()
//		n, remoteAddr, err := sock.ReadFromUdp(frame.Buff[:])
//		if n < 0 || err != nil {
//			continue
//		}
//		logger.Debugf(">>>>>>>>>>>>>>read from p2p data: %v", frame.Buff[:n])
//		if err != nil {
//			logger.Errorf("sock read failed: %v, remoteAddr: %v", err, remoteAddr)
//		}
//
//		frame.Packet = frame.Buff[:n]
//		h, err := util.GetPacketHeader(frame.Packet)
//
//		if err != nil {
//			logger.Debugf("not invalid header: %v", string(frame.Packet))
//		}
//
//		p2pInfo.NodeInfo.P2P = true
//		frame.NetworkId = hex.EncodeToString(h.NetworkId[:])
//		logger.Debugf(">>>>>>>>>>>>>>>>p2p header, networkId: %s", frame.NetworkId)
//		//t.cache.SetCache(frame.NetworkId, p2pInfo.NodeInfo.IP.String(), p2pInfo.NodeInfo)
//		logger.Debugf(">>>>>>>>>>>>>>>>p2p node cached: %v, networkId: %s", p2pInfo.NodeInfo, frame.NetworkId)
//		//加入inbound
//		if h.Flags != option.MsgTypePunchHole {
//			t.Inbound <- frame
//		}
//		logger.Debugf("p2p sock read %d byte, data: %v, remoteAddr: %v", n, frame.Packet[:n], remoteAddr)
//	}
//
//}
