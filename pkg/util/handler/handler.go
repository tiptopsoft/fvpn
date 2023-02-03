package handler

import (
	"context"
)

type StarExecutor interface {
	AddHandler(ctx context.Context, handler Handler) error
}

// Handler is a common handler
type Handler interface {
	Handle(ctx context.Context, udpBytes []byte) error
}

type ChainHandler struct {
	handlers []Handler
}

func NewChainHandler() ChainHandler {
	return ChainHandler{}
}

func (e ChainHandler) AddHandler(ctx context.Context, handler Handler) {
	e.handlers = append(e.handlers, handler)
}

func (e ChainHandler) Handle(ctx context.Context, udpBytes []byte) error {
	for _, h := range e.handlers {
		if err := h.Handle(ctx, udpBytes); err != nil {
			return err
		}
	}

	return nil
}
