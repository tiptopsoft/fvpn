package node

import "sync"

type CacheFunc interface {
	SetPeer(userId, ip string, peer *Peer) error
	GetPeer(userId, ip string) (*Peer, error)
}

type cache struct {
	lock     sync.Mutex
	networks map[string]string         //key: userId value: cidr or ip
	keys     map[string]NoisePublicKey //key: ip or cidr value: NoisePublicKey
	peers    map[NoisePublicKey]*Peer  //key: NoisePublicKey value: peer
}

var (
	_ CacheFunc = (*cache)(nil)
)

func NewCache() CacheFunc {
	return &cache{
		networks: make(map[string]string, 1),
		keys:     make(map[string]NoisePublicKey, 1),
		peers:    make(map[NoisePublicKey]*Peer),
	}
}

func (c *cache) SetPeer(userId, ip string, peer *Peer) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.networks[userId] = ip
	c.keys[ip] = peer.PubKey
	c.peers[peer.PubKey] = peer
	return nil
}

func (c *cache) GetPeer(userId, ip string) (*Peer, error) {
	if userId == "" {
		userId = UCTL.UserId
	}

	key := c.keys[ip]
	peer := c.peers[key]
	return peer, nil
}
