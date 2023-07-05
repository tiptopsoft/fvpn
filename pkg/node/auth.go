package node

import (
	"context"
	"errors"
	"github.com/topcloudz/fvpn/pkg/http"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/util"
	"sync"
)

type UserFunc interface {
	GetUserId() ([]byte, error)
	SetUserId(userId string) error
	SetUserInfo(username, password string) error
}

// User user username password to login, then will receive userId
type user struct {
	lock     sync.Mutex
	Username string
	Password string
	UserId   string
}

func NewUser() UserFunc {
	return &user{}
}

var (
	_    UserFunc = (*user)(nil)
	UCTL          = user{}
)

func (u *user) GetUserId() ([]byte, error) {
	return nil, nil
}

func (u *user) SetUserId(userId string) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	UCTL.UserId = userId
	return nil
}

func (u *user) SetUserInfo(username, password string) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	UCTL.Username = username
	UCTL.Password = password
	return nil
}

// Middleware auth handler to check user login, if not, return an error tell user to login first.
func checkAuth() func(handler Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *packet.Frame) error {
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
