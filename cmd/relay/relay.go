package relay

import (
	"context"
	"fmt"
	"github.com/dezh-tech/immortal/infrastructure/grpc_client/params_keeper"
	"time"

	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/delivery/grpc"
	"github.com/dezh-tech/immortal/delivery/websocket"
	"github.com/dezh-tech/immortal/infrastructure/database"
	grpcclient "github.com/dezh-tech/immortal/infrastructure/grpc_client"
	"github.com/dezh-tech/immortal/infrastructure/meilisearch"
	"github.com/dezh-tech/immortal/infrastructure/metrics"
	"github.com/dezh-tech/immortal/infrastructure/redis"
	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/repository"
)

// Relay keeps all concepts such as server, database and manages them.
type Relay struct {
	config          *config.Config
	websocketServer *websocket.Server
	grpcServer      *grpc.Server
	database        *database.Database
	redis           *redis.Redis
}

// NewRelay creates a new relay.
func New(cfg *config.Config) (*Relay, error) {
	db, err := database.Connect(cfg.Database)
	if err != nil {
		return nil, err
	}
	m := metrics.New()

	meili := meilisearch.New(cfg.Meili)

	r, err := redis.New(cfg.Redis)
	if err != nil {
		return nil, err
	}

	keeper := params_keeper.ParametersKeeper{
		Handler:         &cfg.Handler,
		WebsocketServer: &cfg.WebsocketServer,
	}

	c, err := grpcclient.New(cfg.GRPCClient.Endpoint, cfg.GRPCClient, keeper)
	if err != nil {
		return nil, err
	}

	resp, err := c.RegisterService(context.Background(), fmt.Sprint(cfg.GRPCServer.Port),
		cfg.GRPCClient.Region)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("cant register to master: %s", *resp.Message)
	}

	c.SetID(resp.Token)

	params, err := c.GetParameters(context.Background())
	if err != nil {
		return nil, err
	}

	err = cfg.LoadParameters(params)
	if err != nil {
		return nil, err
	}

	h := repository.New(&cfg.Handler, db, meili, c)

	ws, err := websocket.New(&cfg.WebsocketServer, h, m, r, c)
	if err != nil {
		return nil, err
	}

	gs := grpc.New(&cfg.GRPCServer, r, db, time.Now())

	return &Relay{
		config:          cfg,
		websocketServer: ws,
		database:        db,
		redis:           r,
		grpcServer:      gs,
	}, nil
}

// Start runs the relay and its children.
func (r *Relay) Start(shutdownch chan struct{}) chan error {
	logger.Info("starting the relay")

	errCh := make(chan error, 2)

	go func() {
		if err := r.websocketServer.Start(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := r.grpcServer.Start(shutdownch); err != nil {
			errCh <- err
		}
	}()

	return errCh
}

// Stop shutdowns the relay and its children gracefully.
func (r *Relay) Stop() error {
	logger.Info("stopping the relay")

	if err := r.websocketServer.Stop(); err != nil {
		return err
	}

	if err := r.grpcServer.Stop(); err != nil {
		return err
	}

	if err := r.database.Stop(); err != nil {
		return err
	}

	if err := r.redis.Close(); err != nil {
		return err
	}

	return nil
}
