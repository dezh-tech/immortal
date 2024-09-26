package http

import (
	"encoding/json"
	"net/http"

	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/database"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	config *config.Config
	db     *database.Database
}

func New(cfg *config.Config, db *database.Database) (*Server, error) {
	r := mux.NewRouter()
	s := &Server{
		router: r,
		config: cfg,
		db:     db,
	}

	r.HandleFunc("/nip11", s.nip11Handler)
	r.HandleFunc("/health", s.healthHandler)

	return s, nil
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (s *Server) Start() error {
	// todo::: read from config.
	return http.ListenAndServe("127.0.0.1:8080", s.router)
}
