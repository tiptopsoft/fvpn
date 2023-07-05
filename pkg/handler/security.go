package handler

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/util"
)

func Decode() func(Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
			if frame.FrameType == util.MsgTypePacket {

			}
			return next.Handle(ctx, frame)
		})
	}
}

// Middleware Encrypt use exchangeKey
func Encode() func(Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
			if frame.FrameType == util.MsgTypePacket {
				//key := manager.GetKey(frame.SrcAddr.IP.String())
				//logger.Debugf("executing eecode buff.")
				//newBuff, err := key.Cipher.Encode(frame.Packet)
				//if err != nil {
				//	return err
				//}
				//frame.Clear()
				//copy(frame.Packet, newBuff)
			}
			return next.Handle(ctx, frame)
		})
	}
}
