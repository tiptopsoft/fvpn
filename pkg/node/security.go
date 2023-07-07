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
				start := packet.HeaderBuffSize
				cache := ctx.Value("cache").(CacheFunc)
				peer, err := cache.GetPeer(handler.UCTL.UserId, frame.SrcIP.String())
				if err != nil {
					return errors.New("peer not found")
				}

				logger.Debugf("data before decode: %v", frame.Buff[start:frame.Size])
				decoded, err := peer.GetCodec().Decode(frame.Buff[start:frame.Size])
				if err != nil {
					return err
				}
				frame.Clear()
				copy(frame.Packet[0:start], frame.Buff[0:start])
				copy(frame.Packet[start:], decoded)
				frame.Size = len(decoded) + packet.HeaderBuffSize
				logger.Debugf("data after decode: %v", frame.Packet[:frame.Size])
			}
			frame.Packet = frame.Buff[:frame.Size]
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

				logger.Debugf("data before encode: %v", frame.Buff[:frame.Size])
				encoded, err := peer.GetCodec().Encode(frame.Buff[:frame.Size])
				if err != nil {
					return err
				}
				frame.Clear()
				copy(frame.Packet, encoded)
				frame.Size = len(encoded)
				logger.Debugf("data after encode: %v", frame.Packet[:frame.Size])
				if err != nil {
					return err
				}
			}
			return next.Handle(ctx, frame)
		})
	}
}
