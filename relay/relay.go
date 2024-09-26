package relay

import (
	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/handler"
	"github.com/dezh-tech/immortal/server"
)

// Relay keeps all concepts such as server, database and manages them.
type Relay struct {
	config   config.Config
	server   *server.Server
	database *database.Database
}

// NewRelay creates a new relay.
func New(cfg *config.Config) (*Relay, error) {
	db, err := database.Connect(cfg.DatabaseConf)
	if err != nil {
		return nil, err
	}

	err = cfg.LoadParameters(db)
	if err != nil {
		return nil, err
	}

	h := handler.New(db, cfg.Parameters.Handler)

	s, err := server.New(cfg.ServerConf, h)
	if err != nil {
		return nil, err
	}

	return &Relay{
		config:   *cfg,
		server:   s,
		database: db,
	}, nil
}

// Start runs the relay and its children.
func (r *Relay) Start() error {
	return r.server.Start()
}

// Stop shutdowns the relay and its children gracefully.
func (r *Relay) Stop() error {
	if err := r.server.Stop(); err != nil {
		return err
	}

	return r.database.Stop()
}
