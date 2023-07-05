package node

import (
	"github.com/topcloudz/fvpn/pkg/handler"
	"sync"
)

type CacheFunc interface {
	SetPeer(userId, ip string, peer *Peer) error
	GetPeer(userId, ip string) (*Peer, error)
}

type cache struct {
	lock  sync.Mutex
	peers map[string]PeerMap //userId: map[ip]*Peer
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
	return nil
}

func (c *cache) GetPeer(userId, ip string) (*Peer, error) {
	if userId == "" {
		userId = handler.UCTL.UserId
	}

	peerMap := c.peers[userId]
	peer := peerMap[ip]
	return peer, nil
}
