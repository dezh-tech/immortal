package server

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/dezh-tech/immortal/types/message"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

// Server represents a websocket serer which keeps track of client connections and handle them.
type Server struct {
	knownEvents *bloom.BloomFilter
	config      Config
	conns       map[*websocket.Conn]clientState
	mu          sync.RWMutex
}

func New(cfg Config) *Server {
	seb := bloom.NewWithEstimates(cfg.KnownBloomSize, 0.9)
	f, err := os.Open(cfg.BloomBackupPath)
	if err == nil {
		w := bufio.NewReader(f)
		_, _ = seb.ReadFrom(w)
	}

	return &Server{
		config:      cfg,
		knownEvents: seb,
		conns:       make(map[*websocket.Conn]clientState),
		mu:          sync.RWMutex{},
	}
}

// Start strats a new server instance.
func (s *Server) Start() error {
	http.Handle("/", s)
	err := http.ListenAndServe(net.JoinHostPort(s.config.Bind, //nolint
		strconv.Itoa(int(s.config.Port))), nil)

	return err
}

// handleWS is WebSocket handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.mu.Lock()
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
			delete(s.conns, conn)
			s.mu.Lock()

			break
		}

		msg := message.ParseMessage(buf)
		if msg == nil {
			_ = conn.WriteMessage(1, message.MakeNotice("error: can't parse message."))

			continue
		}

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
	// TODO::: loadfrom database and sent in first query based on limit.
	// TODO::: return EOSE.
	// TODO::: return EVENT messages.
	// TODO::: use mu properly.

	msg, ok := m.(*message.Req)
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice("error: can't parse REQ message."))

		return
	}

	s.mu.Lock()

	client, ok := s.conns[conn]
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
			conn.RemoteAddr())))

		return
	}

	client.RLock()
	client.subs[msg.SubscriptionID] = msg.Filters
	client.RUnlock()

	s.mu.Unlock()
}

// handleEvent handles new incoming EVENT messages from client.
func (s *Server) handleEvent(conn *websocket.Conn, m message.Message) {
	msg, ok := m.(*message.Event)
	if !ok {
		okm := message.MakeOK(false,
			"",
			"error: can't parse EVENT message.",
		)

		_ = conn.WriteMessage(1, okm)

		return
	}

	eID := msg.Event.GetRawID()

	s.mu.RLock()

	if s.knownEvents.Test(eID[:]) {
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

		return
	}

	if !msg.Event.Kind.IsEphemeral() {
		// send to be written in database.
		return // TODO::: remove me.
	}

	s.knownEvents.Add(eID[:])

	_ = conn.WriteMessage(1, message.MakeOK(true, msg.Event.ID, ""))

	for conn, client := range s.conns {
		go func() {
			client.Lock()
			for id, filters := range client.subs {
				if !filters.Match(msg.Event) {
					return
				}
				_ = conn.WriteMessage(1, message.MakeEvent(id, msg.Event))
			}
			client.Unlock()
		}()
	}

	s.mu.RUnlock()
}

// handleClose handles new incoming CLOSE messages from client.
func (s *Server) handleClose(ws *websocket.Conn, m message.Message) {
	msg, ok := m.(*message.Close)
	if !ok {
		_ = ws.WriteMessage(1, message.MakeNotice("error: can't parse CLOSE message."))

		return
	}

	s.mu.Lock()

	client, ok := s.conns[ws]
	if !ok {
		_ = ws.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
			ws.RemoteAddr())))

		return
	}

	client.Lock()
	delete(client.subs, msg.String())
	client.Unlock()

	s.mu.Unlock()
}

// Stop shutdowns the server gracefully.
func (s *Server) Stop() error {
	s.mu.Lock()

	for wsConn, client := range s.conns {
		client.Lock()
		// close all subscriptions.
		for id := range client.subs {
			delete(client.subs, id)

			err := wsConn.WriteMessage(1, message.MakeClosed(id, "error: shutdowning the relay."))
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

	f, err := os.Create(s.config.BloomBackupPath)
	if err != nil {
		return fmt.Errorf("error: creating new file for blooms: %s", err.Error())
	}

	w := bufio.NewWriter(f)
	_, err = s.knownEvents.WriteTo(w)
	if err != nil {
		return fmt.Errorf("error: writing bloom filter to disck: %s", err.Error())
	}

	s.mu.Unlock()

	return nil
}
