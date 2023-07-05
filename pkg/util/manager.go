package util

import (
	"sync"
)

type KeyManager struct {
	lock     sync.Mutex
	NodeKeys map[string]*NodeKey //string is peer appId
}

type NodeKey struct {
	AppId string
}

func (k KeyManager) GetKey(appId string) *NodeKey {
	nodeKey := k.NodeKeys[appId]
	if nodeKey != nil {
		return k.NodeKeys[appId]
	}

	return nil
}

func (k KeyManager) SetKey(appId string, key *NodeKey) {
	k.lock.Lock()
	defer k.lock.Unlock()
	k.NodeKeys[appId] = key
}
