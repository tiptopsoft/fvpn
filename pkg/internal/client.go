package internal

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/pack"
)

type Client interface {
	Request(ctx context.Context, packet pack.Packet)
}
