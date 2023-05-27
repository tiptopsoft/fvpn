package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	nativehttp "github.com/topcloudz/fvpn/pkg/http"
	"github.com/topcloudz/fvpn/pkg/tuntap"
	"github.com/topcloudz/fvpn/pkg/util"
	"io"
	"net/http"
)

func (n *Node) runHttpServer() error {
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

			if err := json.NewEncoder(w).Encode(req); err != nil {
				w.WriteHeader(500)
				return
			}

			logger.Infof("get result ip: %s, mask: %s", req.Ip, req.Mask)
			tap, err := tuntap.New(tuntap.TAP, req.Ip, req.Mask, req.NetworkId)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			n.tun.CacheDevice(req.NetworkId, tap)
			err = util.SendRegister(tap, n.relaySocket)
			if err != nil {
				return
			}
			go n.tun.ReadFromTun(context.Background(), req.NetworkId)
			go n.tun.WriteToUdp()

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
			if err = n.RunJoinNetwork(body.NetworkId); err != nil {
				w.WriteHeader(500)
				return
			}
		}
		w.WriteHeader(200)
	})
	logger.Debugf("node started at: :%d", DefaultPort)
	return s.Start(fmt.Sprintf(":%d", DefaultPort))
}

//func (n *Node) addTuns(networkId string) error {
//	tun, err := tuntap.GetTuntap(networkId)
//	if err != nil {
//		return err
//	}
//
//	n.tuns.Store(networkId, tun)
//	return nil
//}

// GetNetworkIds get network ids when node starting, so can monitor the traffic on the device.
func (n *Node) GetNetworkIds() ([]string, error) {
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
