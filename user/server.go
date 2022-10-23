package user

import (
	"github.com/gin-gonic/gin"
	"github.com/interstellar-cloud/star/option"
)

type Server struct {
	Config *option.Config
}

func (s *Server) Start() error {

	var engine = gin.Default()
	engine.POST("register", register())
	engine.POST("users", users())
	engine.GET("user/:id", getUser())
	err := engine.Run(s.Config.Listen)
	if err != nil {
		return err
	}
	return nil
}

// register register user
func register() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// uesrs user list
func users() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// getUser get user by user id
func getUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
