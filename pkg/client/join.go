package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	httputil "github.com/topcloudz/fvpn/pkg/nativehttp"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"io"
	"net/http"
)

const (
	userUrl = "https://www.efvpn.com"
)

func (n *Node) RunJoinNetwork(netId string) error {
	logger.Infof("start to join %s", netId)
	//user nativehttp to get networkId config
	type body struct {
		userId    string
		networkId string
	}

	request := body{userId: "1"}
	buff, err := json.Marshal(request)
	if err != nil {
		return errors.New("invalid body")
	}
	destUrl := fmt.Sprintf("%s/api/v1/users/user/%s/network/%s/join", userUrl, request.userId, netId)
	resp, err := http.Post(destUrl, "application/json", bytes.NewReader(buff))
	respBuff, err := io.ReadAll(resp.Body)
	var networkResp struct {
		NetworkId string
		Ip        string
		Mask      string
	}

	fmt.Println(string(respBuff))

	var httpResponse httputil.HttpResponse
	err = json.Unmarshal(respBuff, &httpResponse)
	if err != nil {
		panic(err)
	}

	fmt.Println(httpResponse)

	bs, err := json.Marshal(httpResponse.Result)
	if err = json.Unmarshal(bs, &networkResp); err != nil {
		return err
	}

	logger.Infof("get result ip: %s, mask: %s", networkResp.Ip, networkResp.Mask)
	err = tuntap.New(tuntap.TAP, networkResp.Ip, networkResp.Mask, networkResp.NetworkId)
	if err != nil {
		return err
	}

	//request to fvpn to tell that network has been created.
	fvpnUrl := fmt.Sprintf("%s", "http://localhost:6663/api/v1/join")

	req := body{
		networkId: netId,
	}
	buff1, _ := json.Marshal(req)
	resultBuff, err := httputil.Post(fvpnUrl, bytes.NewReader(buff1))

	if err != nil || resultBuff == nil {
		return err
	}

	logger.Infof("join network %s success, resp: %s", networkResp.NetworkId, string(resultBuff))

	return nil
}

func (n *Node) RunLeaveNetwork(networkId string) error {

	return nil
}
