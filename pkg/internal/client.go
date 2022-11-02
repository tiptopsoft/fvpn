package internal

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/packet"
)

type Client interface {
	Request(ctx context.Context, packet packet.Frame)
}
