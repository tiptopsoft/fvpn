package option

const (
	TAP_REGISTER = iota
	TAP_REGISTER_ACK
	TAP_MESSAGE
	TAP_LIST_EDGE_STAR
	TAP_BROADCAST
	TAP_UNREGISTER
)

var (
	Version    uint8  = 1
	DefaultTTL uint8  = 100
	IPV4       uint16 = 0x01
	IPV6       uint16 = 0x02
)
