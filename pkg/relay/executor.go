package relay

//
//import (
//	"context"
//	"encoding/hex"
//	"github.com/topcloudz/fvpn/pkg/handler"
//	"github.com/topcloudz/fvpn/pkg/option"
//	"github.com/topcloudz/fvpn/pkg/packet"
//	"github.com/topcloudz/fvpn/pkg/packet/handshake"
//	"github.com/topcloudz/fvpn/pkg/packet/header"
//	"github.com/topcloudz/fvpn/pkg/packet/notify"
//	notifyack "github.com/topcloudz/fvpn/pkg/packet/notify/ack"
//	"github.com/topcloudz/fvpn/pkg/packet/register"
//	"github.com/topcloudz/fvpn/pkg/security"
//	"github.com/topcloudz/fvpn/pkg/util"
//)
//
//func (r *RegServer) ReadFromUdp() {
//	logger.Infof("start a udp loop")
//	for {
//		ctx := context.Background()
//		frame := packet.NewFrame()
//		n, addr, err := r.socket.ReadFromUDP(frame.Buff[:])
//		frame.Packet = frame.Buff[:n]
//		logger.Debugf("Read from udp %d byte, data: %v", n, frame.Packet)
//
//		packetHeader, err := util.GetPacketHeader(frame.Packet[:12])
//		if err != nil {
//			logger.Errorf("get header falied. %v", err)
//			continue
//		}
//		networkId := hex.EncodeToString(packetHeader.NetworkId[:])
//		frame.Size = n
//		frame.FrameType = packetHeader.Flags
//		frame.SrcAddr = addr
//		//frame.PubKey = hex.EncodeToString(packetHeader.PubKey[:])
//		frame.NetworkId = networkId
//		if err != nil || n < 0 {
//			logger.Warnf("no data exists")
//			continue
//		}
//		err = r.readHandler.Handle(ctx, frame)
//		if err != nil {
//			logger.Errorf(err.Error())
//			continue
//		}
//		r.Outbound <- frame
//	}
//}
//
//func (r *RegServer) WriteToUdp() {
//	logger.Infof("start a udp write loop")
//	for {
//		pkt := <-r.Outbound
//		ctx := context.Background()
//		r.writeHandler.Handle(ctx, pkt)
//	}
//}
//
//func (r *RegServer) writeUdpHandler() handler.HandlerFunc {
//	return func(ctx context.Context, pkt *packet.Frame) error {
//		frameType := pkt.FrameType
//		switch frameType {
//		case option.MsgTypePacket:
//			frameHeader, err := util.GetFrameHeader(pkt.Packet[12:]) //why is 12, because we add our header in, header length is 12
//			if err != nil {
//				logger.Debugf("get header failed, dest ip: %s", frameHeader.DestinationIP.String())
//			}
//			//
//			nodeInfo, err := r.cache.GetNodeInfo(pkt.NetworkId, frameHeader.DestinationIP.String())
//			if nodeInfo == nil || err != nil {
//				logger.Debugf("could not found destitation, destIP: %s", frameHeader.DestinationIP.String())
//			} else {
//				logger.Infof("packet will relay to: %v", nodeInfo.Addr)
//				r.socket.WriteToUDP(pkt.Packet[:], nodeInfo.Addr)
//			}
//
//			break
//		case option.MsgTypeRegisterAck:
//			r.socket.WriteToUDP(pkt.Packet, pkt.SrcAddr)
//			break
//		case option.MsgTypeQueryPeer:
//			logger.Debugf("query nodes result: %v, write to: %v", pkt.Packet, pkt.SrcAddr)
//			_, err := r.socket.WriteToUDP(pkt.Packet, pkt.SrcAddr)
//			if err != nil {
//				logger.Errorf("write query to dest failed: %v", err)
//			}
//			break
//		case option.MsgTypeNotify:
//			//write to dest
//			np, err := notify.Decode(pkt.Packet[:])
//			logger.Debugf("got notify packet: %v, destAddr: %s, networkId: %s", pkt.Packet[:], np.DestAddr.String(), pkt.NetworkId)
//			if err != nil {
//				logger.Errorf("invalid notify packet: %v", err)
//			}
//
//			nodeInfo, err := r.cache.GetNodeInfo(pkt.NetworkId, np.DestAddr.String())
//			if nodeInfo == nil || err != nil {
//				logger.Errorf("node not on line, err: %v", err)
//				break
//			}
//
//			r.socket.WriteToUDP(pkt.Packet, nodeInfo.Addr)
//
//		case option.MsgTypeNotifyAck:
//			//write to dest
//			np, err := notifyack.Decode(pkt.Packet[:])
//			logger.Debugf("got notify ack packet: %v, destAddr: %s, networkId: %s", pkt.Packet[:], np.DestAddr.String(), pkt.NetworkId)
//			if err != nil {
//				logger.Errorf("invalid notify ack packet: %v", err)
//			}
//
//			nodeInfo, err := r.cache.GetNodeInfo(pkt.NetworkId, np.DestAddr.String())
//			if nodeInfo == nil || err != nil {
//				logger.Errorf("node not on line, err: %v", err)
//				break
//			}
//
//			r.socket.WriteToUDP(pkt.Packet, nodeInfo.Addr)
//
//		case option.HandShakeMsgType:
//			key := r.manager.GetKey(pkt.SrcAddr.IP.String())
//			handPkt := handshake.NewPacket("")
//			handPkt.PubKey = key.PubKey
//			buff, err := handshake.Encode(handPkt)
//			if err != nil {
//				logger.Errorf("invalid handshake packet")
//				return err
//			}
//			r.socket.WriteToUDP(buff, pkt.SrcAddr)
//
//		}
//		return nil
//	}
//}
//
//// serverUdpHandler  core self handler
//func (r *RegServer) serverUdpHandler() handler.HandlerFunc {
//	return func(ctx context.Context, frame *packet.Frame) error {
//		data := frame.Packet[:]
//		switch frame.FrameType {
//
//		case option.MsgTypeRegisterSuper:
//			regPkt, err := register.Decode(frame.Packet)
//			if err != nil {
//				logger.Errorf("register failed, err:%v", err)
//				return err
//			}
//			err = r.registerAck(frame.SrcAddr, regPkt.SrcMac, regPkt.SrcIP, frame.NetworkId)
//			h, err := header.NewHeader(option.MsgTypeRegisterAck, frame.NetworkId)
//			if err != nil {
//				logger.Errorf("build resp failed. err: %v", err)
//			}
//			f, _ := header.Encode(h)
//			frame.Packet = f
//			break
//		case option.MsgTypeQueryPeer:
//			peers, size, err := getPeerInfo(r.cache.GetNodes())
//			logger.Debugf("server peers: (%v), size: (%v)", peers, size)
//			if err != nil {
//				logger.Errorf("get peers from server failed. err: %v", err)
//			}
//
//			f, err := peerAckBuild(peers, size, frame.NetworkId)
//			if err != nil {
//				logger.Errorf("get peer ack from server failed. err: %v", err)
//			}
//
//			frame.Packet = f
//			break
//		case option.MsgTypePacket:
//			logger.Infof("server got forward packet size:%d, data: %v", frame.Size, data)
//		case option.MsgTypeNotify:
//			logger.Debugf("notify frame packet: %v", frame.Packet[:])
//		case option.MsgTypeNotifyAck:
//			logger.Debugf("notify ack frame packet: %v", frame.Packet[:])
//		case option.HandShakeMsgType:
//			handPkt, err := handshake.Decode(frame.Packet)
//			if err != nil {
//				logger.Errorf("invalid handshake packet: %v", err)
//				return err
//			}
//			//frame.PubKey = hex.EncodeToString(handPkt.PubKey[:])
//			key := r.manager.GetKey(frame.SrcAddr.IP.String())
//			if key == nil {
//
//				privateKey, err := security.NewPrivateKey()
//				if err != nil {
//					return err
//				}
//				pubKey := privateKey.NewPubicKey()
//				nodeKey := &util.NodeKey{
//					PrivateKey: privateKey,
//					PubKey:     pubKey,
//				}
//				nodeKey.Cipher = security.NewCipher(privateKey, handPkt.PubKey)
//				r.manager.SetKey(frame.SrcAddr.IP.String(), nodeKey)
//			}
//		}
//
//		return nil
//	}
//}
