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
	server.POST(PREFIX+"join", n.joinNet())
	//server.POST(PREFIX+"join", n.joinNet())
	//err := server.Run(":6663")
	return server.Start(":6663")
}

func (n *Node) joinNet() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		var req util.JoinRequest
		err := ctx.Bind(&req)

		if err != nil {
			return ctx.JSON(500, util.HttpError(err.Error()))
		}

		if req.Network != "" {
			err = n.netCtl.JoinNet(util.UCTL.UserId, req.Network)
			if err != nil {
				return ctx.JSON(500, util.HttpError(err.Error()))
			}
		} else if req.IP != "" {
			n.netCtl.JoinIP(util.UCTL.UserId, req.IP)
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
