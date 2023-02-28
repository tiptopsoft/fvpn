package errors

import "errors"

var (
	ErrNotImplemented = errors.New("not implement yet")
	ErrUnsupported    = errors.New("unsupported")
	ErrUnknow         = errors.New("unknown")
	ErrGetMac         = errors.New("invalid mac addr")
	ErrPacket         = errors.New("invalid packet")
	NoSuchInterface   = errors.New("route ip+net: no such network interface")
	ErrInvalieCIDR    = errors.New("invalid cidr")
)

func New(msg string) error {
	return errors.New(msg)
}
