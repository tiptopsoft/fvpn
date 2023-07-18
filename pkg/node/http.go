package node

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/topcloudz/fvpn/pkg/util"
	"net/http"
)

var (
	PREFIX = "/api/v1/"
)

func (n *Node) HttpServer() error {
	server := gin.Default()
	//server := echo.New()
	//server.Use(checkAuth())
	//server.GET("/", hello)
	server.POST(PREFIX+"join", n.joinNet())

	return server.Run(n.cfg.HttpListenStr())
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (n *Node) joinNet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req util.JoinRequest
		err := ctx.Bind(&req)

		if err != nil {
			ctx.JSON(500, util.HttpError(err.Error()))
			return
		}

		if req.CIDR != "" {
			err = n.netCtl.JoinNet(util.UCTL.UserId, req.CIDR)
			if err != nil {
				ctx.JSON(500, util.HttpError(err.Error()))
				return
			}
		} else {
			ctx.JSON(500, util.HttpError("cidr is nil"))
			return
		}

		resp := &util.JoinResponse{
			IP:   n.device.IPToString(),
			Name: n.device.Name(),
		}
		ctx.JSON(200, util.HttpOK(resp))
	}
}

func leaveNet() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		return nil
	}
}

func checkAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return nil
		}
	}
}
