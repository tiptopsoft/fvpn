package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	http2 "github.com/topcloudz/fvpn/pkg/http"
	"io"
	"net/http"
)

func main() {

	destUrl := fmt.Sprintf("%s/api/v1/users/user/%s/network/%s/join", "https://www.efvpn.com", "1", "c04d6b84fd4fc978")
	type body struct {
		userId string
	}

	request := body{userId: "1"}
	buff, err := json.Marshal(request)

	resp, err := http.Post(destUrl, "application/json", bytes.NewBuffer(buff))
	if err != nil {
		panic(err)
	}

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

	bs, _ := json.Marshal(httpResponse.Result)

	json.Unmarshal(bs, &networkResp)

	fmt.Println(networkResp.Ip, networkResp.Mask, networkResp.NetworkId)
}
