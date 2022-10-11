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
	"errors"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"sync"
)

type local struct {
	lock  sync.Mutex
	peers map[string]PeerMap //userId: map[cidr]*Peer
}

type PeerMap map[string]*Peer

func newLocal() Interface {
	return &local{
		peers: make(map[string]PeerMap, 1),
	}
}

var (
	_ Interface = (*local)(nil)
)

func (c *local) SetPeer(userId, ip string, peer *Peer) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	peerMap := c.peers[userId]
	if peerMap == nil {
		peerMap = make(PeerMap, 1)
		c.peers[userId] = peerMap
	}
	peerMap[ip] = peer
	//print
	//for ip, p := range peerMap {
	//	logger.Debugf("========================peer in cache,ip: [%v], peer: [%v], cipher: %v ", ip, p.endpoint.DstToString(), p.cipher)
	//}
	//every add a peer will print current peers in cache
	return nil
}

func (c *local) GetPeer(userId, ip string) (*Peer, error) {
	if userId == "" {
		userId = util.UCTL.UserId
	}

	peerMap := c.peers[userId]
	peer := peerMap[ip]

	// if peer not exists use relay
	//if peer == nil {
	//	return relayPeer, nil
	//}
	if peer == nil {
		return nil, errors.New("peer is nil")
	}
	return peer, nil
}

func (c *local) ListPeers(userId string) PeerMap {
	return c.peers[userId]
}
