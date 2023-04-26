package nativehttp

import (
	"encoding/json"
	"github.com/topcloudz/fvpn/pkg/cache"
	"github.com/topcloudz/fvpn/pkg/log"
	"net/http"
)

var (
	logger = log.Log()
)

type HttpServer struct {
	cache *cache.Cache
}

func New(cache *cache.Cache) HttpServer {
	return HttpServer{cache: cache}
}

func (hs HttpServer) Start() error {
	http.HandleFunc("/cache/list", func(w http.ResponseWriter, r *http.Request) {
		n := hs.cache
		if err := json.NewEncoder(w).Encode(n); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	logger.Infof("nativehttp listen at: %s", ":4001")
	return http.ListenAndServe(":4001", nil)
}
