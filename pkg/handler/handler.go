package handler

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/packet"
)

type Handler interface {
	Handle(ctx context.Context, frame *packet.Frame) error
}

type HandlerFunc func(context.Context, *packet.Frame) error

func (f HandlerFunc) Handle(ctx context.Context, frame *packet.Frame) error {
	return f(ctx, frame)
}

type Middleware func(Handler) Handler

// Chain wrap middleware in order execute
func Chain(middlewares ...Middleware) func(Handler) Handler {
	return func(h Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}

		return h
	}
}

func WithMiddlewares(handler Handler, middlewares ...Middleware) Handler {
	return Chain(middlewares...)(handler)
}
