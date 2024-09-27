package http

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/types/nip11"
	"github.com/gorilla/mux"
)

type Server struct {
	router   *mux.Router
	nip11Doc *nip11.RelayInformationDocument
	db       *database.Database
	config   Config
}

func New(cfg Config, rid *nip11.RelayInformationDocument, db *database.Database) (*Server, error) {
	r := mux.NewRouter()
	s := &Server{
		router:   r,
		nip11Doc: rid,
		db:       db,
		config:   cfg,
	}

	r.HandleFunc("/nip11", s.nip11Handler).Methods("GET")
	r.HandleFunc("/health", s.healthHandler).Methods("GET")

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
