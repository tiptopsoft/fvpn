package node

import (
	"context"
)

// Middleware auth handler to check user login, if not, return an error tell user to login first.
func AuthCheck() func(handler Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, frame *Frame) error {
			//username, password, err := util.GetUserInfo()
			//if err != nil {
			//	return err
			//}
			//
			//client := http.New("http://211.159.225.186:443")
			//req := new(http.LoginRequest)
			//req.Username = username
			//req.Password = password
			//loginResp, err := client.Login(*req)
			//if err != nil {
			//	return errors.New("user should login first")
			//}
			//
			//if loginResp.Token == "" {
			//	return errors.New("token is nil, please login again")
			//}

			return next.Handle(ctx, frame)
		})
	}
}
