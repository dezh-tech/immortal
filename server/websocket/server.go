package websocket

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/dezh-tech/immortal/handler"
	"github.com/dezh-tech/immortal/metrics"
	"github.com/dezh-tech/immortal/relay/redis"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/dezh-tech/immortal/types/message"
	"github.com/dezh-tech/immortal/types/nip11"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool { return true },
	}

	//go:embed landing.html
	landingTempl []byte
)

// Server represents a websocket serer which keeps track of client connections and handle them.
type Server struct {
	mu sync.RWMutex

	config   Config
	conns    map[*websocket.Conn]clientState
	handlers *handler.Handler
	nip11Doc *nip11.RelayInformationDocument
	metrics  *metrics.Metrics
	redis    *redis.Redis
}

func New(cfg Config, nip11Doc *nip11.RelayInformationDocument,
	h *handler.Handler, m *metrics.Metrics, r *redis.Redis,
) (*Server, error) {
	return &Server{
		config:   cfg,
		conns:    make(map[*websocket.Conn]clientState),
		mu:       sync.RWMutex{},
		nip11Doc: nip11Doc,
		handlers: h,
		metrics:  m,
		redis:    r,
	}, nil
}

// Start starts a new server instance.
func (s *Server) Start() error {
	log.Println("websocket server started successfully...")

	http.Handle("/", s)
	err := http.ListenAndServe(net.JoinHostPort(s.config.Bind, //nolint
		strconv.Itoa(int(s.config.Port))), nil)

	return err
}

// handleWS is WebSocket handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") == "application/nostr+json" {
		w.Header().Set("Content-Type", "application/nostr+json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(s.nip11Doc) //nolint

		return
	}

	if r.Header.Get("Upgrade") == "" {
		t := template.New("webpage")
		t, err := t.Parse(string(landingTempl))
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)

			return
		}

		err = t.Execute(w, s.nip11Doc)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)

			return
		}

		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.mu.Lock()

	log.Println("new websocket connection: ", conn.RemoteAddr().String())
	s.metrics.Connections.Inc()

	s.conns[conn] = clientState{
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

		switch msg.Type() {
		case "REQ":
			go s.handleReq(conn, msg)

		case "EVENT":
			go s.handleEvent(conn, msg)

		case "CLOSE":
			go s.handleClose(conn, msg)
		}
	}
}

// handleReq handles new incoming REQ messages from client.
func (s *Server) handleReq(conn *websocket.Conn, m message.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer measureLatency(s.metrics.RequestLatency)()

	status := success
	defer func() {
		s.metrics.RequestsTotal.WithLabelValues(status).Inc()
	}()

	msg, ok := m.(*message.Req)
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice("error: can't parse REQ message."))

		status = parseFail

		return
	}

	if len(msg.Filters) >= s.config.Limitation.MaxFilters {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: max limit of filters is: %d",
			s.config.Limitation.MaxFilters)))

		status = limitsFail

		return
	}

	if s.config.Limitation.MaxSubidLength <= len(msg.SubscriptionID) {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: max limit of sub id is: %d",
			s.config.Limitation.MaxSubidLength)))

		status = limitsFail

		return
	}

	client, ok := s.conns[conn]
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't find connection %s",
			conn.RemoteAddr())))

		status = serverFail

		return
	}

	if len(client.subs) >= s.config.Limitation.MaxSubscriptions {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: max limit of subs is: %d",
			s.config.Limitation.MaxSubscriptions)))

		status = limitsFail

		return
	}

	res, err := s.handlers.HandleReq(msg.Filters)
	if err != nil {
		log.Println(err.Error())
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't process REQ message: %s", err.Error())))
		status = databaseFail

		return
	}

	for _, e := range res {
		msg := message.MakeEvent(msg.SubscriptionID, &e)
		_ = conn.WriteMessage(1, msg)
	}

	_ = conn.WriteMessage(1, message.MakeEOSE(msg.SubscriptionID))

	client.Lock()
	s.metrics.Subscriptions.Inc()
	client.subs[msg.SubscriptionID] = msg.Filters
	client.Unlock()
}

// handleEvent handles new incoming EVENT messages from client.
func (s *Server) handleEvent(conn *websocket.Conn, m message.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer measureLatency(s.metrics.EventLatency)()

	status := success
	defer func() {
		s.metrics.EventsTotal.WithLabelValues(status).Inc()
	}()

	msg, ok := m.(*message.Event)
	if !ok {
		okm := message.MakeOK(false,
			"",
			"error: can't parse EVENT message.",
		)

		_ = conn.WriteMessage(1, okm)
		status = parseFail

		return
	}

	if len(msg.Event.Content) > s.config.Limitation.MaxContentLength {
		okm := message.MakeOK(false,
			"",
			fmt.Sprintf("error: max limit of message length is %d", s.config.Limitation.MaxContentLength),
		)

		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	eID := msg.Event.GetRawID()

	qCtx, cancel := context.WithTimeout(context.Background(), s.redis.QueryTimeout)
	defer cancel()

	exists, err := s.redis.Client.BFExists(qCtx, s.redis.BloomName, eID[:]).Result()
	if err != nil {
		log.Printf("error: checking bloom filter: %s", err.Error())
	}

	if exists {
		okm := message.MakeOK(true, msg.Event.ID, "")
		_ = conn.WriteMessage(1, okm)

		return
	}

	if !msg.Event.IsValid(eID) {
		okm := message.MakeOK(false,
			msg.Event.ID,
			"invalid: id or sig is not correct.",
		)

		_ = conn.WriteMessage(1, okm)

		status = invalidFail

		return
	}

	if !msg.Event.Kind.IsEphemeral() {
		err := s.handlers.HandleEvent(msg.Event)
		if err != nil {
			okm := message.MakeOK(false,
				msg.Event.ID,
				err.Error(),
			)

			_ = conn.WriteMessage(1, okm)

			status = serverFail

			return
		}
		_ = conn.WriteMessage(1, message.MakeOK(true, msg.Event.ID, ""))
	}

	_, err = s.redis.Client.BFAdd(qCtx, s.redis.BloomName, eID[:]).Result()
	if err != nil {
		log.Printf("error: checking bloom filter: %s", err.Error())
	}

	// todo::: can we run goroutines per client?
	for conn, client := range s.conns {
		client.Lock()
		for id, filters := range client.subs {
			if !filters.Match(msg.Event) {
				continue
			}
			_ = conn.WriteMessage(1, message.MakeEvent(id, msg.Event))
		}
		client.Unlock()
	}
}

// handleClose handles new incoming CLOSE messages from client.
func (s *Server) handleClose(conn *websocket.Conn, m message.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := m.(*message.Close)
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice("error: can't parse CLOSE message."))

		return
	}

	client, ok := s.conns[conn]
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
			conn.RemoteAddr())))

		return
	}

	client.Lock()
	s.metrics.Subscriptions.Dec()
	delete(client.subs, msg.String())
	client.Unlock()
}

// Stop shutdowns the server gracefully.
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("stopping websocket server...")

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
