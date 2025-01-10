package websocket

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dezh-tech/immortal/infrastructure/redis"
	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/pkg/utils"
	"github.com/dezh-tech/immortal/types/message"
	"github.com/gorilla/websocket"
	gredis "github.com/redis/go-redis/v9"
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

	qCtx, cancel := context.WithTimeout(context.Background(), s.redis.QueryTimeout)
	defer cancel()

	pipe := s.redis.Client.Pipeline()

	bloomCheckCmd := pipe.BFExists(qCtx, s.redis.BloomFilterName, eID[:])

	var whiteListCheckCmd *gredis.BoolCmd

	if s.config.Limitation.RestrictedWrites {
		whiteListCheckCmd = pipe.CFExists(qCtx, s.redis.WhiteListFilterName, msg.Event.PublicKey)
	}

	blackListCheckCmd := pipe.CFExists(qCtx, s.redis.BlackListFilterName, msg.Event.PublicKey)

	_, err := pipe.Exec(qCtx)
	if err != nil {
		logger.Error("checking bloom filter", "err", err.Error())
	}

	exists, err := bloomCheckCmd.Result()
	if err != nil {
		okm := message.MakeOK(false, msg.Event.ID, "error: internal error")
		_ = conn.WriteMessage(1, okm)

		status = serverFail

		return
	}

	if exists {
		okm := message.MakeOK(false, msg.Event.ID, "duplicate: this event is already received.")
		_ = conn.WriteMessage(1, okm)

		return
	}

	isBlackListed, err := blackListCheckCmd.Result()
	if err != nil {
		okm := message.MakeOK(false, msg.Event.ID, "error: internal error")
		_ = conn.WriteMessage(1, okm)

		status = serverFail

		return
	}
	if isBlackListed {
		okm := message.MakeOK(false, msg.Event.ID, "blocked: pubkey is blocked, contact support for more details.")
		_ = conn.WriteMessage(1, okm)

		status = limitsFail

		return
	}

	if s.config.Limitation.RestrictedWrites {
		isWhiteListed, err := whiteListCheckCmd.Result()
		if err != nil {
			okm := message.MakeOK(false, msg.Event.ID, "error: internal error")
			_ = conn.WriteMessage(1, okm)

			status = serverFail

			return
		}

		if !isWhiteListed {
			okm := message.MakeOK(false, msg.Event.ID, "restricted: not allowed to write.")
			_ = conn.WriteMessage(1, okm)

			status = limitsFail

			return
		}
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

	accepted, authFail, failType, resp := checkLimitations(client, s.redis, *s.config.Limitation, *msg)
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
	}

	_ = conn.WriteMessage(1, message.MakeOK(true, msg.Event.ID, ""))

	_, err = s.redis.Client.BFAdd(qCtx, s.redis.BloomFilterName, eID[:]).Result()
	if err != nil {
		logger.Info("adding event to bloom filter", "err", err.Error(), msg.Event.ID)
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

func checkLimitations(c clientState, r *redis.Redis,
	limits Limitation, msg message.Event) (bool, bool,
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
