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

package node

import (
	"github.com/tiptopsoft/fvpn/pkg/packet/handshake"
	"github.com/tiptopsoft/fvpn/pkg/security"
	"github.com/tiptopsoft/fvpn/pkg/tun"
	"github.com/tiptopsoft/fvpn/pkg/util"
)

type Interface interface {
	SetPeer(userId, ip string, peer *Peer) error
	GetPeer(userId, ip string) (*Peer, error)
	ListPeers(userId string) PeerMap
}

func NewCache(driver string) Interface {
	if driver == "" {
		driver = "local"
	}

	switch driver {
	case "local":
		return newLocal()
	}

	return nil
}

func CachePeers(privateKey security.NoisePrivateKey, frame *Frame, cache Interface, mode int, bind Bind, device tun.Device) (*Peer, error) {
	hpkt, err := handshake.Decode(util.HandShakeMsgTypeAck, frame.Buff)
	if err != nil {
		logger.Errorf("invalid handshake packet: %v", err)
		return nil, err
	}
	uid := frame.UidString()
	srcIP := frame.SrcIP.String()
	logger.Debugf("got remote peer: %v, pubKey: %v", srcIP, hpkt.PubKey)

	p := NewPeer(uid, srcIP, hpkt.PubKey, cache, mode, bind, device)
	p.SetIP(srcIP)
	ep := NewEndpoint(frame.RemoteAddr.String())
	p.SetEndpoint(ep)
	p.SetCodec(security.New(privateKey, hpkt.PubKey))
	err = cache.SetPeer(uid, srcIP, p)
	p.Start()

	if err != nil {
		return nil, err
	}

	return p, nil
}
