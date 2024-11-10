package websocket

import (
	"context"
	"fmt"
	"log"

	"github.com/dezh-tech/immortal/types/message"
	"github.com/dezh-tech/immortal/utils"
	"github.com/gorilla/websocket"
)

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

	client, ok := s.conns[conn]
	if !ok {
		_ = conn.WriteMessage(1, message.MakeOK(false,
			"",
			fmt.Sprintf("error: can't find connection %s",
				conn.RemoteAddr())))

		status = serverFail

		return
	}

	if s.config.Limitation.AuthRequired && !*client.isKnown {
		client.challenge = utils.GenerateChallenge(10)
		authm := message.MakeAuth(client.challenge)

		okm := message.MakeOK(false,
			"",
			"auth-required: we only accept events from authenticated users.",
		)

		_ = conn.WriteMessage(1, okm)

		_ = conn.WriteMessage(1, authm)
		status = authFail

		return
	}

	if msg.Event.IsProtected() && msg.Event.PublicKey != *client.pubkey {
		client.challenge = utils.GenerateChallenge(10)
		authm := message.MakeAuth(client.challenge)

		okm := message.MakeOK(false,
			"",
			"auth-required: this event may only be published by its author.",
		)

		_ = conn.WriteMessage(1, okm)

		_ = conn.WriteMessage(1, authm)
		status = authFail

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

	if len(msg.Event.Content) > s.config.Limitation.MaxContentLength {
		okm := message.MakeOK(false,
			"",
			fmt.Sprintf("error: max limit of content length is %d", s.config.Limitation.MaxContentLength),
		)

		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	if msg.Event.Difficulty() < s.config.Limitation.MinPowDifficulty {
		okm := message.MakeOK(false,
			"",
			fmt.Sprintf("error: min pow required is %d", s.config.Limitation.MinPowDifficulty),
		)

		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	if len(msg.Event.Tags) < s.config.Limitation.MaxEventTags {
		okm := message.MakeOK(false,
			"",
			fmt.Sprintf("error: max limit of tags count is %d", s.config.Limitation.MaxEventTags),
		)

		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	if msg.Event.CreatedAt < s.config.Limitation.CreatedAtLowerLimit ||
		msg.Event.CreatedAt > s.config.Limitation.CreatedAtUpperLimit {
		okm := message.MakeOK(false,
			"",
			fmt.Sprintf("error: created at must be as least %d and at most %d",
				s.config.Limitation.CreatedAtLowerLimit, s.config.Limitation.CreatedAtUpperLimit),
		)

		_ = conn.WriteMessage(1, okm)

		status = limitsFail

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
