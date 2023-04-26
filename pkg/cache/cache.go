package cache

import (
	"errors"
	"github.com/topcloudz/fvpn/pkg/socket"
	"golang.org/x/sys/unix"
	"net"
	"sync"
)

type Cache struct {
	local sync.Map
}

func New() *Cache {
	return &Cache{local: sync.Map{}}
}

// NodeInfo 节点注册到registry时，应保存device ip, NATHost, NATPort
type NodeInfo struct {
	Socket    socket.Interface
	NetworkId string
	Addr      unix.Sockaddr //natip , natport
	MacAddr   net.HardwareAddr
	IP        net.IP
	Port      uint16
	P2P       bool
}

var LocalCache sync.Map

func (c *Cache) SetCache(mac string, node *NodeInfo) {
	LocalCache.Store(mac, node)
}

func (c *Cache) GetNodeInfo(mac string) (*NodeInfo, error) {
	node, b := LocalCache.Load(mac)
	if !b {
		return nil, errors.New("get NodeInfo from LocalCache failed")
	}
	return node.(*NodeInfo), nil
}

func (c *Cache) GetNodes() (nodes []*NodeInfo) {
	c.local.Range(func(key, value any) bool {
		nodes = append(nodes, value.(*NodeInfo))
		return true
	})

	return
}
