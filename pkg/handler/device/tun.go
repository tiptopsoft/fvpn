package device

import (
	"context"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
)

func Handle() handler.HandlerFunc {
	return func(ctx context.Context, frame *packet.Frame) error {
		networkId := ctx.Value("networkId").(string)
		h, _ := header.NewHeader(option.MsgTypePacket, networkId)
		headerBuff, err := header.Encode(h)
		if err != nil {
			return err
		}
		idx := 0
		newPacket := make([]byte, 2048)
		idx = packet.EncodeBytes(newPacket, headerBuff, idx)
		idx = packet.EncodeBytes(newPacket, frame.Buff[:frame.Size], idx)

		frame.Packet = newPacket[:idx]
		//frame.NetworkId = networkId
		frame.FrameType = option.MsgTypePacket
		frame.Type = option.PacketFromTap

		return nil

	}
}
