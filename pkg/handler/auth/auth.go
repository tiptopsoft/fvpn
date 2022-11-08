package auth

import "context"

type StarAuth struct {
	Type     int
	Username string
	Password string
	Token    string
}

type AuthHandler struct {
}

func (ah *AuthHandler) Handle(ctx context.Context) error {
	return nil
}
