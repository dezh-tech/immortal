package server

import (
	"fmt"

	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/message"
	"github.com/dezh-tech/immortal/utils"
	"github.com/gorilla/websocket"
)

func (s *Server) handleAuth(conn *websocket.Conn, m message.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer measureLatency(s.metrics.AuthLatency)()

	status := success
	defer func() {
		s.metrics.AuthsTotal.WithLabelValues(status).Inc()
	}()

	msg, ok := m.(*message.Auth)
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice("error: can't parse AUTH message."))
		status = parseFail

		return
	}

	client, ok := s.conns[conn]
	if !ok {
		_ = conn.WriteMessage(1, message.MakeNotice(fmt.Sprintf("error: can't find connection %s.",
			conn.RemoteAddr())))
		status = serverFail

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

	relayURL, err := utils.ParseURL(relay)
	if err != nil {
		_ = conn.WriteMessage(1, message.MakeNotice("error: invalid auth event."))
		status = parseFail

		return
	}

	if !msg.Event.IsValid(msg.Event.GetRawID()) && msg.Event.Kind != types.KindClientAuthentication &&
		client.challenge != challenge && s.config.URL.Scheme != relayURL.Scheme || s.config.URL.Host != relayURL.Host ||
		s.config.URL.Path != relayURL.Path {
		_ = conn.WriteMessage(1, message.MakeNotice("error: invalid auth event."))
		status = invalidFail

		return
	}

	*client.isKnown = true
	*client.pubkey = msg.Event.PublicKey

	_ = conn.WriteMessage(1, message.MakeOK(true, msg.Event.ID, ""))
	status = success
}
