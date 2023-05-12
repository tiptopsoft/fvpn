package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/http"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"github.com/topcloudz/fvpn/pkg/util"
	"os"
	"path/filepath"
	"strings"
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
			//check login
			homedir, err := os.UserHomeDir()
			if err != nil {
				return errors.New("user not logon")
			}

			path := filepath.Join(homedir, "./fvpn/config.json")
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			decoder := json.NewDecoder(file)

			var resp option.Login
			err = decoder.Decode(&resp)
			if err != nil {
				return err
			}

			values := strings.Split(resp.Auth, ":")
			username := values[0]
			password, err := util.Base64Decode(values[1])
			if err != nil {
				return err
			}

			client := http.New("https://www.efvpn.com")
			req := new(http.LoginRequest)
			req.Username = username
			req.Password = password
			loginResp, err := client.Login(*req)
			if err != nil {
				return errors.New("user should login first")
			}

			if loginResp.Token == "" {
				return errors.New("token is nil, please relogin")
			}

			return next.Handle(ctx, frame)
		})
	}
}
