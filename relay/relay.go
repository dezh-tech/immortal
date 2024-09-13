package relay

import (
	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/server"
)

// Relay keeps all concepts such as server, database and manages them.
type Relay struct {
	config config.Config
	server *server.Server
}

// NewRelay creates a new relay.
func NewRelay(cfg config.Config) *Relay {
	return &Relay{
		config: cfg,
		server: server.NewServer(cfg.ServerConf),
	}
}

// Start runs the relay and its childs.
func (r *Relay) Start() error {
	return r.server.Start()
}

// Stop shutdowns the relay and its childs gracefully.
func (r *Relay) Stop() error {
	return r.server.Stop()
}
