package http

import (
	"encoding/json"
	"errors"
	"github.com/dghubble/sling"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/util"
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
	IP        string `json:"deviceIp"`
	Mask      string `json:"mask"`
	NetworkId string `json:"networkId"`
}

func (c *Client) JoinNetwork(req JoinRequest) (*JoinResponse, error) {
	resp := new(Response)
	//First, read the config.json to get username and password to get token
	username, password, err := util.GetUserInfo()
	if err != nil {
		return nil, err
	}

	loginRequest := LoginRequest{
		Username: username,
		Password: password,
	}

	tokenResp, err := c.Tokens(loginRequest)
	if err != nil {
		return nil, err
	}

	c.sling.New().Post("/api/v1/network/join").BodyJSON(req).Set("token", tokenResp.Token).Receive(resp, resp)
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}

	buff, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}

	var response JoinResponse
	err = json.Unmarshal(buff, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *Client) Login(req LoginRequest) (*LoginResponse, error) {
	resp := new(Response)
	c.sling.New().Post("api/v1/users/login").BodyJSON(req).Receive(resp, resp)
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}

	var tokenResp LoginResponse
	buff, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buff, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (c *Client) Tokens(req LoginRequest) (*LoginResponse, error) {
	resp := new(Response)
	c.sling.New().Post("api/v1/tokens").BodyJSON(req).Receive(resp, resp)
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}

	var tokenResp LoginResponse
	buff, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buff, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (c *Client) Logout(req LoginRequest) (*LoginResponse, error) {
	resp := new(Response)
	c.sling.New().Post("api/v1/logout").BodyJSON(req).Receive(resp, resp)
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}
	return resp.Result.(*LoginResponse), nil
}
