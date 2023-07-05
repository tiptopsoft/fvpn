package node

type NetworkIdFn interface {
	AddPeer(ip string, peer *Peer) error
	FindPeer(ip string) (*Peer, error)
}

var (
	_ NetworkIdFn  = (*networkId)(nil)
	_ NetManagerFn = (*NetManager)(nil)
)

type NetManagerFn interface {
	AddNetwork(cidr string, id NetworkIdFn) error
	FindNetwork(cidr string) NetworkIdFn
}

type NetManager struct {
	networks map[string]NetworkIdFn
}

func (nm *NetManager) AddNetwork(cidr string, id NetworkIdFn) error {
	nm.networks[cidr] = id
	return nil
}

func (nm *NetManager) FindNetwork(cidr string) NetworkIdFn {
	return nm.networks[cidr]
}

type networkId struct {
	cidr  string // 10.0.0.1/24
	peers map[string]*Peer
}

func NewNetworkId() NetworkIdFn {
	return &networkId{}
}

func (n *networkId) AddPeer(ip string, peer *Peer) error {
	return nil
}

func (n *networkId) FindPeer(ip string) (peer *Peer, err error) {
	peer = n.peers[ip]
	if peer == nil {
		err = ErrNotFound
	}
	return
}
