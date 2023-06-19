package tunnel

import "sync"

// Manager tunnel manager
type Manager struct {
	lock      sync.Mutex
	tunnels   map[string]*Tunnel // map addr->tunnel p2p tunnels
	notifyMap map[string]bool
}

func NewManager() *Manager {
	return &Manager{
		tunnels:   make(map[string]*Tunnel, 1),
		notifyMap: make(map[string]bool, 1),
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

func (m *Manager) GetNotifyStatus(dest string) bool {
	return m.notifyMap[dest]
}

func (m *Manager) SetNotifyStatus(dest string, status bool) {
	m.notifyMap[dest] = status
}
