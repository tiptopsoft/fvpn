package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func main() {
	mx := http.NewServeMux()

	mx.HandleFunc("/api/v1/join", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
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
					status    int
					networkId string `json:"NetworkId"`
				}

				req := new(body)

				err = json.Unmarshal(buff, &req)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte("invalid network"))
					return
				}

				//TODO 把networkId加进来
				//if err := n.addTuns(req.NetworkId); err != nil {
				//	w.WriteHeader(500)
				//}
				req.status = 200

				if err := json.NewEncoder(w).Encode(req); err != nil {
					w.WriteHeader(500)
					return
				}

				w.WriteHeader(200)
			}

		}

	})

	http.ListenAndServe(":4000", mx)
}
