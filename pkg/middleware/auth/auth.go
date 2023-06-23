package auth

import (
	"context"
	"errors"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/http"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/util"
)

type StarAuth struct {
	Type     int
	Username string
	Password string
	Token    string
}

type AuthHandler struct {
	packet header.Header
}

// Middleware auth handler to check user login, if not, return an error tell user to login first.
func Middleware() func(handler handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
			username, password, err := util.GetUserInfo()
			if err != nil {
				return err
			}

			client := http.New("http://211.159.225.186:443")
			req := new(http.LoginRequest)
			req.Username = username
			req.Password = password
			loginResp, err := client.Login(*req)
			if err != nil {
				return errors.New("user should login first")
			}

			if loginResp.Token == "" {
				return errors.New("token is nil, please login again")
			}

			return next.Handle(ctx, frame)
		})
	}
}
