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

package relay

import (
	"context"
	"github.com/tiptopsoft/fvpn/pkg/node"
	"github.com/tiptopsoft/fvpn/pkg/security"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"net"
	"sync"

	"github.com/tiptopsoft/fvpn/pkg/log"
)

var (
	logger = log.Log()
)

// RegServer use as server
type RegServer struct {
	*util.RegistryCfg
	conn         *net.UDPConn
	cache        node.CacheFunc
	ws           sync.WaitGroup
	readHandler  node.Handler
	writeHandler node.Handler
	queue        struct {
		outBound *node.OutBoundQueue
		inBound  *node.InBoundQueue
	}

	key struct {
		privateKey security.NoisePrivateKey
		pubKey     security.NoisePublicKey
	}

	//every node has it's own key
	appIds map[string]string
}

func (r *RegServer) Start(address string) error {
	var err error
	r.queue.outBound = node.NewOutBoundQueue()
	r.queue.inBound = node.NewInBoundQueue()
	if r.key.privateKey, err = security.NewPrivateKey(); err != nil {
		return err
	}
	if err = r.start(address); err != nil {
		return err
	}

	r.readHandler = node.WithMiddlewares(r.serverUdpHandler(), node.Decode())
	r.writeHandler = node.WithMiddlewares(r.writeUdpHandler(), node.Encode())
	r.cache = node.NewCache()
	r.ws.Wait()
	return nil
}

// Peer register cache for net, and for user create client
func (r *RegServer) start(address string) error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.IPv4zero, Port: 4000})
	if err != nil {
		return err
	}
	r.conn = conn
	logger.Debugf("server start at: %s", address)

	//nums := runtime.NumCPU()
	//for i := 0; i < nums/2; i++ {
	r.ws.Add(1)
	go r.RoutineInBound(1)
	go r.RoutineOutBound(2)
	//}

	go r.ReadFromUdp()
	return nil
}

func (r *RegServer) PutPktToOutbound(frame *node.Frame) {
	r.queue.outBound.PutPktToOutbound(frame)
}

//func (r *RegServer) GetPktFromOutbound() *packet.Frame {
//	return r.queue.outBound.GetPktFromOutbound()
//}

func (r *RegServer) PutPktToInbound(frame *node.Frame) {
	r.queue.inBound.PutPktToInbound(frame)
}

func (r *RegServer) RoutineInBound(id int) {
	defer r.ws.Done()
	logger.Debugf("start routine %d to handle incomming udp packets", id)
	for {
		select {
		case pkt := <-r.queue.inBound.GetPktFromInbound():
			r.handleInPackets(pkt, id)
		default:

		}

	}
}

func (r *RegServer) handleInPackets(pkt *node.Frame, id int) {
	//pkt.Lock()
	defer func() {
		logger.Debugf("handing in packet success in %d routine finished", id)
		//defer pkt.Unlock()
	}()

	err := r.readHandler.Handle(pkt.Context(), pkt)
	if err != nil {
		logger.Error(err)
	}
}

func (r *RegServer) RoutineOutBound(id int) {
	logger.Debugf("start route %d to handle outgoing udp packets", id)
	for {
		select {
		case pkt := <-r.queue.outBound.GetPktFromOutbound():
			r.handleOutPackets(context.Background(), pkt, id)
		default:

		}
	}
}

func (r *RegServer) handleOutPackets(ctx context.Context, pkt *node.Frame, id int) {
	//pkt.Lock()
	defer func() {
		logger.Debugf("handing out packet success in %d routine finished", id)
	}()

	var err error
	switch pkt.FrameType {
	case util.MsgTypePacket:
		peer, err := r.cache.GetPeer(pkt.UidString(), pkt.DstIP.String())
		if err != nil || peer == nil {
			logger.Errorf("peer %v is not found", pkt.DstIP.String())
		}

		logger.Debugf("write packet to peer %v: ", peer)

		pkt.RemoteAddr = peer.GetEndpoint().DstIP()

		ctx = context.WithValue(ctx, "peer", peer)
		err = r.writeHandler.Handle(ctx, pkt)
	default:
		err = r.writeHandler.Handle(ctx, pkt)
	}

	if err != nil {
		logger.Error(err)
	}
}
