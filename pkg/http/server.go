package http

import (
	"net/http"
)

type Server struct {
	*http.ServeMux
}

func NewServer() Server {
	return Server{
		http.NewServeMux(),
	}
}

func (s Server) HandlerFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}

func (s Server) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}
