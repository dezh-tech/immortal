package http

import (
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	router *mux.Router
	config Config
}

func New(cfg Config) (*Server, error) {
	r := mux.NewRouter()

	s := &Server{
		router: r,
		config: cfg,
	}

	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	return s, nil
}

func (s *Server) Start() error {
	return http.ListenAndServe(net.JoinHostPort(s.config.Bind, //nolint
		strconv.Itoa(int(s.config.Port))), s.router)
}
