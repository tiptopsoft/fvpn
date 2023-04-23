package encrypt

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/handler"
)

func Middleeare() func(handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, buff []byte) error {
			return next.Handle(ctx, buff)
		})
	}
}
