package middleware

import "github.com/topcloudz/fvpn/pkg/handler"

type Middleware func(handler.Handler) handler.Handler

// Chain wrap middleware in order execute
func Chain(middlewares ...Middleware) func(handler.Handler) handler.Handler {
	return func(h handler.Handler) handler.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}

		return h
	}
}

func WithMiddlewares(handler handler.Handler, middlewares ...Middleware) handler.Handler {
	return Chain(middlewares...)(handler)
}
