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
	AddHandler(ctx context.Context, handler ...Handler)
}

// Handler is a common handler
type Handler interface {
	Handle(ctx context.Context, p packet.Packet) error
}

type StarServer struct {
	Protocol Protocol
}
