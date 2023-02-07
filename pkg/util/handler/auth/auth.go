package auth

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/util/packet/common"
)

type StarAuth struct {
	Type     int
	Username string
	Password string
	Token    string
}

type AuthHandler struct {
	packet common.CommonPacket
}

func (ah *AuthHandler) Handle(ctx context.Context, udpBytes []byte) error {
	return nil
}
