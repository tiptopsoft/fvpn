package handler

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/addr"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/packet/notify"
	"github.com/topcloudz/fvpn/pkg/packet/peer"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
	"log"
	"sync"
	"time"
)

type Tun struct {
	socket        socket.Interface // underlay
	relayAddr     *unix.SockaddrInet4
	p2pSocket     sync.Map //p2psocket
	device        map[string]*tuntap.Tuntap
	Inbound       chan *packet.Frame //used from udp
	Outbound      chan *packet.Frame //used for tun
	QueryBound    chan *packet.Frame
	RegisterBound chan *packet.Frame
	P2PBound      chan *P2PNode
	//P2pChannel    chan *P2PSocket
	cache      *cache.Cache
	tunHandler Handler
	udpHandler Handler
	p2pNode    sync.Map // ip target -> socket
}

type P2PNode struct {
	NodeInfo *cache.NodeInfo
	Frame    *packet.Frame
}

func NewTun(tunHandler, udpHandler Handler, s socket.Interface, relayAddr *unix.SockaddrInet4) *Tun {
	tun := &Tun{
		Inbound:       make(chan *packet.Frame, 10000),
		Outbound:      make(chan *packet.Frame, 10000),
		QueryBound:    make(chan *packet.Frame, 10000),
		RegisterBound: make(chan *packet.Frame, 10000),
		P2PBound:      make(chan *P2PNode, 10000),
		//P2pChannel:    make(chan *P2PSocket, 100),
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
				err := tun.sendQueryPeer(id)
				if err != nil {
					logger.Errorf("send query nodes failed. %v", err)
				}
				err = tun.SendRegister(tun.device[id])
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

func (t *Tun) WriteToUdp(pkt *packet.Frame) error {
	packetHeader, err := util.GetPacketHeader(pkt.Packet[:])
	if err != nil {
		return errors.New("buff not encoded by fvpn")
	}

	logger.Debugf("pkt type: %v", packetHeader.Flags)
	if pkt.Type == option.PacketFromTap {

		frameHeader, err := util.GetFrameHeader(pkt.Packet[12:]) //why 12? because packer.Header length is 12.
		logger.Debugf("packet will be write to : mac: %s, ip: %s, content: %v", frameHeader.DestinationAddr, frameHeader.DestinationIP.String(), pkt.Packet)
		if err != nil {
			return err
		}

		//target
		target, err := t.cache.GetNodeInfo(pkt.NetworkId, frameHeader.DestinationIP.String())
		if err != nil {
			err := t.sendQueryPeer(pkt.NetworkId)
			if err != nil {
				logger.Errorf("%v", err)
			}
			return err
		}

		if target.NatType == option.SymmetricNAT {
			//use relay server
			logger.Debugf("use relay server to connect to: %v", target.IP.String())
			_, err := t.socket.Write(pkt.Packet[:])
			if err != nil {
				return err
			}
		} else if target.P2P {
			logger.Debugf("use p2p to connect to: %v, remoteAddr: %v, sock: %v", target.IP, target.Addr, target.Socket)
			if _, err := target.Socket.Write(pkt.Packet); err != nil {
				logger.Errorf("send p2p data failed. %v", err)
			}
		} else {
			//同时进行punch hole
			ip := frameHeader.DestinationIP.String()
			node, err := t.cache.GetNodeInfo(pkt.NetworkId, ip)
			if err != nil {
				logger.Errorf("node has not been query back. %v", err)
			}

			if v, ok := t.p2pNode.Load(node.IP.String()); !ok || v == nil {
				pNode := &P2PNode{
					NodeInfo: node,
					Frame:    pkt,
				}
				t.P2PBound <- pNode
				//go t.p2pHole(pkt, node)
				t.p2pNode.Store(node.IP.String(), node)
			}

			//同时通过relay server发送数据
			t.socket.Write(pkt.Packet[:])
		}
	} else {
		t.socket.Write(pkt.Packet)
	}

	return nil
}

// p2pRunner
func (t *Tun) p2pRunner(sock socket.Interface, node *cache.NodeInfo) {
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
					//go t.p2pHole(frame, frame.Target)
					pNode := &P2PNode{
						NodeInfo: frame.Target,
						Frame:    frame,
					}
					t.P2PBound <- pNode
					t.p2pNode.Store(ip.String(), frame.Target)
				}
			}
		}

	}

}

// WriteToDevice write to device from the queue
func (t *Tun) WriteToDevice(pkt *packet.Frame) {
	device := t.device[pkt.NetworkId]
	if device == nil {
		logger.Errorf("invalid network: %s", pkt.NetworkId)
	}
	logger.Debugf("write to device data :%v", pkt.Packet[12:])
	_, err := device.Write(pkt.Packet[12:]) // start 12, because header length 12
	if err != nil {
		logger.Errorf("write to device err: %v", err)
	}
}

func (t *Tun) HandleFrame() {

	for {
		select {
		case pNode := <-t.P2PBound:
			//这里先尝试P2p, 没有P2P使用relay server
			pkt := pNode.Frame
			hBuf, err := util.GetFrameHeader(pNode.Frame.Packet[12:]) //why 12? because header length is 12.
			if err != nil {
				logger.Errorf("%v", err)
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

			logger.Debugf("build a notify packet to: %v, data: %v", ip, buff)
			frame := packet.NewFrame()
			frame.Packet = buff[:]
			frame.FrameType = option.MsgTypePunchHole
			frame.Type = option.PacketFromUdp
			t.Outbound <- frame

			//punch hole
			sock := socket.NewSocket(6061)
			err = sock.Connect(pkt.Self.Addr)
			if err != nil {
				logger.Errorf("%v", err)
			}

			//open session, node-> remote addr
			holePacket, _ := header.NewHeader(option.MsgTypePunchHole, pkt.NetworkId)
			buff, _ = header.Encode(holePacket)
			_, err = sock.Write(buff)
			if err != nil {
				logger.Errorf("open hole failed: %v", err)
			}
		case pkt := <-t.QueryBound:
			t.socket.Write(pkt.Packet)
		case pkt := <-t.RegisterBound:
			t.socket.Write(pkt.Packet)
		case pkt := <-t.Inbound:
			t.WriteToDevice(pkt)
		case pkt := <-t.Outbound:
			err := t.WriteToUdp(pkt)
			if err != nil {
				logger.Errorf("write to udp failed: %v", err)
			}
		}
	}

}

// register register a node to center.
func (t *Tun) SendRegister(tun *tuntap.Tuntap) error {
	var err error
	//header, err := packet.NewHeader(option.MsgTypeRegisterSuper, tun.NetworkId)
	srcMac, srcIP, err := addr.GetMacAddrAndIPByDev(tun.Name)
	if err != nil {
		return err
	}

	if srcIP == nil {
		return errors.New("device ip not set")
	}
	regPkt := register.NewPacket(tun.NetworkId, srcMac, srcIP)
	copy(regPkt.SrcMac[:], tun.MacAddr)
	if err != nil {
		return err
	}

	data, err := register.Encode(regPkt)
	if err != nil {
		return err
	}
	frame := packet.NewFrame()
	frame.Packet = data
	frame.FrameType = option.MsgTypeRegister
	frame.FrameType = option.PacketFromUdp

	t.RegisterBound <- frame
	return nil
}

func (t *Tun) sendQueryPeer(networkId string) error {
	pkt := peer.NewPacket(networkId)
	data, err := peer.Encode(pkt)
	if err != nil {
		log.Printf("query data failed: %v", err)
	}

	frame := packet.NewFrame()
	frame.Packet = data
	frame.FrameType = option.MsgTypeQueryPeer
	frame.Type = option.PacketFromUdp
	t.QueryBound <- frame

	return nil
}

type P2PSocket struct {
	Socket   socket.Interface
	NodeInfo *cache.NodeInfo
}
