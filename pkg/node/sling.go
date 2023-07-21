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

type ClientManager struct {
	ConsoleClient *client
	LocalClient   *client
}

func NewClient(base string) *client {
	return &client{
		sling: sling.New().Client(http.DefaultClient).Base(base),
	}
}

func NewManager(cfg *util.ClientConfig) *ClientManager {
	return &ClientManager{
		ConsoleClient: NewClient(cfg.ControlUrl()),
		LocalClient:   NewClient(cfg.HostUrl()),
	}
}

func (c *ClientManager) JoinNetwork(networkId string) (*util.JoinResponse, error) {
	resp := new(util.Response)
	//First, read the config.json to get username and password to get token
	info, err := util.GetLocalInfo()
	if err != nil {
		return nil, err
	}

	loginRequest := util.LoginRequest{
		Username: info.Username,
		Password: info.Password,
	}

	tokenResp, err := c.ConsoleClient.Tokens(loginRequest)
	if err != nil {
		return nil, err
	}

	req := new(util.JoinRequest)
	req.NetWorkId = networkId

	c.ConsoleClient.sling.New().Post("/api/v1/network/join").BodyJSON(req).Set("token", tokenResp.Token).Receive(resp, resp)
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

	req.CIDR = response.CIDR
	joinResp, err := c.JoinLocalFvpn(*req)
	if err != nil {
		return nil, err
	}

	return joinResp, nil
}

func (c *client) LeaveNetwork() error {
	return nil
}

// JoinLocalFvpn call fvpn to create device handle traffic
func (c *ClientManager) JoinLocalFvpn(req util.JoinRequest) (*util.JoinResponse, error) {
	resp := new(util.Response)
	c.LocalClient.sling.New().Post("/api/v1/join").BodyJSON(req).Receive(&resp, &resp)
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

	result.CIDR = req.CIDR
	return &result, nil
}

func (c *client) Login(req util.LoginRequest) (*util.LoginResponse, error) {
	resp := new(util.Response)
	c.sling.New().Post("api/v1/users/login").BodyJSON(req).Receive(resp, resp)
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}

	var tokenResp util.LoginResponse
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

func (c *client) Tokens(req util.LoginRequest) (*util.LoginResponse, error) {
	resp := new(util.Response)
	c.sling.New().Post("api/v1/tokens").BodyJSON(req).Receive(resp, resp)
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}

	var tokenResp util.LoginResponse
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

func (c *client) Logout(req util.LoginRequest) (*util.LoginResponse, error) {
	resp := new(util.Response)
	c.sling.New().Post("api/v1/logout").BodyJSON(req).Receive(resp, resp)
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}
	return resp.Result.(*util.LoginResponse), nil
}

func (c *client) Init(appId string) (*util.InitResponse, error) {
	resp := new(util.Response)
	//First, read the config.json to get username and password to get token
	info, err := util.GetLocalInfo()
	if err != nil {
		return nil, err
	}

	loginRequest := util.LoginRequest{
		Username: info.Username,
		Password: info.Password,
	}

	err = util.UCTL.SetUserId(info.UserId)
	if err != nil {
		return nil, err
	}
	tokenResp, err := c.Tokens(loginRequest)
	if err != nil {
		return nil, err
	}

	c.sling.New().Post("/api/v1/network/init/"+appId).Set("token", tokenResp.Token).Receive(resp, resp)

	var initResp util.InitResponse
	buff, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buff, &initResp)
	if err != nil {
		return nil, err
	}

	return &initResp, nil
}
