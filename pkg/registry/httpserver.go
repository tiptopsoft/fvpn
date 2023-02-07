package registry

import (
	"github.com/gin-gonic/gin"
	"net"
)

type RegistryServer struct {
	*RegStar
}

func (s RegistryServer) Start() error {

	var router = gin.Default()
	router.GET("registry/list", s.list())
	router.GET("registry/addr/list", s.addrList())

	err := router.Run(s.HttpListen)
	if err != nil {
		return err
	}
	return nil
}

func (s RegistryServer) list() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, s.Peers)
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
