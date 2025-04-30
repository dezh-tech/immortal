package grpc

import (
	"context"
	"time"

	"github.com/dezh-tech/immortal"
	rpb "github.com/dezh-tech/immortal/delivery/grpc/gen"
)

type healthServer struct {
	*Server
}

func newHealthServer(server *Server) *healthServer {
	return &healthServer{
		Server: server,
	}
}

func (s healthServer) Status(ctx context.Context, _ *rpb.StatusRequest) (*rpb.StatusResponse, error) {
	services := make([]*rpb.Service, 0)

	redisStatus := rpb.Status_CONNECTED
	redisMessage := ""

	if err := s.redis.Client.Ping(ctx).Err(); err != nil {
		redisStatus = rpb.Status_DISCONNECTED
		redisMessage = err.Error()
	}

	redis := rpb.Service{
		Name:    "redis",
		Status:  redisStatus,
		Message: redisMessage,
	}

	services = append(services, &redis)

	mongoStatus := rpb.Status_CONNECTED
	mongoMessage := ""

	if err := s.database.Client.Ping(ctx, nil); err != nil {
		mongoStatus = rpb.Status_DISCONNECTED
		mongoMessage = err.Error()
	}

	mongo := rpb.Service{
		Name:    "mongo",
		Status:  mongoStatus,
		Message: mongoMessage,
	}

	services = append(services, &mongo)

	return &rpb.StatusResponse{
		Uptime:   int64(time.Since(s.startTime).Seconds()),
		Version:  immortal.StringVersion(),
		Services: services,
	}, nil
}
