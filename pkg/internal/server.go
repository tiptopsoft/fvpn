package internal

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/packet"
)

type Protocol int

const (
	TCP Protocol = iota
	UDP
)

type Server interface {
	Start(port int) error
	Stop() error
}

type StarFunc interface {
	AddHandler(handler ...StarHandler)
}

// StarHandler is a common handler
type StarHandler interface {
	Handle(ctx context.Context, p packet.Packet) error
}

type StarServer struct {
	Protocol Protocol
}
