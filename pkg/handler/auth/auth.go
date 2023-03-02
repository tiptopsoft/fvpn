package auth

import (
	"context"
	"github.com/interstellar-cloud/star/pkg/packet"
)

type StarAuth struct {
	Type     int
	Username string
	Password string
	Token    string
}

type AuthHandler struct {
	packet packet.Header
}

func (ah *AuthHandler) Handle(ctx context.Context, buff []byte) error {
	return nil
}
