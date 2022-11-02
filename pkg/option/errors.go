package option

import "errors"

var (
	ErrNotImplemented = errors.New("not implement yet")
	ErrUnsupported    = errors.New("unsupported")
	ErrUnknow         = errors.New("unknown")
	ErrGetMac         = errors.New("invalid mac addr")
)
