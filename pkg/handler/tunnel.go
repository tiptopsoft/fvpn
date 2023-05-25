package handler

import (
	"context"
	"encoding/hex"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
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
	cache      *cache.Cache
	tunHandler Handler
	udpHandler Handler
	NetworkId  string
	p2pNode    sync.Map
}

func NewTun(tunHandler, udpHandler Handler, socket socket.Interface) *Tun {
	tun := &Tun{
		Inbound:    make(chan *packet.Frame, 10000),
		Outbound:   make(chan *packet.Frame, 10000),
		QueryBound: make(chan *packet.Frame, 10000),
		P2PBound:   make(chan *cache.NodeInfo, 10000),
		device:     make(map[string]*tuntap.Tuntap),
		cache:      cache.New(),
		tunHandler: tunHandler,
		udpHandler: udpHandler,
	}
	tun.socket = socket
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
		v, ok := t.p2pNode.Load(ip.String())
		if !ok {
			//get node base info
			info := &cache.NodeInfo{
				NatType: util.NatType,
			}
			info.IP = ip
			info.Port = 6061
			t.p2pNode.Store(ip.String(), info)
			frame.Self = info
		} else {
			frame.Self = v.(*cache.NodeInfo)
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
			if err := t.addQueryRemoteNodes(pkt.NetworkId); err != nil {
				logger.Errorf("add query task failed: %v", err)
			}
			continue
		}
		if target.NatType == option.SymmetricNAT {
			//use relay server
			logger.Debugf("use relay server to connect to: %v", target.IP.String())
			t.socket.Write(pkt.Packet[:])
		} else if target.P2P {
			logger.Debugf("use p2p to connect to: %v", target.IP.String())
			target.Socket.WriteToUdp(pkt.Packet, target.Addr)
		} else {
			//build a notifypacket
			np := notify.NewPacket(pkt.NetworkId)
			np.Addr = target.IP
			np.Port = target.Port
			np.DestAddr = header.DestinationIP
			np.NatType = util.NatType

			buff, err := notify.Encode(np)
			logger.Debugf("send a notify packet to: %v, data: %v", header.DestinationIP.String(), buff)
			if err != nil {
				logger.Errorf("build notify packet failed: %v", err)
			}

			//write to notify
			t.socket.Write(buff[:])

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

		n, err := t.socket.Read(frame.Buff[:])
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
			t.P2PBound <- frame.Self
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
func (t *Tun) addQueryRemoteNodes(networkId string) error {
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
		if node.Status {
			logger.Debugf("node %v already in queue", node)
			continue
		}
		address := node.Addr
		//p2pSocket := socket.NewSocket(6061)
		p2pSocket := socket.NewSocket(6061)
		err := p2pSocket.Connect(address)
		if err != nil {
			logger.Errorf("init p2p failed. address: %v, err: %v", address, err)
			continue
		}

		//open session, node-> remote addr
		err = p2pSocket.WriteToUdp([]byte("hello"), address)
		if err != nil {
			logger.Errorf("open hole failed: %v", err)
		}
		addr := address.(*unix.SockaddrInet4)
		logger.Infof(">>>>>>>>>>>>>>>>>>>>>punch message addr: %v natip: %v, natport: %d, ip: %v, port: %v", address, addr.Addr, addr.Port, node.IP, node.Port)
		node.Status = true
		node.P2P = true
		t.cache.SetCache(node.NetworkId, node.IP.String(), node)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		taskCh := make(chan int)
		go func() {

			for {
				frame := packet.NewFrame()
				n, remoteAddr, err := p2pSocket.ReadFromUdp(frame.Buff[:])
				logger.Debugf("read from p2p data: %v", frame.Buff[:n])
				if err != nil {
					logger.Errorf("sock read failed: %v, remoteAddr: %v", err, remoteAddr)
				}

				frame.Packet = frame.Buff[:n]
				h, err := util.GetPacketHeader(frame.Packet)
				if err != nil {
					logger.Warnf("this may be a hole msg: %v", string(frame.Packet))
					continue
				}

				frame.NetworkId = hex.EncodeToString(h.NetworkId[:])
				//加入inbound
				t.Inbound <- frame
				logger.Debugf("p2p sock read %d byte, data: %v, remoteAddr: %v", n, frame.Packet[:n], remoteAddr)
				taskCh <- 1
			}
		}()

		select {
		case <-taskCh:
			//设置为P2P
			logger.Infof("use p2p>>>>>>>>>>>>>>>>>>")
			break
		case <-ctx.Done():
			node.Status = false
			logger.Infof("p2p connect timeout")
		}

	}
}
