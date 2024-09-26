package relay

import (
	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/handler"
	"github.com/dezh-tech/immortal/server/http"
	"github.com/dezh-tech/immortal/server/websocket"
)

// Relay keeps all concepts such as server, database and manages them.
type Relay struct {
	config          config.Config
	webscoketServer *websocket.Server
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

	ws, err := websocket.New(cfg.WebsocketServer, h)
	if err != nil {
		return nil, err
	}

	hs, err := http.New(cfg, db)
	if err != nil {
		return nil, err
	}

	return &Relay{
		config:          *cfg,
		webscoketServer: ws,
		database:        db,
		httpServer:      hs,
	}, nil
}

// Start runs the relay and its children.
func (r *Relay) Start() error {
	go func() error {
		if err := r.webscoketServer.Start(); err != nil {
			return err
		}

		return nil
	}()

	return r.httpServer.Start()
}

// Stop shutdowns the relay and its children gracefully.
func (r *Relay) Stop() error {
	if err := r.webscoketServer.Stop(); err != nil {
		return err
	}

	return r.database.Stop()
}
