// Copyright 2023 Tiptopsoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package device

import (
	"encoding/json"
	"errors"
	"github.com/dghubble/sling"
	"github.com/tiptopsoft/fvpn/pkg/util"
	"net/http"
)

//type Interface interface {
//	ListNodes(userId string)
//}

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

func NewManager(cfg *util.NodeCfg) *ClientManager {
	return &ClientManager{
		ConsoleClient: NewClient(cfg.ControlUrl()),
		LocalClient:   NewClient(cfg.HostUrl()),
	}
}

func (c *ClientManager) JoinNetwork(networkId string) (*JoinResponse, error) {
	resp := new(Response)
	//First, read the config.json to get username and password to get token
	info, err := util.GetLocalUserInfo()
	if err != nil {
		return nil, err
	}

	loginRequest := LoginRequest{
		Username: info.Username,
		Password: info.Password,
	}

	tokenResp, err := c.ConsoleClient.Tokens(loginRequest)
	if err != nil {
		return nil, err
	}

	req := new(JoinRequest)
	req.NetWorkId = networkId

	c.ConsoleClient.sling.New().Post("/api/v1/network/join").BodyJSON(req).Set("token", tokenResp.Token).Receive(resp, resp)
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

	req.CIDR = response.CIDR
	joinResp, err := c.JoinLocalFvpn(*req)
	if err != nil {
		return nil, err
	}

	return joinResp, nil
}

func (c *ClientManager) LeaveNetwork(networkId string) (*LeaveResponse, error) {
	resp := new(Response)
	//First, read the config.json to get username and password to get token
	info, err := util.GetLocalUserInfo()
	if err != nil {
		return nil, err
	}

	loginRequest := LoginRequest{
		Username: info.Username,
		Password: info.Password,
	}

	tokenResp, err := c.ConsoleClient.Tokens(loginRequest)
	if err != nil {
		return nil, err
	}

	req := new(LeaveRequest)
	req.NetWorkId = networkId

	c.ConsoleClient.sling.New().Post("/api/v1/network/join").BodyJSON(req).Set("token", tokenResp.Token).Receive(resp, resp)
	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}

	buff, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}

	var response LeaveResponse
	err = json.Unmarshal(buff, &response)
	if err != nil {
		return nil, err
	}

	req.CIDR = response.CIDR
	leaveResp, err := c.LeaveFvpnLocal(*req)
	if err != nil {
		return nil, err
	}

	return leaveResp, nil
}

// JoinLocalFvpn call fvpn to create device handle traffic
func (c *ClientManager) JoinLocalFvpn(req JoinRequest) (*JoinResponse, error) {
	resp := new(Response)
	c.LocalClient.sling.New().Post("/api/v1/join").BodyJSON(req).Receive(&resp, &resp)
	if resp.Code == 500 {
		return nil, errors.New(resp.Message)
	}

	jsonStr, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, errors.New("invalid result")
	}

	var result JoinResponse
	err = json.Unmarshal(jsonStr, &result)
	if err != nil {
		return nil, err
	}

	result.CIDR = req.CIDR
	return &result, nil
}

func (c *ClientManager) LeaveFvpnLocal(req LeaveRequest) (*LeaveResponse, error) {
	resp := new(Response)
	c.LocalClient.sling.New().Post("/api/v1/leave").BodyJSON(req).Receive(&resp, &resp)
	if resp.Code == 500 {
		return nil, errors.New(resp.Message)
	}

	jsonStr, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, errors.New("invalid result")
	}

	var result LeaveResponse
	err = json.Unmarshal(jsonStr, &result)
	if err != nil {
		return nil, err
	}

	result.CIDR = req.CIDR

	return &result, nil
}

func (c *client) Login(req LoginRequest) (*LoginResponse, error) {
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

func (c *client) Tokens(req LoginRequest) (*LoginResponse, error) {
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

func (c *client) Init(appId string) (*InitResponse, error) {
	resp := new(Response)
	//First, read the config.json to get username and password to get token
	info, err := util.GetLocalUserInfo()
	if err != nil {
		return nil, err
	}

	loginRequest := LoginRequest{
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

	var initResp InitResponse
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

// Status check node status
func (c *client) Status() (*StatusResponse, error) {
	resp := new(Response)
	_, err := c.sling.New().Get("/api/v1/status").Receive(resp, resp)
	if err != nil {
		return nil, errors.New("check status failed, please check if fvpn is running or not")
	}
	var statusResp StatusResponse
	buff, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buff, &statusResp)
	if err != nil {
		return nil, err
	}

	return &statusResp, nil
}

func (c *client) Stop() (*StopResponse, error) {
	resp := new(Response)
	_, err := c.sling.New().Post("/api/v1/stop").Receive(resp, resp)
	if err != nil {
		return nil, errors.New("stop failed, please check if fvpn is running or not")
	}
	var stopResp StopResponse
	buff, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buff, &stopResp)
	if err != nil {
		return nil, err
	}

	return &stopResp, nil
}
