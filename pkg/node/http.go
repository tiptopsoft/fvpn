package node

import (
	echo "github.com/labstack/echo/v4"
	"github.com/topcloudz/fvpn/pkg/util"
)

var (
	PREFIX = "api/v1/"
)

func (n *Node) HttpServer() error {
	//server := gin.Default()
	server := echo.New()
	server.Use(checkAuth())
	server.POST(PREFIX+"join", n.joinNet())
	return server.Start(n.cfg.HttpListenStr())
}

func (n *Node) joinNet() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		var req util.JoinRequest
		err := ctx.Bind(&req)

		if err != nil {
			return ctx.JSON(500, util.HttpError(err.Error()))
		}

		if req.NetWorkId != "" {
			err = n.netCtl.JoinNet(util.UCTL.UserId, req.NetWorkId)
			if err != nil {
				return ctx.JSON(500, util.HttpError(err.Error()))
			}
		}

		resp := &util.JoinResponse{
			IP:   n.device.IPToString(),
			Name: n.device.Name(),
		}
		return ctx.JSON(200, util.HttpOK(resp))
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
