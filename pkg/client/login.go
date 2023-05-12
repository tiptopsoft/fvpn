package client

import (
	"encoding/json"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/http"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/util"
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
	homeDir, err := os.UserHomeDir()
	path := filepath.Join(homeDir, ".fvpn/config.json")
	_, err = os.Stat(path)
	var file *os.File
	if os.IsNotExist(err) {
		parentDir := filepath.Dir(path)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return err
		}
		file, err = os.Create(path)
	} else {
		file, err = os.Open(path)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)

	b := option.Login{
		Auth: fmt.Sprintf("%s:%s", username, util.StringToBase64(password)),
	}
	return encoder.Encode(b)
}
