package handler

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/packet"
)

type StarHandler interface {
	AddHandler(ctx context.Context, handler ...Handler)
}

// Handler is a common handler
type Handler interface {
	Handle(ctx context.Context, p packet.Packet) error
}
