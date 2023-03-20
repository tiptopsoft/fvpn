package handler

import "context"

type Handler interface {
	Handle(ctx context.Context, buff []byte) error
}

type HandlerFunc func(context.Context, []byte) error

func (f HandlerFunc) Handle(ctx context.Context, buff []byte) error {
	return f(ctx, buff)
}
