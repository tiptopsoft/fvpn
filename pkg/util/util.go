package util

import (
	"encoding/json"
	"errors"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type IPHeader struct {
	SrcIP net.IP
	DstIP net.IP
}

// GetIPFrameHeader return srcIP, destIP
func GetIPFrameHeader(buff []byte) (*IPHeader, error) {
	if len(buff) < packet.IPBuffSize {
		return nil, errors.New("invalid ip frame")
	}

	h := new(IPHeader)
	h.SrcIP = net.IPv4(buff[12], buff[13], buff[14], buff[15])
	h.DstIP = net.IPv4(buff[16], buff[17], buff[18], buff[19])
	return h, nil
}

func GetPacketHeader(buff []byte) (header.Header, error) {
	if len(buff) < packet.HeaderBuffSize {
		return header.Header{}, errors.New("not invalid packer")
	}
	h, err := header.Decode(buff[:packet.HeaderBuffSize])
	if err != nil {
		return header.Header{}, err
	}
	return h, nil
}

func GetUserInfo() (string, string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	path := filepath.Join(homedir, ".fvpn/config.json")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return "", "", errors.New("please login")
	}
	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}

	decoder := json.NewDecoder(file)

	var resp Login
	err = decoder.Decode(&resp)
	if err != nil {
		return "", "", err
	}

	values := strings.Split(resp.Auth, ":")
	username := values[0]
	password, err := Base64Decode(values[1])
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}
