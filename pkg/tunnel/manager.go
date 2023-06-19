package tunnel

import (
	"sync"
)

var l sync.Mutex

// Manager tunnel manager
type Manager struct {
	lock      sync.Mutex
	tunnels   map[string]*Tunnel // map addr->tunnel p2p tunnels
	notifyMap map[string]*PortPair
}

func NewManager() *Manager {
	return &Manager{
		tunnels:   make(map[string]*Tunnel, 1),
		notifyMap: make(map[string]*PortPair, 1),
	}
}

func (m *Manager) GetTunnel(dest string) *Tunnel {
	return m.tunnels[dest]
}

func (m *Manager) SetTunnel(dest string, t *Tunnel) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.tunnels[dest] = t
}

func (m *Manager) GetNotifyPortPair(dest string) *PortPair {
	return m.notifyMap[dest]
}

func (m *Manager) SetNotifyPortPair(dest string, pkt *PortPair) {
	l.Lock()
	if m.notifyMap[dest] != nil {
		return
	}
	m.notifyMap[dest] = pkt
	defer l.Unlock()
}
