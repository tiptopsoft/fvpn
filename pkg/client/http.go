package client

import (
	"encoding/json"
	"fmt"
	nativehttp "github.com/topcloudz/fvpn/pkg/http"
	"github.com/topcloudz/fvpn/pkg/tuntap"
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

			var body struct {
				status    int
				networkId string
			}

			err = json.Unmarshal(buff, &body)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("invalid network"))
				return
			}

			//TODO 把networkId加进来
			logger.Infof("join network %s success", body.networkId)
			if err := n.addTuns(body.networkId); err != nil {
				w.WriteHeader(500)
			}
			body.status = 200

			if err := json.NewEncoder(w).Encode(body); err != nil {
				w.WriteHeader(500)
				return
			}

			w.WriteHeader(200)
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
	logger.Debugf("node start at: %d", DefaultPort)
	return s.Start(fmt.Sprintf(":%d", DefaultPort))
}

func (n *Node) addTuns(networkId string) error {
	tun, err := tuntap.GetTuntap(networkId)
	if err != nil {
		return err
	}

	n.tuns.Store(networkId, tun)
	return nil
}
