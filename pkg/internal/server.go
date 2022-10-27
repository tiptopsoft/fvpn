package internal

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/packet"
)

type Server interface {
	Start(port int) error
}

type StarFunc interface {
	AddHandler(handler ...StarHandler)
}

// StarHandler is a common handler
type StarHandler interface {
	Handle(ctx context.Context, p packet.Packet) error
}

type StarServer struct {
}
