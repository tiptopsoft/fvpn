package infra

import (
	"github.com/topcloudz/fvpn/pkg/middleware"
	"github.com/topcloudz/fvpn/pkg/middleware/auth"
)

// first :
func Middlewares(params ...bool) []middleware.Middleware {
	var result []middleware.Middleware
	if params[0] {
		result = append(result, auth.Middleware())
	}

	//if params[1] {
	//	result = append(result, encrypt.Middleware())
	//}

	//if params[2] {
	//	//TODO
	//}

	return result
}
