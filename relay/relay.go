package relay

import (
	"context"
	"log"

	"github.com/dezh-tech/immortal/client"
	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/handler"
	"github.com/dezh-tech/immortal/metrics"
	"github.com/dezh-tech/immortal/relay/redis"
	"github.com/dezh-tech/immortal/server"
	"github.com/dezh-tech/immortal/utils"
)

// Relay keeps all concepts such as server, database and manages them.
type Relay struct {
	config          config.Config
	websocketServer *server.Server
	database        *database.Database
	redis           *redis.Redis
}

// NewRelay creates a new relay.
func New(cfg *config.Config) (*Relay, error) {
	db, err := database.Connect(cfg.Database)
	if err != nil {
		return nil, err
	}

	err = cfg.LoadParameters(db)
	if err != nil {
		return nil, err
	}

	h := handler.New(db, cfg.Parameters.Handler)

	m := metrics.New()

	r, err := redis.New(cfg.RedisConf)
	if err != nil {
		return nil, err
	}

	ws, err := server.New(cfg.WebsocketServer, cfg.GetNIP11Documents(), h, m, r)
	if err != nil {
		return nil, err
	}

	c, err := client.NewClient(cfg.Kraken.Endpoint)
	if err != nil {
		return nil, err
	}

	la, err := utils.LocalAddr()
	if err != nil {
		return nil, err
	}

	resp, err := c.Register(context.Background(), la, cfg.Kraken.Heartbeat)
	if err != nil || resp == nil || !resp.Success {
		return nil, err
	}

	return &Relay{
		config:          *cfg,
		websocketServer: ws,
		database:        db,
		redis:           r,
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

	return errCh
}

// Stop shutdowns the relay and its children gracefully.
func (r *Relay) Stop() error {
	log.Println("stopping relay...")

	if err := r.websocketServer.Stop(); err != nil {
		return err
	}

	return r.database.Stop()
}
