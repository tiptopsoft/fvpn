package node

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/util"
	"sync"
)

type CacheFunc interface {
	SetPeer(userId, ip string, peer *Peer) error
	GetPeer(userId, ip string) (*Peer, error)
	ListPeers(userId string) PeerMap
}

type cache struct {
	lock  sync.Mutex
	peers map[string]PeerMap //userId: map[cidr]*Peer
}

var (
	_ CacheFunc = (*cache)(nil)
)

type PeerMap map[string]*Peer

func NewCache() CacheFunc {
	return &cache{
		peers: make(map[string]PeerMap, 1),
	}
}

func (c *cache) SetPeer(userId, ip string, peer *Peer) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	peerMap := c.peers[userId]
	if peerMap == nil {
		peerMap = make(PeerMap, 1)
		c.peers[userId] = peerMap
	}
	peerMap[ip] = peer
	//print
	for ip, p := range peerMap {
		logger.Debugf("========================peer in cache,ip: [%v], peer: [%v]", ip, p.endpoint.DstToString())
	}
	//every add a peer will print current peers in cache
	return nil
}

func (c *cache) GetPeer(userId, ip string) (*Peer, error) {
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

func (c *cache) ListPeers(userId string) PeerMap {
	return c.peers[userId]
}
