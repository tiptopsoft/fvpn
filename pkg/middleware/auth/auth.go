package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/interstellar-cloud/star/pkg/handler"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/packet"
	"io"
	"net/http"
	"time"
)

var (
	logger      = log.Log()
	defaultUser = "http://fvpn.user.internal"
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

type AuthResponse struct {
	Valid    bool
	ExpireAt time.Time
	Message  string
}

func Middleware() func(handler handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return handler.HandlerFunc(func(ctx context.Context, buff []byte) error {
			//TODO impl authcation
			//user url:
			req, _ := http.NewRequest(http.MethodPost, defaultUser, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {

			}
			defer resp.Body.Close()
			buf, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Errorf("check auth failed. %v", err)
				return errors.New("check auth failed")
			}

			var res AuthResponse
			err = json.Unmarshal(buf, res)
			if err != nil {
				return err
			}

			if !res.Valid {
				logger.Errorf("user is invalid")
				return errors.New("user is invalid")
			}
			return next.Handle(ctx, buff)
		})
	}
}
