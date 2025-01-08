package websocket

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	grpcclient "github.com/dezh-tech/immortal/infrastructure/grpc_client"
	"github.com/dezh-tech/immortal/infrastructure/metrics"
	"github.com/dezh-tech/immortal/infrastructure/redis"
	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/repository"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/dezh-tech/immortal/types/message"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

// Server represents a websocket serer which keeps track of client connections and handle them.
type Server struct {
	mu sync.RWMutex

	config  Config
	conns   map[*websocket.Conn]clientState
	handler *repository.Handler
	metrics *metrics.Metrics
	redis   *redis.Redis
	grpc    grpcclient.IClient
}

func New(cfg Config, h *repository.Handler, m *metrics.Metrics,
	r *redis.Redis, grpc grpcclient.IClient,
) (*Server, error) {
	return &Server{
		config:  cfg,
		conns:   make(map[*websocket.Conn]clientState),
		mu:      sync.RWMutex{},
		handler: h,
		metrics: m,
		redis:   r,
		grpc:    grpc,
	}, nil
}

// Start starts a new server instance.
func (s *Server) Start() error {
	go s.checkExpiration()

	addr := net.JoinHostPort(s.config.Bind,
		strconv.Itoa(int(s.config.Port)))

	logger.Info("websocket server started", "listen", addr)

	http.Handle("/", s)
	err := http.ListenAndServe(addr, nil) //nolint

	return err
}

// handleWS is WebSocket handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.mu.Lock()

	logger.Debug("incoming websocket connection",
		"addr", conn.RemoteAddr().String())

	s.metrics.Connections.Inc()

	known := false
	pubkey := ""

	s.conns[conn] = clientState{
		pubkey:  &pubkey,
		isKnown: &known,
		subs:    make(map[string]filter.Filters),
		RWMutex: &sync.RWMutex{},
	}

	s.mu.Unlock()

	s.readLoop(conn)
}

// readLoop reads incoming messages from a client and answer to them.
func (s *Server) readLoop(conn *websocket.Conn) {
	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			logger.Debug("failed to read form connection", "conn",
				conn.RemoteAddr().String(), "err", err.Error())

			// clean up closed connection.
			s.mu.Lock()

			s.metrics.Connections.Dec()

			client, ok := s.conns[conn]
			if ok {
				s.metrics.Subscriptions.Sub(float64(len(client.subs)))
			}

			delete(s.conns, conn)

			s.mu.Unlock()

			break
		}

		msg, err := message.ParseMessage(buf)
		if err != nil {
			_ = conn.WriteMessage(1, message.MakeNotice(
				fmt.Sprintf("error: can't parse message: %s",
					err.Error())))

			continue
		}

		s.metrics.MessagesTotal.Inc()

		logger.Debug("incoming message", "conn",
			conn.RemoteAddr().String(), "msg", msg.String())

		switch msg.Type() {
		case "REQ":
			go s.handleReq(conn, msg)

		case "EVENT":
			go s.handleEvent(conn, msg)

		case "CLOSE":
			go s.handleClose(conn, msg)

		case "AUTH":
			go s.handleAuth(conn, msg)
		}
	}
}

// Stop shutdowns the server gracefully.
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Info("stopping websocket server")

	for wsConn, client := range s.conns {
		client.Lock()
		// close all subscriptions.
		for id := range client.subs {
			delete(client.subs, id)

			err := wsConn.WriteMessage(1, message.MakeClosed(id, "error: shutdown the relay."))
			if err != nil {
				return fmt.Errorf("error: closing subscription: %s, connection: %s, error: %s",
					id, wsConn.RemoteAddr(), err.Error())
			}
		}

		// close connection.
		delete(s.conns, wsConn)
		err := wsConn.Close()
		if err != nil {
			return fmt.Errorf("error: closing connection: %s, error: %s",
				wsConn.RemoteAddr(), err.Error())
		}

		client.Unlock()
	}

	return nil
}
