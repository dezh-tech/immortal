package websocket

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dezh-tech/immortal/types/message"
	"github.com/dezh-tech/immortal/utils"
	"github.com/gorilla/websocket"
)

// handleEvent handles new incoming EVENT messages from client.
// todo::: too much complexity.
func (s *Server) handleEvent(conn *websocket.Conn, m message.Message) { //nolint
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

	eID := msg.Event.GetRawID()
	pubkey := msg.Event.PublicKey

	qCtx, cancel := context.WithTimeout(context.Background(), s.redis.QueryTimeout)
	defer cancel()

	pipe := s.redis.Client.Pipeline()

	bloomCheckCmd := pipe.BFExists(qCtx, s.redis.BloomFilterName, eID[:])

	// TODO::: check config to enable filter checks
	whiteListCheckCmd := pipe.CFExists(qCtx, s.redis.WhiteListFilterName, pubkey)
	blackListCheckCmd := pipe.CFExists(qCtx, s.redis.BlackListFilterName, pubkey)

	_, err := pipe.Exec(qCtx)
	if err != nil {
		log.Printf("error: checking filters: %s", err.Error())
	}

	exists, err := bloomCheckCmd.Result()
	if err != nil {
		okm := message.MakeOK(false, msg.Event.ID, "error: internal error")
		_ = conn.WriteMessage(1, okm)

		status = serverFail

		return
	}
	if exists {
		okm := message.MakeOK(true, msg.Event.ID, "")
		_ = conn.WriteMessage(1, okm)

		return
	}

	notAllowedToWrite, err := blackListCheckCmd.Result()
	if err != nil {
		okm := message.MakeOK(false, msg.Event.ID, "error: internal error")
		_ = conn.WriteMessage(1, okm)

		status = serverFail

		return
	}
	if notAllowedToWrite {
		okm := message.MakeOK(false, msg.Event.ID, "blocked: pubkey is blocked, contact support for more details.")
		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	allowedToWrite, err := whiteListCheckCmd.Result()
	if err != nil {
		okm := message.MakeOK(false, msg.Event.ID, "error: internal error")
		_ = conn.WriteMessage(1, okm)

		status = serverFail

		return
	}
	if !allowedToWrite {
		okm := message.MakeOK(false, msg.Event.ID, "restricted: not allowed to write.")
		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	client, ok := s.conns[conn]
	if !ok {
		_ = conn.WriteMessage(1, message.MakeOK(false,
			msg.Event.ID,
			fmt.Sprintf("error: can't find connection %s",
				conn.RemoteAddr())))

		status = serverFail

		return
	}

	if s.config.Limitation.AuthRequired && !*client.isKnown {
		client.challenge = utils.GenerateChallenge(10)
		authm := message.MakeAuth(client.challenge)

		okm := message.MakeOK(false,
			msg.Event.ID,
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
			msg.Event.ID,
			"auth-required: this event may only be published by its author.",
		)

		_ = conn.WriteMessage(1, authm)

		_ = conn.WriteMessage(1, okm)

		status = authFail

		return
	}

	expirationTag := msg.Event.Tags.GetValue("expiration")

	if expirationTag != "" {
		expiration, err := strconv.ParseInt(expirationTag, 10, 64)
		if err != nil {
			okm := message.MakeOK(false,
				msg.Event.ID,
				fmt.Sprintf("invalid: expiration tag %s.", expirationTag),
			)

			_ = conn.WriteMessage(1, okm)

			status = invalidFail

			return
		}

		if time.Now().Unix() >= expiration {
			okm := message.MakeOK(false,
				msg.Event.ID,
				fmt.Sprintf("invalid: this event was expired in %s.", time.Unix(expiration, 0).String()),
			)

			_ = conn.WriteMessage(1, okm)

			status = invalidFail

			return
		}

		if err := s.redis.AddDelayedTask("expiration_events",
			fmt.Sprintf("%s:%d", msg.Event.ID, msg.Event.Kind), time.Until(time.Unix(expiration, 0))); err != nil {
			okm := message.MakeOK(false,
				msg.Event.ID, "error: can't add event to expiration queue.",
			)

			_ = conn.WriteMessage(1, okm)

			status = invalidFail

			return
		}
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

	if len(msg.Event.Content) > int(s.config.Limitation.MaxContentLength) {
		okm := message.MakeOK(false,
			msg.Event.ID,
			fmt.Sprintf("error: max limit of content length is %d", s.config.Limitation.MaxContentLength),
		)

		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	if msg.Event.Difficulty() < int(s.config.Limitation.MinPowDifficulty) {
		okm := message.MakeOK(false,
			msg.Event.ID,
			fmt.Sprintf("error: min pow required is %d", s.config.Limitation.MinPowDifficulty),
		)

		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	if len(msg.Event.Tags) > int(s.config.Limitation.MaxEventTags) {
		okm := message.MakeOK(false,
			msg.Event.ID,
			fmt.Sprintf("error: max limit of tags count is %d", s.config.Limitation.MaxEventTags),
		)

		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	if msg.Event.CreatedAt < s.config.Limitation.CreatedAtLowerLimit ||
		msg.Event.CreatedAt > s.config.Limitation.CreatedAtUpperLimit {
		okm := message.MakeOK(false,
			msg.Event.ID,
			fmt.Sprintf("error: created at must be as least %d and at most %d",
				s.config.Limitation.CreatedAtLowerLimit, s.config.Limitation.CreatedAtUpperLimit),
		)

		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	if !msg.Event.Kind.IsEphemeral() {
		err := s.handler.HandleEvent(msg.Event)
		if err != nil {
			okm := message.MakeOK(false,
				msg.Event.ID,
				"error: can't write event to database.",
			)

			_ = conn.WriteMessage(1, okm)

			status = serverFail

			return
		}
		_ = conn.WriteMessage(1, message.MakeOK(true, msg.Event.ID, ""))
	}

	_, err = s.redis.Client.BFAdd(qCtx, s.redis.BloomFilterName, eID[:]).Result()
	if err != nil {
		log.Printf("error: adding event to bloom filter.")
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
