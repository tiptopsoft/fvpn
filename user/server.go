package user

import (
	"github.com/gin-gonic/gin"
	"github.com/interstellar-cloud/star/option"
)

type Server struct {
	Config *option.Config
	db     *Db
}

func (s *Server) Start() error {

	db := &Db{
		Config: s.Config,
	}
	if err := db.Init(); err != nil {
		return err
	}
	s.db = db
	var engine = gin.Default()
	engine.POST("register", s.register())
	engine.GET("users", s.users())
	engine.GET("user/:id", s.getUser())
	err := engine.Run(s.Config.Listen)
	if err != nil {
		return err
	}
	return nil
}

// register register user
func (s Server) register() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var u User
		if err := ctx.ShouldBind(&u); err != nil {
			ctx.JSON(500, "invalid body")
			return
		}

		if err := u.Create(s.db); err != nil {
			ctx.JSON(500, "failed to register.")
			return
		}

		ctx.JSON(200, "success")
	}
}

// uesrs user list
func (s Server) users() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var u User

		if users, err := u.ListUser(s.db); err != nil {
			ctx.JSON(500, "failed to register.")
			return
		} else {
			ctx.JSON(200, users)
		}
	}
}

// getUser get user by user id
func (s Server) getUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var u User
		var err error
		if u, err = u.Get(s.db); err != nil {
			ctx.JSON(500, "failed to register.")
			return
		} else {
			ctx.JSON(200, u)
		}
	}
}
