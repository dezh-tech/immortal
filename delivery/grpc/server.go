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
	"google.golang.org/grpc"
)

type Server struct {
	ctx       context.Context
	cancel    context.CancelFunc
	config    *Config
	listener  net.Listener
	grpc      *grpc.Server
	StartTime time.Time
	DB        *database.Database
	Redis     *redis.Redis
	keeper    ParametersKeeper
}

func New(conf *Config, r *redis.Redis, db *database.Database, st time.Time, keeper ParametersKeeper) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:       ctx,
		cancel:    cancel,
		config:    conf,
		StartTime: st,
		Redis:     r,
		DB:        db,
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

	rpb.RegisterHealthServiceServer(grpcServer, healthServer)
	rpb.RegisterShutdownServiceServer(grpcServer, shutdownServer)
	rpb.RegisterParametersServiceServer(grpcServer, paramsServer)

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
