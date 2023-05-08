package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/http"
	"os"
	"path/filepath"
)

func (n *Node) Login(username, password string) error {
	client := http.New(userUrl)
	req := new(http.LoginRequest)
	req.Username = username
	req.Password = password
	resp, err := client.Login(*req)
	if err != nil {
		return err
	}

	//登陆成功
	logger.Infof("login success. token:%s", resp.Token)
	//write to local
	path := filepath.Join("~/.fvpn/config.json")
	file, err := os.Open(path)
	defer file.Close()
	encoder := json.NewEncoder(file)
	type body struct {
		Auth string `json:"auth"`
	}

	b := body{
		Auth: fmt.Sprintf("%s:%s", username, stringToBase64(password)),
	}
	return encoder.Encode(b)
}

func stringToBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) (string, error) {
	buff, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	return string(buff), nil
}
