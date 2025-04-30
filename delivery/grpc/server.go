package grpc

import (
	"context"
	"net"
	"strconv"
	"time"

	rpb "github.com/dezh-tech/immortal/delivery/grpc/gen"
	"github.com/dezh-tech/immortal/infrastructure/database"
	"github.com/dezh-tech/immortal/infrastructure/redis"
	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/repository"
	"google.golang.org/grpc"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc

	config   Config
	grpc     *grpc.Server
	database database.Database
	redis    *redis.Redis
	handler  repository.Handler

	keeper ParametersKeeper

	listener  net.Listener
	startTime time.Time
}

func New(conf Config, r *redis.Redis, db database.Database,
	handler repository.Handler, st time.Time, keeper ParametersKeeper,
) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:       ctx,
		cancel:    cancel,
		config:    conf,
		startTime: st,
		redis:     r,
		database:  db,
		handler:   handler,
		keeper:    keeper,
	}
}

func (s *Server) Start(shutdownch chan struct{}) error {
	listener, err := net.Listen("tcp", net.JoinHostPort(s.config.Bind, //nolint
		strconv.Itoa(int(s.config.Port))))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor())

	healthServer := newHealthServer(s)
	shutdownServer := newShutdownServer(s, shutdownch)
	paramsServer := newParamsServer(s)
	migrationServer := newMigrationServer(s)

	rpb.RegisterHealthServer(grpcServer, healthServer)
	rpb.RegisterShutdownServer(grpcServer, shutdownServer)
	rpb.RegisterParametersServer(grpcServer, paramsServer)
	rpb.RegisterMigrationServer(grpcServer, migrationServer)

	s.listener = listener
	s.grpc = grpcServer

	logger.Info("gRPC server started successfully", "listen", listener.Addr().String())

	if err := s.grpc.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	logger.Info("stopping gRPC server")

	s.cancel()
	s.grpc.Stop()

	return nil
}
