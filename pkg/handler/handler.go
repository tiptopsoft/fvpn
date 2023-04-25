package handler

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/packet"
)

var (
	logger = log.Log()
)

type Handler interface {
	Handle(ctx context.Context, frame *packet.Frame) error
}

type HandlerFunc func(context.Context, *packet.Frame) error

func (f HandlerFunc) Handle(ctx context.Context, frame *packet.Frame) error {
	return f(ctx, frame)
}
