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
				offset := packet.HeaderBuffSize
				buff := frame.Packet[offset:frame.Size]
				cache := ctx.Value("cache").(CacheFunc)
				peer, err := cache.GetPeer(handler.UCTL.UserId, frame.SrcIP.String())
				if err != nil {
					return errors.New("peer not found")
				}

				logger.Debugf("data before decode: %v", buff)
				decoded, err := peer.GetCodec().Decode(buff)
				if err != nil {
					return err
				}
				frame.Clear()
				copy(frame.Packet[0:offset], frame.Buff[0:offset])
				copy(frame.Packet[offset:], decoded)
				frame.Size = len(decoded) + offset
				logger.Debugf("data after decode: %v", frame.Packet[:frame.Size])
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
				offset := packet.HeaderBuffSize
				buff := frame.Packet[offset:frame.Size]
				cache := ctx.Value("cache").(CacheFunc)
				peer, err := cache.GetPeer(handler.UCTL.UserId, frame.DstIP.String())
				if err != nil {
					return errors.New("peer not found")
				}

				logger.Debugf("data before encode: %v", buff)
				encoded, err := peer.GetCodec().Encode(buff)
				if err != nil {
					return err
				}
				copy(frame.Packet[offset:], encoded)
				frame.Size = offset + len(encoded)
				logger.Debugf("data after encode: %v", frame.Packet[:frame.Size])
				if err != nil {
					return err
				}
			}
			return next.Handle(ctx, frame)
		})
	}
}
