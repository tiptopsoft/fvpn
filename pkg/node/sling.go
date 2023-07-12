package node

import (
	"encoding/json"
	"errors"
	"github.com/dghubble/sling"
	"github.com/topcloudz/fvpn/pkg/util"
	"net/http"
)

type Interface interface {
	ListNodes(userId string)
}

type client struct {
	sling *sling.Sling
}

func NewClient(base string) *client {
	return &client{
		sling: sling.New().Client(http.DefaultClient).Base(base),
	}
}

func (c *client) JoinNetwork(req util.JoinRequest) (*util.JoinResponse, error) {
	resp := new(util.Response)
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

	var response util.JoinResponse
	err = json.Unmarshal(buff, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *client) LeaveNetwork() error {
	return nil
}

// JoinLocalFvpn call fvpn to create device handle traffic
func (c *client) JoinLocalFvpn(req util.JoinRequest) (*util.JoinResponse, error) {
	resp := new(util.Response)
	c.sling.New().Post("/api/v1/join").BodyJSON(req).Receive(&resp, &resp)
	if resp.Code == 500 {
		return nil, errors.New(resp.Message)
	}

	jsonStr, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, errors.New("invalid result")
	}

	var result util.JoinResponse
	err = json.Unmarshal(jsonStr, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *client) Login(req LoginRequest) (*LoginResponse, error) {
	resp := new(util.Response)
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

func (c *client) Tokens(req LoginRequest) (*LoginResponse, error) {
	resp := new(util.Response)
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

func (c *client) Logout(req LoginRequest) (*LoginResponse, error) {
	resp := new(util.Response)
	c.sling.New().Post("api/v1/logout").BodyJSON(req).Receive(resp, resp)
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}
	return resp.Result.(*LoginResponse), nil
}
