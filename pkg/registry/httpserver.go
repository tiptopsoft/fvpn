package registry

import (
	"github.com/gin-gonic/gin"
	"github.com/interstellar-cloud/star/pkg/util/option"
	"net"
)

type RegistryServer struct {
	*option.RegConfig
}

func (s RegistryServer) Start(addr string) error {

	var router = gin.Default()
	router.GET("registry/list", s.list())
	router.GET("registry/addr/list", s.addrList())

	err := router.Run(addr)
	if err != nil {
		return err
	}
	return nil
}

func (s RegistryServer) list() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp = make(map[string]string)
		m.Range(func(key, value any) bool {
			resp[key.(string)] = value.(*net.UDPAddr).String()
			return true
		})

		ctx.JSON(200, resp)
	}
}

func (s RegistryServer) addrList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var resp = make(map[string]string)
		socketMap.Range(func(key, value any) bool {
			resp[key.(string)] = value.(*net.UDPAddr).String()
			return true
		})

		ctx.JSON(200, resp)
	}
}
