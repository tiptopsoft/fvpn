package codec

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/security"
)

func Decode(cipher security.CipherFunc) func(handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {

			newBuff, err := cipher.Decode(frame.Packet[12:])
			if err != nil {
				return err
			}
			copy(frame.Packet[12:], newBuff)
			return next.Handle(ctx, frame)
		})
	}
}

// Middleware Encrypt use exchangeKey
func Encode(cipher security.CipherFunc) func(handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
			newBuff, err := cipher.Encode(frame.Packet)
			if err != nil {
				return err
			}
			frame.Clear()
			copy(frame.Packet, newBuff)
			return next.Handle(ctx, frame)
		})
	}
}
