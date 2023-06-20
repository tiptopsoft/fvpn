package infra

import (
	"github.com/topcloudz/fvpn/pkg/middleware"
	"github.com/topcloudz/fvpn/pkg/middleware/auth"
	"github.com/topcloudz/fvpn/pkg/middleware/codec"
	"github.com/topcloudz/fvpn/pkg/security"
)

// first :
func Middlewares(cipher security.CipherFunc) []middleware.Middleware {
	var result []middleware.Middleware
	result = append(result, auth.Middleware())
	result = append(result, codec.Decode(cipher))
	return result
}
