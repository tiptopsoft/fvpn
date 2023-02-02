package user

import (
	"github.com/gin-gonic/gin"
	"github.com/interstellar-cloud/star/pkg/util/option"
)

type UserServer struct {
	Config *option.Config
	db     *Db
}

func (s UserServer) Start(port int) error {

	db := &Db{
		Config: s.Config,
	}
	if err := db.Init(); err != nil {
		return err
	}
	s.db = db
	var router = gin.Default()
	router.POST("registry", s.register())
	router.GET("users", s.users())
	router.GET("user/:id", s.getUser())

	router.GET("/registry/list", s.getResource())

	err := router.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}

// register register user
func (s UserServer) register() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var u User
		if err := ctx.ShouldBind(&u); err != nil {
			ctx.JSON(500, "invalid body")
			return
		}

		if err := u.Create(s.db); err != nil {
			ctx.JSON(500, "failed to registry.")
			return
		}

		ctx.JSON(200, "success")
	}
}

// uesrs user list
func (s UserServer) users() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var u User

		if users, err := u.ListUser(s.db); err != nil {
			ctx.JSON(500, "failed to registry.")
			return
		} else {
			ctx.JSON(200, users)
		}
	}
}

// getUser get user by user id
func (s UserServer) getUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var u User
		var err error
		if u, err = u.Get(s.db); err != nil {
			ctx.JSON(500, "failed to registry.")
			return
		} else {
			ctx.JSON(200, u)
		}
	}
}

func (s UserServer) getResource() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
