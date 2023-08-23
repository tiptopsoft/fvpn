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
	"time"
)

type local struct {
	lock    sync.Mutex
	peers   map[string]PeerMap //userId: map[Cidr]*Peer
	timeMap sync.Map
}

const (
	expire = 1 * time.Minute
)

type PeerMap map[string]*Value

type Value struct {
	Cidr string
	Peer *Peer
	Time time.Time
}

func newValue(cidr string, peer *Peer) *Value {
	return &Value{
		Cidr: cidr,
		Peer: peer,
		Time: time.Now(),
	}
}

func newLocal() Interface {

	local := &local{
		peers: make(map[string]PeerMap, 1),
	}
	go func() {
		local.checkExpire()
	}()
	return local
}

func (c *local) checkExpire() {
	timer := time.NewTimer(time.Minute * 2)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			for _, peerMap := range c.peers {
				for key, value := range peerMap {
					if time.Now().Sub(value.Time) > 0 {
						peerMap[key] = nil
					}
				}
			}
			timer.Reset(time.Minute * 2)
		}
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
	peerMap[ip] = newValue(ip, peer)
	//expt := time.Now().Add(expire)
	//c.timeMap.Store(ip, expire)
	//time.AfterFunc(expire, func() {
	//
	//})
	//print
	//for ip, p := range peerMap {
	//	logger.Debugf("========================Peer in cache,ip: [%v], Peer: [%v], cipher: %v ", ip, p.endpoint.DstToString(), p.cipher)
	//}
	//every add a Peer will print current peers in cache
	return nil
}

func (c *local) GetPeer(userId, ip string) (*Peer, error) {
	if userId == "" {
		userId = util.UCTL.UserId
	}

	peerMap := c.peers[userId]
	value := peerMap[ip]

	// if Peer not exists use relay
	//if Peer == nil {
	//	return relayPeer, nil
	//}
	if value == nil {
		return nil, errors.New("Peer is nil")
	}
	return value.Peer, nil
}

func (c *local) ListPeers(userId string) PeerMap {
	return c.peers[userId]
}
