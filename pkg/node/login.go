package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/topcloudz/fvpn/pkg/util"
	"os"
	"path/filepath"
)

func Login(username, password string, cfg *util.ClientConfig) error {
	client := NewClient(cfg.ControlUrl())
	req := new(util.LoginRequest)
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
		file, err = os.OpenFile(path, os.O_RDWR, 0755)
	}
	defer file.Close()
	var local util.LocalConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&local)
	if err != nil {
		return errors.New("login failed")
	}

	appId := local.AppId
	if appId == "" {
		appId = local.AppId
	}

	encoder := json.NewEncoder(file)

	b := util.LocalConfig{
		Auth:  fmt.Sprintf("%s:%s", username, util.StringToBase64(password)),
		AppId: local.AppId,
	}
	return encoder.Encode(b)
}
