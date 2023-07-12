package node

import (
	"github.com/gin-gonic/gin"
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/util"
)

var (
	PREFIX = "api/v1/"
)

func (n *Node) HttpServer() error {
	server := gin.Default()
	server.POST(PREFIX+"join", n.joinNet())
	err := server.Run(":6663")
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) joinNet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req util.JoinRequest
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(500, util.HttpError(err.Error()))
			return
		}

		if req.Network != "" {
			err = n.netCtl.JoinNet(handler.UCTL.UserId, req.Network)
			if err != nil {
				ctx.JSON(500, util.HttpError(err.Error()))
				return
			}
		} else if req.IP != "" {
			n.netCtl.JoinIP(handler.UCTL.UserId, req.IP)
		}

		resp := &util.JoinResponse{
			IP:   n.device.IPToString(),
			Name: n.device.Name(),
		}
		ctx.JSON(200, util.HttpOK(resp))
	}
}

func leaveNet() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}
