package util

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/topcloudz/fvpn/pkg/log"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/packet/header"
	"golang.org/x/sys/unix"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
)

var (
	logger = log.Log()
)

type FrameHeader struct {
	DestinationAddr net.HardwareAddr
	SourceAddr      net.HardwareAddr
	SourceIP        net.IP
	DestinationIP   net.IP
	EtherType       uint16
}

func GetAddress(address string, port int) (unix.SockaddrInet4, error) {
	ad, err := netip.ParseAddr(address)
	return unix.SockaddrInet4{
		Port: port,
		Addr: ad.As4(),
	}, err
}

// GetFrameHeader return dest mac, dest ip, if data provide is null, error return
func GetFrameHeader(buff []byte) (*FrameHeader, error) {
	if len(buff) == 0 {
		return nil, errors.New("no data exists")
	}
	header := parseHeader(buff)

	//是ARP
	if header.EtherType == 0x0806 {
		logger.Debugf("this is an arp frame")
		header.SourceIP = net.IPv4(buff[34], buff[35], buff[36], buff[37])
		header.DestinationIP = net.IPv4(buff[38], buff[39], buff[40], buff[41])
	}

	//IP
	if header.EtherType == 0x0800 {
		logger.Debugf("this is an IP frame")
		header.SourceIP = net.IPv4(buff[26], buff[27], buff[28], buff[29])
		header.DestinationIP = net.IPv4(buff[30], buff[31], buff[32], buff[33])
	}

	logger.Debugf("recevice header is: %v", header)
	return header, nil
}

func GetPacketHeader(buff []byte) (header.Header, error) {
	h, err := header.Decode(buff[:12])
	if err != nil {
		return header.Header{}, err
	}
	return h, nil
}

func parseHeader(buf []byte) *FrameHeader {
	header := new(FrameHeader)
	var hd net.HardwareAddr
	hd = buf[0:6]
	header.DestinationAddr = hd
	hd = buf[6:12]
	header.SourceAddr = hd
	header.EtherType = binary.BigEndian.Uint16(buf[12:14])
	return header
}

func GetUserInfo() (string, string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	path := filepath.Join(homedir, "./fvpn/config.json")
	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}

	decoder := json.NewDecoder(file)

	var resp option.Login
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
