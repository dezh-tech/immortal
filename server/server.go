package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/dezh-tech/immortal/types/filter"
	"github.com/dezh-tech/immortal/types/message"
	"golang.org/x/net/websocket"
)

// TODO::: replace with https://github.com/coder/websocket.
// TODO::: replace `log` with main logger.

// Server represents a nostr relay which keeps track of client connections and handle them.
type Server struct {
	config    Config
	conns     map[*websocket.Conn]map[string]filter.Filters
	connsLock sync.RWMutex
}

func NewServer(cfg Config) *Server {
	return &Server{
		config:    cfg,
		conns:     make(map[*websocket.Conn]map[string]filter.Filters),
		connsLock: sync.RWMutex{},
	}
}

// Start strats a new server instance.
func (s *Server) Start() error {
	http.Handle("/", websocket.Handler(s.handleWS))
	err := http.ListenAndServe(net.JoinHostPort(s.config.Bind,
		strconv.Itoa(int(s.config.Port))), nil) //nolint

	return err
}

// handleWS is WebSocket handler.
func (s *Server) handleWS(ws *websocket.Conn) {
	s.connsLock.Lock()
	s.conns[ws] = make(map[string]filter.Filters)
	s.connsLock.Unlock()

	s.readLoop(ws)
}

// readLoop reads incoming messages from a client and answer to them.
func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			continue
		}

		msg := message.ParseMessage(buf[:n])
		if msg == nil {
			_, _ = ws.Write(message.MakeNotice("error: can't parse message."))

			continue
		}

		switch msg.Type() {
		case "REQ":
			go s.handleReq(ws, msg)

		case "EVENT":
			go s.handleEvent(ws, msg)

		case "CLOSE":
			s.handleClose(ws, msg)
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
		_, _ = ws.Write(message.MakeNotice("error: can't parse REQ message."))

		return
	}

	s.connsLock.Lock()
	defer s.connsLock.Unlock()

	subs, ok := s.conns[ws]
	if !ok {
		_, _ = ws.Write(message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
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

		_, _ = ws.Write(okm)

		return
	}

	if !msg.Event.IsValid() {
		okm := message.MakeOK(false,
			msg.SubscriptionID,
			"invalid: id or sig is not correct.",
		)

		_, _ = ws.Write(okm)

		return
	}

	_, _ = ws.Write(message.MakeOK(true, msg.SubscriptionID, ""))

	for conn, subs := range s.conns {
		for id, filters := range subs {
			if !filters.Match(msg.Event) {
				return
			}
			_, _ = conn.Write(message.MakeEvent(id, msg.Event))
		}
	}
}

// handleClose handles new incoming CLOSE messages from client.
func (s *Server) handleClose(ws *websocket.Conn, m message.Message) {
	msg, ok := m.(*message.Close)
	if !ok {
		_, _ = ws.Write(message.MakeNotice("error: can't parse CLOSE message."))

		return
	}

	s.connsLock.Lock()
	defer s.connsLock.Unlock()

	conn, ok := s.conns[ws]
	if !ok {
		_, _ = ws.Write(message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
			ws.RemoteAddr())))

		return
	}

	delete(conn, msg.String())
	_, _ = ws.Write(message.MakeClosed(msg.String(), "ok: closed successfully."))
}

// Stop shutdowns the server gracefully.
func (s *Server) Stop() error {
	s.connsLock.Lock()
	defer s.connsLock.Unlock()

	for wsConn, subs := range s.conns {
		// close all subscriptions.
		for id := range subs {
			_, _ = wsConn.Write(message.MakeClosed(id, "error: shutdowning the relay."))
		}

		// close connection.
		_ = wsConn.Close()
	}

	return nil
}
