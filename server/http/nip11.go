package http

import (
	"net/http"
)

func (s *Server) nip11Handler(w http.ResponseWriter, r *http.Request) {
	s.respondWithJSON(w, 200, s.nip11Doc)
}
