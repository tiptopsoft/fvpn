package compress

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/handler"
)

func Middleeare() func(handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, buff []byte) error {
			return next.Handle(ctx, buff)
		})
	}
}
