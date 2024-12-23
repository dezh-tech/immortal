package relay

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dezh-tech/immortal/client"
	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/handler"
	"github.com/dezh-tech/immortal/metrics"
	"github.com/dezh-tech/immortal/relay/redis"
	"github.com/dezh-tech/immortal/server/grpc"
	"github.com/dezh-tech/immortal/server/websocket"
)

// Relay keeps all concepts such as server, database and manages them.
type Relay struct {
	config          config.Config
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

	r, err := redis.New(cfg.RedisConf)
	if err != nil {
		return nil, err
	}

	c, err := client.NewClient(cfg.Kraken.Endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.RegisterService(context.Background(), fmt.Sprint(cfg.GRPCServer.Port),
		cfg.Kraken.Region, cfg.Kraken.Heartbeat)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("cant register to master: %s", *resp.Message)
	}

	params, err := c.GetConfig(context.Background(), resp.Token)
	if err != nil {
		return nil, err
	}

	err = cfg.LoadParameters(params)
	if err != nil {
		return nil, err
	}

	h := handler.New(db, cfg.Handler)

	ws, err := websocket.New(cfg.WebsocketServer, h, m, r)
	if err != nil {
		return nil, err
	}

	gs := grpc.New(&cfg.GRPCServer, r, db, time.Now())

	return &Relay{
		config:          *cfg,
		websocketServer: ws,
		database:        db,
		redis:           r,
		grpcServer:      gs,
	}, nil
}

// Start runs the relay and its children.
func (r *Relay) Start() chan error {
	log.Println("relay started successfully...")
	errCh := make(chan error, 2)

	go func() {
		if err := r.websocketServer.Start(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := r.grpcServer.Start(); err != nil {
			errCh <- err
		}
	}()

	return errCh
}

// Stop shutdowns the relay and its children gracefully.
func (r *Relay) Stop() error {
	log.Println("stopping relay...")

	if err := r.websocketServer.Stop(); err != nil {
		return err
	}

	if err := r.grpcServer.Stop(); err != nil {
		return err
	}

	return r.database.Stop()
}
