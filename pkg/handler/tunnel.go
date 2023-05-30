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
	"github.com/topcloudz/fvpn/pkg/packet/notify/ack"
	"github.com/topcloudz/fvpn/pkg/packet/peer"
	"github.com/topcloudz/fvpn/pkg/packet/register"
	"github.com/topcloudz/fvpn/pkg/socket"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"sync"
	"time"
)

type Tun struct {
	socket        socket.Socket // underlay
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
	Socket   socket.Socket
}

func NewTun(tunHandler, udpHandler Handler, s socket.Socket, relayAddr *unix.SockaddrInet4) *Tun {
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

		////find self, send self to remote to tell remote to connect to self
		//nodeInfo, err := t.GetSelf(networkId)
		//if err != nil {
		//	logger.Errorf("%v", err)
		//} else {
		//	frame.Self = nodeInfo
		//}
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
		ip := frameHeader.DestinationIP.String()
		target, err := t.cache.GetNodeInfo(pkt.NetworkId, ip)
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
			//同时通过relay server发送数据
			t.socket.Write(pkt.Packet[:])

			//同时进行punch hole
			go t.sendNotifyMessage(pkt.NetworkId, t.relayAddr, ip, option.MsgTypeNotify)

		}
	} else {
		t.socket.Write(pkt.Packet)
	}

	return nil
}

func (t *Tun) sendNotifyMessage(networkId string, address unix.Sockaddr, ip string, flag uint16) {
	srcIP := t.device[networkId].IP
	if _, ok := t.p2pSocket.Load(ip); !ok {
		//新建一个client
		logger.Debugf("========will create a new socket for p2p connection for: %v", ip)
		newSocket := socket.NewSocket(0)
		err := newSocket.Connect(address)
		if err != nil {
			logger.Errorf("connect to registry failed:%v", err)
			return
		}

		var buff []byte
		if flag == option.MsgTypeNotify {
			pkt := notify.NewPacket(networkId)
			pkt.DestAddr = net.ParseIP(ip)
			addr, err := newSocket.LocalAddr()
			if err != nil {
				return
			}
			pkt.SourceIP = srcIP
			pkt.Port = uint16(addr.Port)
			buff, err = notify.Encode(pkt)
		}

		if flag == option.MsgTypeNotifyAck {
			pkt := ack.NewPacket(networkId)
			pkt.DestAddr = net.ParseIP(ip)
			localAddr, err := newSocket.LocalAddr()
			if err != nil {
				return
			}
			pkt.SourceIP = srcIP
			pkt.Port = uint16(localAddr.Port)
			buff, err = ack.Encode(pkt)
		}

		if err != nil {
			logger.Errorf("encode notify failed: %v", err)
			return
		}
		//发送notify message
		_, err = newSocket.Write(buff)
		if err != nil {
			logger.Errorf("write notify packet failed: %v", err)
			return
		}

		t.p2pSocket.Store(ip, newSocket)
	}

}

// p2pRunner
func (t *Tun) p2pRunner(sock socket.Socket, node *cache.NodeInfo) {
	for {
		logger.Debugf("start a p2p runner......")
		frame := packet.NewFrame()
		n, err := sock.Read(frame.Buff[:])
		if n < 0 || err != nil {
			break
		}

		logger.Debugf(">>>>>>>>>>>>>>read from p2p data: %v", frame.Buff[:n])
		if err != nil {
			logger.Errorf("sock read failed: %v", err)
		}

		frame.Packet = frame.Buff[:n]
		h, err := util.GetPacketHeader(frame.Packet)

		if err != nil {
			logger.Debugf("not invalid header: %v", string(frame.Packet))
		}

		frame.NetworkId = hex.EncodeToString(h.NetworkId[:])

		logger.Debugf(">>>>>>>>>>>>>>>>p2p header, networkId: %s", frame.NetworkId)

		//加入inbound
		if h.Flags != option.MsgTypePunchHole {
			t.Inbound <- frame
		}
		logger.Debugf("p2p sock read %d byte, data: %v", n, frame.Packet[:n])
	}
}

func (t *Tun) GetSocket(targetIP string) socket.Socket {
	v, b := t.p2pSocket.Load(targetIP)
	if !b {
		return socket.Socket{}
	}

	return v.(socket.Socket)
}

func (t *Tun) SaveSocket(target string, s socket.Socket) {
	t.p2pSocket.Store(target, s)
}

func (t *Tun) ReadFromUdp() {
	logger.Debugf("start a udp loop socket is: %v", t.socket)

	for {
		if !t.socket.Run {
			logger.Debugf("exit origin reader from udp")
			break
		}
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
		switch frame.FrameType {
		case option.MsgTypePacket:
			t.Inbound <- frame
		case option.MsgTypeNotify:

			//send a ack,  open hole
			ch := make(chan int)
			go func() {
				t.sendNotifyMessage(frame.NetworkId, t.relayAddr, frame.Target.IP.String(), option.MsgTypeNotifyAck)
				ch <- 1
			}()

			//also punch hole
			go func() {
				<-ch
				t.bidirectionHole(frame, err)
			}()

		case option.MsgTypeNotifyAck:
			//go t.p2pHole(frame, frame.Target)
			if t.bidirectionHole(frame, err) {
				return
			}
		}

	}

}

func (t *Tun) bidirectionHole(frame *packet.Frame, err error) bool {
	ip := frame.Target.IP.String()
	s, ok := t.p2pSocket.Load(ip)
	if !ok {
		logger.Errorf("notify ack get socket failed: %v", err)
		return true
	}
	sock := s.(socket.Socket)
	pNode := &P2PNode{
		NodeInfo: frame.Target,
		Frame:    frame,
		Socket:   sock,
	}
	t.P2PBound <- pNode
	return false
}

// WriteToDevice write to device from the queue
func (t *Tun) WriteToDevice(pkt *packet.Frame) error {
	device := t.device[pkt.NetworkId]
	if device == nil {
		return errors.New("invalid network: " + pkt.NetworkId)
	}
	_, err := device.Write(pkt.Packet[12:]) // start 12, because header length 12
	if err != nil {
		return err
	}

	return nil
}

func (t *Tun) HandleFrame() {

	for {
		select {
		case pNode := <-t.P2PBound:
			//这里先尝试P2p, 没有P2P使用relay server
			var err error
			pkt := pNode.Frame
			//同时进行punch hole
			sock := pNode.Socket
			//open session, node-> remote addr
			holePacket, _ := header.NewHeader(option.MsgTypePunchHole, pkt.NetworkId)
			buff, _ := header.Encode(holePacket)
			logger.Debugf(">>>>>>> punching hole to: %v, socket is: %v", pNode.NodeInfo.Addr, sock)
			if err != nil {
				logger.Errorf("open hole failed: %v", err)
			}

			_, err = sock.Write(buff)
			if err != nil {
				logger.Errorf("send punch hole failed: %v", err)
				return
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
					go t.p2pRunner(sock, pNode.NodeInfo)
					ch <- 1
				}
			}()

			select {
			case v := <-ch:
				if v == 1 {
					pNode.NodeInfo.P2P = true
					pNode.NodeInfo.Socket = sock
					t.cache.SetCache(pkt.NetworkId, pNode.NodeInfo.IP.String(), pNode.NodeInfo)
					logger.Debugf("punch hole success")
				} else {
					logger.Debugf("punch hole failed.")
				}
			case <-timeout.Done():
				logger.Debugf("punch hole failed.")
			}

		case pkt := <-t.QueryBound:
			t.socket.Write(pkt.Packet)
		case pkt := <-t.RegisterBound:
			t.socket.Write(pkt.Packet)
		case pkt := <-t.Inbound:
			err := t.WriteToDevice(pkt)
			if err != nil {
				logger.Errorf("write to device failed: %v", err)
			}
			logger.Debugf("write to device data :%v", pkt.Packet[12:])
		case pkt := <-t.Outbound:
			err := t.WriteToUdp(pkt)
			if err != nil {
				logger.Errorf("write to udp failed: %v", err)
			}
		}
	}

}

// SendRegister register register a node to center.
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

	t.Outbound <- frame
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
	t.Outbound <- frame

	return nil
}
