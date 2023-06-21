package util

import (
	"github.com/topcloudz/fvpn/pkg/security"
	"sync"
)

type KeyManager struct {
	lock     sync.Mutex
	NodeKeys map[string]*NodeKey //string is peer appId
}

type NodeKey struct {
	PrivateKey security.NoisePrivateKey
	PubKey     security.NoisePublicKey
	SharedKey  security.NoiseSharedKey
	Cipher     security.CipherFunc
}

func (k KeyManager) GetKey(ip string) *NodeKey {
	nodeKey := k.NodeKeys[ip]
	if nodeKey != nil {
		return k.NodeKeys[ip]
	}

	return nil
}

func (k KeyManager) SetKey(ip string, key *NodeKey) {
	k.lock.Lock()
	defer k.lock.Unlock()
	k.NodeKeys[ip] = key
}
