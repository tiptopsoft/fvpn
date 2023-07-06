package node

import (
	"context"
	"errors"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/util"
)

func Decode() func(handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
			if frame.FrameType == util.MsgTypePacket {
				cache := ctx.Value("cache").(CacheFunc)
				peer, err := cache.GetPeer(handler.UCTL.UserId, frame.DstIP.String())
				if err != nil {
					return errors.New("peer not found")
				}

				frame.Packet, err = peer.GetCodec().Decode(frame.Buff[:])
				if err != nil {
					return err
				}
			}
			return next.Handle(ctx, frame)
		})
	}
}

// Encode Middleware Encrypt use exchangeKey
func Encode() func(handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
			if frame.FrameType == util.MsgTypePacket {
				cache := ctx.Value("cache").(CacheFunc)
				peer, err := cache.GetPeer(handler.UCTL.UserId, frame.DstIP.String())
				if err != nil {
					return errors.New("peer not found")
				}

				frame.Packet, err = peer.GetCodec().Encode(frame.Buff[:])
				if err != nil {
					return err
				}
			}
			return next.Handle(ctx, frame)
		})
	}
}
