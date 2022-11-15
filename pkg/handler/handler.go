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

type Executor struct {
	handlers []Handler
}

func NewExecutor() Executor {
	return Executor{}
}

func (e Executor) AddHandler(ctx context.Context, handler Handler) {
	e.handlers = append(e.handlers, handler)
}

func (e Executor) Execute(ctx context.Context, udpBytes []byte) error {
	for _, h := range e.handlers {
		if err := h.Handle(ctx, udpBytes); err != nil {
			return err
		}
	}

	return nil
}
