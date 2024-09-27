package http

import (
	"context"
	"net/http"
	"runtime"
)

type status string

const (
	statusOK                 status = "OK"
	statusPartiallyAvailable status = "Partially Available"
	statusUnavailable        status = "Unavailable"
	statusTimeout            status = "Timeout during health check"
)

type system struct {
	Version          string `json:"version"`
	GoroutinesCount  int    `json:"goroutines_count"`
	TotalAllocBytes  uint64 `json:"total_alloc_bytes"`
	HeapObjectsCount uint64 `json:"heap_objects_count"`
	AllocBytes       uint64 `json:"alloc_bytes"`
}

type service struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Reason string `json:"reason"`
}

type healthResponse struct {
	Status   status  `json:"status"`
	Database service `json:"databse"`
	System   system  `json:"system"`
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	resp := healthResponse{
		Status: statusOK,
		System: system{
			Version:          runtime.Version(),
			GoroutinesCount:  runtime.NumGoroutine(),
			TotalAllocBytes:  ms.Alloc,
			HeapObjectsCount: ms.HeapObjects,
			AllocBytes:       ms.Alloc,
		},
		Database: service{
			Name: "mongo_db",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.db.QueryTimeout)
	defer cancel()

	resp.Database.Status = true
	if err := s.db.Client.Ping(ctx, nil); err != nil {
		resp.Database.Status = false
		resp.Database.Reason = err.Error()
		resp.Status = statusPartiallyAvailable
	}

	s.respondWithJSON(w, 200, resp)
}
