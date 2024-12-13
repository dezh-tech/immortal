package grpc

import (
	"context"
	"time"

	"github.com/dezh-tech/immortal"
	"github.com/dezh-tech/immortal/server/grpc/gen"
)

type healthServer struct {
	*Server
}

func newHealthServer(server *Server) *healthServer {
	return &healthServer{
		Server: server,
	}
}

func (s healthServer) Status(ctx context.Context, _ *gen.StatusRequest) (*gen.StatusResponse, error) {
	services := make([]*gen.Service, 0)

	redisStatus := gen.Status_CONNECTED
	redisMessage := ""

	if err := s.Redis.Client.Ping(ctx).Err(); err != nil {
		redisStatus = gen.Status_DISCONNECTED
		redisMessage = err.Error()
	}

	redis := gen.Service{
		Name:    "redis",
		Status:  redisStatus,
		Message: redisMessage,
	}

	services = append(services, &redis)

	mongoStatus := gen.Status_CONNECTED
	mongoMessage := ""

	if err := s.DB.Client.Ping(ctx, nil); err != nil {
		mongoStatus = gen.Status_DISCONNECTED
		mongoMessage = err.Error()
	}

	mongo := gen.Service{
		Name:    "mongo",
		Status:  mongoStatus,
		Message: mongoMessage,
	}

	services = append(services, &mongo)

	return &gen.StatusResponse{
		Uptime:   int64(time.Since(s.StartTime).Seconds()),
		Version:  immortal.StringVersion(),
		Services: services,
	}, nil
}
