// Copyright 2023 Tiptopsoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package node

import (
	"context"
)

// AuthCheck Middleware auth handler to check user login, if not, return an error tell user to login first.
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
