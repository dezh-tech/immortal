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
	// Version is the go version.
	Version string `json:"version"`
	// GoroutinesCount is the number of the current goroutines.
	GoroutinesCount int `json:"goroutines_count"`
	// TotalAllocBytes is the total bytes allocated.
	TotalAllocBytes int `json:"total_alloc_bytes"`
	// HeapObjectsCount is the number of objects in the go heap.
	HeapObjectsCount int `json:"heap_objects_count"`
	// TotalAllocBytes is the bytes allocated and not yet freed.
	AllocBytes int `json:"alloc_bytes"`
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
			TotalAllocBytes:  int(ms.Alloc),
			HeapObjectsCount: int(ms.HeapObjects),
			AllocBytes:       int(ms.Alloc),
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
