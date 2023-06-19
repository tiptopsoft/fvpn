package encrypt

import (
	"context"
	"encoding/base64"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/peer"
	"github.com/topcloudz/fvpn/pkg/security"
)

func Decode() func(handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {

			//

			return next.Handle(ctx, frame)
		})
	}
}

// Middleware Encrypt use exchangeKey
func Encode(keytext string) func(handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
			var key [peer.NosiePrivateKeySize]byte
			_, err := base64.StdEncoding.Decode(key[:], []byte(keytext))
			if err != nil {
				return err
			}

			cipher := security.NewCipher()

			newPkt, err := cipher.Encode(key[:], frame.Packet)
			if err != nil {
				return err
			}

			frame.Packet = newPkt
			return next.Handle(ctx, frame)
		})
	}
}
