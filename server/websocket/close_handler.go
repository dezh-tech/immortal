package websocket

import (
	"fmt"

	"github.com/dezh-tech/immortal/types/message"
	"github.com/gorilla/websocket"
)

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
