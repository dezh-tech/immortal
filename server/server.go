package server

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/dezh-tech/immortal/types/filter"
	"github.com/dezh-tech/immortal/types/message"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

// Server represents a websocket serer which keeps track of client connections and handle them.
type Server struct {
	config    Config
	conns     map[*websocket.Conn]map[string]filter.Filters
	connsLock sync.RWMutex
}

func New(cfg Config) *Server {
	return &Server{
		config:    cfg,
		conns:     make(map[*websocket.Conn]map[string]filter.Filters),
		connsLock: sync.RWMutex{},
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

	s.connsLock.Lock()
	s.conns[conn] = make(map[string]filter.Filters)
	s.connsLock.Unlock()

	s.readLoop(conn)
}

// readLoop reads incoming messages from a client and answer to them.
func (s *Server) readLoop(ws *websocket.Conn) {
	for {
		_, buf, err := ws.ReadMessage()
		if err != nil {
			break
		}

		msg := message.ParseMessage(buf)
		if msg == nil {
			_ = ws.WriteMessage(1, message.MakeNotice("error: can't parse message."))

			continue
		}

		switch msg.Type() {
		case "REQ":
			go s.handleReq(ws, msg)

		case "EVENT":
			go s.handleEvent(ws, msg)

		case "CLOSE":
			go s.handleClose(ws, msg)
		}
	}
}

// handleReq handles new incoming REQ messages from client.
func (s *Server) handleReq(ws *websocket.Conn, m message.Message) {
	// TODO::: loadfrom database and sent in first query based on limit.
	// TODO::: return EOSE.
	// TODO::: return EVENT messages.

	msg, ok := m.(*message.Req)
	if !ok {
		_ = ws.WriteMessage(1, message.MakeNotice("error: can't parse REQ message."))

		return
	}

	s.connsLock.Lock()
	defer s.connsLock.Unlock()

	subs, ok := s.conns[ws]
	if !ok {
		_ = ws.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
			ws.RemoteAddr())))

		return
	}

	subs[msg.SubscriptionID] = msg.Filters
}

// handleEvent handles new incoming EVENT messages from client.
func (s *Server) handleEvent(ws *websocket.Conn, m message.Message) {
	msg, ok := m.(*message.Event)
	if !ok {
		okm := message.MakeOK(false,
			"",
			"error: can't parse EVENT message.",
		)

		_ = ws.WriteMessage(1, okm)

		return
	}

	if !msg.Event.IsValid() {
		okm := message.MakeOK(false,
			msg.Event.ID,
			"invalid: id or sig is not correct.",
		)

		_ = ws.WriteMessage(1, okm)

		return
	}

	_ = ws.WriteMessage(1, message.MakeOK(true, msg.Event.ID, ""))

	for conn, subs := range s.conns {
		for id, filters := range subs {
			if !filters.Match(msg.Event) {
				return
			}
			_ = conn.WriteMessage(1, message.MakeEvent(id, msg.Event))
		}
	}
}

// handleClose handles new incoming CLOSE messages from client.
func (s *Server) handleClose(ws *websocket.Conn, m message.Message) {
	msg, ok := m.(*message.Close)
	if !ok {
		_ = ws.WriteMessage(1, message.MakeNotice("error: can't parse CLOSE message."))

		return
	}

	s.connsLock.Lock()
	defer s.connsLock.Unlock()

	conn, ok := s.conns[ws]
	if !ok {
		_ = ws.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
			ws.RemoteAddr())))

		return
	}

	delete(conn, msg.String())
	_ = ws.WriteMessage(1, message.MakeClosed(msg.String(), "ok: closed successfully."))
}

// Stop shutdowns the server gracefully.
func (s *Server) Stop() error {
	s.connsLock.Lock()
	defer s.connsLock.Unlock()

	for wsConn, subs := range s.conns {
		// close all subscriptions.
		for id := range subs {
			_ = wsConn.WriteMessage(1, message.MakeClosed(id, "error: shutdowning the relay."))
		}

		// close connection.
		_ = wsConn.Close()
	}

	return nil
}
