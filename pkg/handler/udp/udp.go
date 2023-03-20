package udp

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/handler"
)

type UdpHandler struct {
}

func New() handler.Handler {
	return UdpHandler{}
}

func (uh UdpHandler) Handle(ctx context.Context, buff []byte) error {
	return nil
}
