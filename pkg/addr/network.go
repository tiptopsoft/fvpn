package addr

import "fmt"

// use networkId
// use 64 bit (8byte) to explain networkId, as use 16 hex digst
// like: 8056c2e21c123456
// address: 5 bytes
// netword controllerId : 3byte
type Network struct {
	Address  []byte //5 byte
	Identity []byte // 3byte
}

func (n Network) NetworkId() string {
	fmt.Sprintf("%s%s", "", string(n.Address))
	return ""
}
