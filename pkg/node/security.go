package node

import (
	"context"
	"errors"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/util"
)

func PeerDecode() func(Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
			if frame.FrameType == util.MsgTypePacket {
				logger.Debugf("executing decode buff.")
				//newBuff, err := cipher.Decode(frame.Packet[12:])
				//if err != nil {
				//	return err
				//}
				//copy(frame.Packet[12:], newBuff)
			}
			return next.Handle(ctx, frame)
		})
	}
}

// Middleware Encrypt use exchangeKey
func PeerEncode() func(Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
			return next.Handle(ctx, frame)
		})
	}
}

func Decode(manager *util.KeyManager) func(Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
			if frame.FrameType == util.MsgTypePacket {
				key := manager.GetKey(frame.SrcAddr.IP.String())
				if key == nil {
					return errors.New("not found cipher")
				}
				logger.Debugf("executing decode buff.")
				//newBuff, err := key.Cipher.Decode(frame.Packet[12:])
				//if err != nil {
				//	return err
				//}
				//copy(frame.Packet[12:], newBuff)
			}
			return next.Handle(ctx, frame)
		})
	}
}

// Middleware Encrypt use exchangeKey
func Encode(manager *util.KeyManager) func(Handler) Handler {
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
