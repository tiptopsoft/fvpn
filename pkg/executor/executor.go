package executor

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/packet"
)

type Executor interface {
	Execute(ctx context.Context, frame *packet.Frame) error
}
