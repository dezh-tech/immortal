package websocket

import (
	"fmt"

	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/message"
	"github.com/gorilla/websocket"
)

func (s *Server) handleAuth(conn *websocket.Conn, m message.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := m.(*message.Auth)
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice("error: can't parse AUTH message."))

		return
	}

	client, ok := s.conns[conn]
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
			conn.RemoteAddr())))

		return
	}

	challenge := ""
	relay := ""

	for _, t := range msg.Event.Tags {
		if len(t) < 2 {
			continue
		}

		if t[0] == "challenge" {
			challenge = t[1]
		}

		if t[0] == "relay" {
			relay = t[1]
		}
	}

	if !msg.Event.IsValid(msg.Event.GetRawID()) && msg.Event.Kind != types.KindClientAuthentication &&
		client.challenge != challenge && s.nip11Doc.URL != relay {
		_ = conn.WriteMessage(1, message.MakeNotice("error: invalid auth event."))

		return
	}

	*client.isKnown = true
	*client.pubkey = msg.Event.PublicKey

	_ = conn.WriteMessage(1, message.MakeOK(true, msg.Event.ID, ""))
}
