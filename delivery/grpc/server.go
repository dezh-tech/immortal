package grpc

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	rpb "github.com/dezh-tech/immortal/delivery/grpc/gen"
	"github.com/dezh-tech/immortal/infrastructure/database"
	"github.com/dezh-tech/immortal/infrastructure/redis"
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
}

func New(conf *Config, r *redis.Redis, db *database.Database, st time.Time) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		ctx:       ctx,
		cancel:    cancel,
		config:    conf,
		StartTime: st,
		Redis:     r,
		DB:        db,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", net.JoinHostPort(s.config.Bind, //nolint
		strconv.Itoa(int(s.config.Port))))
	if err != nil {
		return err
	}

	log.Println("grpc server started...")

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor())

	healthServer := newHealthServer(s)

	rpb.RegisterHealthServiceServer(grpcServer, healthServer)

	s.listener = listener
	s.grpc = grpcServer

	return s.grpc.Serve(listener)
}

func (s *Server) Stop() error {
	s.cancel()

	log.Println("grpc server stopped...")

	s.grpc.Stop()
	if err := s.listener.Close(); err != nil {
		return err
	}

	return nil
}
