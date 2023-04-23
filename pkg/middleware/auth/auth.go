package auth

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/packet"
)

type StarAuth struct {
	Type     int
	Username string
	Password string
	Token    string
}

type AuthHandler struct {
	packet packet.Header
}

func Middleware() func(handler handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, buff []byte) error {
			return next.Handle(ctx, buff)
		})
	}
}
