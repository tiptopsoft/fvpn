package http

import (
	"encoding/json"
	"github.com/interstellar-cloud/star/pkg/log"
	"github.com/interstellar-cloud/star/pkg/node"
	"net/http"
)

var (
	logger = log.Log()s
)

type HttpServer struct {
	cache node.NodesCache
}

func New(cache node.NodesCache) HttpServer {
	return HttpServer{cache: cache}
}

func (hs HttpServer) Start() error {
	http.HandleFunc("/node/list", func(w http.ResponseWriter, r *http.Request) {
		n := hs.cache
		if err := json.NewEncoder(w).Encode(n); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	logger.Infof("http listen at: %s", ":4001")
	return http.ListenAndServe(":4001", nil)
}
