package relay

import (
	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/server"
)

// Relay keeps all concepts such as server, database and manages them.
type Relay struct {
	config   config.Config
	server   *server.Server
	database *database.Database
}

// NewRelay creates a new relay.
func New(cfg config.Config) (*Relay, error) {
	db, err := database.New(cfg.DSN)
	if err != nil {
		return nil, err
	}

	return &Relay{
		config:   cfg,
		server:   server.New(cfg.ServerConf),
		database: db,
	}, nil
}

// Start runs the relay and its children.
func (r *Relay) Start() error {
	return r.server.Start()
}

// Stop shutdowns the relay and its children gracefully.
func (r *Relay) Stop() error {
	return r.server.Stop()
}
