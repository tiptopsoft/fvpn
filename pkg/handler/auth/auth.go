package auth

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/packet/common"
)

type StarAuth struct {
	Type     int
	Username string
	Password string
	Token    string
}

type AuthHandler struct {
	packet common.PacketHeader
}

func (ah *AuthHandler) Handle(ctx context.Context, buff []byte) error {
	return nil
}
