package cache

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/socket"
	"golang.org/x/sys/unix"
	"net"
	"sync"
)

var (
	logger = log.Log()
)

type Cache struct {
	local map[string]*NodeInfo
}

func New() *Cache {
	m := make(map[string]*NodeInfo)
	return &Cache{local: m}
}

// NodeInfo 节点注册到registry时，应保存device ip, NATHost, NATPort
type NodeInfo struct {
	Socket    socket.Interface //natip or innerip
	NetworkId string
	Addr      unix.Sockaddr //natip , natport
	MacAddr   net.HardwareAddr
	IP        net.IP
	Port      uint16
	P2P       bool
	Status    bool // true in queue
	NatType   uint8
}

var LocalCache sync.Map

func (c *Cache) SetCache(networkId, ip string, node *NodeInfo) {
	m, b := LocalCache.Load(networkId)
	if !b {
		c.local[ip] = node
		LocalCache.Store(networkId, c)
	} else {
		s := m.(*Cache)
		s.local[ip] = node
	}
	logger.Debugf("cache %s, ip: %s", networkId, ip)
}

func (c *Cache) GetNodeInfo(networkId, ip string) (*NodeInfo, error) {
	m, b := LocalCache.Load(networkId)
	if !b {
		return nil, errors.New("not networkId " + networkId + " cached")
	}
	s := m.(*Cache)
	node := s.local[ip]
	if node == nil {
		return nil, errors.New("get Self from " + networkId + " LocalCache failed")
	}
	return node, nil
}

// ListNodesByNetworkId list all node in this networkId
func (c *Cache) ListNodesByNetworkId(networkId string) (nodes []*NodeInfo, err error) {
	m, b := LocalCache.Load(networkId)
	if !b {
		return nil, errors.New("not networkId " + networkId + " cached")
	}
	s := m.(*Cache)
	for _, node := range s.local {
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (c *Cache) GetNodes() (nodes []*NodeInfo) {
	for _, value := range c.local {
		nodes = append(nodes, value)
	}
	return
}
