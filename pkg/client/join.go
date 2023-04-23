package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	http2 "github.com/topcloudz/fvpn/pkg/http"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"io"
	"net/http"
)

const (
	userUrl = "https://www.efvpn.com"
)

func (n *Node) RunJoinNetwork(netId string) error {
	logger.Infof("start to join %s", netId)
	//user http to get networkId config
	type body struct {
		userId string
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

	var httpResponse http2.HttpResponse
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
	return tuntap.New(tuntap.TAP, networkResp.Ip, networkResp.Mask, networkResp.NetworkId)
}

func (n *Node) RunLeaveNetwork(networkId string) error {

	return nil
}
