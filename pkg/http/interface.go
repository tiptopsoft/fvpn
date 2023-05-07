package http

import (
	"fmt"
	"github.com/dghubble/sling"
	"github.com/topcloudz/fvpn/pkg/cache"
	"net/http"
)

type Interface interface {
	ListNodes(userId string) []cache.NodeInfo
}

type Client struct {
	sling *sling.Sling
}

func New(base string) *Client {
	return &Client{
		sling: sling.New().Client(http.DefaultClient).Base(base),
	}
}

type JoinRequest struct {
	SrcMac    string `json:"srcMac"`
	NetworkId string `json:"networkId"`
	Ip        string `json:"ip"`
	Mask      string `json:"mask"`
}

type JoinResponse struct {
	IP        string
	Mask      string
	NetworkId string
}

func (c *Client) JoinNetwork(userId, networkId string, req JoinRequest) (*JoinResponse, error) {
	resp := new(Response)
	_, err := c.sling.New().Post(fmt.Sprintf("/api/v1/users/user/%s/network/%s/join", userId, networkId)).BodyJSON(req).ReceiveSuccess(resp)
	if err != nil {
		return nil, err
	}
	return resp.Result.(*JoinResponse), nil
}

func (c *Client) LeaveNetwork() error {
	return nil
}

// JoinLocalFvpn call fvpn to create device handle traffic
func (c *Client) JoinLocalFvpn(req JoinRequest) error {
	resp := new(Response)
	_, err := c.sling.New().Post("/api/v1/join").BodyJSON(req).ReceiveSuccess(resp)
	return err
}
