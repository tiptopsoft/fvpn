package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	nativehttp "github.com/topcloudz/fvpn/pkg/http"
	"github.com/topcloudz/fvpn/pkg/option"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"io"
	"net"
	"net/http"
)

func (p *Peer) runHttpServer() error {
	s := nativehttp.NewServer()
	s.HandlerFunc("/api/v1/join", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//将networkId 写入cache
		switch r.Method {
		case "POST":
			buff, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("invalid network"))
				return
			}

			type body struct {
				Status    int    `json:"status"`
				NetworkId string `json:"NetworkId"`
				Ip        string `json:"ip"`
				Mask      string `json:"mask"`
			}

			req := new(body)

			err = json.Unmarshal(buff, &req)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("invalid network"))
				return
			}

			logger.Infof("request network is: %s", req.NetworkId)

			tap, err := tuntap.New(tuntap.TAP, req.NetworkId)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			if p.devices == nil {
				p.devices[req.NetworkId] = tap
			}

			if err := json.NewEncoder(w).Encode(req); err != nil {
				w.WriteHeader(500)
				return
			}

			logger.Infof("get result ip: %s, mask: %s", req.Ip, req.Mask)

			//get dev
			if err = option.ExecCommand("/bin/sh", "-c", fmt.Sprintf("ifconfig %s %s netmask %s mtu %d up", tap.Name, req.Ip, req.Mask, 1420)); err != nil {
				w.WriteHeader(500)
				return
			}

			// start tap
			tap.IP = net.ParseIP(req.Ip)
			p.devices[req.NetworkId] = tap
			go p.ReadFromTun(tap, req.NetworkId)
			p.SendRegister(tap)
			//give a timer

			go util.AddJob(req.NetworkId, p.sendQueryPeer)

			w.WriteHeader(200)
			logger.Infof("join network %s success", req.NetworkId)
		}

	})

	s.HandlerFunc("/api/v1/leave", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//将networkId 写入cache
		switch r.Method {
		case "POST":
			buff, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("invalid network"))
				return
			}

			var body struct {
				NetworkId string
			}

			err = json.Unmarshal(buff, &body)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("invalid network"))
				return
			}

			//TODO 把networkId加进来
			logger.Infof("join network %s success", body.NetworkId)
			if err = p.RunJoinNetwork(body.NetworkId); err != nil {
				w.WriteHeader(500)
				return
			}
		}
		w.WriteHeader(200)
	})
	logger.Debugf("node started at: :%d", DefaultPort)
	return s.Start(fmt.Sprintf(":%d", DefaultPort))
}

// GetNetworkIds get network ids when node starting, so can monitor the traffic on the device.
func (p *Peer) GetNetworkIds() ([]string, error) {
	var body struct {
		userId string
	}

	body.userId = "1"
	buff, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	destUrl := fmt.Sprintf("%s%s", userUrl, "/api/v1/users/user/1/getJoinNetwork")
	resp, err := nativehttp.Post(destUrl, bytes.NewBuffer(buff))

	if err != nil {
		return nil, err
	}

	var result []string
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
