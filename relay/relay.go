package relay

import (
	"log"

	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/handler"
	"github.com/dezh-tech/immortal/metrics"
	"github.com/dezh-tech/immortal/server/http"
	"github.com/dezh-tech/immortal/server/websocket"
)

// Relay keeps all concepts such as server, database and manages them.
type Relay struct {
	config          config.Config
	websocketServer *websocket.Server
	httpServer      *http.Server
	database        *database.Database
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

	ws, err := websocket.New(cfg.WebsocketServer, cfg.GetNIP11Documents(), h, m)
	if err != nil {
		return nil, err
	}

	hs, err := http.New(cfg.HTTPServer, db)
	if err != nil {
		return nil, err
	}

	return &Relay{
		config:          *cfg,
		websocketServer: ws,
		database:        db,
		httpServer:      hs,
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
		if err := r.httpServer.Start(); err != nil {
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
