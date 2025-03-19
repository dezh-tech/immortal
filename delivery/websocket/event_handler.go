package websocket

import (
	"fmt"
	"github.com/dezh-tech/immortal/delivery/websocket/configs"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dezh-tech/immortal/infrastructure/redis"
	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/pkg/utils"
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/dezh-tech/immortal/types/message"
	"github.com/gorilla/websocket"
)

// handleEvent handles new incoming EVENT messages from client.
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

	if !msg.Event.IsValid(eID) {
		okm := message.MakeOK(false,
			msg.Event.ID,
			"invalid: id or sig is not correct.",
		)

		_ = conn.WriteMessage(1, okm)

		status = invalidFail

		return
	}

	if err := s.redis.CheckAcceptability(s.config.GetLimitation().RestrictedWrites,
		eID[:], msg.Event.PublicKey); err != nil {
		okm := message.MakeOK(false, msg.Event.ID, err.Error())
		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		if strings.HasPrefix(err.Error(), "internal") {
			status = serverFail
		}

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

	accepted, authFail, failType, resp := checkLimitations(client, s.redis, *s.config.GetLimitation(), *msg)
	if !accepted && authFail {
		client.challenge = utils.GenerateChallenge(10)
		authm := message.MakeAuth(client.challenge)

		okm := message.MakeOK(false,
			msg.Event.ID,
			resp,
		)

		_ = conn.WriteMessage(1, okm)

		_ = conn.WriteMessage(1, authm)
		status = failType

		return
	}

	if !accepted {
		okm := message.MakeOK(false,
			msg.Event.ID,
			resp,
		)

		_ = conn.WriteMessage(1, okm)

		status = failType

		return
	}

	if !msg.Event.Kind.IsEphemeral() { //nolint
		if msg.Event.Kind == types.KindEventDeletionRequest {
			if err := s.handler.NIP09Deletion(msg.Event); err != nil {
				okm := message.MakeOK(false,
					msg.Event.ID,
					"error: can't delete requested event(s).",
				)

				_ = conn.WriteMessage(1, okm)

				status = serverFail

				return
			}
		}

		if msg.Event.Kind == types.KindRightToVanish {
			relays := msg.Event.Tags.GetValues("relay")
			if !slices.Contains(relays, "ALL_RELAYS") &&
				!slices.Contains(relays, s.config.GetURL().String()) {
				okm := message.MakeOK(false,
					msg.Event.ID,
					"error: can't execute vanish request.",
				)

				_ = conn.WriteMessage(1, okm)

				status = serverFail

				return
			}

			if err := s.handler.DeleteByFilter(
				&filter.Filter{Authors: []string{msg.Event.PublicKey}, Until: msg.Event.CreatedAt}); err != nil {
				okm := message.MakeOK(false,
					msg.Event.ID,
					"error: can't execute vanish request.",
				)

				_ = conn.WriteMessage(1, okm)

				status = serverFail

				return
			}

			if err := s.handler.DeleteByFilter(
				&filter.Filter{Kinds: []types.Kind{types.KindGiftWrap}, Tags: map[string][]string{
					"p": {msg.Event.PublicKey},
				}, Until: msg.Event.CreatedAt}); err != nil {
				okm := message.MakeOK(false,
					msg.Event.ID,
					"error: can't delete requested event(s).",
				)

				_ = conn.WriteMessage(1, okm)

				status = serverFail

				return
			}
		}

		if err := s.handler.HandleEvent(msg.Event); err != nil {
			okm := message.MakeOK(false,
				msg.Event.ID,
				"error: can't write event to database.",
			)

			_ = conn.WriteMessage(1, okm)

			status = serverFail

			return
		}
	}

	_ = conn.WriteMessage(1, message.MakeOK(true, msg.Event.ID, ""))

	if err := s.redis.AddEventToBloom(eID[:]); err != nil {
		logger.Info("adding event to bloom filter", "err", err.Error(), msg.Event.ID)
	}

	// todo::: can we run goroutines per client?
	for conn, client := range s.conns {
		client.Lock()
		for id, filter := range client.subs {
			if !filter.Match(msg.Event, *client.pubkey) {
				continue
			}
			_ = conn.WriteMessage(1, message.MakeEvent(id, msg.Event))
		}
		client.Unlock()
	}
}

func checkLimitations(c clientState, r *redis.Redis,
	limits configs.Limitation, msg message.Event) (bool, bool,
	string, string,
) {
	if limits.AuthRequired && !*c.isKnown {
		return false, true, authFail, "auth-required: we only accept events from authenticated users."
	}

	if msg.Event.IsProtected() && msg.Event.PublicKey != *c.pubkey {
		return false, true, authFail, "auth-required: this event may only be published by its author."
	}

	expirationTag := msg.Event.Tags.GetValue("expiration")

	if expirationTag != "" {
		expiration, err := strconv.ParseInt(expirationTag, 10, 64)
		if err != nil {
			return false, false, serverFail, fmt.Sprintf("invalid: expiration tag %s.", expirationTag)
		}

		if time.Now().Unix() >= expiration {
			return false, false, invalidFail, fmt.Sprintf("invalid: this event was expired in %s.",
				time.Unix(expiration, 0).String())
		}

		if err := r.AddDelayedTask(expirationTaskListName,
			fmt.Sprintf("%s:%d", msg.Event.ID, msg.Event.Kind), time.Until(time.Unix(expiration, 0))); err != nil {
			return false, false, serverFail, "error: can't add event to expiration queue."
		}
	}

	if len(msg.Event.Content) > int(limits.MaxContentLength) {
		return false, false, limitsFail, fmt.Sprintf("error: max limit of content length is %d", limits.MaxContentLength)
	}

	if msg.Event.Difficulty() < int(limits.MinPowDifficulty) {
		return false, false, limitsFail, fmt.Sprintf("error: min pow required is %d", limits.MinPowDifficulty)
	}

	if len(msg.Event.Tags) > int(limits.MaxEventTags) {
		return false, false, limitsFail, fmt.Sprintf("error: max limit of tags count is %d", limits.MaxEventTags)
	}

	if msg.Event.CreatedAt < limits.CreatedAtLowerLimit ||
		msg.Event.CreatedAt > limits.CreatedAtUpperLimit {
		return false, false, limitsFail, fmt.Sprintf("error: created at must be as least %d and at most %d",
			limits.CreatedAtLowerLimit, limits.CreatedAtUpperLimit)
	}

	return true, false, "", ""
}
