package http

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/dezh-tech/immortal/database"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	router *mux.Router
	db     *database.Database
	config Config
}

func New(cfg Config, db *database.Database) (*Server, error) {
	r := mux.NewRouter()
	s := &Server{
		router: r,
		db:     db,
		config: cfg,
	}

	r.HandleFunc("/health", s.healthHandler).Methods("GET")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	return s, nil
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload) //nolint

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func (s *Server) Start() error {
	return http.ListenAndServe(net.JoinHostPort(s.config.Bind, //nolint
		strconv.Itoa(int(s.config.Port))), s.router)
}
