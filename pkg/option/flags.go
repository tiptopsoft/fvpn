package option

const (
	TAP_REGISTER = iota
	TAP_REGISTER_ACK
	TAP_MESSAGE
	TAP_BROADCAST
	TAP_UNREGISTER
)

var (
	Version    = 1
	DefaultTTL = 2000
)
