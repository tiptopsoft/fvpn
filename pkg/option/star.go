package option

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/packet"
)

type Client interface {
	Request(ctx context.Context, packet packet.Packet)
}

type Server interface {
	Start(port int) error
	Stop() error
}
