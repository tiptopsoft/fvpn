package util

import (
	"encoding/json"
	"errors"
	"github.com/topcloudz/fvpn/pkg/packet"
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

func GetPacketHeader(buff []byte) (packet.Header, error) {
	if len(buff) < packet.HeaderBuffSize {
		return packet.Header{}, errors.New("not invalid packer")
	}
	h, err := packet.Decode(buff[:packet.HeaderBuffSize])
	if err != nil {
		return packet.Header{}, err
	}
	return h, nil
}

func GetLocalInfo() (info *LocalInfo, err error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(homedir, ".fvpn/config.json")
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return nil, errors.New("please login")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)

	var local LocalConfig
	err = decoder.Decode(&local)
	if err != nil {
		return nil, err
	}

	values := strings.Split(local.Auth, ":")
	info.Username = values[0]
	info.Password, err = Base64Decode(values[1])
	info.UserId = local.UserId
	if err != nil {
		return nil, err
	}

	return info, nil
}
