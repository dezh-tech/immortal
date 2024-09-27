package http

import (
	"net/http"
)

func (s *Server) nip11Handler(w http.ResponseWriter, _ *http.Request) {
	s.respondWithJSON(w, 200, s.nip11Doc)
}
