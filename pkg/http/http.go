package http

import (
	"encoding/json"
	"github.com/interstellar-cloud/star/pkg/cache"
	"github.com/interstellar-cloud/star/pkg/log"
	"net/http"
)

var (
	logger = log.Log()
)

type HttpServer struct {
	cache cache.PeersCache
}

func New(cache cache.PeersCache) HttpServer {
	return HttpServer{cache: cache}
}

func (hs HttpServer) Start() error {
	http.HandleFunc("/cache/list", func(w http.ResponseWriter, r *http.Request) {
		n := hs.cache
		if err := json.NewEncoder(w).Encode(n); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	logger.Infof("http listen at: %s", ":4001")
	return http.ListenAndServe(":4001", nil)
}
