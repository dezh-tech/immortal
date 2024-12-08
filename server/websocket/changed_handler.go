package websocket

import (
	"github.com/dezh-tech/immortal/types/message"
	"github.com/gorilla/websocket"
)

// handleChanged handles new incoming CHANGED messages from client.
func (s *Server) handleChanged(conn *websocket.Conn, m message.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := m.(*message.Changed)
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice("error: can't parse CHANGED message."))

		return
	}

	// todo:::
}
