package node

type Interface interface {
	SetPeer(userId, ip string, peer *Peer) error
	GetPeer(userId, ip string) (*Peer, error)
	ListPeers(userId string) PeerMap
}

func NewCache(driver string) Interface {
	if driver == "" {
		driver = "local"
	}

	switch driver {
	case "local":
		return newLocal()
	}

	return nil
}
