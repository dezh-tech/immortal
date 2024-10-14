package http

import (
	"context"
	"net/http"
)

type status string

const (
	statusOK                 status = "OK"
	statusPartiallyAvailable status = "Partially Available"
	statusUnavailable        status = "Unavailable"
	statusTimeout            status = "Timeout during health check"
)

type service struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Reason string `json:"reason"`
}

type healthResponse struct {
	Status   status  `json:"status"`
	Database service `json:"databse"`
}

func (s *Server) healthHandler(w http.ResponseWriter, _ *http.Request) {
	resp := healthResponse{
		Status: statusOK,
		Database: service{
			Name: "mongo_db",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.db.QueryTimeout)
	defer cancel()

	resp.Database.Status = true
	if err := s.db.Client.Ping(ctx, nil); err != nil { //nolint
		resp.Database.Status = false
		resp.Database.Reason = err.Error()
		resp.Status = statusPartiallyAvailable
	}

	s.respondWithJSON(w, 200, resp)
}
